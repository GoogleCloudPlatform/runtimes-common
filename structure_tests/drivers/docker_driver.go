// Copyright 2017 Google Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package drivers

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/types/unversioned"
	"github.com/fsouza/go-dockerclient"
)

type DockerDriver struct {
	cli           docker.Client
	originalImage string
	currentImage  string
	env           map[string]string
}

func NewDockerDriver(image string) Driver {
	newCli, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	return &DockerDriver{
		originalImage: image,
		currentImage:  image,
		cli:           *newCli,
		env:           nil,
	}
}

func (d *DockerDriver) Setup(t *testing.T, envVars []unversioned.EnvVar, fullCommands []unversioned.Command) {
	env := d.processEnvVars(t, envVars)
	for _, cmd := range fullCommands {
		d.currentImage = d.runAndCommit(t, env, cmd)
	}
}

func (d *DockerDriver) ProcessCommand(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string) (string, string, int) {
	var env []string
	for _, envVar := range envVars {
		env = append(env, envVar.Key+"="+envVar.Value)
	}
	stdout, stderr, exitCode := d.exec(t, env, fullCommand)

	if stdout != "" {
		t.Logf("stdout: %s", stdout)
	}
	if stderr != "" {
		t.Logf("stderr: %s", stderr)
	}
	return stdout, stderr, exitCode
}

func (d *DockerDriver) retrieveEnvVar(envVar string) string {
	// since we're only retrieving these during processing, we can cache this on the driver
	if d.env == nil {
		image, err := d.cli.InspectImage(d.currentImage)
		if err != nil {
			return ""
		}

		// convert env to map for processing
		imageEnv := make(map[string]string)
		for _, varPair := range image.Config.Env {
			pair := strings.Split(varPair, "=")
			imageEnv[pair[0]] = pair[1]
		}
		d.env = imageEnv
	}

	return d.env[envVar]
}

func (d *DockerDriver) processEnvVars(t *testing.T, vars []unversioned.EnvVar) []string {
	if len(vars) == 0 {
		return nil
	}

	env := []string{}

	for _, envVar := range vars {
		expandedVal := os.Expand(envVar.Value, d.retrieveEnvVar)
		env = append(env, envVar.Key+"="+expandedVal)
	}
	return env
}

// copies a tar archive starting at the specified path from the image, and returns
// a tar reader which can be used to iterate through its contents and retrieve metadata
func (d *DockerDriver) retrieveTar(t *testing.T, path string) (*tar.Reader, error) {
	// this contains a placeholder command which does not get run, since
	// the client doesn't allow creating a container without a command.
	container, err := d.cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: d.currentImage,
			Cmd:   []string{"NOOP_COMMAND_DO_NOT_RUN"},
		},
		HostConfig:       nil,
		NetworkingConfig: nil,
	})
	if err != nil {
		t.Errorf("Error creating container: %s", err.Error())
		return nil, err
	}

	var b bytes.Buffer
	stream := bufio.NewWriter(&b)

	err = d.cli.DownloadFromContainer(container.ID, docker.DownloadFromContainerOptions{
		OutputStream: stream,
		Path:         path,
	})
	if err != nil {
		t.Errorf("Error when downloading file from container: %s", err.Error())
		return nil, err
	}
	stream.Flush()
	return tar.NewReader(bytes.NewReader(b.Bytes())), nil
}

func (d *DockerDriver) StatFile(t *testing.T, target string) (os.FileInfo, error) {
	reader, err := d.retrieveTar(t, target)
	if err != nil {
		return nil, err
	}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		switch header.Typeflag {
		case tar.TypeDir, tar.TypeReg:
			if filepath.Clean(header.Name) == path.Base(target) {
				return header.FileInfo(), nil
			}
		default:
			continue
		}
	}
	return nil, fmt.Errorf("File %s not found in image", target)
}

func (d *DockerDriver) ReadFile(t *testing.T, target string) ([]byte, error) {
	reader, err := d.retrieveTar(t, target)
	if err != nil {
		return nil, err
	}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if filepath.Clean(header.Name) == path.Base(target) {
				return nil, fmt.Errorf("Cannot read specified path: %s is a directory, not a file", target)
			}
		case tar.TypeReg:
			if filepath.Clean(header.Name) == path.Base(target) {
				var b bytes.Buffer
				stream := bufio.NewWriter(&b)
				io.Copy(stream, reader)
				return b.Bytes(), nil
			}
		default:
			continue
		}
	}
	return nil, fmt.Errorf("File %s not found in image", target)
}

func (d *DockerDriver) ReadDir(t *testing.T, target string) ([]os.FileInfo, error) {
	reader, err := d.retrieveTar(t, target)
	if err != nil {
		return nil, err
	}
	var infos []os.FileInfo
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if header.Typeflag == tar.TypeDir {
			// we only want top level dirs here, no recursion. to get these, remove
			// trailing separator and split on separator. there should only be two parts.
			parts := strings.Split(strings.TrimSuffix(header.Name, string(os.PathSeparator)), string(os.PathSeparator))
			if len(parts) == 2 {
				infos = append(infos, header.FileInfo())
			}
		}
	}
	return infos, nil
}

// This method takes a command (in the form of a list of args), and does the following:
// 1) creates a container, based on the "current latest" image, with the command set as
// the command to run when the container starts
// 2) starts the container
// 3) commits the container with its changes to a new image,
// and sets that image as the new "current image"
func (d *DockerDriver) runAndCommit(t *testing.T, env []string, command []string) string {
	container, err := d.cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        d.currentImage,
			Env:          env,
			Cmd:          command,
			AttachStdout: true,
			AttachStderr: true,
		},
		HostConfig:       nil,
		NetworkingConfig: nil,
	})
	if err != nil {
		t.Errorf("Error creating container: %s", err.Error())
		return ""
	}

	err = d.cli.StartContainer(container.ID, nil)
	if err != nil {
		t.Errorf("Error creating container: %s", err.Error())
	}

	_, err = d.cli.WaitContainer(container.ID)

	image, err := d.cli.CommitContainer(docker.CommitContainerOptions{
		Container: container.ID,
	})

	if err != nil {
		t.Errorf("Error committing container: %s", err.Error())
	}

	d.currentImage = image.ID
	return image.ID
}

func (d *DockerDriver) exec(t *testing.T, env []string, command []string) (string, string, int) {
	// first, start container from the current image
	container, err := d.cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        d.currentImage,
			Env:          env,
			Cmd:          command,
			AttachStdout: true,
			AttachStderr: true,
		},
		HostConfig:       nil,
		NetworkingConfig: nil,
	})
	if err != nil {
		t.Errorf("Error creating container: %s", err.Error())
		return "", "", -1
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	if err = d.cli.StartContainer(container.ID, nil); err != nil {
		t.Errorf("Error creating container: %s", err.Error())
	}

	//TODO(nkubala): look into adding timeout
	exitCode, err := d.cli.WaitContainer(container.ID)
	if err != nil {
		t.Errorf("Error when waiting for container: %s", err.Error())
	}

	if err = d.cli.Logs(docker.LogsOptions{
		Container:    container.ID,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Stdout:       true,
		Stderr:       true,
	}); err != nil {
		t.Errorf("Error retrieving container logs: %s", err.Error())
	}

	return stdout.String(), stderr.String(), exitCode
}

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
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
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

func (d *DockerDriver) Setup(t *testing.T, envVars []unversioned.EnvVar, fullCommand []unversioned.Command) {
	env := d.processEnvVars(t, envVars)
	for _, cmd := range fullCommand {
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

// copies a file from a docker container to the local fs, and returns its path
// caller is responsible for removing this file when finished
func (d *DockerDriver) retrieveFile(t *testing.T, path string) (string, error) {
	// this contains a placeholder command which does not get run, since
	// the client doesn't allow creating a container without a command.
	container, err := d.cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        d.currentImage,
			Cmd:          []string{"NOOP_COMMAND_DO_NOT_RUN"},
			AttachStdout: true,
			AttachStderr: true,
		},
		HostConfig:       nil,
		NetworkingConfig: nil,
	})
	if err != nil {
		t.Errorf("Error creating container: %s", err.Error())
		return "", err
	}

	tmpFile, err := ioutil.TempFile("", "structure_test")
	if err != nil {
		t.Errorf("Error when creating temp file: %s", err.Error())
		return "", err
	}
	stream := bufio.NewWriter(tmpFile)

	err = d.cli.DownloadFromContainer(container.ID, docker.DownloadFromContainerOptions{
		OutputStream: stream,
		Path:         path,
	})
	if err != nil {
		t.Errorf("Error when downloading file from container: %s", err.Error())
		return "", err
	}
	stream.Flush()
	tmpFile.Close()
	return tmpFile.Name(), nil
}

func (d *DockerDriver) StatFile(t *testing.T, path string) (os.FileInfo, error) {
	file, err := d.retrieveFile(t, path)
	if err != nil {
		return nil, err
	}
	defer os.Remove(file)

	f, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, nil
	}
	return f, nil
}

func (d *DockerDriver) ReadFile(t *testing.T, path string) ([]byte, error) {
	file, err := d.retrieveFile(t, path)
	if err != nil {
		return nil, err
	}
	defer os.Remove(file)
	return ioutil.ReadFile(file)
}

func (d *DockerDriver) ReadDir(t *testing.T, path string) ([]os.FileInfo, error) {
	tmpDir, err := d.retrieveFile(t, path)
	defer os.RemoveAll(tmpDir)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadDir(tmpDir)
}

// This method takes a command (in the form of a list of args), and does the following:
// 1) creates a container, based on the "current latest" image, with the command set as
// the command to run when the container starts
// 2) starts the container
// 3) commits the container with its changes to a new image,
// and sets that image as the new "current image"
func (d *DockerDriver) runAndCommit(t *testing.T, env []string, command []string) string {
	shouldRun := true

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

	if shouldRun {
		err = d.cli.StartContainer(container.ID, nil)
		if err != nil {
			t.Errorf("Error creating container: %s", err.Error())
		}

		_, err = d.cli.WaitContainer(container.ID)
	}

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

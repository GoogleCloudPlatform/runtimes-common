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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"context"

	"bufio"

	"regexp"

	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/types/unversioned"
	"github.com/fsouza/go-dockerclient"
)

type DockerDriver struct {
	cli           docker.Client
	originalImage string
	currentImage  string
}

func (d *DockerDriver) Info() string {
	return fmt.Sprintf("DockerDriver:\nOriginalImage: %s\ncurrentImage: %s\ncli: %s\n",
		d.originalImage, d.currentImage, d.cli)
}

func NewDockerDriver(image string) *DockerDriver {
	newCli, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	return &DockerDriver{
		originalImage: image,
		currentImage:  image,
		cli:           *newCli,
	}
}

func (d *DockerDriver) ProcessCommand(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) (string, string, int) {

	if len(fullCommand) == 0 {
		t.Logf("empty command provided: skipping...")
		return "", "", -1
	}
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
	//reset image for next test
	d.currentImage = d.originalImage
	return stdout, stderr, exitCode
}

func (d *DockerDriver) SetEnvVars(t *testing.T, vars []unversioned.EnvVar) []unversioned.EnvVar {
	if len(vars) == 0 {
		return nil
	}

	image, err := d.cli.InspectImage(d.currentImage)
	if err != nil {
		t.Errorf("Error when inspecting image: %s", err.Error())
		return nil
	}

	// convert env to map for easier processing
	imageEnv := make(map[string]string)
	for _, varPair := range image.Config.Env {
		pair := strings.Split(varPair, "=")
		imageEnv[pair[0]] = pair[1]
	}

	before := regexp.MustCompile(".*\\$(.*?):")
	after := regexp.MustCompile(".*:\\$(.*)")

	env := []string{}
	for _, envVar := range vars {
		currentVar := ""
		if match := before.FindStringSubmatch(envVar.Value); match != nil {
			// first entry is the leftmost substring: second entry is the first group
			currentVar = match[1]
		}
		if match := after.FindStringSubmatch(envVar.Value); match != nil {
			currentVar = match[1]
		}
		if currentVar != "" {
			if val, ok := imageEnv[currentVar]; ok {
				env = append(env, envVar.Key+"="+strings.Replace(envVar.Value, "$"+currentVar, val, -1))
			} else {
				t.Errorf("Variable %s not found in image env! Check test config.", currentVar)
			}
		} else {
			env = append(env, envVar.Key+"="+envVar.Value)
		}
	}
	// since these are global envvars, just overwrite the original image
	d.originalImage = d.runAndCommit(t, env, nil)
	d.currentImage = d.originalImage
	return nil
}

func (d *DockerDriver) ResetEnvVars(t *testing.T, vars []unversioned.EnvVar) {
	// since the container will be destroyed after the tests, this is a noop
}

// copies a file from a docker container to the local fs, and returns its path
// caller is responsible for removing this file when finished
func (d *DockerDriver) retrieveFile(t *testing.T, path string, directory bool) (string, error) {
	ctx := context.Background()

	// TODO(nkubala): this contains a hack to get around the fact that
	// the client will not allow creating a container without a command.
	// we should remove this, or at the very least find a better alternative
	// given that not every container is guaranteed to have a shell.
	container, err := d.cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        d.currentImage,
			Cmd:          []string{"NOOP_COMMAND_DO_NOT_RUN"},
			AttachStdout: true,
			AttachStderr: true,
		},
		HostConfig:       nil,
		NetworkingConfig: nil,
		Context:          ctx,
	})
	if err != nil {
		t.Errorf("Error creating container: %s", err.Error())
		return "", err
	}

	if directory {
		return "", nil
	} else {
		tmpFile, err := ioutil.TempFile("", "structure_test")
		if err != nil {
			t.Errorf("Error when creating temp file: %s", err.Error())
			return "", err
		}
		stream := bufio.NewWriter(tmpFile)

		err = d.cli.DownloadFromContainer(container.ID, docker.DownloadFromContainerOptions{
			OutputStream: stream,
			Path:         path,
			Context:      ctx,
		})
		if err != nil {
			t.Errorf("Error when downloading file from container: %s", err.Error())
			return "", err
		}
		stream.Flush()
		tmpFile.Close()
		return tmpFile.Name(), nil
	}
}

func (d *DockerDriver) StatFile(t *testing.T, path string) (os.FileInfo, error) {
	file, err := d.retrieveFile(t, path, false)
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
	file, err := d.retrieveFile(t, path, false)
	if err != nil {
		return nil, err
	}
	defer os.Remove(file)
	return ioutil.ReadFile(file)
}

func (d *DockerDriver) ReadDir(t *testing.T, path string) ([]os.FileInfo, error) {
	// TODO(nkubala): unimplemented
	return nil, nil
}

// This method takes a command (in the form of a list of args), and does the following:
// 1) creates a container, based on the "current latest" image, with the command set as
// the command to run when the container starts
// 2) starts the container
// 3) commits the container with its changes to a new image,
// and sets that image as the new "current image"
func (d *DockerDriver) runAndCommit(t *testing.T, env []string, command []string) string {

	// this is a placeholder command since apparently the client doesnt allow creating
	// a container without a command.
	// TODO(nkubala): figure out how to remove this
	if len(command) == 0 {
		command = []string{"/bin/sh"}
	}

	ctx := context.Background()

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
		Context:          ctx,
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
		Context:   ctx,
	})

	if err != nil {
		t.Errorf("Error committing container: %s", err.Error())
	}

	d.currentImage = image.ID
	return image.ID
}

func (d *DockerDriver) exec(t *testing.T, env []string, command []string) (string, string, int) {
	ctx := context.Background()

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
		Context:          ctx,
	})
	if err != nil {
		t.Errorf("Error creating container: %s", err.Error())
		return "", "", -1
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err = d.cli.StartContainer(container.ID, nil)
	if err != nil {
		t.Errorf("Error creating container: %s", err.Error())
	}

	err = d.cli.PauseContainer(container.ID)

	_, err = d.cli.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Logs:         true,
		Stream:       true,
		Stdout:       true,
		Stderr:       true,
	})

	// since we can't block without losing the exit code, give go some
	// time to attach before unpausing anything in the container
	time.Sleep(1000 * time.Millisecond)

	if err != nil {
		t.Errorf("Error when attaching to container: %s", err.Error())
	}

	err = d.cli.UnpauseContainer(container.ID)

	exitCode, err := d.cli.WaitContainer(container.ID)
	if err != nil {
		t.Errorf("Error when waiting for container: %s", err.Error())
	}

	return stdout.String(), stderr.String(), exitCode
}

func (d *DockerDriver) Setup(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) {
	env := []string{}
	for _, envVar := range envVars {
		env = append(env, envVar.Key+"="+envVar.Value)
	}
	d.currentImage = d.runAndCommit(t, env, fullCommand)
}

func (d *DockerDriver) Teardown(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) {
	// reset to the original image
	d.currentImage = d.originalImage
}

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
	"fmt"
	"os"
	"testing"

	"context"

	"bytes"

	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/types/unversioned"
	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/utils"
	"github.com/fsouza/go-dockerclient"
)

type DockerDriver struct {
	OriginalImage    string
	CurrentImage     string
	CurrentContainer string
	cli              docker.Client
}

func (d DockerDriver) Info() string {
	return fmt.Sprintf("DockerDriver:\nOriginalImage: %s\nCurrentImage: %s\nCurrentContainer: %s\ncli: %s\n",
		d.OriginalImage, d.CurrentImage, d.CurrentContainer, d.cli)
}

func NewDockerDriver(image string) DockerDriver {
	newCli, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	return DockerDriver{
		OriginalImage:    image,
		CurrentImage:     image,
		CurrentContainer: "",
		cli:              *newCli,
	}
}

func (d DockerDriver) ProcessCommand(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) (string, string, int) {

	if len(fullCommand) == 0 {
		t.Logf("empty command provided: skipping...")
		return "", "", -1
	}
	env := []string{}
	if len(envVars) > 0 {
		for _, envVar := range envVars {
			env = append(env, envVar.Key+"="+envVar.Value)
		}
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

func (d DockerDriver) SetEnvVars(t *testing.T, vars []unversioned.EnvVar) []unversioned.EnvVar {
	if len(vars) == 0 {
		return nil
	}
	env := []string{}
	for _, envVar := range vars {
		env = append(env, envVar.Key+"="+envVar.Value)
	}
	d.runAndCommit(t, env, nil)
	return nil
}

func (d DockerDriver) ResetEnvVars(t *testing.T, vars []unversioned.EnvVar) {
	// since the container will be destroyed after the tests, this is a noop
}

func (d DockerDriver) StatFile(path string) (os.FileInfo, error) {
	// TODO(nkubala): unimplemented
	return nil, nil
}

func (d DockerDriver) ReadFile(path string) ([]byte, error) {
	// TODO(nkubala): unimplemented
	return nil, nil
}

func (d DockerDriver) ReadDir(path string) ([]os.FileInfo, error) {
	// TODO(nkubala): unimplemented
	return nil, nil
}

// This method takes a command (in the form of a list of args), and does the following:
// 1) creates a container, based on the "current latest" image, with the command set as
// the command to run when the container starts
// 2) starts the container
// 3) commits the container with its changes to a new image,
// and sets that image as the new "current image"
func (d DockerDriver) runAndCommit(t *testing.T, env []string, command []string) string {
	name := utils.GenerateContainerName()
	t.Logf("container name is %s", name)
	t.Logf("env: %v", env)
	t.Logf("command: %v", command)

	// this is a placeholder command since apparently the client doesnt allow creating
	// a container without a command.
	// TODO(nkubala): figure out how to remove this
	if len(command) == 0 {
		command = []string{"/bin/sh"}
	}

	ctx := context.Background()

	container, err := d.cli.CreateContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
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
	t.Logf("container name: %s", container.Name)

	err = d.cli.StartContainer(name, nil)
	if err != nil {
		t.Errorf("Error creating container: %s", err.Error())
	}

	image, err := d.cli.CommitContainer(docker.CommitContainerOptions{
		Container: name,
		Context:   ctx,
	})

	if err != nil {
		t.Errorf("Error committing container: %s", err.Error())
	}

	d.CurrentImage = image.ID
	return image.ID

	// return name
}

func (d DockerDriver) exec(t *testing.T, env []string, command []string) (string, string, int) {
	ctx := context.Background()

	exec, err := d.cli.CreateExec(docker.CreateExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Env:          env,
		Cmd:          command,
		Container:    d.CurrentContainer,
		Context:      ctx,
	})

	if err != nil {
		t.Errorf("Error when creating exec instance: %s", err.Error())
		return "", "", -1
	}

	execID := exec.ID

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	err = d.cli.StartExec(execID, docker.StartExecOptions{
		OutputStream: stdout,
		ErrorStream:  stderr,
		Detach:       true, //do we want to detach here?
		Tty:          false,
		Context:      ctx,
	})

	if err != nil {
		t.Errorf("Error when starting exec instance: %s", err.Error())
		return "", "", -1
	}

	// TODO(nkubala): do we need to wait here?

	var exitCode int
	for {
		execInspect, err := d.cli.InspectExec(execID)
		if err != nil {
			t.Errorf("Error when inspecting exec: %s", err.Error())
			return "", "", -1
		}
		if execInspect.Running {
			continue
		} else {
			exitCode = execInspect.ExitCode
			break
		}
	}

	return stdout.String(), stderr.String(), exitCode
}

func (d DockerDriver) Setup(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) {
	env := []string{}
	for _, envVar := range envVars {
		env = append(env, envVar.Key+"="+envVar.Value)
	}
	d.runAndCommit(t, env, fullCommand)
}

func (d DockerDriver) Teardown(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) {
	// since the container will be destroyed after the tests, this is a noop
}

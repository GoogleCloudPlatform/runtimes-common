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
	"fmt"
	"os"
	"strings"
	"testing"

	"context"

	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/types/unversioned"
	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerDriver struct {
	OriginalImage    string
	CurrentImage     string
	CurrentContainer string
	cli              client.Client
}

func (d DockerDriver) Info() string {
	return fmt.Sprintf("DockerDriver:\nOriginalImage: %s\nCurrentImage: %s\nCurrentContainer: %s\ncli: %s\n",
		d.OriginalImage, d.CurrentImage, d.CurrentContainer, d.cli)
}

func NewDockerDriver(image string) DockerDriver {
	newCli, err := client.NewEnvClient()
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
	args := []string{}
	if len(envVars) > 0 {
		args = append(args, "--env")
		for _, envVar := range envVars {
			args = append(args, envVar.Key+"="+envVar.Value)
			env = append(env, envVar.Key+"="+envVar.Value)
		}
	}
	args = append(args, d.CurrentContainer)
	// TODO(nkubala): do we still need this?
	if shellMode {
		args = append(args, "/bin/sh", "-c", strings.Join(fullCommand, " "))
	}
	// stdout, stderr, err := d.exec(t, flags)
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
	flags := []string{}
	for _, envVar := range vars {
		flags = append(flags, envVar.Key+"="+envVar.Value)
	}
	d.runAndCommit(t, flags)
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
func (d DockerDriver) runAndCommit(t *testing.T, args []string) string {
	name := utils.GenerateContainerName()
	t.Logf("container name is %s", name)
	t.Logf("args: %v", args)

	ctx := context.Background()
	_, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: d.CurrentImage,
		Cmd:   args, // this command gets run when the container starts
	}, nil, nil, name)
	if err != nil {
		panic(err)
	}

	startoptions := &types.ContainerStartOptions{}
	err = d.cli.ContainerStart(ctx, name, *startoptions)

	commitOptions := &types.ContainerCommitOptions{}
	resp, err := d.cli.ContainerCommit(ctx, name, *commitOptions)
	d.CurrentImage = resp.ID
	return name
}

func (d DockerDriver) exec(t *testing.T, env []string, command []string) (string, string, int) {
	ctx := context.Background()

	config := &types.ExecConfig{
		User:         "root",
		Privileged:   true,
		Tty:          false,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		DetachKeys:   "",
		Env:          env,
		Cmd:          command,
	}

	response, err := d.cli.ContainerExecCreate(ctx, d.CurrentContainer, *config)
	if err != nil {
		panic(err)
	}
	execID := response.ID

	// var stdout, stderr io.Writer

	//TODO(nkubala): figure out how to capture stdout/stderr here
	// stdout = d.cli.out
	// stderr = d.cli.err

	var stdout, stderr string

	resp, err := d.cli.ContainerExecAttach(ctx, execID, *config)
	if err != nil {
		panic(err)
	}

	// TODO: how to unmultiplex this stream?
	// TODO: this should be in a goroutine so it can continue reading while we inspect the proc?
	// see https://godoc.org/github.com/moby/moby/client (ctrl + f "containerattach") and
	// https://github.com/moby/moby/blob/master/client/container_attach.go#L34
	scanner := bufio.NewScanner(resp.Reader)
	for scanner.Scan() {
		stdout = stdout + scanner.Text()
	}

	// TODO: this should be in another goroutine
	// should wrap in a 'while(inspectResp.running)' loop
	var status int
	inspectResp, err := d.cli.ContainerExecInspect(ctx, execID)
	if err != nil {
		status = -1
	}
	status = inspectResp.ExitCode

	resp.Close()
	return stdout, stderr, status
}

func (d DockerDriver) Setup(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) {
	flags := []string{}
	// TODO(nkubala): run the command
	for _, envVar := range envVars {
		flags = append(flags, envVar.Key+"="+envVar.Value)
	}
	d.runAndCommit(t, flags)
}

func (d DockerDriver) Teardown(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) {
	// since the container will be destroyed after the tests, this is a noop
}

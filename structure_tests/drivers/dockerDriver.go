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
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"

	"context"

	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/types/unversioned"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/moby/moby/api/client"
)

type DockerDriver struct {
	OriginalImage    string
	CurrentImage     string
	CurrentContainer string
	cli              client.Client
}

func New(image string) DockerDriver {
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
	flags := []string{}
	if len(envVars) > 0 {
		flags = append(flags, "--env")
		for _, envVar := range envVars {
			flags = append(flags, envVar.Key+"="+envVar.Value)
			env = append(env, envVar.Key+"="+envVar.Value)
		}
	}
	flags = append(flags, d.CurrentContainer)
	if shellMode {
		flags = append(flags, "/bin/sh", "-c", strings.Join(fullCommand, " "))
	}
	// stdout, stderr, err := d.exec(t, flags)
	d.exec(t, env, fullCommand)

	if stdout != "" {
		t.Logf("stdout: %s", stdout)
	}
	if stderr != "" {
		t.Logf("stderr: %s", stderr)
	}
	var exitCode int
	if err != nil {
		if checkOutput {
			// The test might be designed to run a command that exits with an error.
			t.Logf("Error running command: %s. Continuing.", err)
		} else {
			t.Fatalf("Error running setup/teardown command: %s.", err)
		}
		switch err := err.(type) {
		default:
			t.Errorf("Command failed to start! Unable to retrieve error info!")
		case *exec.ExitError:
			exitCode = err.Sys().(syscall.WaitStatus).ExitStatus()
		case *exec.Error:
			// Command started but failed to finish, so we can at least check the stderr
			stderr = err.Error()
		}
	} else {
		exitCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	}
	return stdout, stderr, exitCode
}

func (d DockerDriver) SetEnvVars(t *testing.T, vars []unversioned.EnvVar) []unversioned.EnvVar {
	flags := []string{}
	for _, envVar := range vars {
		flags = append(flags, envVar.Key+"="+envVar.Value)
	}
	newImage := d.runAndCommit(flags)
	t.Logf(newImage)
	return nil
}

func (d DockerDriver) ResetEnvVars(t *testing.T, vars []unversioned.EnvVar) {
	// since the container will be destroyed after the tests, this is a noop
}

func (d DockerDriver) StatFile(path string) (os.FileInfo, error) {
	return nil, nil
}

func (d DockerDriver) ReadFile(path string) ([]byte, error) {
	return nil, nil
}

func (d DockerDriver) runAndCommit(flags []string) string {
	name := "abc1"
	// name := utils.GenerateContainerName()

	ctx := context.Background()
	createResp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: d.CurrentImage,
		Cmd:   flags,
	}, nil, nil, name)
	if err != nil {
		panic(err)
	}
	d.CurrentContainer = name
	return name

	// var cmd *exec.Cmd
	// flags = append([]string{"run", "-itd", "--name", name, d.CurrentImage}, flags...)
	// cmd = exec.Command("docker", flags...)
	// var outbuf, errbuf bytes.Buffer

	// cmd.Stdout = &outbuf
	// cmd.Stderr = &errbuf

	// if err := cmd.Run(); err != nil {
	// 	panic(err)
	// }
	// stdout := outbuf.String()
	// d.CurrentContainer = stdout
	// return stdout
}

func (d DockerDriver) exec(t *testing.T, env []string, command []string) (string, string, error) {

	ctx := context.Background()

	config := &types.ExecConfig{
		User:         "root",
		Privileged:   true,
		Tty:          true,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		DetachKeys:   "",
		Env:          env,
		Cmd:          command,
	}

	response, err := cli.ContainerExecCreate(ctx, d.CurrentContainer, config)

	//TODO: figure out how to capture stdout/stderr here

	// var cmd *exec.Cmd
	// flags = append([]string{"exec", "-itd"}, flags...)
	// cmd = exec.Command("docker", flags...)
	// t.Logf("Executing: %s", cmd.Args)

	// var outbuf, errbuf bytes.Buffer

	// cmd.Stdout = &outbuf
	// cmd.Stderr = &errbuf

	// err := cmd.Run()
	// return outbuf.String(), errbuf.String(), err
}

func (d DockerDriver) Setup(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) (string, string, int) {
	//create args (we already know command will be "docker")
	//send to runAndCommit
}

func (d DockerDriver) Teardown(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) (string, string, int) {

}

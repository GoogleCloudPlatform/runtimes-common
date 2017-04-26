// Copyright 2016 Google Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
)

type StructureTestv1 struct {
	GlobalEnvVars      []EnvVar
	CommandTests       []CommandTestv1
	FileExistenceTests []FileExistenceTestv1
	FileContentTests   []FileContentTestv1
	LicenseTests       []LicenseTestv1
}

func (st StructureTestv1) RunAll(t *testing.T) int {
	originalVars := SetEnvVars(t, st.GlobalEnvVars)
	defer ResetEnvVars(t, originalVars)
	testsRun := 0
	testsRun += st.RunCommandTests(t)
	testsRun += st.RunFileExistenceTests(t)
	testsRun += st.RunFileContentTests(t)
	testsRun += st.RunLicenseTests(t)
	return testsRun
}

func (st StructureTestv1) RunCommandTests(t *testing.T) int {
	counter := 0
	for _, tt := range st.CommandTests {
		t.Run(tt.LogName(), func(t *testing.T) {
			validateCommandTestV1(t, tt)
			for _, setup := range tt.Setup {
				ProcessCommand(t, tt.EnvVars, setup, tt.ShellMode, false)
			}

			stdout, stderr, exitcode := ProcessCommand(t, tt.EnvVars, tt.Command, tt.ShellMode, true)
			CheckOutput(t, tt, stdout, stderr, exitcode)

			for _, teardown := range tt.Teardown {
				ProcessCommand(t, tt.EnvVars, teardown, tt.ShellMode, false)
			}
			counter++
		})
	}
	return counter
}

func (st StructureTestv1) RunFileExistenceTests(t *testing.T) int {
	counter := 0
	for _, tt := range st.FileExistenceTests {
		t.Run(tt.LogName(), func(t *testing.T) {
			validateFileExistenceTestV1(t, tt)
			var err error
			var info os.FileInfo
			if tt.IsDirectory {
				info, err = os.Stat(tt.Path)
			} else {
				info, err = os.Stat(tt.Path)
			}
			if tt.ShouldExist && err != nil {
				if tt.IsDirectory {
					t.Errorf("Directory %s should exist but does not!", tt.Path)
				} else {
					t.Errorf("File %s should exist but does not!", tt.Path)
				}
			} else if !tt.ShouldExist && err == nil {
				if tt.IsDirectory {
					t.Errorf("Directory %s should not exist but does!", tt.Path)
				} else {
					t.Errorf("File %s should not exist but does!", tt.Path)
				}
			}
			if tt.Permissions != "" {
				perms := info.Mode()
				if perms.String() != tt.Permissions {
					t.Errorf("%s has incorrect permissions. Expected: %s, Actual: %s", tt.Path, tt.Permissions, perms.String())
				}
			}
			counter++
		})
	}
	return counter
}

func (st StructureTestv1) RunFileContentTests(t *testing.T) int {
	counter := 0
	for _, tt := range st.FileContentTests {
		t.Run(tt.LogName(), func(t *testing.T) {
			validateFileContentTestV1(t, tt)
			actualContents, err := ioutil.ReadFile(tt.Path)
			if err != nil {
				t.Errorf("Failed to open %s. Error: %s", tt.Path, err)
			}

			contents := string(actualContents[:])

			var errMessage string
			for _, s := range tt.ExpectedContents {
				errMessage = "Expected string " + s + " not found in file contents!"
				compileAndRunRegex(s, contents, t, errMessage, true)
			}
			for _, s := range tt.ExcludedContents {
				errMessage = "Excluded string " + s + " found in file contents!"
				compileAndRunRegex(s, contents, t, errMessage, false)
			}
			counter++
		})
	}
	return counter
}

func (st StructureTestv1) RunLicenseTests(t *testing.T) int {
	for num, tt := range st.LicenseTests {
		t.Run(tt.LogName(num), func(t *testing.T) {
			checkLicenses(t, tt)
		})
		return 1
	}
	return 0
}

// given an array of command parts, construct a full command and execute it against the
// current environment. a list of environment variables can be passed to be set in the
// environment before the command is executed. additionally, a boolean flag is passed
// to specify whether or not we care about the output of the command.
func ProcessCommand(t *testing.T, envVars []EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) (string, string, int) {
	var cmd *exec.Cmd
	if len(fullCommand) == 0 {
		t.Logf("empty command provided: skipping...")
		return "", "", -1
	}
	var command string
	var flags []string
	if shellMode {
		command = "/bin/sh"
		flags = []string{"-c", strings.Join(fullCommand, " ")}
	} else {
		command = fullCommand[0]
		flags = fullCommand[1:]
	}
	originalVars := SetEnvVars(t, envVars)
	defer ResetEnvVars(t, originalVars)
	if len(flags) > 0 {
		cmd = exec.Command(command, flags...)
	} else {
		cmd = exec.Command(command)
	}

	if checkOutput {
		t.Logf("Executing: %s", cmd.Args)
	} else {
		t.Logf("Executing setup/teardown: %s", cmd.Args)
	}

	var outbuf, errbuf bytes.Buffer

	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout := outbuf.String()
	if stdout != "" {
		t.Logf("stdout: %s", stdout)
	}
	stderr := errbuf.String()
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

// given a list of environment variable key/value pairs, set these in the current environment.
// also, keep track of the previous values of these vars to reset after test execution.
func SetEnvVars(t *testing.T, vars []EnvVar) []EnvVar {
	var originalVars []EnvVar
	for _, env_var := range vars {
		originalVars = append(originalVars, EnvVar{env_var.Key, os.Getenv(env_var.Key)})
		if err := os.Setenv(env_var.Key, os.ExpandEnv(env_var.Value)); err != nil {
			t.Fatalf("error setting env var: %s", err)
		}
	}
	return originalVars
}

func ResetEnvVars(t *testing.T, vars []EnvVar) {
	for _, env_var := range vars {
		var err error
		if env_var.Value == "" {
			// if the previous value was empty string, the variable did not
			// exist in the environment; unset it
			err = os.Unsetenv(env_var.Key)
		} else {
			// otherwise, set it back to its previous value
			err = os.Setenv(env_var.Key, env_var.Value)
		}
		if err != nil {
			t.Fatalf("error resetting env var: %s", err)
		}
	}
}

func CheckOutput(t *testing.T, tt CommandTestv1, stdout string, stderr string, exitCode int) {
	for _, errStr := range tt.ExpectedError {
		errMsg := fmt.Sprintf("Expected string '%s' not found in error!", errStr)
		compileAndRunRegex(errStr, stderr, t, errMsg, true)
	}
	for _, errStr := range tt.ExcludedError {
		errMsg := fmt.Sprintf("Excluded string '%s' found in error!", errStr)
		compileAndRunRegex(errStr, stderr, t, errMsg, false)
	}
	for _, outStr := range tt.ExpectedOutput {
		errMsg := fmt.Sprintf("Expected string '%s' not found in output!", outStr)
		compileAndRunRegex(outStr, stdout, t, errMsg, true)
	}
	for _, outStr := range tt.ExcludedOutput {
		errMsg := fmt.Sprintf("Excluded string '%s' found in output!", outStr)
		compileAndRunRegex(outStr, stdout, t, errMsg, false)
	}
	if tt.ExitCode != exitCode {
		t.Errorf("Test %s exited with incorrect error code! Expected: %d, Actual: %d", tt.Name, tt.ExitCode, exitCode)
	}
}

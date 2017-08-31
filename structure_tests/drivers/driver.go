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
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/types/unversioned"
)

type Driver interface {
	Setup(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
		shellMode bool, checkOutput bool)

	Teardown(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
		shellMode bool, checkOutput bool)

	// given an array of command parts, construct a full command and execute it against the
	// current environment. a list of environment variables can be passed to be set in the
	// environment before the command is executed. additionally, a boolean flag is passed
	// to specify whether or not we care about the output of the command.
	ProcessCommand(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
		shellMode bool, checkOutput bool) (string, string, int)

	// given a list of environment variable key/value pairs, set these in the current environment.
	// also, keep track of the previous values of these vars to reset after test execution.
	SetEnvVars(t *testing.T, vars []unversioned.EnvVar) []unversioned.EnvVar

	ResetEnvVars(t *testing.T, vars []unversioned.EnvVar)

	StatFile(path string) (os.FileInfo, error)

	ReadFile(path string) ([]byte, error)

	ReadDir(path string) ([]os.FileInfo, error)
}

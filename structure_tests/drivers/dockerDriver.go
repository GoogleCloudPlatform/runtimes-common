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

type DockerDriver struct {
}

func (d DockerDriver) ProcessCommand(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string,
	shellMode bool, checkOutput bool) (string, string, int) {
	return "", "", 0
}

func (d DockerDriver) SetEnvVars(t *testing.T, vars []unversioned.EnvVar) []unversioned.EnvVar {
	return nil
}

func (d DockerDriver) ResetEnvVars(t *testing.T, vars []unversioned.EnvVar) {

}

func (d DockerDriver) StatFile(path string) (os.FileInfo, error) {
	return nil, nil
}

func (d DockerDriver) ReadFile(path string) ([]byte, error) {
	return nil, nil
}

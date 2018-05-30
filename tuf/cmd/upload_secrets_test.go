/*
Copyright 2018 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/config"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/deployer"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/testutil"
)

type MockDeployer struct {
	returnError bool
}

var deployerError bool

func NewMockDeployer() (deployer.DeployTool, error) {
	return &MockDeployer{
		returnError: deployerError,
	}, nil
}

func (d MockDeployer) UpdateSecrets(tufConfig config.TUFConfig, rootKeyFile string, targetKeyFile string, snapshotKeyFile string) error {
	if d.returnError {
		return errors.New("Some err")
	}
	return nil
}

func (d MockDeployer) GenerateMetadata(tufConfig config.TUFConfig, root string,
	target string, snapshot string, oldroot []byte) error {
	return nil
}

var uploadSecretsTC = []struct {
	name        string
	config      string
	deployErr   bool
	expectedErr error
	emptyArgs   bool
}{
	// Since flag variables are set, we need run this always as first test case.
	{"emptyArgs", testutil.MarshalledTUFConfig(), false,
		errors.New("Please specify atleast on secret to upload"), true},
	{"updateSecretsSuccess", testutil.MarshalledTUFConfig(), false, nil, false},
	{"invalidConfig", "invalidYaml", false, errors.New("yaml: unmarshal errors"), false},
	{"deployError", testutil.MarshalledTUFConfig(), true, errors.New("Some err"), false},
}

func TestUpdateSecrets(t *testing.T) {
	for _, tc := range uploadSecretsTC {
		t.Run(tc.name, func(t *testing.T) {
			tmpdir, err := ioutil.TempDir("", "upload_")
			if err != nil {
				panic(fmt.Sprintf("Cannot run tests due to %v", err))
			}
			defer os.Remove(tmpdir)

			tufConfig := testutil.CreateAndWriteFile(tmpdir, "tufConfig.yaml", tc.config)
			if !tc.emptyArgs {
				tmpfile := testutil.CreateAndWriteFile(tmpdir, "encrtypedKey.json", "")
				RootCommand.SetArgs([]string{"upload-secrets",
					"--config", tufConfig,
					"--root-key", tmpfile,
					"--target-key", tmpfile,
					"--snapshot-key", tmpfile})
			} else {
				RootCommand.SetArgs([]string{"upload-secrets",
					"--config", tufConfig})
			}
			deployerError = tc.deployErr
			DefaultDeployTool = NewMockDeployer
			err = RootCommand.Execute()

			if !testutil.IsErrorEqualOrContains(err, tc.expectedErr) {
				t.Fatalf("Expected Err: %v\nGot: %v", tc.expectedErr, err)
			}

		})
	}
}

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
package ctc_lib

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigCommand(t *testing.T) {
	Version = "1.0.1"
	ConfigFile = "testdata/testConfig.json"
	testCommand := ContainerToolCommand{
		ContainerToolCommandBase: &ContainerToolCommandBase{
			Command: &cobra.Command{
				Use: "Hello Command",
			},
			Phase: "test",
		},
		Output: nil,
		RunO: func(command *cobra.Command, args []string) (interface{}, error) {
			return nil, nil
		},
	}

	testCommand.SetArgs([]string{"config"})
	Execute(&testCommand)
	expectedConfig := &ConfigOutput{
		Config: map[string]interface{}{
			"message":     "echo", // Make sure user defined config are also returned
			"updatecheck": "true", // inhertited from the Default Config
			"logdir":      "/tmp", // This overrides the default Config
		},
	}
	if reflect.DeepEqual(ConfigCommand.Output, expectedConfig) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expectedConfig, ConfigCommand.Output)
	}
}

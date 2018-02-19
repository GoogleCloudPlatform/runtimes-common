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
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/container_command_output"

	"github.com/spf13/cobra"
)

type TestInterface struct {
	testField string
}

func TestContainerToolCommandVersion(t *testing.T) {
	var testCommand = ContainerToolCommand{
		Command: &cobra.Command{
			Use: "test Command",
			Run: func(command *cobra.Command, args []string) {
				fmt.Println("Example test")
			},
		},
		Phase:         "test",
		CommandOutput: container_command_output.ContainerCommandOutput{},
		Version:       "1.0.1",
	}
	testCommand.SetArgs([]string{"version"})
	testCommand.Execute()
	//	var expected string = "1.0"
	if testCommand.CommandOutput.Version != "1.0" {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", "1.0", testCommand.CommandOutput.Version)
	}
}

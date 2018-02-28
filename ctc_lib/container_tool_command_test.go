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
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/version"
	"github.com/spf13/cobra"
)

type TestInterface struct {
	Salutation string
	Name       string
}

var Salutation string
var Name string

var testCommand = ContainerToolCommand{
	Command: &cobra.Command{
		Use:  "Hello Command",
		RunE: RunCommand,
	},
	Phase:   "test",
	Version: "1.0.1",
	Output:  &TestInterface{},
}

func RunCommand(command *cobra.Command, args []string) error {
	fmt.Println("Running Hello World Command")
	if Name == "" {
		return errors.New("Please supply Name Argument")
	}
	testOutput := TestInterface{
		Salutation: Salutation,
		Name:       Name,
	}
	WriteObject(command, testOutput)
	return nil
}

func setup() {
	testCommand.Flags().StringVarP(&Salutation, "salutation", "s", "", "Salutation")
	testCommand.Flags().StringVarP(&Name, "name", "n", "", "Name")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestContainerToolCommandVersion(t *testing.T) {
	testCommand.SetArgs([]string{"version"})
	output := testCommand.ExecuteO()
	var expectedOutput = version.VersionOutput{
		Version: "1.0.1",
	}
	if reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expectedOutput, output)
	}
}

func TestContainerToolCommandHelp(t *testing.T) {
	testCommand.SetArgs([]string{"help"})
	testCommand.Execute()
	if "1" != "HELP STRING" {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", "HELP STRING", "!")
	}
}

func TestContainerToolCommandOutput(t *testing.T) {
	testCommand.SetArgs([]string{"--name=Sparks", "--salutation=Mr."})
	output := testCommand.ExecuteO()
	var expectedOutput = TestInterface{
		Salutation: "Mr.",
		Name:       "Sparks",
	}
	if reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expectedOutput, testCommand.Output)
	}
}

func TestContainerToolCommandOutputError(t *testing.T) {
	testCommand.Execute()
	// if testCommand.OutputBuffer.String() != "1.0.1" {
	// 	t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", "1.0.1", testCommand.OutputBuffer)
	// }
}

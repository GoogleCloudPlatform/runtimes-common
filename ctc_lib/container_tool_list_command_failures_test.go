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
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

var Channel1 chan interface{}

func addInts() {
	for i := 0; i < 2; i++ {
		Channel1 <- i
	}
	// Make sure to close the stream.
	close(Channel1)
}

func RunStreamFailureCommand(command *cobra.Command, args []string) {
	// Run the method which writes to the stream
	go addInts()
}

//Error in Summary Template should return Error.
func TestContainerToolCommandListFailureOutput(t *testing.T) {
	Channel1 = make(chan interface{}, 1)
	testCommand := ContainerToolListCommand{
		ContainerToolCommandBase: &ContainerToolCommandBase{
			Command: &cobra.Command{
				Use: "Failure command",
			},
			Phase:           "test",
			DefaultTemplate: "{{.}}",
		},
		OutputList:      make([]interface{}, 0),
		StreamO:         RunStreamFailureCommand,
		Stream:          Channel1,
		SummaryTemplate: "{{.fail}}",
		TotalO: func(list []interface{}) (interface{}, error) {
			return len(list), nil
		},
	}
	err := ExecuteE(&testCommand)
	if err == nil {
		t.Errorf("Expected Error however command executed successfully")
	}
}

//Error in Command and in Summary Template should return CommandError
func TestContainerToolCommandListFailureOutput1(t *testing.T) {
	Channel1 = make(chan interface{}, 1)
	testCommand := ContainerToolListCommand{
		ContainerToolCommandBase: &ContainerToolCommandBase{
			Command: &cobra.Command{
				Use: "Failure command",
			},
			Phase:           "test",
			DefaultTemplate: "{{.fail}}",
		},
		OutputList:      make([]interface{}, 0),
		StreamO:         RunStreamFailureCommand,
		Stream:          Channel1,
		SummaryTemplate: "{{.fail}}",
		TotalO: func(list []interface{}) (interface{}, error) {
			return len(list), nil
		},
	}
	err := ExecuteE(&testCommand)
	expected := "can't evaluate field fail in type int"
	if strings.Contains(expected, err.Error()) {
		t.Errorf("Expected to Contain : \n %q \nGot:\n %q\n", expected, err.Error())
	}
}

//Error when calculating Total should return the error
func TestContainerToolCommandListFailureOutput2(t *testing.T) {
	defer SetExitOnError(true)
	Channel1 = make(chan interface{}, 1)
	testCommand := ContainerToolListCommand{
		ContainerToolCommandBase: &ContainerToolCommandBase{
			Command: &cobra.Command{
				Use: "Failure command",
			},
			Phase:           "test",
			DefaultTemplate: "{{.}}",
		},
		OutputList:      make([]interface{}, 0),
		StreamO:         RunStreamFailureCommand,
		Stream:          Channel1,
		SummaryTemplate: "{{.}}",
		TotalO: func(list []interface{}) (interface{}, error) {
			return len(list), errors.New("Did not return empty list")
		},
	}
	SetExitOnError(false)
	err := ExecuteE(&testCommand)
	expected := "Did not return empty list"
	if expected != err.Error() {
		t.Errorf("Expected to Contain : \n %q \nGot:\n %q\n", expected, err.Error())
	}
}

//Error when calculating Total should return the error
func TestContainerToolCommandListFailureOutput3(t *testing.T) {
	Channel1 = make(chan interface{}, 1)
	testCommand := ContainerToolListCommand{
		ContainerToolCommandBase: &ContainerToolCommandBase{
			Command: &cobra.Command{
				Use: "Failure command",
			},
			Phase:           "test",
			DefaultTemplate: "{{.}}",
		},
		OutputList:      make([]interface{}, 0),
		StreamO:         RunStreamFailureCommand,
		Stream:          Channel1,
		SummaryTemplate: "{{.}}",
		TotalO: func(list []interface{}) (interface{}, error) {
			return len(list), errors.New("Did not return empty list")
		},
	}
	testCommand.SetArgs([]string{"--jsonOutput"})
	err := ExecuteE(&testCommand)
	expected := "Did not return empty list"
	if expected != err.Error() {
		t.Errorf("Expected to Contain : \n %q \nGot:\n %q\n", expected, err.Error())
	}
}

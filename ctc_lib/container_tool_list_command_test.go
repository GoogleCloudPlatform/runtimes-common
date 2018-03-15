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
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/cobra"
)

type TestListOutput struct {
	Name string
}

var LName string

func RunListCommand(command *cobra.Command, args []string) ([]interface{}, error) {
	Log.Debug("Running Hello World Command")
	var testOutputs []TestListOutput
	for _, name := range strings.Split(LName, ",") {
		testOutputs = append(testOutputs, TestListOutput{
			Name: name,
		})
	}
	s := make([]interface{}, len(testOutputs))
	for i, v := range testOutputs {
		s[i] = v
	}
	return s, nil
}

func TestContainerToolCommandListOutput(t *testing.T) {
	testCommand := ContainerToolListCommand{
		ContainerToolCommandBase: &ContainerToolCommandBase{
			Command: &cobra.Command{
				Use: "Hello Command",
			},
			Phase: "test",
		},
		OutputList: make([]interface{}, 0),
		RunO:       RunListCommand,
	}
	testCommand.Flags().StringVarP(&LName, "name", "n", "", "Comma Separated list of Name")
	var OutputBuffer bytes.Buffer
	testCommand.Command.SetOutput(&OutputBuffer)
	testCommand.SetArgs([]string{"--name=John,Jane"})
	Execute(&testCommand)
	var expectedOutput = []TestListOutput{
		{Name: "John"},
		{Name: "Jane"},
	}
	s := make([]interface{}, len(expectedOutput))
	for i, v := range expectedOutput {
		s[i] = v
	}
	if !reflect.DeepEqual(s, testCommand.OutputList) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", s, testCommand.OutputList)
	}
}

func TestContainerToolCommandLogging(t *testing.T) {
	var hook *test.Hook
	testCommand := ContainerToolListCommand{
		ContainerToolCommandBase: &ContainerToolCommandBase{
			Command: &cobra.Command{
				Use: "Hello Command",
				PersistentPreRun: func(c *cobra.Command, args []string) {
					_, hook = test.NewNullLogger()
					Log.AddHook(hook)
				},
			},
			Phase: "test",
		},
		OutputList: make([]interface{}, 0),
		RunO:       RunListCommand,
	}
	testCommand.Flags().StringVarP(&LName, "name", "n", "", "Comma Separated list of Name")
	testCommand.SetArgs([]string{"--name=John,Jane", "--loglevel=debug"})
	Execute(&testCommand)
	assert.Equal(t, len(hook.Entries), 1)
	assert.Equal(t, hook.LastEntry().Message, "Running Hello World Command")
}

func TestContainerToolCommandHandlePanicLogging(t *testing.T) {
	defer SetExitOnError(true)
	var hook *test.Hook
	testCommand := ContainerToolListCommand{
		ContainerToolCommandBase: &ContainerToolCommandBase{
			Command: &cobra.Command{
				Use: "kill",
				PersistentPreRun: func(c *cobra.Command, args []string) {
					_, hook = test.NewNullLogger()
					Log.AddHook(hook)
				},
			},
			Phase: "test",
		},
		OutputList: make([]interface{}, 0),
		RunO: func(command *cobra.Command, args []string) ([]interface{}, error) {
			Log.Panic("Please dont kill me")
			return nil, nil
		},
	}
	SetExitOnError(false)
	Execute(&testCommand)
	assert.Equal(t, hook.LastEntry().Message, "Please dont kill me")
}
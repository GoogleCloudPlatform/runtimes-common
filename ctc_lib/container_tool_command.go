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
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/help"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/version"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/writer"
	"github.com/spf13/cobra"
)

type ContainerToolCommand struct {
	*cobra.Command
	Phase                string
	Version              string
	HelpText             string
	Output               interface{}
	subCommandsOutputMap map[string]interface{}
	args                 []string
}

type ContainerToolCommandList struct {
	*cobra.Command
	Phase      string
	Version    string
	OutputList []interface{}
}

func (ctc *ContainerToolCommand) init() {
	// Add version subcommand
	ctc.AddCommand(version.NewVersionCommand(ctc.Version,
		ctc.Command.Name()).Command, &version.VersionOutput{})
	// Add help subcommand
	ctc.AddCommand(help.NewHelpCommand(ctc.HelpText,
		ctc.Command.Name()).Command, &help.HelpOutput{})

	// Set up Root Command
	ctc.Command.SetHelpTemplate(help.HelpTemplate)
}

func (ctc *ContainerToolCommand) AddCommand(command *cobra.Command, output interface{}) {
	ctc.Command.AddCommand(command)
	ctc.SetSubCommandsOutputMap(command.Name(), output)
}

func (ctc *ContainerToolCommand) SetArgs(args []string) {
	ctc.args = args
	ctc.Command.SetArgs(args)
}

func (ctc *ContainerToolCommand) SetSubCommandsOutputMap(name string, output interface{}) {
	if ctc.subCommandsOutputMap == nil {
		ctc.subCommandsOutputMap = make(map[string]interface{})
	}
	ctc.subCommandsOutputMap[name] = output
}

func (ctc *ContainerToolCommand) GetSubCommandsOutputMap() map[string]interface{} {
	return ctc.subCommandsOutputMap
}

func (ctc *ContainerToolCommand) Execute() error {
	ctc.init()
	fmt.Println("Execute called")
	return ctc.Command.Execute()
}

func (ctc *ContainerToolCommand) ExecuteO() interface{} {
	ctc.init()
	ctcWriter := writer.NewCTCBuffer(ctc.Output)
	ctc.Command.SetOutput(ctcWriter)
	err := ctc.Command.Execute()
	if err != nil {
		return err
	}

	targetCommand, _, _ := ctc.Command.Find(ctc.args)
	output := ctc.GetSubCommandsOutputMap()[targetCommand.Name()]
	err = json.Unmarshal(ctcWriter.OutputBuffer.Bytes(), output)
	if err != nil {
		return err
	}
	return output
}

func WriteObject(cmd *cobra.Command, obj interface{}) {
	jsonEncoded, _ := json.Marshal(obj)
	cmd.Print(string(jsonEncoded))
}

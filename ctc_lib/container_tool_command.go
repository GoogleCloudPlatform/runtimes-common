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

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/help"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/version"
	"github.com/spf13/cobra"
)

type ContainerToolCommand struct {
	*cobra.Command
	Phase    string
	Version  string
	HelpText string
	Output   interface{}
	args     []string
	RunO     func(command *cobra.Command, args []string) (interface{}, error)
}

type ContainerToolCommandList struct {
	*cobra.Command
	Phase      string
	Version    string
	OutputList []interface{}
	RunO       func(ctc *ContainerToolCommand, args []string) (interface{}, error)
}

// Define all Flags
var Template string

func (ctc *ContainerToolCommand) init() {
	ctc.AddSubCommands()
	ctc.AddFlags()
}

func (ctc *ContainerToolCommand) AddSubCommands() {
	// Add version subcommand
	ctc.Command.AddCommand(
		version.NewVersionCommand(ctc.Version, ctc.Command.Name()).Command)
	// Add help subcommand
	ctc.Command.AddCommand(
		help.NewHelpCommand(ctc.HelpText, ctc.Command.Name()).Command)

	// Set up Root Command
	ctc.Command.SetHelpTemplate(help.HelpTemplate)
}

func (ctc *ContainerToolCommand) AddFlags() {
	// Add template Flag
	ctc.Command.Flags().StringVarP(&Template, "template", "t", "{{.}}", "Output format")

}

func (ctc *ContainerToolCommand) Execute() error {
	ctc.init()
	if (ctc.Command.Run != nil || ctc.Command.RunE != nil) && ctc.RunO != nil {
		errors.New("Cannot provide both Command.Run and RunO implementation" +
			"Either implement Command.Run implementation or RunO implemetation")
	}
	cobraRun := func(c *cobra.Command, args []string) {
		obj, _ := ctc.RunO(c, args)
		fmt.Println(obj, Template)
	}

	ctc.Command.Run = cobraRun
	return ctc.Command.Execute()
}

// This function is used for Testing.
func (ctc *ContainerToolCommand) SetArgs(args []string) {
	ctc.args = args
	ctc.Command.SetArgs(args)
}

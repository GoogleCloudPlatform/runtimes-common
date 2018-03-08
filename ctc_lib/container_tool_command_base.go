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
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/flags"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/util"
	"github.com/spf13/cobra"
)

type ContainerToolCommandBase struct {
	*cobra.Command
	Phase string
}

func (ctb *ContainerToolCommandBase) init() {
	ctb.AddSubCommands()
	ctb.AddFlags()
}

func (ctb *ContainerToolCommandBase) AddSubCommands() {
	// Add version subcommand
	ctb.AddCommand(VersionCommand)

	// Set up Root Command
	ctb.Command.SetHelpTemplate(HelpTemplate)
}

func (ctb *ContainerToolCommandBase) AddCommand(command *ContainerToolCommand) {
	cobraRun := func(c *cobra.Command, args []string) {
		obj, _ := command.RunO(c, args)
		command.Output = obj
		util.ExecuteTemplate(flags.TemplateString, command.Output, ctb.OutOrStdout())
	}
	command.Run = cobraRun
	ctb.Command.AddCommand(command.Command)
}

func (ctb *ContainerToolCommandBase) AddCommandList(command *ContainerToolListCommand) {
	cobraRun := func(c *cobra.Command, args []string) {
		obj, _ := command.RunO(c, args)
		command.OutputList = obj
		util.ExecuteTemplate(flags.TemplateString, command.OutputList, ctb.OutOrStdout())
	}
	command.Run = cobraRun
	ctb.Command.AddCommand(command.Command)
}

func (ctb *ContainerToolCommandBase) AddFlags() {
	// Add template Flag
	ctb.PersistentFlags().StringVarP(&flags.TemplateString, "template", "t", "{{.}}", "Output format")
}

// These functions are used for Testing.
func (ctb *ContainerToolCommandBase) SetArgs(args []string) {
	ctb.Command.SetArgs(args)
}

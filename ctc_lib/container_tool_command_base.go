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

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/flags"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/help"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/sub_command"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/util"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/version"
	"github.com/spf13/cobra"
)

type ContainerToolCommandBase struct {
	*cobra.Command
	Phase   string
	Version string
}

func (ctb *ContainerToolCommandBase) init() {
	ctb.AddSubCommands()
	ctb.AddFlags()
}

func (ctb *ContainerToolCommandBase) AddSubCommands() {
	// Add version subcommand
	ctb.AddCommand(version.NewVersionCommand(ctb.Version, ctb.Command.Name()).ContainerToolSubCommand)

	// Set up Root Command
	ctb.Command.SetHelpTemplate(help.HelpTemplate)
}

func (ctb *ContainerToolCommandBase) AddCommand(subCommand *sub_command.ContainerToolSubCommand) {
	cobraRun := func(c *cobra.Command, args []string) {
		obj, _ := subCommand.RunO(c, args)
		subCommand.Output = obj
		util.ExecuteTemplate(flags.TemplateString, subCommand.Output, ctb.OutOrStdout())
	}
	subCommand.Command.Run = cobraRun
	ctb.Command.AddCommand(subCommand.Command)
}

func (ctb *ContainerToolCommandBase) AddFlags() {
	// Add template Flag
	ctb.PersistentFlags().StringVarP(&flags.TemplateString, "template", "t", "{{.}}", "Output format")

	// Add NoExit Flag used for testing
	ctb.PersistentFlags().BoolVar(&flags.NoExit, "noexit", false,
		"Do not Exit with Status 1 on Error. Used for Testing")
	ctb.PersistentFlags().MarkHidden("noexit")
}

func (ctb *ContainerToolCommandBase) CheckValidCommand(runODefined bool) error {
	if (ctb.Run != nil || ctb.RunE != nil) && runODefined {
		return errors.New("Cannot provide both Command.Run and RunO implementation." +
			"\nEither implement Command.Run implementation or RunO implemetation")
	}
	return nil
}

// This function is used for Testing.
func (ctb *ContainerToolCommandBase) SetArgs(args []string) {
	ctb.Command.SetArgs(args)
}

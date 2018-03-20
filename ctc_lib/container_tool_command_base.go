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
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/constants"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/flags"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ContainerToolCommandBase struct {
	*cobra.Command
	Phase           string
	DefaultTemplate string
}

func (ctb *ContainerToolCommandBase) GetCommand() *cobra.Command {
	return ctb.Command
}

func (ctb *ContainerToolCommandBase) SetRun(cobraRun func(c *cobra.Command, args []string)) {
	ctb.Run = cobraRun
}

func (ctb *ContainerToolCommandBase) Init() {
	toolName = ctb.Name()
	cobra.OnInitialize(initConfig, initLogging)
	ctb.AddFlags()
	ctb.AddSubCommands()
}

func (ctb *ContainerToolCommandBase) AddSubCommands() {
	// Add version subcommand
	ctb.AddCommand(VersionCommand)

	// Set up Root Command
	ctb.Command.SetHelpTemplate(HelpTemplate)
}

func (ctb *ContainerToolCommandBase) AddCommand(command CLIInterface) {
	cobraRun := func(c *cobra.Command, args []string) {
		command.PrintO(c, args)
	}
	command.SetRun(cobraRun)
	ctb.Command.AddCommand(command.GetCommand())
}

func (ctb *ContainerToolCommandBase) AddFlags() {
	// Add template Flag
	ctb.PersistentFlags().StringVarP(&flags.TemplateString, "template", "t", constants.EmptyTemplate, "Output format")
	ctb.PersistentFlags().VarP(types.NewLogLevel(constants.DefaultLogLevel, &flags.LogLevel), "loglevel", "l", "LogLevel")
	ctb.PersistentFlags().BoolVarP(&flags.UpdateCheck, "updateCheck", "u", true, "Run Update Check")
	viper.BindPFlag("updateCheck", ctb.PersistentFlags().Lookup("updateCheck"))

	// Also Log to StdOut
	ctb.PersistentFlags().BoolVar(&flags.AlsoLogToStdOut, "alsoLogToStdOut", false, "Also Log to Std Out")
}

func (ctb *ContainerToolCommandBase) ReadTemplateFromFlagOrCmdDefault() string {
	if flags.TemplateString == constants.EmptyTemplate && ctb.DefaultTemplate != "" {
		return ctb.DefaultTemplate
	}
	return flags.TemplateString
}

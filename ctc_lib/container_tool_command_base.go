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
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	// Init Logger and set the Output to StdOut or Error
	Log = log.New()
	// Initialize Global flags in PrePersistentRun Hook.
	ctb.initGlobalFlags()
	ctb.AddSubCommands()
	ctb.AddFlags()

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
	ctb.PersistentFlags().StringVarP(&flags.TemplateString, "template", "t", emptyTemplate, "Output format")
	ctb.PersistentFlags().VarP(types.NewLogLevel(defaultLogLevel, &flags.LogLevel), "loglevel", "l", "LogLevel")
}

func (ctb *ContainerToolCommandBase) ReadTemplateFromFlagOrCmdDefault() string {
	if flags.TemplateString == emptyTemplate && ctb.DefaultTemplate != "" {
		return ctb.DefaultTemplate
	}
	return flags.TemplateString
}

// The command line flags are not parsed until Execute is called.
// Hence all logic to init global flags should go in this function.
func (ctb *ContainerToolCommandBase) initGlobalFlags() {
	init_persisted_flags := func(cmd *cobra.Command, args []string) {
		Log.SetLevel(flags.LogLevel.Level)
	}
	if ctb.PersistentPreRun == nil && ctb.PersistentPreRunE == nil {
		ctb.PersistentPreRun = init_persisted_flags
		return
	}
	// Cobra run PersistentPreRunE first if it is set.
	if ctb.PersistentPreRunE != nil {
		existingPreRun := ctb.PersistentPreRunE
		ctb.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			init_persisted_flags(cmd, args)
			return existingPreRun(cmd, args)
		}
	} else {
		existingPreRun := ctb.PersistentPreRun
		ctb.PersistentPreRun = func(cmd *cobra.Command, args []string) {
			init_persisted_flags(cmd, args)
			existingPreRun(cmd, args)
		}
	}
}

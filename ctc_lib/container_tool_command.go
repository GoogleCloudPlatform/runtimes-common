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

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/flags"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/help"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/sub_command"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/util"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/version"
	"github.com/spf13/cobra"
)

type ContainerToolCommand struct {
	*cobra.Command
	Phase    string
	Version  string
	HelpText string
	Output   interface{}
	RunO     func(command *cobra.Command, args []string) (interface{}, error)
}

type ContainerToolListCommand struct {
	*cobra.Command
	Phase      string
	Version    string
	OutputList []interface{}
	RunO       func(command *cobra.Command, args []string) (interface{}, error)
}

func (ctc *ContainerToolCommand) init() {
	ctc.AddSubCommands()
	ctc.AddFlags()
}

func (ctc *ContainerToolCommand) AddSubCommands() {
	// Add version subcommand
	ctc.AddCommand(version.NewVersionCommand(ctc.Version, ctc.Command.Name()).ContainerToolSubCommand)
	// Add help subcommand
	ctc.AddCommand(help.NewHelpCommand(ctc.HelpText, ctc.Command.Name()).ContainerToolSubCommand)

	// Set up Root Command
	ctc.Command.SetHelpTemplate(help.HelpTemplate)
}

func (ctc *ContainerToolCommand) AddCommand(subCommand *sub_command.ContainerToolSubCommand) {
	cobraRun := func(c *cobra.Command, args []string) {
		obj, _ := subCommand.RunO(c, args)
		subCommand.Output = obj
		util.ExecuteTemplate(flags.TemplateString, subCommand.Output, ctc.OutOrStdout())
	}
	subCommand.Command.Run = cobraRun
	ctc.Command.AddCommand(subCommand.Command)
}

func (ctc *ContainerToolCommand) AddFlags() {
	// Add template Flag
	ctc.Command.PersistentFlags().StringVarP(&flags.TemplateString, "template", "t", "{{.}}", "Output format")

	// Add NoExit Flag used for testing
	ctc.Command.PersistentFlags().BoolVar(&flags.NoExit, "noexit", false,
		"Do not Exit with Status 1 on Error. Used for Testing")
	ctc.Command.Flags().MarkHidden("noexit")

}

func (ctc *ContainerToolCommand) Execute() (err error) {
	defer errRecover(&err)

	ctc.init()
	fmt.Println(ctc.PersistentFlags().GetBool("noexit"))
	if (ctc.Command.Run != nil || ctc.Command.RunE != nil) && ctc.RunO != nil {
		panic("Cannot provide both Command.Run and RunO implementation." +
			"\nEither implement Command.Run implementation or RunO implemetation")
	}
	cobraRun := func(c *cobra.Command, args []string) {
		obj, _ := ctc.RunO(c, args)
		ctc.Output = obj
		util.ExecuteTemplate(flags.TemplateString, obj, ctc.OutOrStdout())
	}

	ctc.Command.Run = cobraRun
	err = ctc.Command.Execute()
	if err != nil {
		panic(err)
	}
	return err
}

// This function is used for Testing.
func (ctc *ContainerToolCommand) SetArgs(args []string) {
	ctc.Command.SetArgs(args)
}

// errRecover is the handler that turns panics into returns from the top
// level of Parse.
func errRecover(errp *error) {
	if e := recover(); e != nil {
		// TODO: Change this to Log.Error once Logging is introduced.
		fmt.Println(e, flags.NoExit)
		fmt.Println(reflect.TypeOf(e))
		if !flags.NoExit {
			os.Exit(1)
		}
		*errp = errors.New(fmt.Sprintf("%v", e))
	}
}

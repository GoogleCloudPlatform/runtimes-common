/*
Copyright 2018 Google, Inc. All rights reserved.
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

package help

import (
	"strings"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/sub_command"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type HelpCommand struct {
	*sub_command.ContainerToolSubCommand
}

type HelpOutput struct {
	HelpText string
}

var HelpTemplate = `
{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

func NewHelpCommand(helpText string, toolname string) *HelpCommand {
	var command = &cobra.Command{
		Use:   "help [command]",
		Short: "Help about the command",
		RunE: func(c *cobra.Command, args []string) error {
			cmd, args, e := c.Root().Find(args)
			if cmd == nil || e != nil || len(args) > 0 {
				return errors.Errorf("unknown help topic: %v", strings.Join(args, " "))
			}

			helpFunc := cmd.HelpFunc()
			helpFunc(cmd, args)
			return nil
		},
	}
	var helpCommand = &HelpCommand{
		ContainerToolSubCommand: &sub_command.ContainerToolSubCommand{
			Command: command,
			Output:  &HelpOutput{},
		},
	}
	return helpCommand
}

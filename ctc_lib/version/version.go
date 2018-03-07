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

package version

import (
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/sub_command"
	"github.com/spf13/cobra"
)

type VersionCommand struct {
	*sub_command.ContainerToolSubCommand
}

type VersionOutput struct {
	Version string
}

func NewVersionCommand(version string, toolname string) *VersionCommand {
	var command = &cobra.Command{
		Use:   "version",
		Short: "Print the version of " + toolname,
		Long:  `Print the version of ` + toolname,
		Args:  cobra.ExactArgs(0),
	}
	var versionCommand = &VersionCommand{
		ContainerToolSubCommand: &sub_command.ContainerToolSubCommand{
			Command: command,
			Output:  &VersionOutput{},
			RunO: func(command *cobra.Command, args []string) (interface{}, error) {
				var versionOutput = VersionOutput{
					Version: version,
				}
				return versionOutput, nil
			},
		},
	}

	return versionCommand
}

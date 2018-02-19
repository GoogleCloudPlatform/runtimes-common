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
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/container_command_output"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/version"
	"github.com/spf13/cobra"
)

type ContainerToolCommand struct {
	*cobra.Command
	Phase         string
	Version       string
	CommandOutput container_command_output.ContainerCommandOutput
}

type ContainerToolCommandList struct {
	*cobra.Command
	Phase             string
	Version           string
	CommandOutputList container_command_output.ContainerCommandOutputList
}

func (ctc ContainerToolCommand) init() {
	// Add version subcommand
	ctc.Command.AddCommand(version.NewVersionCommand(ctc.Version,
		ctc.Command.Name()).Command)
	// Add help subcommand

}

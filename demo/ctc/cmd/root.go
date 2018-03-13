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

package cmd

import (
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib"
	"github.com/spf13/cobra"
)

type RootCommandOutput struct {
	Message string
}

var Message string

var RootCommand = &ctc_lib.ContainerToolCommand{
	ContainerToolCommandBase: &ctc_lib.ContainerToolCommandBase{
		Command: &cobra.Command{
			Use:   "echo",
			Short: "Echo's Message",
		},
		Phase: "test",
	},
	Output: &RootCommandOutput{},
	RunO: func(command *cobra.Command, args []string) (interface{}, error) {
		return RootCommandOutput{
			Message: Message,
		}, nil
	},
}

func Execute() {
	ctc_lib.Version = "1.0.1"
	ctc_lib.Execute(RootCommand)
}

func init() {
	RootCommand.Flags().StringVarP(&Message, "message", "m", "text", "Message to Echo")
}

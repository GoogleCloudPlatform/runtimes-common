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
	"errors"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib"
	"github.com/spf13/cobra"
)

// Flags
var tufConfigFilename string

// RootCommand
var RootCommand = &ctc_lib.ContainerToolCommand{
	ContainerToolCommandBase: &ctc_lib.ContainerToolCommandBase{
		Command: &cobra.Command{
			Use: "TUF command",
		},
		Phase:           "test",
		DefaultTemplate: "{{.}}",
	},
	RunO: func(command *cobra.Command, args []string) (interface{}, error) {
		return nil, errors.New("Please Specify a Sub Command")
	},
}

func init() {
	RootCommand.PersistentFlags().StringVar(&tufConfigFilename, "config", "tuf.yaml", "File name for Tool config")
	RootCommand.AddCommand(GenerateKeyCommand)
	RootCommand.AddCommand(UploadSecretsCommand)
}

func Execute() {
	ctc_lib.Version = "0.0.1"
	ctc_lib.Execute(RootCommand)
}

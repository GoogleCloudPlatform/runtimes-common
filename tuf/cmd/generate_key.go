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
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/scheme"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
	"github.com/spf13/cobra"
)

var cryptoScheme types.CryptoScheme
var filename string

// Command To upload Secrets to GCS store.
var GenerateKeyCommand = &ctc_lib.ContainerToolCommand{
	ContainerToolCommandBase: &ctc_lib.ContainerToolCommandBase{
		Command: &cobra.Command{
			Use: "generate-secret key value pair",
			RunE: func(command *cobra.Command, args []string) error {
				var err error
				switch cryptoScheme.Scheme {
				case types.ECDSA256:
					err = scheme.NewECDSA().Store(filename)
				default:
					err = scheme.NewECDSA().Store(filename)
				}
				return err
			},
		},
		Phase: "test",
	},
}

func init() {
	GenerateKeyCommand.Flags().Var(types.NewCryptoScheme("ECDSA256", &cryptoScheme), "scheme", "Generate Public/Private key pair and store it")
	GenerateKeyCommand.Flags().StringVar(&filename, "file", "keys.json", "File name to store the secret in json format")
}

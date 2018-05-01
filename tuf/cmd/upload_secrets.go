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
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/config"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/deployment"

	"github.com/spf13/cobra"
)

// Flags
var rootKey = "~/root.json"
var targetKey = "~/target.json"
var snapshotKey = "~/snapshot.json"

// Command To upload Secrets to GCS store and Renegerate all the MetaData.
var UploadSecretsCommand = &ctc_lib.ContainerToolCommand{
	ContainerToolCommandBase: &ctc_lib.ContainerToolCommandBase{
		Command: &cobra.Command{
			Use: "Prototype GCS Update Keys to Google Cloud Storage.",
		},
		Phase:           "test",
		DefaultTemplate: "{{.}}",
	},
	RunO: func(command *cobra.Command, args []string) (interface{}, error) {
		tufConfig, err := config.ReadConfig(tufConfigFilename)
		if err != nil {
			return nil, err
		}
		deployment.UpdateSecrets(tufConfig, rootKey, targetKey, snapshotKey)
		return nil, nil
	},
}

func init() {
	RootCommand.Flags().StringVar(&rootKey, "root-key", "", "Secret key.json for Root role")
	RootCommand.Flags().StringVar(&targetKey, "target-key", "", "Secret key.json for Snapshot role")
	RootCommand.Flags().StringVar(&snapshotKey, "snapshot-key", "", "Secret key.json for Target role")
}

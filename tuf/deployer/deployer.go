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
package deployer

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/config"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/gcs"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/kms"
)

type KeyPair struct {
	Public  string
	Private string
}

type DeployTool interface {
	UpdateSecrets(config.TUFConfig, string, string, string) error
}

type Deployer struct {
	KmsService kms.KMService
	Storage    gcs.Storage
}

func New() (DeployTool, error) {
	kmsService, err := kms.New()
	if err != nil {
		return nil, err
	}
	gcsClient, err := gcs.New()
	if err != nil {
		return nil, err
	}
	return &Deployer{
		KmsService: kmsService,
		Storage:    gcsClient,
	}, nil
}

func (d *Deployer) UpdateSecrets(tufConfig config.TUFConfig, rootKeyFile string, targetKeyFile string, snapshotKeyFile string) error {
	errorStr := make([]string, 0)
	if rootKeyFile != "" {
		errorStr = append(errorStr, d.uploadSecret(rootKeyFile, tufConfig, config.RootSecretFileName).Error())
	}
	if targetKeyFile != "" {
		errorStr = append(errorStr, d.uploadSecret(targetKeyFile, tufConfig, config.TargetSecretFileName).Error())
	}
	if snapshotKeyFile != "" {
		errorStr = append(errorStr, d.uploadSecret(snapshotKeyFile, tufConfig, config.SnapshotSecretFileName).Error())
	}
	if len(errorStr) > 0 {
		// Exit if there were errors uploading secrets.
		return fmt.Errorf("Encountered following errors %s", strings.Join(errorStr, "\n"))
	}

	// TODO Generate all the Metadata.

	// TODO Write Consistent Snapshots
	return nil
}

func (d *Deployer) uploadSecret(file string, tufConfig config.TUFConfig, name string) error {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	encyptedResponse, err := d.KmsService.Encrypt(kms.CryptoKeyFromConfig(tufConfig), string(text))
	tmpFile, errWrite := ioutil.TempFile("", "key")
	defer os.Remove(tmpFile.Name())
	if errWrite != nil {
		return err
	}
	ioutil.WriteFile(tmpFile.Name(), []byte(encyptedResponse.Ciphertext), os.ModePerm)
	tmpFile.Close()

	_, _, err = d.Storage.Upload(tufConfig.GCSProjectID, tufConfig.GCSBucketID, name, tmpFile)
	return err
}

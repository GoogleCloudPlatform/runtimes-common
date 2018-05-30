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

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/config"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/gcs"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/kms"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/metadata/v1"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/scheme"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

type DeployTool interface {
	UpdateSecrets(config.TUFConfig, string, string, string) error
	GenerateMetadata(tufConfig config.TUFConfig, root string, target string, snapshot string, oldroot []byte) (
		*v1.RootMetadata, *v1.TargetMetadata, *v1.SnapshotMetadata, error)
}

type Deployer struct {
	KmsService kms.KMService
	Storage    gcs.Store
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
	var oldRootSecretBytes []byte
	var err error
	if rootKeyFile != "" {
		// If root secret is changed, first download the old root secret.
		// We need to sign the metadata with old as well as new root secret.
		oldRootSecretBytes, err = d.Storage.Download(tufConfig.GCSBucketID, config.RootSecretFileName)
		if err != nil && err != storage.ErrObjectNotExist {
			// The old root file exists but there was an error reading it. This is fatal hence return error
			return err
		}
		err = d.uploadSecret(rootKeyFile, tufConfig, config.RootSecretFileName)
		if err != nil {
			errorStr = append(errorStr, err.Error())
		}
	}

	if targetKeyFile != "" {
		err := d.uploadSecret(targetKeyFile, tufConfig, config.TargetSecretFileName)
		if err != nil {
			errorStr = append(errorStr, err.Error())
		}
	}
	if snapshotKeyFile != "" {
		err := d.uploadSecret(snapshotKeyFile, tufConfig, config.SnapshotSecretFileName)
		if err != nil {
			errorStr = append(errorStr, err.Error())
		}
	}
	if len(errorStr) > 0 {
		// Exit if there were errors uploading secrets.
		return fmt.Errorf("Encountered following errors %s", strings.Join(errorStr, "\n"))
	}

	d.GenerateMetadata(tufConfig, rootKeyFile, targetKeyFile, snapshotKeyFile, oldRootSecretBytes)

	// TODO Write Consistent Snapshots
	return nil
}

func (d *Deployer) uploadSecret(file string, tufConfig config.TUFConfig, name string) error {
	textBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	encyptedResponse, err := d.KmsService.Encrypt(kms.CryptoKeyFromConfig(tufConfig), string(textBytes))
	tmpFile, errWrite := ioutil.TempFile("", "key")
	defer os.Remove(tmpFile.Name())
	if errWrite != nil {
		return err
	}
	ioutil.WriteFile(tmpFile.Name(), []byte(encyptedResponse.Ciphertext), os.ModePerm)
	_, _, err = d.Storage.Upload(tufConfig.GCSBucketID, name, tmpFile)
	return err
}

func (d *Deployer) GenerateMetadata(tufConfig config.TUFConfig, root string,
	target string, snapshot string, oldroot []byte) (*v1.RootMetadata, *v1.TargetMetadata, *v1.SnapshotMetadata, error) {
	targetKey, err := d.readSecret(tufConfig.GCSBucketID, config.TargetSecretFileName, target)
	if err != nil {
		return nil, nil, nil, err
	}
	snapshotKey, err := d.readSecret(tufConfig.GCSBucketID, config.SnapshotSecretFileName, snapshot)
	if err != nil {
		return nil, nil, nil, err
	}
	rootKey, err := d.readSecret(tufConfig.GCSBucketID, config.RootSecretFileName, root)
	if err != nil {
		return nil, nil, nil, err
	}
	// Once you have all the secrets, populate the structs.
	rootMetadata := v1.PopulateRootMetadata(rootKey, snapshotKey, targetKey)
	return &rootMetadata, &v1.TargetMetadata{}, &v1.SnapshotMetadata{}, nil
}

func (d *Deployer) readSecret(bucketID string, key string, secretFile string) (types.Scheme, error) {
	if secretFile != "" {
		// Read secret from the target file.
		return scheme.Read(secretFile)
	}
	// Download the secret from GCS.
	bytes, err := d.Storage.Download(bucketID, key)
	if err != nil {
		return nil, err
	}
	return scheme.ReadBytes(bytes)
}

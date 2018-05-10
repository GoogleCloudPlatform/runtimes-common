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

package integrationtest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/cmd"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/config"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/gcs"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/kms"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/testutil"
)

var rootSecret = fmt.Sprintf("This is root file secret number %d", rand.Int())
var targetSecret = fmt.Sprintf("This is target file secret number %d", rand.Int())
var snapshotSecret = fmt.Sprintf("This is snapshot file secret number %d", rand.Int())

func TestUploadSecretsCommand(t *testing.T) {
	rootFile, _ := ioutil.TempFile("", "rawSecret1.json")
	defer os.Remove(rootFile.Name())
	ioutil.WriteFile(rootFile.Name(), []byte(rootSecret), 644)

	targetFile, _ := ioutil.TempFile("", "rawSecret2.json")
	defer os.Remove(targetFile.Name())
	ioutil.WriteFile(targetFile.Name(), []byte(targetSecret), 644)

	snapshotFile, _ := ioutil.TempFile("", "rawSecret3.json")
	defer os.Remove(snapshotFile.Name())
	ioutil.WriteFile(snapshotFile.Name(), []byte(snapshotSecret), 644)

	tufConfig, _ := ioutil.TempFile("", "tufConfig.yaml")
	defer os.Remove(tufConfig.Name())

	buf, err := yaml.Marshal(&testutil.IntegrationTufConfig)
	if err != nil {
		t.Fatalf("Error while writing config %v", err)
	}
	ioutil.WriteFile(tufConfig.Name(), buf, 644)

	prevCred := os.Getenv(testutil.GoogleCredConstant)
	os.Setenv(testutil.GoogleCredConstant, os.Getenv(testutil.TufIntegrationConstant))

	defer os.Setenv(testutil.GoogleCredConstant, prevCred)

	cmd.RootCommand.SetArgs([]string{"upload-secrets",
		"--config", tufConfig.Name(),
		"--root-key", rootFile.Name(),
		"--target-key", targetFile.Name(),
		"--snapshot-key", snapshotFile.Name()})

	err = cmd.RootCommand.Execute()

	if err != nil {
		t.Fatalf("Unexpected Err: %v", err)
	}

	// Download Files from GCS
	rootFileEncrypt, _ := ioutil.TempFile("", "encryptSecret1.json")
	defer os.Remove(rootFileEncrypt.Name())

	targetFileEncrypt, _ := ioutil.TempFile("", "encryptSecret2.json")
	defer os.Remove(targetFileEncrypt.Name())

	snapshotFileEncrpyt, _ := ioutil.TempFile("", "encryptSecret3.json")
	defer os.Remove(snapshotFile.Name())

	err = downloadAndVerifyFiles(
		testutil.IntegrationTufConfig,
		rootFileEncrypt.Name(),
		targetFileEncrypt.Name(),
		snapshotFileEncrpyt.Name(),
	)

	if err != nil {
		t.Fatalf("Unexpected Error %v", err)
	}

}

func downloadAndVerifyFiles(tufConfig config.TUFConfig, rootFile string, targetFile string, snapshotFile string) error {
	errorStrings := make([]string, 0)
	gcsService, err := gcs.New()
	if err != nil {
		return err
	}
	defer cleanAllStorage(gcsService, tufConfig.GCSBucketID)
	errorStrings = appendErrorIfExists(errorStrings, downloadFile(gcsService, tufConfig.GCSBucketID, config.RootSecretFileName, rootFile))
	errorStrings = appendErrorIfExists(errorStrings, downloadFile(gcsService, tufConfig.GCSBucketID, config.TargetSecretFileName, targetFile))
	errorStrings = appendErrorIfExists(errorStrings, downloadFile(gcsService, tufConfig.GCSBucketID, config.SnapshotSecretFileName, snapshotFile))

	// Decrypt the file and see if its same as.
	kmsService, err := kms.New()
	if err != nil {
		return err
	}
	errorStrings = appendErrorIfExists(errorStrings, decryptFile(kmsService, tufConfig, rootFile, rootSecret))
	errorStrings = appendErrorIfExists(errorStrings, decryptFile(kmsService, tufConfig, targetFile, targetSecret))
	errorStrings = appendErrorIfExists(errorStrings, decryptFile(kmsService, tufConfig, snapshotFile, snapshotSecret))
	if len(errorStrings) > 0 {
		return errors.New(strings.Join(errorStrings, "\n"))
	}
	return nil
}

func downloadFile(gcsService *gcs.GCSStore, bucketID string, key string, dest string) error {
	err := gcsService.Download(bucketID, key, dest)
	if err != nil {
		return err
	}
	return nil
}

func decryptFile(kmsService *kms.KMS, tufConfig config.TUFConfig, decryptFile string, plainTextExp string) error {
	bytes, _ := ioutil.ReadFile(decryptFile)
	plainText, decryptErr := kmsService.Decrypt(kms.CryptoKeyFromConfig(tufConfig), string(bytes))
	if decryptErr != nil {
		return decryptErr
	} else if plainText != plainTextExp {
		return errors.New(fmt.Sprintf("Expected: %v\nGot: %v", plainTextExp, plainText))
	}
	return nil
}

func cleanAllStorage(gcsService *gcs.GCSStore, bucketID string) {
	err1 := gcsService.Delete(bucketID, config.RootSecretFileName)
	err2 := gcsService.Delete(bucketID, config.TargetSecretFileName)
	err3 := gcsService.Delete(bucketID, config.SnapshotSecretFileName)
	if err1 != nil || err2 != nil || err3 != nil {
		panic(fmt.Sprintf("Error cleaning buckts %v, %v, %v", err1, err2, err3))
	}
}

func appendErrorIfExists(errString []string, err error) []string {
	if err != nil {
		return append(errString, err.Error())
	}
	return errString
}

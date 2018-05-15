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

var testFiles = make([]string, 0)

func createAndWriteFile(filename string, text string) string {
	tmpFile, err := ioutil.TempFile("", filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot run tests due to %v", err))
	}
	testFiles = append(testFiles, tmpFile.Name())
	if text != "" {
		ioutil.WriteFile(tmpFile.Name(), []byte(text), 644)
	}
	return tmpFile.Name()
}

func cleanUpFiles() {
	for _, file := range testFiles {
		os.Remove(file)
	}
}

func TestUploadSecretsCommand(t *testing.T) {
	rootFile := createAndWriteFile("rawSecret1.json", rootSecret)
	targetFile := createAndWriteFile("rawSecret1.json", targetSecret)
	snapshotFile := createAndWriteFile("rawSecret1.json", snapshotSecret)

	buf, err := yaml.Marshal(&testutil.IntegrationTufConfig)
	if err != nil {
		t.Fatalf("Error while writing config %v", err)
	}
	tufConfig := createAndWriteFile("tufConfig.yaml", string(buf))

	defer cleanUpFiles()

	cmd.RootCommand.SetArgs([]string{"upload-secrets",
		"--config", tufConfig,
		"--root-key", rootFile,
		"--target-key", targetFile,
		"--snapshot-key", snapshotFile})

	err = cmd.RootCommand.Execute()

	if err != nil {
		t.Fatalf("Unexpected Err: %v", err)
	}

	err = downloadAndVerifySecrets(testutil.IntegrationTufConfig)

	if err != nil {
		t.Fatalf("Unexpected Error %v", err)
	}

}

func downloadAndVerifySecrets(tufConfig config.TUFConfig) error {
	errorStrings := make([]string, 0)
	gcsService, err := gcs.New()
	if err != nil {
		return err
	}
	defer cleanAllStorage(gcsService, tufConfig.GCSBucketID)
	rootBytes, err := downloadFile(gcsService, tufConfig.GCSBucketID, config.RootSecretFileName)
	errorStrings = appendErrorIfExists(errorStrings, err)
	targetBytes, err := downloadFile(gcsService, tufConfig.GCSBucketID, config.TargetSecretFileName)
	errorStrings = appendErrorIfExists(errorStrings, err)
	snapshotBytes, err := downloadFile(gcsService, tufConfig.GCSBucketID, config.SnapshotSecretFileName)
	errorStrings = appendErrorIfExists(errorStrings, err)

	// Decrypt the file and see if its same as.
	kmsService, err := kms.New()
	if err != nil {
		return err
	}
	errorStrings = appendErrorIfExists(errorStrings, decryptFile(kmsService, tufConfig, rootBytes, rootSecret))
	errorStrings = appendErrorIfExists(errorStrings, decryptFile(kmsService, tufConfig, targetBytes, targetSecret))
	errorStrings = appendErrorIfExists(errorStrings, decryptFile(kmsService, tufConfig, snapshotBytes, snapshotSecret))
	if len(errorStrings) > 0 {
		return errors.New(strings.Join(errorStrings, "\n"))
	}
	return nil
}

func downloadFile(gcsService *gcs.GCSStore, bucketID string, key string) ([]byte, error) {
	bytes, err := gcsService.Download(bucketID, key)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func decryptFile(kmsService *kms.KMS, tufConfig config.TUFConfig, decryptBytes []byte, plainTextExp string) error {
	plainText, decryptErr := kmsService.Decrypt(kms.CryptoKeyFromConfig(tufConfig), string(decryptBytes))
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

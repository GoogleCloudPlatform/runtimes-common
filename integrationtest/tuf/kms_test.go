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
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/kms"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/testutil"
)

func TestKMSIntegration(t *testing.T) {

	kmsService, err := kms.New()

	if err != nil {
		t.Fatalf("Failed creating KMS client. %v", err)
	}
	encryptResp, err := kmsService.Encrypt(kms.CryptoKeyFromConfig(testutil.IntegrationTufConfig), "this is secret")
	if err != nil {
		t.Fatalf("Unexpected error when encrypting. %v", err)
	}
	plainText, decryptErr := kmsService.Decrypt(kms.CryptoKeyFromConfig(testutil.IntegrationTufConfig), encryptResp.Ciphertext)
	if decryptErr != nil {
		t.Fatalf("Unexpected error. %v", decryptErr)
	}
	if plainText != "this is secret" {
		t.Fatalf("Expected: this is secret\nGot: %v", plainText)
	}

}

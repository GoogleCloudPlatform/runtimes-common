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
package kms

import (
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/config"
)

func TestCryptoKeyFromConfig(t *testing.T) {
	tufConfig := config.TUFConfig{
		KMSProjectID: "project",
		KMSLocation:  "global",
		KeyRingID:    "testring",
		CryptoKeyID:  "key",
	}
	cryptoKey := CryptoKeyFromConfig(tufConfig)
	if cryptoKey.Name() != "projects/project/locations/global/keyRings/testring/cryptoKeys/key" {
		t.Fatalf(`\nExpected: \n\t projects/project/locations/global/keyRings/testring/cryptoKeys/key
Got: \n\t%s`, cryptoKey.Name())
	}
}

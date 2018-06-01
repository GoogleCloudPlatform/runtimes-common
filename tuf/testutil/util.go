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

package testutil

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/config"
)

var TestTUFConfig = config.TUFConfig{
	GCSProjectID:      "testProjectID",
	KMSProjectID:      "testKmsProjectID",
	KMSLocation:       "global",
	KeyRingID:         "testKeyRing",
	CryptoKeyID:       "testKey",
	GCSBucketID:       "testBucket",
	RootThreshold:     1,
	TargetThreshold:   1,
	SnapshotThreshold: 1,
	SpecVersion:       1,
	Targets:           []string{"file1.txt", "file22.txt"},
}

var IntegrationTufConfig = config.TUFConfig{
	KMSProjectID: "gcp-runtimes",
	KMSLocation:  "global",
	KeyRingID:    "tuftestkeyring",
	CryptoKeyID:  "testkey",
	GCSProjectID: "gcp-runtimes",
	GCSBucketID:  "tuf-integration",
}

func IsErrorEqualOrContains(err error, subErr error) bool {
	if err == nil && subErr == nil {
		return true // Return true Both of them are nil
	} else if err == nil || subErr == nil {
		return false // Return false if either of them are nil
	} else if strings.Contains(err.Error(), subErr.Error()) {
		return true // Return true if Messages are equal
	}
	return false // Return false
}

func CreateAndWriteFile(dir string, filename string, text string) string {
	tmpFile, err := ioutil.TempFile(dir, filename)
	if err != nil {
		panic(fmt.Sprintf("Cannot run tests due to %v", err))
	}

	if text != "" {
		ioutil.WriteFile(tmpFile.Name(), []byte(text), 644)
	}
	return tmpFile.Name()
}

func MarshalledTUFConfig() string {
	bytes, err := yaml.Marshal(TestTUFConfig)
	if err != nil {
		panic(fmt.Sprintf("Cannot run tests due to %v", err))
	}
	return string(bytes)
}

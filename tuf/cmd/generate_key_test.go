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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/testutil"
)

var testCases = []struct {
	name             string
	scheme           string
	expectedErr      error
	shouldFileExists bool
}{
	{"ecdsaSchemeSuccess", "ECDSA256", nil, true},
	{"invalidScheme", "ECDSA1", errors.New("not a valid CryptoScheme"), false},
	{"notImplementedScheme", "RSA256", errors.New("Not Implemented Yet"), false},
}

func TestGenerateKey(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpfile, _ := ioutil.TempFile("", "testKey.json")
			defer os.Remove(tmpfile.Name())
			RootCommand.SetArgs([]string{"generate-secret", "--scheme", tc.scheme, "--file", tmpfile.Name()})
			err := RootCommand.Execute()
			if !testutil.IsErrorEqualOrContains(err, tc.expectedErr) {
				t.Fatalf("Expected Err: %v\nGot: %v", tc.expectedErr, err)
			}
			// Check if output file exists.
			_, err = os.Stat(filename)
			if !os.IsNotExist(err) != tc.shouldFileExists {
				t.Fatalf("Expected file to exist: %t\nGot %v",
					tc.shouldFileExists, err)
			}
			// Check if Json file can be read properly.
			if tc.shouldFileExists {
				raw, err := ioutil.ReadFile(tmpfile.Name())
				if err != nil {
					t.Fatalf("Error while reading the test report %v", err)
				}
				c := map[string]interface{}{}
				json.Unmarshal(raw, &c)
				fmt.Println(c)
				if s := c["KeyType"]; s != tc.scheme {
					t.Fatalf("Expected File to contain Scheme: %s.\nGot: %s", s, tc.scheme)
				}
			}
		})
	}
}

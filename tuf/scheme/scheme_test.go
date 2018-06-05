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
package scheme

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/testutil"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

var testReadTC = []struct {
	name          string
	inputKey      types.Scheme
	expectedError error
}{
	{"testSuccess", NewECDSA(), nil},
	{"testFail", &testutil.TestKey{PrivateKey: "secret", KeyType: "test"},
		fmt.Errorf("Could not parse key test")},
}

func TestReadFile(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "test_read_")
	defer os.RemoveAll(tmpdir)
	for _, tc := range testReadTC {
		tmpfile, _ := ioutil.TempFile(tmpdir, tc.name)
		t.Run(tc.name, func(t *testing.T) {
			tc.inputKey.Store(tmpfile.Name())
			readKey, err := Read(tmpfile.Name())
			if !reflect.DeepEqual(tc.expectedError, err) {
				t.Fatalf("\nExpected Error %v \nGot %v", tc.expectedError.Error(), err.Error())
			}
			if tc.expectedError == nil && !reflect.DeepEqual(tc.inputKey, readKey) {
				t.Fatalf("Expected Key %v \n Got %v", tc.inputKey, readKey)
			}
		})
	}
}

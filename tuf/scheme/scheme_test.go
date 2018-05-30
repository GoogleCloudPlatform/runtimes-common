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
	"crypto"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

type TempKey struct {
	PrivateKey string
	KeyType    string
}

func (tmp *TempKey) Store(filename string) error {
	keyJson, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, keyJson, 0644)
}

func (tmp *TempKey) Sign(rand io.Reader, priv *crypto.PrivateKey, hash []byte) (r, s *big.Int, err error) {
	return nil, nil, nil
}

func (tmp *TempKey) Verify(pub *crypto.PublicKey, hash []byte, r, s *big.Int) bool {
	return true
}

func (tmp *TempKey) GetPublicKey() string {
	return "Test-Public"
}

func (tmp *TempKey) GetKeyId() types.KeyId {
	return "Test"
}

func (tmp *TempKey) GetKeyIdHashAlgo() []types.HashAlgo {
	return []types.HashAlgo{"sha256"}
}

func (tmp *TempKey) GetScheme() types.KeyScheme {
	return "TEST_KEY"
}

var testReadTC = []struct {
	name              string
	inputKey          types.Scheme
	expectedError     error
	expectedOutputNil bool
}{
	{"testSuccess", NewECDSA(),
		nil,
		false,
	},
	{"testFail", &TempKey{PrivateKey: "temp", KeyType: "test"},
		fmt.Errorf("Could not parse key test"),
		true,
	},
}

func TestReadFile(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("", "test_read_")
	defer os.RemoveAll(tmpdir)
	for _, tc := range testReadTC {
		tmpfile, _ := ioutil.TempFile(tmpdir, tc.name)
		t.Run(tc.name, func(t *testing.T) {
			tc.inputKey.Store(tmpfile.Name())
			readKey, err := Read(tmpfile.Name())
			if tc.expectedError != nil && err != nil && tc.expectedError.Error() != err.Error() {
				t.Fatalf("\nExpected Error %v \nGot %v", tc.expectedError.Error(), err.Error())
			} else if (tc.expectedError != nil && err == nil) || (tc.expectedError == nil && err != nil) {
				t.Fatalf("\nExpected Error %v \nGot %v", tc.expectedError, err)
			}
			if !tc.expectedOutputNil && !reflect.DeepEqual(tc.inputKey, readKey) {
				t.Fatalf("Expected Key %v \n Got %v", tc.inputKey, readKey)
			}
			if tc.expectedOutputNil && readKey != nil {
				t.Fatalf("Expected nil \n Got %v", readKey)
			}
		})
	}
}

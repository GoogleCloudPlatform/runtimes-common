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
	"bytes"
	"encoding/gob"
	"errors"
	"reflect"
	"strings"
	"testing"
)

type testinterface struct {
	Foo    string `json:"foo"`
	FooMap map[string]string
}

var ecdsaTC = []struct {
	name                 string
	metadata             interface{}
	verifyMetadata       interface{}
	expectedSignErr      error
	expectedVerifyResult bool
}{
	{"success",
		testinterface{
			Foo: "foostring",
			FooMap: map[string]string{
				"1": "one",
				"2": "two",
			},
		}, testinterface{
			Foo: "foostring",
			FooMap: map[string]string{
				"1": "one",
				"2": "two",
			},
		}, nil, true,
	},
	{"success_empty_interface", testinterface{}, testinterface{}, nil, true},
	{"success_native_string", "abc", "abc", nil, true},
	{"success_native_int", 1, 1, nil, true},
	{"fail_native_int", 1, 2, nil, false},
	{"fail",
		testinterface{Foo: "foostring"},
		testinterface{Foo: "barstring"}, nil, false,
	},
	{"sign_err", nil, nil, errors.New("gob: cannot encode nil value"), false},
}

func TestSignVerify(t *testing.T) {
	for _, tc := range ecdsaTC {
		t.Run(tc.name, func(t *testing.T) {
			ecdsaKey := NewECDSA()
			sig, err := ecdsaKey.Sign(tc.metadata)
			if !reflect.DeepEqual(tc.expectedSignErr, err) {
				t.Fatalf("Expected a error while signing %v\n Got %v", tc.expectedSignErr, err)
			}
			if tc.expectedSignErr == nil {
				var buf bytes.Buffer
				err := gob.NewEncoder(&buf).Encode(tc.verifyMetadata)
				if err != nil {
					t.Fatalf("Cannot Verify due to %v", err)
				}
				actualVerifyResult := ecdsaKey.Verify(string(buf.Bytes()), sig)
				if tc.expectedVerifyResult != actualVerifyResult {
					t.Fatalf("%s : Expected Verify result to be %t. \n Got %t", tc.name, tc.expectedVerifyResult, actualVerifyResult)
				}
			}
		})
	}
}

func TestEncodeDecode(t *testing.T) {
	ecdsaKey := NewECDSA()
	privateKey, publicKey, err := ecdsaKey.encode()
	if err != nil {
		t.Fatalf("Did not expect encoding error %v", err)
	}
	if !strings.Contains(privateKey, "BEGIN PRIVATE KEY") {
		t.Fatalf("Expected to container BEGIN PRIVATE KEY substring in %v", privateKey)
	}
	if !strings.Contains(publicKey, "BEGIN PUBLIC KEY") {
		t.Fatalf("Expected to container BEGING PUBLIC KEY substring in %v", privateKey)
	}
	// Make sure we can decode the same key using pem encoded private and public key.
	ecdsaDecodedKey := &ECDSA{}
	err = ecdsaDecodedKey.decode(privateKey)
	if err != nil {
		t.Fatalf("Did not expect decoding error %v", err)
	}
	if !reflect.DeepEqual(ecdsaDecodedKey, ecdsaKey) {
		t.Fatalf("Expected %v\nGot %v", ecdsaKey, ecdsaDecodedKey)
	}
}

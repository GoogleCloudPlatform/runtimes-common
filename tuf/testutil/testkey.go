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
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

type TestKey struct {
	PrivateKey string
	KeyType    string
	PublicKey  string
	SignStr    string
	KeyId      string
}

func (tmp *TestKey) Store(filename string) error {
	keyJson, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, keyJson, 0644)
}

func (tmp *TestKey) Sign(singedMetadata interface{}) (string, error) {
	singedMetadataStr, ok := singedMetadata.(string)
	if !ok {
		return "", errors.New("Signing Error")
	}
	return tmp.SignStr + singedMetadataStr, nil
}

func (tmp *TestKey) Verify(signingstring string, signature string) bool {
	return signature == signingstring
}

func (tmp *TestKey) GetPublicKey() string {
	return tmp.PublicKey
}

func (tmp *TestKey) GetKeyId() types.KeyId {
	return types.KeyId(tmp.KeyId)
}

func (tmp *TestKey) GetKeyIdHashAlgo() []types.HashAlgo {
	return []types.HashAlgo{"sha256"}
}

func (tmp *TestKey) GetScheme() types.KeyScheme {
	return types.KeyScheme(tmp.KeyType)
}

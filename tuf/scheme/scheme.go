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
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/constants"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

type SchemeKey struct {
	PrivateKey string          `json:"PrivateKey"`
	PublicKey  string          `json:"PublicKey"`
	KeyType    types.KeyScheme `json:"KeyType"`
}

func Read(filename string) (types.Scheme, error) {
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ReadBytes(text)
}

func ReadBytes(text []byte) (types.Scheme, error) {
	var schemeKey SchemeKey
	err := json.Unmarshal(text, &schemeKey)
	if err != nil {
		return nil, err
	}
	switch schemeKey.KeyType {
	case constants.ECDSA256Scheme:
		ecdsaKey := &ECDSA{}
		ecdsaKey.decode(schemeKey.PrivateKey)
		return ecdsaKey, err
	}
	return nil, fmt.Errorf("Could not parse key %v", schemeKey.KeyType)
}

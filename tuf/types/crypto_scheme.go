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

package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	ECDSA256 KeyScheme = "ECDSA256"
	RSA256   KeyScheme = "RSA256"
)

var VALID_CRYPTO_SCHEMES = map[KeyScheme]bool{
	ECDSA256: true,
	RSA256:   false, // Not implemented.
}

var CryptoSchemes []KeyScheme
var ImplementedSchemes []KeyScheme

type CryptoScheme struct {
	Scheme KeyScheme
}

func (scheme *CryptoScheme) String() string {
	return string(scheme.Scheme)
}

func (scheme *CryptoScheme) Set(s string) error {
	keyScheme := KeyScheme(s)
	value, ok := VALID_CRYPTO_SCHEMES[keyScheme]
	if ok && value {
		scheme.Scheme = keyScheme
		return nil
	}
	if !ok {
		return fmt.Errorf(`%s is not a valid CryptoScheme.
		Please Provide one of %s`, s, JoinKeyScheme(CryptoSchemes, ", "))
	}
	return fmt.Errorf(`%s is not a Not Implemented Yet!
		Please Provide one of %s`, s, JoinKeyScheme(ImplementedSchemes, ", "))
}

func (scheme *CryptoScheme) Type() string {
	return "types.CryptoScheme"
}

func NewCryptoScheme(val KeyScheme, p *CryptoScheme) *CryptoScheme {
	value, ok := VALID_CRYPTO_SCHEMES[val]
	if ok && value {
		*p = CryptoScheme{
			Scheme: val,
		}
		return p
	}
	return nil
}

func (scheme *CryptoScheme) Store(filename string) error {
	schemeJson, err := json.Marshal(scheme)
	if err != nil {
		return fmt.Errorf("Error while marshalling json %s", err.Error())
	}
	return ioutil.WriteFile(filename, schemeJson, 0644)
}

func init() {
	for k, v := range VALID_CRYPTO_SCHEMES {
		CryptoSchemes = append(CryptoSchemes, k)
		if v {
			ImplementedSchemes = append(ImplementedSchemes, k)
		}
	}
}

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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/big"
)

type ECDSA struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  crypto.PublicKey
}

func NewECDSA() *ECDSA {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil
	}
	return &ECDSA{
		PrivateKey: privateKey,
		PublicKey:  privateKey.Public(),
	}
}

func (ecdsa *ECDSA) Sign(rand io.Reader, priv *crypto.PrivateKey, hash []byte) (r, s *big.Int, err error) {
	return nil, nil, nil

}

func (ecdsa *ECDSA) Verify(pub *crypto.PublicKey, hash []byte, r, s *big.Int) bool {
	return true
}

func (ecdsa *ECDSA) Store(filename string) error {
	keyJson, err := json.Marshal(ecdsa)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, keyJson, 0644)
}

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
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

type ECDSA struct {
	*ecdsa.PrivateKey
	KeyType types.KeyScheme
}

func NewECDSA() *ECDSA {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil
	}
	return &ECDSA{
		PrivateKey: privateKey,
		KeyType:    types.ECDSA256,
	}
}

func (ecdsa *ECDSA) encode() (string, string, error) {
	x509Encoded, err := x509.MarshalECPrivateKey(ecdsa.PrivateKey)
	if err != nil {
		return "", "", err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, err := x509.MarshalPKIXPublicKey(&ecdsa.PrivateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return string(pemEncoded), string(pemEncodedPub), nil
}

func (ecdsaKey *ECDSA) decode(pemEncoded string) error {
	block, _ := pem.Decode([]byte(pemEncoded))
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	ecdsaKey.PrivateKey = privateKey
	ecdsaKey.KeyType = types.ECDSA256
	return nil
}

func (ecdsa *ECDSA) Store(filename string) error {
	privateKey, publicKey, err := ecdsa.encode()
	schemeKey := SchemeKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		KeyType:    ecdsa.KeyType,
	}
	jsonBytes, err := json.Marshal(schemeKey)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, jsonBytes, 0644)
}

func (ecdsa *ECDSA) Sign(rand io.Reader, priv *crypto.PrivateKey, hash []byte) (r, s *big.Int, err error) {
	return nil, nil, nil
}

func (ecdsa *ECDSA) Verify(pub *crypto.PublicKey, hash []byte, r, s *big.Int) bool {
	return true
}

func (ecdsa *ECDSA) GetPublicKey() string {
	_, publicKey, _ := ecdsa.encode()
	return publicKey
}

func (ecdsa *ECDSA) GetKeyId() types.KeyId {
	var bytes = sha256.Sum256([]byte(ecdsa.GetPublicKey()))
	var b = bytes[0:len(bytes)]
	return types.KeyId(fmt.Sprintf("%x", b))
}

func (ecdsa *ECDSA) GetKeyIdHashAlgo() []types.HashAlgo {
	return []types.HashAlgo{"sha256"}
}

func (ecdsa *ECDSA) GetScheme() types.KeyScheme {
	return types.ECDSA256
}

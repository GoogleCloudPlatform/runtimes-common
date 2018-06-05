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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/constants"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

type ECDSA struct {
	*ecdsa.PrivateKey
	KeyType types.KeyScheme
}

type ecdsaSignature struct {
	R, S *big.Int
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

func (ecdsaKey *ECDSA) encode() (string, string, error) {
	x509Encoded, err := x509.MarshalECPrivateKey(ecdsaKey.PrivateKey)
	if err != nil {
		return "", "", err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, err := x509.MarshalPKIXPublicKey(&ecdsaKey.PrivateKey.PublicKey)
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

func (ecdsaKey *ECDSA) Store(filename string) error {
	privateKey, publicKey, err := ecdsaKey.encode()
	schemeKey := SchemeKey{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		KeyType:    ecdsaKey.KeyType,
	}
	jsonBytes, err := json.Marshal(schemeKey)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, jsonBytes, 0644)
}

func (ecdsaKey *ECDSA) Sign(singedMetadata interface{}) (string, error) {
	// Convert singedMetadata to bytes.
	b, err := json.Marshal(singedMetadata)
	if err != nil {
		return "", err
	}

	// Calculate hash of string using SHA256 algo because only first 32 bytes are verified.
	sha256Sum := sha256.Sum256(b)
	r, s, err := ecdsa.Sign(rand.Reader, ecdsaKey.PrivateKey, sha256Sum[0:])
	if err != nil {
		return "", err
	}
	// Use asn1.Marshall as json.Marshal cannot unmarshall big ints.
	out, marshalError := asn1.Marshal(ecdsaSignature{r, s})
	return hex.EncodeToString(out), marshalError
}

func (ecdsaKey *ECDSA) Verify(signingstring string, signature string) bool {
	// Decode the hex String.
	decSignatureString, err := hex.DecodeString(signature)

	if err != nil {
		return false
	}
	es := ecdsaSignature{}
	_, err = asn1.Unmarshal([]byte(decSignatureString), &es)
	if err != nil {
		return false
	}
	// Calculate hash of string using SHA256 algo
	sha256Sum := sha256.Sum256([]byte(signingstring))
	// Verify the signature
	return ecdsa.Verify(&ecdsaKey.PublicKey, sha256Sum[0:], es.R, es.S)
}

func (ecdsaKey *ECDSA) GetPublicKey() string {
	_, publicKey, _ := ecdsaKey.encode()
	return publicKey
}

func (ecdsaKey *ECDSA) GetKeyId() types.KeyId {
	var bytes = sha256.Sum256([]byte(ecdsaKey.GetPublicKey()))
	var b = bytes[0:]
	return types.KeyId(fmt.Sprintf("%x", b))
}

func (ecdsaKey *ECDSA) GetKeyIdHashAlgo() []types.HashAlgo {
	return []types.HashAlgo{constants.SHA256}
}

func (ecdsaKey *ECDSA) GetScheme() types.KeyScheme {
	return types.ECDSA256
}

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
package kms

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/config"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

type KMService interface {
	Encrypt(CryptoKey, string) (*cloudkms.EncryptResponse, error)
	Decrypt(CryptoKey, string) (string, error)
}

type KMS struct {
	Service *cloudkms.Service
}

func New() (*KMS, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	kmsService, err := cloudkms.New(client)
	if err != nil {
		return nil, err
	}
	return &KMS{
		Service: kmsService,
	}, nil
}

type CryptoKey struct {
	ProjectID  string
	LocationID string
	KeyRingID  string
	KeyName    string
}

func (cryptoKey *CryptoKey) Name() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		cryptoKey.ProjectID, cryptoKey.LocationID, cryptoKey.KeyRingID, cryptoKey.KeyName)
}

func CryptoKeyFromConfig(tufConfig config.TUFConfig) CryptoKey {
	return CryptoKey{
		ProjectID:  tufConfig.KMSProjectID,
		LocationID: tufConfig.KMSLocation,
		KeyRingID:  tufConfig.KeyRingID,
		KeyName:    tufConfig.CryptoKeyID,
	}
}

func (kms *KMS) Encrypt(cryptoKey CryptoKey, text string) (*cloudkms.EncryptResponse, error) {
	encryptRequest := &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString([]byte(text)),
	}
	return kms.Service.Projects.Locations.KeyRings.CryptoKeys.Encrypt(cryptoKey.Name(), encryptRequest).Do()
}

func (kms *KMS) Decrypt(cryptoKey CryptoKey, cipherText string) (string, error) {
	decryptRequest := &cloudkms.DecryptRequest{
		Ciphertext: cipherText,
	}
	decryptResp, err := kms.Service.Projects.Locations.KeyRings.CryptoKeys.Decrypt(cryptoKey.Name(), decryptRequest).Do()
	if err != nil {
		return "", err
	}
	bytes, err := base64.StdEncoding.DecodeString(decryptResp.Plaintext)
	return string(bytes), err
}

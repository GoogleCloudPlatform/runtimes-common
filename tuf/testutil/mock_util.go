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
	"io"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/kms"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

// Mock Store
type MockStore struct {
}

func (s *MockStore) Upload(projectId string, bucket string, name string, r io.Reader) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	return nil, nil, nil
}

func (gcs *MockStore) Download(projectId string, bucketId string, objectName string, destPath string) error {
	return nil
}

// Mock KMS
type MockKMS struct {
}

func (kms *MockKMS) Encrypt(cryptoKey kms.CryptoKey, text string) (*cloudkms.EncryptResponse, error) {
	return nil, nil
}

func (kms *MockKMS) Decrypt(cryptoKey kms.CryptoKey, cipherText string) (string, error) {
	return "", nil
}

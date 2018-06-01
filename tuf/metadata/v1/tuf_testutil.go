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

package v1

import (
	"errors"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

func (tmf *TestTargetMetadataFetcher) FetchFileWithAttributes(fileName string) ([]byte, map[string]string, error) {
	if fileName == "" {
		return nil, nil, errors.New("Could not fetch file")
	} else if fileName == "AttrError" {
		return []byte(fileName), nil, errors.New("Could not get attributes")
	}
	return []byte(fileName), map[string]string{
		"name":       fileName,
		"permission": "0644",
	}, nil
}

type TestTargetMetadataFetcher struct {
}

func NewTestTargetMetadataFetcher(filename string) *TestTargetMetadataFetcher {
	return &TestTargetMetadataFetcher{}
}

func (tmf *TestTargetMetadataFetcher) FetchTargetMetadata(filename string, algos []types.HashAlgo) (Target, error) {
	target := Target{}
	fileBytes, attr, err := tmf.FetchFileWithAttributes(filename)
	if err != nil {
		return target, err
	}
	target.Custom = attr
	if filename == "HashError" {
		return target, errors.New("hash error")
	}
	target.Hashes = make([]string, len(algos))
	for i, algo := range algos {
		target.Hashes[i] = string(algo)
	}
	target.Length = len(fileBytes)
	return target, nil
}

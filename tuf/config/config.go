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
package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type TUFConfig struct {
	GCSProjectID      string `yaml: "gcsProjectID"`
	KMSProjectID      string `yaml: "kmsProjectID"`
	KMSLocation       string `yaml: "kmlLocation"`
	KeyRingID         string `yaml: "kmsKeyringID"`
	CryptoKeyID       string `yaml: "cryptoKey"`
	GCSBucketID       string `yaml: "gcsBucketID"`
	RootThreshold     int    `yaml: "rootThreshold"`
	SnapshotThreshold int    `yaml: "snapshotThreshold"`
	TargetThreshold   int    `yaml: "targetThreshold"`
	Targets           []string
}

const (
	RootSecretFileName     = "encrypted-root.key"
	TargetSecretFileName   = "encrypted-target.key"
	SnapshotSecretFileName = "encrypted-snapshot.key"
	TimelineSecretFileName = "encrypted-timeline.key"
)

func ReadConfig(filename string) (TUFConfig, error) {
	buf, err := ioutil.ReadFile(filename)
	tufConfig := TUFConfig{}
	if err != nil {
		return tufConfig, err
	}
	err = yaml.Unmarshal(buf, &tufConfig)
	return tufConfig, err
}

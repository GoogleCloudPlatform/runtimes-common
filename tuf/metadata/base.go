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
package metadata

import (
	"time"
)

type KeyType string

const (
	RSAKey KeyType = "rsa"
)

type HashAlgo string

const (
	SHA256 HashAlgo = "sha256"
	SHA512 HashAlgo = "sha512"
)

type KeyScheme string

const (
	RSA KeyScheme = "rsa"
	ED  KeyScheme = "ed25519"
)

type RoleType string

const (
	RootRole      RoleType = "root"
	TargetRole    RoleType = "target"
	SnapshotRole  RoleType = "snapshot"
	TimestampRole RoleType = "timestamp"
)

type KeyId string

type Signature struct {
	KeyId string
	Sig   string
}

type BaseSigned struct {
	Type        RoleType
	Expires     time.Time
	Version     int
	SpecVersion int
}

type RootSigned struct {
	*BaseSigned
	ConsistentSnapshot bool
	Keys               []KeysDef
	Roles              []Role
}

type KeysDef struct {
	KeyId               KeyId
	KeyidHashAlgorithms []HashAlgo
	Keytype             KeyType
	PrivateKey          string
	Scheme              KeyScheme
}

type Role struct {
	Type      RoleType
	Keyids    []KeyId
	Threshold int
}

type SnapshotSigned struct {
	*BaseSigned
	Meta SnapshotMeta
}

type SnapshotMeta struct {
	Type    RoleType
	Version int
}

type TargetSigned struct {
	*BaseSigned
	Targets []Target
}

type Target struct {
	Filename string
	Custom   interface{}
	Hashes   []string
	Length   int
}

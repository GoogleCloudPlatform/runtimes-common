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
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/metadata"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

type RootMetadata struct {
	Signatures []metadata.Signature `json:"signatures"`
	Signed     RootSigned           `json:"signed"`
}

type RootSigned struct {
	*metadata.BaseSigned
	ConsistentSnapshot bool                    `json:"consistent_snapshot"`
	Keys               map[types.KeyId]KeysDef `json:"keys"`
	Roles              map[types.RoleType]Role `json:"roles"`
}

type KeysDef struct {
	KeyidHashAlgorithms []types.HashAlgo  `json:"keyid_hash_algorithms"`
	Keytype             types.KeyType     `json:"keytype"`
	Val                 map[string]string `json:"keyval"`
	Scheme              types.KeyScheme   `json:"scheme"`
}

type Role struct {
	Keyids    []types.KeyId `json:"keyids"`
	Threshold int           `json:"threshold"`
}

type SnapshotMetadata struct {
	Signatures []metadata.Signature `json:"signatures"`
	Signed     SnapshotSigned       `json:"signed"`
}

type SnapshotSigned struct {
	*metadata.BaseSigned
	Meta map[string]MetaVersionInfo `json:"meta"`
}

type MetaVersionInfo struct {
	Version int `json:"version"`
}

type TargetMetadata struct {
	Signatures []metadata.Signature `json:"signatures"`
	Signed     TargetSigned         `json:"signed"`
}
type TargetSigned struct {
	*metadata.BaseSigned
	Targets map[string]Target `json:"targets"`
}

type Target struct {
	Custom map[string]string `json:"custom"`
	Hashes []string          `json:"hashes"`
	Length int               `json:"length"`
}

type Metadata struct {
	RootMetadata     RootMetadata
	TargetMetadata   TargetMetadata
	SnapshotMetadata SnapshotMetadata
}

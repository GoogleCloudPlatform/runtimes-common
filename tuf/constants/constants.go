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
package constants

import (
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

const (
	RootType                              = "root"
	TargetType                            = "target"
	SnapshotType                          = "snapshot"
	RootSecretFileKey                     = "encrypted-root.key"
	TargetSecretFileKey                   = "encrypted-target.key"
	SnapshotSecretFileKey                 = "encrypted-snapshot.key"
	TimelineSecretFileKey                 = "encrypted-timeline.key"
	RSAKey                types.KeyType   = "rsa"
	ECDSA256Key           types.KeyType   = "ecdsa256"
	SHA256                types.HashAlgo  = "sha256"
	SHA512                types.HashAlgo  = "sha512"
	RootRole              types.RoleType  = "root"
	TargetRole            types.RoleType  = "target"
	SnapshotRole          types.RoleType  = "snapshot"
	TimestampRole         types.RoleType  = "timestamp"
	ECDSA256Scheme        types.KeyScheme = "ECDSA256"
	RSA256Scheme          types.KeyScheme = "RSA256"
)

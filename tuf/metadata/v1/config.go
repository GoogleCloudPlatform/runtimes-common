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
	"time"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/constants"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/metadata"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

type RootMetadata struct {
	Signatures []metadata.Signature `json:"signatures"`
	Signed     metadata.RootSigned  `json:"signed"`
}

type SnapshotMetadata struct {
	Signatures []metadata.Signature    `json:"signatures"`
	Signed     metadata.SnapshotSigned `json:"signed"`
}

type TargetMetadata struct {
	Signatures []metadata.Signature  `json:"signatures"`
	Signed     metadata.TargetSigned `json:"signed"`
}

func PopulateRootMetadata(rootKey types.Scheme, targetKey types.Scheme, snapshotKey types.Scheme) RootMetadata {
	rootMetadata := RootMetadata{}
	rootMetadata.Signed = metadata.RootSigned{
		BaseSigned: &metadata.BaseSigned{
			Type:    constants.RootType,
			Expires: time.Now().AddDate(10, 0, 0), // 10 years later.
		},
		ConsistentSnapshot: false,
		Keys: map[types.KeyId]metadata.KeysDef{
			rootKey.GetKeyId(): metadata.KeysDef{
				KeyidHashAlgorithms: rootKey.GetKeyIdHashAlgo(),
				Scheme:              rootKey.GetScheme(),
				Val: map[string]string{
					"public": rootKey.GetPublicKey(),
				},
			},
		},
	}
	rootMetadata.Signatures = []metadata.Signature{
		metadata.Signature{
			KeyId: rootKey.GetKeyId(),
			//	Sig:   rootKey.Sign(rootMetadata.Signed),
		},
	}
	return rootMetadata
}

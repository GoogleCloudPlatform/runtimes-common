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
	"fmt"
	"time"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/config"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/constants"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/metadata"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

var HASHALGOS = []types.HashAlgo{
	types.HashAlgo("sha256"),
}

type TUF struct {
	TufConfig       config.TUFConfig
	RootSecrets     []types.Scheme
	OldRootSecrets  []types.Scheme
	TargetSecrets   []types.Scheme
	SnapshotSecrets []types.Scheme
}

func (tuf *TUF) PopulateRootMetadata() (RootMetadata, error) {
	rootMetadata := RootMetadata{}
	rootMetadata.Signed = RootSigned{
		BaseSigned:         tuf.populateBaseSigned(constants.RootType),
		ConsistentSnapshot: false, // TODO: Verify logic for consistent snapshot.
	}
	rootMetadata.Signed.Roles = tuf.populateRoles()
	rootMetadata.Signed.Keys = populateKeys(tuf.RootSecrets)
	var err error
	rootMetadata.Signatures, err = sign(append(tuf.RootSecrets, tuf.OldRootSecrets...), rootMetadata.Signed)
	if err != nil {
		return rootMetadata, err
	}
	return rootMetadata, nil
}

func (tuf *TUF) PopulateTargetMetadata(targetFetcher TargetMetadataFetcher) (TargetMetadata, error) {
	targetMetadata := TargetMetadata{}
	targetMetadata.Signed = TargetSigned{
		BaseSigned: tuf.populateBaseSigned(constants.TargetType),
	}
	var err error
	targetMetadata.Signed.Targets, err = populateTargets(targetFetcher, tuf.TufConfig.Targets)
	if err != nil {
		return targetMetadata, err
	}
	targetMetadata.Signatures, err = tuf.SignTargetMetadata(targetMetadata.Signed)
	if err != nil {
		return targetMetadata, err
	}
	return targetMetadata, nil
}

func (tuf *TUF) PopulateSnapshotMetadata() (SnapshotMetadata, error) {
	snaphotMetadata := SnapshotMetadata{}
	snaphotMetadata.Signed = SnapshotSigned{
		BaseSigned: tuf.populateBaseSigned(constants.SnapshotRole),
	}
	var err error
	snaphotMetadata.Signed.Meta = populateMeta(tuf.getVersion())
	snaphotMetadata.Signatures, err = tuf.SignSnapshotMetadata(snaphotMetadata.Signed)
	if err != nil {
		return snaphotMetadata, err
	}
	return snaphotMetadata, nil
}

func (tuf *TUF) populateRoles() map[types.RoleType]Role {
	var roles = map[types.RoleType]Role{}
	roles[constants.RootType] = addKeyIds(tuf.RootSecrets, tuf.TufConfig.RootThreshold)
	roles[constants.TargetType] = addKeyIds(tuf.TargetSecrets, tuf.TufConfig.TargetThreshold)
	roles[constants.SnapshotType] = addKeyIds(tuf.SnapshotSecrets, tuf.TufConfig.SnapshotThreshold)
	return roles
}

func addKeyIds(secrets []types.Scheme, threshold int) Role {
	role := Role{}
	keyids := make([]types.KeyId, len(secrets))
	for i, secret := range secrets {
		keyids[i] = secret.GetKeyId()
	}
	role.Keyids = keyids
	role.Threshold = threshold
	return role
}

func populateKeys(secrets []types.Scheme) map[types.KeyId]KeysDef {
	var keys = map[types.KeyId]KeysDef{}
	for _, secret := range secrets {
		keys[secret.GetKeyId()] = KeysDef{
			KeyidHashAlgorithms: secret.GetKeyIdHashAlgo(),
			Scheme:              secret.GetScheme(),
			Keytype:             types.KeyType(secret.GetScheme()),
			Val: map[string]string{
				"public": secret.GetPublicKey(),
			},
		}
	}
	return keys
}

func (tuf *TUF) SignTargetMetadata(metadata interface{}) ([]metadata.Signature, error) {
	return sign(tuf.TargetSecrets, metadata)
}

func (tuf *TUF) SignRootMetadata(metadata interface{}) ([]metadata.Signature, error) {
	return sign(tuf.RootSecrets, metadata)
}

func (tuf *TUF) SignSnapshotMetadata(metadata interface{}) ([]metadata.Signature, error) {
	return sign(tuf.SnapshotSecrets, metadata)
}

func sign(secrets []types.Scheme, signed interface{}) ([]metadata.Signature, error) {
	signatures := make([]metadata.Signature, len(secrets))
	for i, secret := range secrets {
		sig, err := secret.Sign(signed)
		if err != nil {
			return nil, err
		}
		signatures[i] = metadata.Signature{
			KeyId: secret.GetKeyId(),
			Sig:   sig,
		}
	}
	return signatures, nil
}

func (tuf *TUF) populateBaseSigned(metadataType types.RoleType) *metadata.BaseSigned {
	return &metadata.BaseSigned{
		Type: metadataType,
		// TODO (tejaldesai): Read this from TUFConfig.TargetSecretExpiryDuration
		Expires:     time.Now().AddDate(10, 0, 0), // 10 years later.
		SpecVersion: tuf.TufConfig.SpecVersion,
		Version:     tuf.getVersion(),
	}
}

func populateTargets(targetFetcher TargetMetadataFetcher, targetsConfigured []string) (map[string]Target, error) {
	targets := map[string]Target{}
	var err error
	for _, configTarget := range targetsConfigured {
		targets[configTarget], err = targetFetcher.FetchTargetMetadata(configTarget, HASHALGOS)
		if err != nil {
			return nil, err
		}
	}
	return targets, nil
}

func populateMeta(version int) map[string]MetaVersionInfo {
	// TODO (tejaldesai): Look at previous meta and current version to the
	meta := map[string]MetaVersionInfo{} // Fetch previous Metadata
	// Add entries for root and target metadata for current version
	currentVersionJson := fmt.Sprintf("%d.json", version)
	consistentVersion := MetaVersionInfo{
		Version: version,
	}
	meta["root"+currentVersionJson] = consistentVersion
	meta["target"+currentVersionJson] = consistentVersion
	return meta
}

func (tuf *TUF) getVersion() int {
	// TODO (tejaldesai): Calculate Version from previously published metadata files.
	return 1
}

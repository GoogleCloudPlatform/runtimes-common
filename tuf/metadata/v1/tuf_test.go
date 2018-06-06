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
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/tuf/metadata"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/testutil"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/types"
)

func TestRootMetadata(t *testing.T) {
	tuf := TUF{testutil.TestTUFConfig,
		[]types.Scheme{},
		[]types.Scheme{},
		[]types.Scheme{}, []types.Scheme{}}
	rootMetadata, err := tuf.PopulateRootMetadata()
	if err != nil {
		t.Fatalf("Did not expect error %v", err)
	}
	expectedBaseSigned := &metadata.BaseSigned{
		Type:        "root",
		SpecVersion: 1,
		Version:     1,
	}
	// TODO (tejaldesai): Fix this by reading time form config.
	expectedBaseSigned.Expires = rootMetadata.Signed.BaseSigned.Expires
	if !reflect.DeepEqual(rootMetadata.Signed.BaseSigned, expectedBaseSigned) {
		t.Fatalf("Expected %v\n Got %v", rootMetadata.Signed.BaseSigned, expectedBaseSigned)
	}
}

var tuf = TUF{
	TufConfig: testutil.TestTUFConfig,
	RootSecrets: []types.Scheme{
		&testutil.TestKey{
			PrivateKey: "secret1",
			KeyType:    "test-scheme1",
			PublicKey:  "public1",
			KeyId:      "idroot1",
			SignStr:    "secret1",
		},
		&testutil.TestKey{
			PrivateKey: "secret2",
			KeyType:    "test-scheme2",
			PublicKey:  "public2",
			KeyId:      "idroot2",
			SignStr:    "secret2",
		},
	},
	TargetSecrets: []types.Scheme{
		&testutil.TestKey{
			PrivateKey: "target",
			KeyType:    "test-scheme1",
			PublicKey:  "targetpublic",
			KeyId:      "idtarget",
			SignStr:    "targetsecret",
		},
	},
	SnapshotSecrets: []types.Scheme{
		&testutil.TestKey{
			PrivateKey: "snapshot",
			KeyType:    "test-scheme1",
			PublicKey:  "snapshotpublic",
			KeyId:      "idsnapshot",
			SignStr:    "snapshotsecret",
		},
	},
}

func TestPopulateRoles(t *testing.T) {
	expectedResult := map[types.RoleType]Role{
		"root": {
			Keyids: []types.KeyId{
				types.KeyId("idroot1"),
				types.KeyId("idroot2"),
			},
			Threshold: 1,
		},
		"target": {
			Keyids: []types.KeyId{
				types.KeyId("idtarget"),
			},
			Threshold: 1,
		},
		"snapshot": {
			Keyids: []types.KeyId{
				types.KeyId("idsnapshot"),
			},
			Threshold: 1,
		},
	}
	actual := tuf.populateRoles()
	if !reflect.DeepEqual(expectedResult, actual) {
		t.Fatalf("Expected %v\n Got : %v", expectedResult, actual)
	}
}

func TestPopulateKeys(t *testing.T) {
	expectedResult := map[types.KeyId]KeysDef{
		types.KeyId("idroot1"): {
			KeyidHashAlgorithms: []types.HashAlgo{types.HashAlgo("sha256")},
			Val: map[string]string{
				"public": "public1",
			},
			Keytype: types.KeyType("test-scheme1"),
			Scheme:  types.KeyScheme("test-scheme1"),
		},
		types.KeyId("idroot2"): {
			KeyidHashAlgorithms: []types.HashAlgo{types.HashAlgo("sha256")},
			Val: map[string]string{
				"public": "public2",
			},
			Keytype: types.KeyType("test-scheme2"),
			Scheme:  types.KeyScheme("test-scheme2"),
		},
	}
	actual := populateKeys(tuf.RootSecrets)
	if !reflect.DeepEqual(expectedResult, actual) {
		t.Fatalf("Expected:\n%v\n Got :\n%v", expectedResult, actual)
	}
}

var signTC = []struct {
	name         string
	signMetadata interface{}
	secrets      string
	expectedErr  error
	expectedRes  []metadata.Signature
}{
	{"signRootSuccess", "str", "root", nil,
		[]metadata.Signature{
			{
				KeyId: types.KeyId("idroot1"),
				Sig:   "secret1str",
			},
			{
				KeyId: types.KeyId("idroot2"),
				Sig:   "secret2str",
			},
		},
	},
	{"signTargetSuccess", "str", "target", nil,
		[]metadata.Signature{
			{
				KeyId: types.KeyId("idtarget"),
				Sig:   "targetsecretstr",
			},
		},
	},
	{"signSnapshotSuccess", "str", "snapshot", nil,
		[]metadata.Signature{
			{
				KeyId: types.KeyId("idsnapshot"),
				Sig:   "snapshotsecretstr",
			},
		},
	},
	{"singErr", 1, "snapshot", errors.New("Signing Error"), nil},
}

func TestSign(t *testing.T) {
	for _, tc := range signTC {
		t.Run(tc.name, func(t *testing.T) {
			var actual []metadata.Signature
			var err error
			switch tc.secrets {
			case "root":
				actual, err = tuf.SignRootMetadata(tc.signMetadata)
			case "target":
				actual, err = tuf.SignTargetMetadata(tc.signMetadata)
			default:
				actual, err = tuf.SignSnapshotMetadata(tc.signMetadata)

			}
			if !reflect.DeepEqual(tc.expectedErr, err) {
				t.Fatalf("Expected Error:\n%v\n Got :\n%v", tc.expectedErr, err)
			}
			if !reflect.DeepEqual(tc.expectedRes, actual) {
				t.Fatalf("Expected:\n%v\n Got :\n%v", tc.expectedRes, actual)
			}
		})
	}
}

var populateTargetsTC = []struct {
	name        string
	targets     []string
	expectedRes map[string]Target
	err         error
}{
	{"success", []string{"file1.txt", "file22.txt"},
		map[string]Target{
			"file1.txt": {
				Custom: map[string]string{
					"name":       "file1.txt",
					"permission": "0644",
				},
				Length: 9,
				Hashes: []string{"sha256"},
			},
			"file22.txt": {
				Custom: map[string]string{
					"name":       "file22.txt",
					"permission": "0644",
				},
				Length: 10,
				Hashes: []string{"sha256"},
			},
		}, nil},
	{"fetcherror", []string{""}, nil, errors.New("Could not fetch file")},
	{"attrerror", []string{"AttrError"}, nil, errors.New("Could not get attributes")},
	{"hasherror", []string{"HashError"}, nil, errors.New("hash error")},
}

func TestPopulateTargets(t *testing.T) {
	for _, tc := range populateTargetsTC {
		t.Run(tc.name, func(t *testing.T) {
			tf := TestTargetMetadataFetcher{}
			actual, err := populateTargets(&tf, tc.targets)
			if !reflect.DeepEqual(tc.err, err) {
				t.Fatalf("Expected Error:\n%v\n Got :\n%v", tc.err, err)
			}
			if tc.err == nil {
				if !reflect.DeepEqual(tc.expectedRes, actual) {
					t.Fatalf("Expected:\n%v\n Got :\n%v", tc.expectedRes, actual)
				}
			}
		})
	}
}

var metaTC = []struct {
	name        string
	version     int
	expectedRes map[string]MetaVersionInfo
}{
	{"one", 1, map[string]MetaVersionInfo{
		"root1.json":   {Version: 1},
		"target1.json": {Version: 1},
	}},
	{"two", 2, map[string]MetaVersionInfo{
		"root2.json":   {Version: 2},
		"target2.json": {Version: 2},
	}},
}

func TestMeta(t *testing.T) {
	for _, tc := range metaTC {
		t.Run(tc.name, func(t *testing.T) {
			actual := populateMeta(tc.version)
			if !reflect.DeepEqual(tc.expectedRes, actual) {
				t.Fatalf("Expected:\n%v\n Got :\n%v", tc.expectedRes, actual)
			}
		})
	}
}

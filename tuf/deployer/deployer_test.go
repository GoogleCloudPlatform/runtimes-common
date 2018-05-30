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
package deployer

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"reflect"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/metadata/v1"
	"github.com/GoogleCloudPlatform/runtimes-common/tuf/testutil"
)

type TestStore struct{}

func (ts *TestStore) Upload(string, string, io.Reader) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	return nil, nil, nil
}

func (ts *TestStore) Download(string, string) ([]byte, error) {
	return nil, nil
}

func NewTestDeployer() (DeployTool, error) {
	return &Deployer{
		KmsService: nil,
		Storage:    &TestStore{},
	}, nil
}

func TestGenerateMetadata(t *testing.T) {
	//tmpdir, err := ioutil.TempDir("", "generate_")
	//if err != nil {
	//	panic(fmt.Sprintf("Cannot run tests due to %v", err))
	//}
	// fmt.Println(tmpdir)
	// defer os.Remove(tmpdir)
	d, err := NewTestDeployer()
	if err != nil {
		t.Fatalf("fail due to %s", err)
	}

	root, _, _, err := d.GenerateMetadata(testutil.TestTUFConfig, "testdata/testRootSecret.key", "testdata/testTargetSecret.key",
		"testdata/testSnapshotSecret.key", []byte{})

	if err != nil {
		t.Fatalf("Expected nil, got %v", err)
	}

	expectedRootMetadata := v1.RootMetadata{}
	rootBytes, _ := ioutil.ReadFile("testdata/testRoot.json")
	json.Unmarshal(rootBytes, &expectedRootMetadata)
	root.Signed.Expires = expectedRootMetadata.Signed.Expires
	if !reflect.DeepEqual(root, &expectedRootMetadata) {
		t.Fatalf("Expected RootMetadata %v \n Got %v", &expectedRootMetadata, root)
	}
}

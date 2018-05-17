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
package gcs

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

type GCSStore struct {
	Client  *storage.Client
	Context context.Context
}

func New() (*GCSStore, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCSStore{
		Client:  client,
		Context: ctx,
	}, nil
}

func (gcs *GCSStore) Upload(bucket string, name string, r io.Reader) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {

	bh := gcs.Client.Bucket(bucket)
	// Next check if the bucket exists
	if _, err := bh.Attrs(gcs.Context); err != nil {
		// TODO: Create a new bucket with read permissions for "project-team-<projectId>"
		return nil, nil, errors.New(fmt.Sprintf("Error checking for bucket. %v", err))
	}

	obj := bh.Object(name)
	w := obj.NewWriter(gcs.Context)
	if _, err := io.Copy(w, r); err != nil {
		return nil, nil, err
	}
	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	attrs, err := obj.Attrs(gcs.Context)
	return obj, attrs, err
}

func (gcs *GCSStore) Download(bucketId string, objectName string) ([]byte, error) {
	bh := gcs.Client.Bucket(bucketId)

	rc, err := bh.Object(objectName).NewReader(gcs.Context)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return slurp, nil
}

func (gcs *GCSStore) Delete(bucketId string, objectName string) error {
	o := gcs.Client.Bucket(bucketId).Object(objectName)
	if err := o.Delete(gcs.Context); err != nil {
		return err
	}
	return nil
}

/*
Copyright 2018 Google LLC
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

package cmd

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"

	"github.com/containers/image/docker"

	"testing"
)

var files = map[string]string{
	"foo":     "baz",
	"bar/baz": "bat",
}

func setupTar(t *testing.T) string {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error setting up tempfile: %s", err)
	}
	defer f.Close()

	tarballPath = f.Name()
	t.Log(tarballPath)

	tw := tar.NewWriter(f)
	defer tw.Close()
	for p, c := range files {
		hdr := &tar.Header{
			Name: p,
			Mode: 0600,
			Size: int64(len(c)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("Error writing test tar: %s", err)
		}
		if _, err := tw.Write([]byte(c)); err != nil {
			t.Fatalf("Error writing test tar: %s", err)
		}
	}
	return f.Name()
}

func checkTar(t *testing.T) {
	ref, err := docker.ParseReference("//" + image)
	if err != nil {
		t.Fatalf("Error pulling built image: %s", err)
	}

	img, err := ref.NewImage(nil)
	if err != nil {
		t.Fatalf("Error pulling built image: %s", err)
	}

	imgSrc, err := ref.NewImageSource(nil)
	if err != nil {
		t.Fatalf("Error pulling built image: %s", err)
	}

	layers := img.LayerInfos()
	appendedLayerID := layers[len(layers)-1]
	blob, _, err := imgSrc.GetBlob(appendedLayerID)
	if err != nil {
		t.Fatalf("Error fetching appended layer: %s", err)
	}

	gzr, err := gzip.NewReader(blob)
	if err != nil {
		t.Fatalf("Error decompressing appended layer: %s", err)
	}
	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error reading tar: %s", err)
		}
		c, ok := files[hdr.Name]
		if !ok {
			t.Fatalf("Unexpected file in layer: %s", hdr.Name)
		}
		b, err := ioutil.ReadAll(tr)
		if err != nil {
			t.Fatalf("Error reading file from layer: %s", err)
		}
		if string(b) != c {
			t.Fatalf("Files do not match: %s %s", b, c)
		}
	}
}

func TestAppend(t *testing.T) {
	baseImage = "gcr.io/google-appengine/debian9:latest"
	image = "gcr.io/gcp-runtimes/test/appender:latest"
	tarballPath = setupTar(t)
	defer os.RemoveAll(tarballPath)

	rootCmd.Run(nil, nil)

	checkTar(t)

}

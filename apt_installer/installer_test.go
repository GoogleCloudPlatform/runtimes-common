// Copyright 2017 Google Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var BUF_SIZE = 1024

var test_dir string
var err error

func TestParseNormal(t *testing.T) {
	dockerfile, _ := ioutil.TempFile(test_dir, "TestParseNormal")
	config := createInstaller(dockerfile.Name(), "test_configs/normal.yaml")
	if len(config.AptPackages.Packages) == 0 {
		t.Errorf("Packages not parsed correctly!")
	} else if len(config.AptPackages.PPAs) == 0 {
		t.Errorf("PPAs not parsed correctly!")
	}
}

func TestParseNoPackages(t *testing.T) {
	dockerfile, _ := ioutil.TempFile(test_dir, "TestParseNoPackages")
	config := createInstaller(dockerfile.Name(), "test_configs/no_packages.yaml")
	if len(config.AptPackages.Packages) != 0 {
		t.Errorf("Extra packages included on parse!")
	} else if len(config.AptPackages.PPAs) == 0 {
		t.Errorf("PPAs not parsed correctly!")
	}
}

func TestParseNoPpas(t *testing.T) {
	dockerfile, _ := ioutil.TempFile(test_dir, "TestParseNoPpas")
	config := createInstaller(dockerfile.Name(), "test_configs/no_ppas.yaml")
	if len(config.AptPackages.Packages) == 0 {
		t.Errorf("Packages not parsed correctly!")
	} else if len(config.AptPackages.PPAs) != 0 {
		t.Errorf("Extra PPAs included on parse!")
	}
}

func TestParseEmpty(t *testing.T) {
	dockerfile, _ := ioutil.TempFile(test_dir, "TestParseEmpty")
	config := createInstaller(dockerfile.Name(), "test_configs/empty.yaml")
	if len(config.AptPackages.Packages) != 0 {
		t.Errorf("Extra packages included on parse!")
	} else if len(config.AptPackages.PPAs) != 0 {
		t.Errorf("Extra PPAs included on parse!")
	}
}

func TestDockerfileNormal(t *testing.T) {
	dockerfile, err := ioutil.TempFile(test_dir, "TestDockerfileNormal")
	installer := RuntimeConfig{AptPackageHolder{[]string{"ppa1", "ppa2"}, []string{"package1", "package2"}}, dockerfile.Name()}
	err = generateDockerfile(installer)
	if err != nil {
		t.Errorf(err.Error())
	}
	_checkFileContains(t, dockerfile, "package1", "package2", "add-apt-repository", "ppa1", "ppa2")
}

func TestDockerfileNoPackages(t *testing.T) {
	dockerfile, err := ioutil.TempFile(test_dir, "TestDockerfileNoPackages")
	installer := RuntimeConfig{AptPackageHolder{[]string{"ppa1", "ppa2"}, []string{}}, dockerfile.Name()}
	err = generateDockerfile(installer)
	if err != nil {
		t.Errorf(err.Error())
	}
	_checkFileContains(t, dockerfile, "add-apt-repository", "ppa1", "ppa2")
}

func TestDockerfileNoPpas(t *testing.T) {
	dockerfile, err := ioutil.TempFile(test_dir, "TestDockerfileNoPpas")
	installer := RuntimeConfig{AptPackageHolder{[]string{}, []string{"package1", "package2"}}, dockerfile.Name()}
	err = generateDockerfile(installer)
	if err != nil {
		t.Errorf(err.Error())
	}
	_checkFileContains(t, dockerfile, "package1", "package2")
	_checkFileOmits(t, dockerfile, "add-apt-repository")
}

func TestDockerfileEmpty(t *testing.T) {
	dockerfile, err := ioutil.TempFile(test_dir, "TestDockerfileEmpty")
	installer := RuntimeConfig{AptPackageHolder{[]string{}, []string{}}, dockerfile.Name()}
	err = generateDockerfile(installer)
	if err != nil {
		t.Errorf(err.Error())
	}

	buf := make([]byte, BUF_SIZE)
	n, err := dockerfile.Read(buf)
	if err != nil {
		t.Errorf("Error when reading Dockerfile: %s", err)
	}
	var d_str string
	if n > 0 {
		d_str = string(buf[:n])
	}

	if d_str != "\n" {
		t.Errorf("Dockerfile contains installation artifacts with no ppas/packages specified!")
	}
}

func _checkFileContains(t *testing.T, dockerfile *os.File, contents ...string) {
	buf := make([]byte, BUF_SIZE)
	n, err := dockerfile.Read(buf)
	if err != nil {
		t.Errorf("Error when reading Dockerfile: %s", err)
	}
	var d_str string
	if n > 0 {
		d_str = string(buf[:n])
	}
	for _, s := range contents {
		if !strings.Contains(d_str, s) {
			t.Errorf("String %s not found in Dockerfile", s)
		}
	}
}

func _checkFileOmits(t *testing.T, dockerfile *os.File, omit_strings ...string) {
	buf := make([]byte, BUF_SIZE)
	t.Logf("Dockerfile location: %s", dockerfile.Name())
	n, err := dockerfile.Read(buf)
	if err != nil && err != io.EOF {
		t.Errorf("Error when reading Dockerfile: %s", err)
	}
	var d_str string
	if n > 0 {
		d_str = string(buf[:n])
	}
	for _, s := range omit_strings {
		if strings.Contains(d_str, s) {
			t.Errorf("String %s found in Dockerfile", s)
		}
	}
}

func TestMain(m *testing.M) {
	test_dir, err = ioutil.TempDir(".", "installer_test.")
	if err != nil {
		fmt.Printf("Error when creating temp_dir: %s", err)
		os.Exit(1)
	}

	code := m.Run()
	os.RemoveAll(test_dir)
	os.Exit(code)
}

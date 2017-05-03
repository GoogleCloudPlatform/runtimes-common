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
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

var testBaseImage = "gcr.io/gcp-runtimes/check_if_tag_exists"

func TestParseNormal(t *testing.T) {
	config := createInstaller(testBaseImage, "test_configs/normal.yaml")
	if len(config.AptPackages.Packages) == 0 {
		t.Errorf("Packages not parsed correctly!")
	} else if len(config.AptPackages.PPAs) == 0 {
		t.Errorf("PPAs not parsed correctly!")
	}
}

func TestParseNoPackages(t *testing.T) {
	config := createInstaller(testBaseImage, "test_configs/no_packages.yaml")
	if len(config.AptPackages.Packages) != 0 {
		t.Errorf("Extra packages included on parse!")
	} else if len(config.AptPackages.PPAs) == 0 {
		t.Errorf("PPAs not parsed correctly!")
	}
}

func TestParseNoPpas(t *testing.T) {
	config := createInstaller(testBaseImage, "test_configs/no_ppas.yaml")
	if len(config.AptPackages.Packages) == 0 {
		t.Errorf("Packages not parsed correctly!")
	} else if len(config.AptPackages.PPAs) != 0 {
		t.Errorf("Extra PPAs included on parse!")
	}
}

func TestParseEmpty(t *testing.T) {
	config := createInstaller(testBaseImage, "test_configs/empty.yaml")
	if len(config.AptPackages.Packages) != 0 {
		t.Errorf("Extra packages included on parse!")
	} else if len(config.AptPackages.PPAs) != 0 {
		t.Errorf("Extra PPAs included on parse!")
	}
}

func TestDockerfileNormal(t *testing.T) {
	installer := Installer{AptPackageHolder{[]string{"ppa1", "ppa2"}, []string{"package1", "package2"}}, "foobar"}
	err, build_dir := generateDockerfile(installer)
	defer os.RemoveAll(build_dir)
	if err != nil {
		t.Errorf(err.Error())
	}
	_checkFileContains(t, build_dir, "package1", "package2", "add-apt-repository", "ppa1", "ppa2")
}

func TestDockerfileNoPackages(t *testing.T) {
	installer := Installer{AptPackageHolder{[]string{"ppa1", "ppa2"}, []string{}}, "foobar"}
	err, build_dir := generateDockerfile(installer)
	defer os.RemoveAll(build_dir)
	if err != nil {
		t.Errorf(err.Error())
	}
	_checkFileContains(t, build_dir, "add-apt-repository", "ppa1", "ppa2")
}

func TestDockerfileNoPpas(t *testing.T) {
	installer := Installer{AptPackageHolder{[]string{}, []string{"package1", "package2"}}, "foobar"}
	err, build_dir := generateDockerfile(installer)
	defer os.RemoveAll(build_dir)
	if err != nil {
		t.Errorf(err.Error())
	}
	_checkFileContains(t, build_dir, "package1", "package2")
	_checkFileOmits(t, build_dir, "add-apt-repository")
}

func TestDockerfileEmpty(t *testing.T) {
	installer := Installer{AptPackageHolder{[]string{}, []string{}}, "foobar"}
	err, build_dir := generateDockerfile(installer)
	defer os.RemoveAll(build_dir)
	if err != nil {
		t.Errorf(err.Error())
	}
	dockerfile, err := ioutil.ReadFile(path.Join(build_dir, "/Dockerfile"))
	if err != nil {
		t.Errorf("Error when reading Dockerfile: %s", err)
	}
	d_str := string(dockerfile)
	if d_str != "FROM foobar\n\n" {
		t.Errorf("Dockerfile contains installation artifacts with no ppas/packages specified!")
	}
}

func _checkFileContains(t *testing.T, build_dir string, contents ...string) {
	dockerfile, err := ioutil.ReadFile(path.Join(build_dir, "/Dockerfile"))
	if err != nil {
		t.Errorf("Error when reading Dockerfile: %s", err)
	}
	d_str := string(dockerfile)
	for _, s := range contents {
		if !strings.Contains(d_str, s) {
			t.Errorf("String %s not found in Dockerfile", s)
		}
	}
}

func _checkFileOmits(t *testing.T, build_dir string, omit_strings ...string) {
	dockerfile, err := ioutil.ReadFile(path.Join(build_dir, "/Dockerfile"))
	if err != nil {
		t.Errorf("Error when reading Dockerfile: %s", err)
	}
	d_str := string(dockerfile)
	for _, s := range omit_strings {
		if strings.Contains(d_str, s) {
			t.Errorf("String %s found in Dockerfile", s)
		}
	}
}

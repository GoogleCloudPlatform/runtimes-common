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
	"bufio"
	"flag"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var INTERMEDIATE_IMAGE_TAG = "apt:intermediate"

func generateDockerfile(installer Installer) (int, string) {
	var err error

	build_dir, err := ioutil.TempDir(".", "intermediate_base.")
	if err != nil {
		log.Printf("Error when creating temp_dir: %s", err)
		return 1, ""
	}

	f, err := os.Create(build_dir + "/Dockerfile")
	if err != nil {
		log.Printf("Error opening file for writing: %s", err)
		return 1, build_dir
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	err = INSTALL_TMPL.Execute(w, installer)
	if err != nil {
		log.Printf("Error when executing template: %s", err)
		return 1, build_dir
	}
	for _, ppa := range installer.AptPackages.PPAs {
		err = PPA_TMPL.Execute(w, PpaHolder{ppa})
		if err != nil {
			log.Printf("Error when executing template: %s", err)
			return 1, build_dir
		}
	}
	err = APT_TMPL.Execute(w, installer.AptPackages)
	if err != nil {
		log.Printf("Error when executing template: %s", err)
		return 1, build_dir
	}
	err = REMOVE_TOOLS_TMPL.Execute(w, installer.AptPackages)
	if err != nil {
		log.Printf("Error when executing template: %s", err)
		return 1, build_dir
	}

	w.Flush()
	return 0, build_dir
}

func doBuild(build_dir string) int {
	docker_flags := []string{"build", "--no-cache", "-t", INTERMEDIATE_IMAGE_TAG, build_dir}
	cmd := exec.Command("docker", docker_flags...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return 1
	}
	return 0
}

var baseImage, configFile string

func main() {
	init_templates()
	flag.StringVar(&configFile, "yaml", "/workspace/app.yaml",
		"path to the .yaml file containing packages to install.")
	flag.StringVar(&baseImage, "base", "",
		"base runtime image to install packages on.")
	flag.Parse()

	if baseImage == "" {
		log.Fatalf("Please provide base image.")
	}

	configContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error when reading config file: %s", err)
	}

	var config Installer

	if err := yaml.Unmarshal(configContents, &config); err != nil {
		log.Fatalf("Error when unmarshaling yaml file: %s", err)
	}
	config.BaseImage = baseImage

	gen_code, build_dir := generateDockerfile(config)
	if gen_code != 0 {
		if build_dir != "" {
			os.RemoveAll(build_dir)
		}
		os.Exit(gen_code)
	}
	build_code := doBuild(build_dir)
	os.RemoveAll(build_dir)
	os.Exit(build_code)
}

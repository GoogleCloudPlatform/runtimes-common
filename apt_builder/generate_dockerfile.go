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
	"errors"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var INTERMEDIATE_IMAGE_TAG = "apt:intermediate"

func generateDockerfile(config RuntimeConfig) error {
	var err error

	f, err := os.OpenFile(config.Dockerfile, os.O_RDONLY|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return errors.New(fmt.Sprintf("Error when reading Dockerfile %s: %s",
									  config.Dockerfile, err))
	}

	defer f.Close()
	w := bufio.NewWriter(f)

	if err := DOCKERFILE_TMPL.Execute(w, config); err != nil {
		return errors.New(fmt.Sprintf("Error when executing template: %s", err))
	}

	w.Flush()
	return nil
}

func doBuild(build_dir string) error {
	docker_flags := []string{"build", "--no-cache", "-t", INTERMEDIATE_IMAGE_TAG, build_dir}
	cmd := exec.Command("docker", docker_flags...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

var dockerfile, configFile string

func main() {
	init_templates()
	flag.StringVar(&configFile, "yaml", "/workspace/app.yaml",
		"path to the .yaml file containing packages to install.")
	flag.StringVar(&dockerfile, "dockerfile", "",
		"path to the Dockerfile for the application.")
	flag.Parse()

	if dockerfile == "" {
		log.Fatalf("Please provide path to Dockerfile")
	}

	configContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error when reading config file: %s", err)
	}

	var config RuntimeConfig

	if err := yaml.Unmarshal(configContents, &config); err != nil {
		log.Fatalf("Error when unmarshaling yaml file: %s", err)
	}
	config.Dockerfile = dockerfile

	err = generateDockerfile(config)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

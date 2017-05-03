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
	"path"
	"text/template"
)

var INTERMEDIATE_IMAGE_TAG = "apt:intermediate"
var DOCKERFILE_TMPL *template.Template

func generateDockerfile(installer Installer) (error, string) {
	var err error

	build_dir, err := ioutil.TempDir(".", "intermediate_base.")
	if err != nil {
		return errors.New(fmt.Sprintf("Error when creating temp_dir: %s", err)), ""
	}

	f, err := os.Create(path.Join(build_dir, "/Dockerfile"))
	if err != nil {
		return errors.New(fmt.Sprintf("Error opening file for writing: %s", err)), build_dir
	}

	defer f.Close()
	w := bufio.NewWriter(f)

	if err := DOCKERFILE_TMPL.Execute(w, installer); err != nil {
		return errors.New(fmt.Sprintf("Error when executing template: %s", err)), build_dir
	}

	w.Flush()
	return nil, build_dir
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

var baseImage, configFile string

func main() {
	var err error

	DOCKERFILE := `FROM {{.BaseImage}}
{{if or .AptPackages.Packages .AptPackages.PPAs}}
RUN apt-get update && apt-get install -y --force-yes \
    apt-utils software-properties-common python-software-properties \
    {{range $ppa := .AptPackages.PPAs}}
    {{"&& add-apt-repository -y "}}{{$ppa}} \{{end}}
    {{if .AptPackages.Packages}}
    && apt-get update && apt-get install -y --force-yes \
    {{range $pkg := .AptPackages.Packages}}   {{$pkg }} \
    {{end}}{{end}}
    && apt-get remove -y --force-yes software-properties-common \
       python-software-properties apt-utils \
    && apt-get autoremove -y --force-yes \
    && apt-get clean -y --force-yes
{{end}}
`

	DOCKERFILE_TMPL = template.Must(template.New("DOCKERFILE").Parse(DOCKERFILE))

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

	err, build_dir := generateDockerfile(config)
	if err != nil {
		log.Printf(err.Error())
		if build_dir != "" {
			os.RemoveAll(build_dir)
		}
		os.Exit(1)
	}
	err = doBuild(build_dir)
	os.RemoveAll(build_dir)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

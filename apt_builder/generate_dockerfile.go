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
	"io/ioutil"
	"log"
	"os"
	"github.com/ghodss/yaml"
)

func generateDockerfile(installer Installer) {
	var err error

    f, err := os.Create("./test/Dockerfile")
    if err != nil {
        log.Fatalf("Error opening file for writing: %s", err)
    }

    defer f.Close()

    w := bufio.NewWriter(f)

	err = INSTALL_TMPL.Execute(w, installer)
	if err != nil {
		log.Fatalf("Error when executing template: %s", err)
	}
    for _, ppa := range installer.AptPackages.PPAs {
        err = PPA_TMPL.Execute(w, PpaHolder{ppa})
        if err != nil {
            log.Fatalf("Error when executing template: %s", err)
        }
    }
	// err = PPA_TMPL.Execute(w, installer)
	// if err != nil {
	// 	log.Fatalf("Error when executing template: %s", err)
	// }
	err = APT_TMPL.Execute(w, installer.AptPackages)
	if err != nil {
		log.Fatalf("Error when executing template: %s", err)
	}
	err = REMOVE_TOOLS_TMPL.Execute(w, installer.AptPackages)
	if err != nil {
		log.Fatalf("Error when executing template: %s", err)
	}

    w.Flush()
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

	generateDockerfile(config)
}

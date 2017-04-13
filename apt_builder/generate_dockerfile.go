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
	// "errors"
	"flag"
	// "io/ioutil"
	"log"
	"os"
	// "regexp"
	// "strings"
	// "testing"
	"text/template"

	// "github.com/ghodss/yaml"
)

var INSTALL_TMPL, PPA_TMPL, APT_TMPL, REMOVE_TOOLS_TMPL *template.Template


func init_templates() {
	var err error

	INSTALL_TOOLS := `FROM {{.BaseImage}}
	RUN apt-get update && apt-get install -y --force-yes \\
	    software-properties-common python-software-properties \\`

	INSTALL_TMPL, err = template.New("INSTALL_TOOLS").Parse(INSTALL_TOOLS)
	if err != nil { log.Fatalf("Error creating template: %s", err) }

	// PPA_ADD := `
 //    	&& add-apt-repository -y {ppa} \\`

 //    PPA_TMPL, err := template.New("PPA_ADD").Parse(PPA_ADD)
	// if err != nil { log.Fatalf("Error creating template: %s", err) }

	// APT_INSTALL := `
	//     && apt-get install -y --force-yes \\
	//     {package_list} \\`

	// APT_TMPL, err := template.New("APT_INSTALL").Parse(APT_INSTALL)
	// if err != nil { log.Fatalf("Error creating template: %s", err) }

	// REMOVE_TOOLS := `
	//     && apt-get remove -y --force-yes software-properties-common \\
	//     python-software-properties \\
	//     && apt-get autoremove -y --force-yes \\
	//     && apt-get clean -y --force-yes

	// `

	// REMOVE_TOOLS_TMPL, err := template.New("REMOVE_TOOLS").Parse(REMOVE_TOOLS)
	// if err != nil { log.Fatalf("Error creating template: %s", err) }
}

type Installer struct {
	BaseImage string
}

func generateDockerfile() {
	installer := Installer{"foobar"}
	// TODO: change ostream to file
	err := INSTALL_TMPL.Execute(os.Stdout, installer)
	if err != nil {
		log.Fatalf("Error when executing template: %s", err)
	}
}


var configFile string

func main() {
	init_templates()
	flag.StringVar(&configFile, "yaml", "",
				   "path to the .yaml file containing test definitions.")
	flag.Parse()

	if configFile == "" {
		os.Exit(1)
	}
	log.Printf(configFile)
	generateDockerfile()
}

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
	"log"
	"text/template"
)

var DOCKERFILE_TMPL *template.Template

func init_templates() {
	var err error

	DOCKERFILE := `FROM {{.BaseImage}}
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
`

	DOCKERFILE_TMPL, err = template.New("DOCKERFILE").Parse(DOCKERFILE)
	if err != nil {
		log.Fatalf("Error creating template: %s", err)
	}
}

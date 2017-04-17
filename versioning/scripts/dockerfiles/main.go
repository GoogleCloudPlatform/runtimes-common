/*
Command line tool for updating Dockerfiles based on versions.yaml.
*/
package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/GoogleCloudPlatform/runtimes-common/versioning/versions"
)

func renderDockerfile(version versions.Version, tmpl template.Template) []byte {
	var result bytes.Buffer
	tmpl.Execute(&result, version)
	return result.Bytes()
}

func writeDockerfile(version versions.Version, data []byte) {
	path := filepath.Join(version.Dir, "Dockerfile")
	err := ioutil.WriteFile(path, data, 0644)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	templateDirPtr := flag.String("template_dir", "templates", "Path to directory containing Dockerfile.template")

	var spec versions.Spec
	spec = versions.LoadVersions("versions.yaml")

	templatePath := filepath.Join(*templateDirPtr, "Dockerfile.template")
	templateData, err := ioutil.ReadFile(templatePath)
	templateString := string(templateData)
	check(err)

	tmpl, err := template.
		New("dockerfileTemplate").
		Parse(templateString)
	check(err)

	for _, version := range spec.Versions {
		data := renderDockerfile(version, *tmpl)
		writeDockerfile(version, data)
	}
}

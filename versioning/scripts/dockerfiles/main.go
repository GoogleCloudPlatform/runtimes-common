/*
Command line tool for updating Dockerfiles based on versions.yaml.
*/
package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
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
	verifyPtr := flag.Bool("verify", false, "Verify dockerfiles")
	flag.Parse()

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

	if *verifyPtr {
		foundDockerfile := make(map[string]bool)
		failureCount := 0

		for _, version := range spec.Versions {
			data := renderDockerfile(version, *tmpl)

			path := filepath.Join(version.Dir, "Dockerfile")
			dockerfile, err := ioutil.ReadFile(path)
			check(err)

			foundDockerfile[path] = true

			if string(dockerfile) == string(data) {
				log.Printf("%s: OK", path)
			} else {
				failureCount++
				log.Printf("%s: FAIL", path)
			}
		}

		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			check(err)
			if info.Name() == "Dockerfile" && !info.IsDir() && !foundDockerfile[path] {
				failureCount++
				log.Printf("%s: UNIDENTIFIED", path)
			}
			return nil
		})
		check(err)

		os.Exit(failureCount)
	} else {
		for _, version := range spec.Versions {
			data := renderDockerfile(version, *tmpl)
			writeDockerfile(version, data)
		}
	}
}

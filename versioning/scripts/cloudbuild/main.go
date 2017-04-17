/*
Command line tool for generating a Cloud Build yaml file based on versions.yaml.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/runtimes-common/versioning/versions"
)

type cloudBuildOptions struct {
	// Whether to restrict to a particular set of Dockerfile directories.
	// If empty, all directories are used.
	Directories []string

	// Whether to run tests as part of the build.
	RunTests bool

	// Whether to require that image tags do not already exist in the repo.
	RequireNewTags bool

	// Whether to push to all declared tags
	FirstTagOnly bool

	// Optional timeout duration. If not specified, the Cloud Builder default timeout is used.
	TimeoutSeconds int
}

const cloudBuildTemplateString = `
steps:
{{- if .RequireNewTags }}
# Check if tags exist.
{{- range .Images }}
  - name: gcr.io/gcp-runtimes/check_if_tag_exists
    args:
      - 'python'
      - '/main.py'
      - '--image={{ . }}'
{{- end }}
{{- end }}

# Build images
{{- range .ImageBuilds }}
  - name: gcr.io/cloud-builders/docker
    args:
      - 'build'
      - '--tag={{ .Tag }}'
      - '{{ .Directory }}'
{{- end }}

# Run tests
{{- range .ImageBuilds }}
{{- $primary := .Tag }}
{{- range .Tests }}
  - name: gcr.io/gcp-runtimes/structure_test
    args:
      - '--image'
      - '{{ $primary }}'
      - '--config'
      - '{{ . }}'
{{- end }}
{{- end }}

# Add alias tags
{{- range .ImageBuilds }}
{{- $primary := .Tag }}
{{- range .Aliases }}
  - name: gcr.io/cloud-builders/docker
    args:
      - 'tag'
      - '{{ $primary }}'
      - '{{ . }}'
{{- end }}
{{- end }}

images:
{{- range .AllImages }}
  - '{{ . }}'
{{- end }}

{{- if not eq .TimeoutSeconds 0 }}
timeout: {{ .TimeoutSeconds }}s
{{- end }}
`

const testsDir = "tests"
const testJsonSuffix = "_test.json"
const testYamlSuffix = "_test.yaml"

type imageBuildTemplateData struct {
	Directory string
	Tag       string
	Aliases   []string
	Tests     []string
}

type cloudBuildTemplateData struct {
	RequireNewTags bool
	ImageBuilds    []imageBuildTemplateData
	AllImages      []string
	TimeoutSeconds int
}

func newCloudBuildTemplateData(
	registry string, spec versions.Spec, options cloudBuildOptions) cloudBuildTemplateData {
	data := cloudBuildTemplateData{}
	data.RequireNewTags = options.RequireNewTags

	// Determine the set of directories to operate on.
	dirs := make(map[string]bool)
	if len(options.Directories) > 0 {
		for _, d := range options.Directories {
			dirs[d] = true
		}
	} else {
		for _, v := range spec.Versions {
			dirs[v.Dir] = true
		}
	}

	// Extract tests to run.
	var tests []string
	if options.RunTests {
		if info, err := os.Stat(testsDir); err == nil && info.IsDir() {
			files, err := ioutil.ReadDir(testsDir)
			check(err)
			for _, f := range files {
				if strings.HasSuffix(f.Name(), testJsonSuffix) || strings.HasSuffix(f.Name(), testYamlSuffix) {
					tests = append(tests, fmt.Sprintf("/workspace/tests/%s", f.Name()))
				}
			}
		}
	}

	// Extract a list of full image names to build.
	for _, v := range spec.Versions {
		if !dirs[v.Dir] {
			continue
		}
		var images []string
		for _, t := range v.Tags {
			image := fmt.Sprintf("%v/%v:%v", registry, v.Repo, t)
			images = append(images, image)
			if options.FirstTagOnly {
				break
			}
		}
		data.AllImages = append(data.AllImages, images...)
		data.ImageBuilds = append(
			data.ImageBuilds, imageBuildTemplateData{v.Dir, images[0], images[1:], tests})
	}

	data.TimeoutSeconds = options.TimeoutSeconds
	return data
}

func renderCloudBuildConfig(
	registry string, spec versions.Spec, options cloudBuildOptions) string {
	data := newCloudBuildTemplateData(registry, spec, options)
	tmpl, _ := template.
		New("cloudBuildTemplate").
		Parse(cloudBuildTemplateString)
	var result bytes.Buffer
	tmpl.Execute(&result, data)
	return result.String()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	registryPtr := flag.String("registry", "gcr.io/$PROJECT_ID", "Registry, e.g: 'gcr.io/my-project'")
	dirsPtr := flag.String("dirs", "", "Comma separated list of Dockerfile dirs to use.")
	testsPtr := flag.Bool("tests", true, "Run tests.")
	newTagsPtr := flag.Bool("new_tags", false, "Require that image tags do not already exist.")
	firstTagOnly := flag.Bool("first_tag", false, "Build only the first per version.")
	timeoutPtr := flag.Int("timeout", 0, "Timeout in seconds. If not set, the default Cloud Build timeout is used.")
	flag.Parse()

	if *registryPtr == "" {
		log.Fatalf("--registry flag is required")
	}

	if strings.Contains(*registryPtr, ":") {
		*registryPtr = strings.Replace(*registryPtr, ":", "/", 1)
	}

	var dirs []string
	if *dirsPtr != "" {
		dirs = strings.Split(*dirsPtr, ",")
	}

	spec := versions.LoadVersions("versions.yaml")
	options := cloudBuildOptions{dirs, *testsPtr, *newTagsPtr, *firstTagOnly, *timeoutPtr}
	result := renderCloudBuildConfig(*registryPtr, spec, options)
	fmt.Println(result)
}

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
	"reflect"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/runtimes-common/versioning/versions"
	"strconv"
)

type Format struct {
	Indent string
	Count  int
}

func parseSpecToFormat(c map[string]string) Format {
	fIndent := " "
	fCount := 4
	if val, ok := c["template_retry_intendt_type"]; ok {
		fIndent = val
	}
	if val, ok := c["template_retry_intendt_length"]; ok {
		count, err := strconv.Atoi(val)
		check(err)
		fCount = count
	}
	return Format{Indent: fIndent, Count: fCount}
}

func (f Format) indent(s string) string {
	temp := strings.Split(s, "\n")
	str := ""
	for index, line := range temp {
		if index > 0 {
			str = str + "\n"
		}
		trimmed := strings.TrimSpace(line)
		diff := len(line) - len(trimmed)
		prefix := strings.Repeat(f.Indent, diff*f.Count)
		str = str + prefix + trimmed
	}
	return str
}

func renderDockerfile(version versions.Version, tmpl template.Template) []byte {
	var result bytes.Buffer
	tmpl.Execute(&result, version)
	return result.Bytes()
}

func renderRetry(pack versions.Package, tmpl template.Template) string {
	var result bytes.Buffer
	tmpl.Execute(&result, pack)
	return string(result.Bytes())
}

func writeDockerfile(version versions.Version, data []byte) {
	path := filepath.Join(version.Dir, "Dockerfile")
	// Delete first to make sure file is created with the right mode.
	deleteIfFileExists(path)
	err := ioutil.WriteFile(path, data, 0644)
	check(err)
}

func findFilesToCopy(templateDir string, callback func(path string, fileInfo os.FileInfo)) {
	filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		check(err)
		if strings.HasSuffix(info.Name(), ".template") || info.IsDir() {
			return nil
		}
		path, err = filepath.Rel(templateDir, path)
		check(err)
		callback(path, info)
		return nil
	})
}

func copyFiles(version versions.Version, templateDir string) {
	findFilesToCopy(templateDir, func(path string, fileInfo os.FileInfo) {
		data, err := ioutil.ReadFile(filepath.Join(templateDir, path))
		check(err)

		target := filepath.Join(version.Dir, path)
		// Delete first to make sure file is created with the right mode.
		deleteIfFileExists(target)
		err = ioutil.WriteFile(target, data, fileInfo.Mode())
		check(err)
	})
}

func deleteIfFileExists(path string) {
	if fileInfo, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("File %s exists but cannot be stat'ed", path)
		}
	} else {
		if fileInfo.IsDir() {
			log.Fatalf("%s is unexpectedly a directory", path)
		}
		err = os.Remove(path)
		check(err)
	}
}

func verifyDockerfiles(spec versions.Spec, tmpl template.Template) (failureCount int) {
	foundDockerfile := make(map[string]bool)
	failureCount = 0
	warningCount := 0

	for _, version := range spec.Versions {
		data := renderDockerfile(version, tmpl)

		path := filepath.Join(version.Dir, "Dockerfile")
		dockerfile, err := ioutil.ReadFile(path)
		check(err)

		foundDockerfile[path] = true

		if string(dockerfile) == string(data) {
			log.Printf("%s: OK", path)
		} else {
			failureCount++
			log.Printf("%s: FAILED", path)
		}
	}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		check(err)
		if info.Name() == "Dockerfile" && !info.IsDir() && !foundDockerfile[path] {
			warningCount++
			log.Printf("%s: UNIDENTIFIED (warning)", path)
		}
		return nil
	})
	check(err)

	if failureCount == 0 && warningCount > 0 {
		log.Print("Dockerfile verification completed: PASSED (with warnings)")
	} else if failureCount == 0 {
		log.Print("Dockerfile verification completed: PASSED")
	} else {
		log.Print("Dockerfile verification completed: FAILED")
	}

	return
}

func verifyCopiedFiles(spec versions.Spec, templateDir string) (failureCount int) {
	failureCount = 0
	for _, version := range spec.Versions {
		findFilesToCopy(templateDir, func(path string, sourceFileInfo os.FileInfo) {
			failureCount++

			source := filepath.Join(templateDir, path)
			target := filepath.Join(version.Dir, path)
			targetFileInfo, err := os.Stat(target)
			if err != nil {
				log.Printf("%s is expected but cannot be stat'ed", target)
				return
			}

			// Check mode for owner only.
			sourcePerm := os.FileMode(sourceFileInfo.Mode().Perm() & 0700)
			targetPerm := os.FileMode(targetFileInfo.Mode().Perm() & 0700)
			if sourcePerm != targetPerm {
				log.Printf("%s has wrong file mode %v, expected %v", target, targetPerm, sourcePerm)
				return
			}

			expected, err := ioutil.ReadFile(source)
			check(err)
			actual, err := ioutil.ReadFile(filepath.Join(version.Dir, path))
			if err != nil {
				log.Printf("%s is expected but cannot be read", target)
				return
			}

			if !reflect.DeepEqual(expected, actual) {
				log.Printf("%s content is different from its template", target)
				return
			}

			log.Printf("%s: OK", target)

			failureCount--
		})
	}

	if failureCount == 0 {
		log.Print("Copied files verification completed: PASSED")
	} else {
		log.Print("Copied files verification completed: FAILED")
	}

	return
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const templateFrom = `FROM {{ .From }}
`

const templateRetry = `found='' && \
 for server in \
  pool.sks-keyservers.net \
  na.pool.sks-keyservers.net \
  eu.pool.sks-keyservers.net \
  oc.pool.sks-keyservers.net \
  ha.pool.sks-keyservers.net \
  hkp://p80.pool.sks-keyservers.net:80 \
  hkp://keyserver.ubuntu.com:80 \
  pgp.mit.edu \
 ; do \
  {{ $.RetryCommand }} \
   && found=yes && break; \
 done; \
 test -n "$found"`

func main() {
	templateDirPtr := flag.String("template_dir", "templates", "Path to directory containing Dockerfile.template and any other files to copy over")
	verifyPtr := flag.Bool("verify_only", false, "Verify dockerfiles")
	flag.Parse()

	var spec versions.Spec
	spec = versions.LoadVersions("versions.yaml")

	templatePath := filepath.Join(*templateDirPtr, "Dockerfile.template")
	templateData, err := ioutil.ReadFile(templatePath)
	templateString := string(templateData)
	check(err)
	if spec.CheckExtension("TEMPLATE_FROM") {
		templateString = templateFrom + templateString
	}
	if spec.CheckExtension("TEMPLATE_RETRY") {
		templateFormat := parseSpecToFormat(spec.Config)
		retryTmpl, err := template.
			New("retryTemplate").
			Parse(templateFormat.indent(templateRetry))
		check(err)
		for index, element := range spec.Versions {
			for pName, pValues := range element.Packages {
				if pValues.RetryCommand != "" {
					pValues.RetryCommand = renderRetry(pValues, *retryTmpl)
					spec.Versions[index].Packages[pName] = pValues
				}
			}
		}
	}
	tmpl, err := template.
		New("dockerfileTemplate").
		Parse(templateString)
	check(err)

	if *verifyPtr {
		failureCount := verifyDockerfiles(spec, *tmpl)
		failureCount += verifyCopiedFiles(spec, *templateDirPtr)
		os.Exit(failureCount)
	} else {
		for _, version := range spec.Versions {
			data := renderDockerfile(version, *tmpl)
			writeDockerfile(version, data)
			copyFiles(version, *templateDirPtr)
		}
	}
}

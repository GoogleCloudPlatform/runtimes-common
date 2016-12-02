package structure_tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

type problem struct {
	name  string
	issue string
}

func TestImageLicenses(t *testing.T) {
	whitelist := [0]string{}
	blacklist := [4]string{"AGPL", "WTFPL", "AFFERO GENERAL PUBLIC LICENSE", "DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE"}
	root := "/usr/share/doc"
	packages, err := ioutil.ReadDir("/usr/share/doc")
	var problems []problem
	if err != nil {
		t.Fatalf("%s", err)
	}
	for _, p := range packages {
		if !p.IsDir() {
			continue
		}

		// Skip over packages in the whitelist
		for _, w := range whitelist {
			if w == p.Name() {
				continue
			}
		}
		// If package doesn't have copyright file, add it to list of problematic packages
		licenseFile := path.Join(root, p.Name(), "copyright")
		//fmt.Println(licenseFile)
		_, err := os.Stat(licenseFile)
		if err != nil {
			problems = append(problems, problem{p.Name(), err.Error()})
			continue
		}
		// Read through the copyright file and make sure don't have an unauthorized license
		license, err := ioutil.ReadFile(licenseFile)
		if err != nil {
			problems = append(problems, problem{p.Name(), err.Error()})
			continue
		}
		contents := strings.ToUpper(string(license))
		for _, b := range blacklist {
			if strings.Contains(contents, b) {
				problems = append(problems, problem{p.Name(), "invalid license"})
				break
			}
		}
	}

	if len(problems) > 0 {
		for _, p := range problems {
			fmt.Println(p)
		}
		t.Fatalf("The above packages require attention")
	}
}

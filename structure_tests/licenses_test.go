package structure_tests

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

var (
	// Whitelist is the list of packages that we want to automatically pass this
	// check even if it would normally fail for one reason or another.
	whitelist = []string{}

	// Blacklist is the set of words that, if contained in a license file, should cause a failure.
	// This will most likely just be names of unsupported licenses.
	blacklist = []string{"AGPL", "WTFPL", "AFFERO GENERAL PUBLIC LICENSE", "DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE"}
)

func TestImageLicenses(t *testing.T) {
	root := "/usr/share/doc"
	packages, err := ioutil.ReadDir("/usr/share/doc")
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
		_, err := os.Stat(licenseFile)
		if err != nil {
			t.Errorf("Error reading license file for %s: %s", p.Name(), err.Error())
			continue
		}
		// Read through the copyright file and make sure don't have an unauthorized license
		license, err := ioutil.ReadFile(licenseFile)
		if err != nil {
			t.Errorf("Error reading license file for %s: %s", p.Name(), err.Error())
			continue
		}
		contents := strings.ToUpper(string(license))
		for _, b := range blacklist {
			if strings.Contains(contents, b) {
				t.Errorf("Invalid license for %s, license contains %s", p.Name(), b)
				break
			}
		}
	}
}

package differs

import (
	"bufio"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type PackageDiff struct {
	Image1      string
	Packages1   []string
	Image2      string
	Packages2   []string
	VersionDiff []VDiff
}

type VDiff struct {
	Package  string
	Version1 string
	Version2 string
}

func output(diff PackageDiff) error {
	const master = `Packages found only in {{.Image1}}:{{block "list" .Packages1}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}
Packages found only in {{.Image2}}:{{block "list2" .Packages2}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}
Version differences:{{"\n"}}	(Package: {{.Image1}}{{"\t\t"}}{{.Image2}}){{range .VersionDiff}}
	{{.Package}}: {{.Version1}}	{{.Version2}}
	{{end}}`

	funcs := template.FuncMap{"join": strings.Join}

	masterTmpl, err := template.New("master").Funcs(funcs).Parse(master)
	if err != nil {
		log.Fatal(err)
	}

	if err := masterTmpl.Execute(os.Stdout, diff); err != nil {
		log.Fatal(err)
	}
	return nil
}

// AptDiff compares the packages installed by apt-get.
func AptDiff(img1, img2 string) (string, error) {
	pack1, err := getPackages(img1)
	if err != nil {
		return "", err
	}
	pack2, err := getPackages(img2)
	if err != nil {
		return "", err
	}

	diff := diffMaps(pack1, pack2)
	diff.Image1 = img1
	diff.Image2 = img2
	output(diff)
	return "", nil
}

func diffMaps(map1, map2 map[string]string) PackageDiff {
	diff1 := []string{}
	diff2 := []string{}
	versionDiff := []VDiff{}
	for key1, value1 := range map1 {
		value2, ok := map2[key1]
		if !ok {
			diff1 = append(diff1, key1+":"+value1)
		} else if value2 != value1 {
			versionDiff = append(versionDiff, VDiff{key1, value1, value2})
			delete(map2, key1)
		} else {
			delete(map2, key1)
		}
	}
	for key2, value2 := range map2 {
		diff2 = append(diff2, key2+":"+value2)
	}
	diff := PackageDiff{Packages1: diff1, Packages2: diff2, VersionDiff: versionDiff}
	return diff
}

func getPackages(path string) (map[string]string, error) {
	packages := make(map[string]string)

	var layerStems []string

	layers, err := ioutil.ReadDir(path)
	if err != nil {
		return packages, err
	}
	for _, layer := range layers {
		layerStems = append(layerStems, filepath.Join(path, layer.Name(), "layer/var/lib/dpkg/status"))
	}

	for _, statusFile := range layerStems {
		if _, err := os.Stat(statusFile); err == nil {
			if file, err := os.Open(statusFile); err == nil {
				// make sure it gets closed
				defer file.Close()

				var currPackage string
				// create a new scanner and read the file line by line
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := strings.Split(scanner.Text(), ": ")
					if len(line) == 2 {
						key := line[0]
						value := line[1]
						if key == "Package" {
							currPackage = value
						}
						if key == "Version" {
							packages[currPackage] = value
						}
					}

				}

			} else {
				return packages, err
			}
		} else {
			// status file does not exist in this layer
			continue
		}
	}
	return packages, nil
}

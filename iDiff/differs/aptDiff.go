package differs

import (
	"bufio"
	"html/template"
	"log"
	"os"
	"reflect"
	"runtimes-common/iDiff/utils"
	"strings"

	"github.com/golang/glog"
)

// PackageDiff stores the difference information between two images.
type PackageDiff struct {
	Image1    string
	Packages1 []string
	Image2    string
	Packages2 []string
	InfoDiff  []Info
}

// Info stores the information for one package in two different images.
type Info struct {
	Package string
	Info1   PackageInfo
	Info2   PackageInfo
}

// PackageInfo stores the specific metadata about a package.
type PackageInfo struct {
	Version string
	Size    string
}

func output(diff PackageDiff) error {
	const master = `Packages found only in {{.Image1}}:{{block "list" .Packages1}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}
Packages found only in {{.Image2}}:{{block "list2" .Packages2}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}
Version differences:{{"\n"}}	(Package:	{{.Image1}}{{"\t\t"}}{{.Image2}}){{range .InfoDiff}}
	{{.Package}}:	{{.Info1.Version}}	{{.Info2.Version}}
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

func (pi PackageInfo) string() string {
	return pi.Version
}

func diffMaps(map1, map2 map[string]PackageInfo) PackageDiff {
	diff1 := []string{}
	diff2 := []string{}
	infoDiff := []Info{}
	for key1, value1 := range map1 {
		value2, ok := map2[key1]
		if !ok {
			diff1 = append(diff1, key1+":"+value1.string())
		} else if !reflect.DeepEqual(value2, value1) {
			infoDiff = append(infoDiff, Info{key1, value1, value2})
			delete(map2, key1)
		} else {
			delete(map2, key1)
		}
	}
	for key2, value2 := range map2 {
		diff2 = append(diff2, key2+":"+value2.string())
	}
	diff := PackageDiff{Packages1: diff1, Packages2: diff2, InfoDiff: infoDiff}
	return diff
}

func getPackages(path string) (map[string]PackageInfo, error) {
	packages := make(map[string]PackageInfo)
	layerStems, err := utils.BuildLayerTargets(path, "layer/var/lib/dpkg/status")
	if err != nil {
		return packages, err
	}
	for _, statusFile := range layerStems {
		if _, err := os.Stat(statusFile); err != nil {
			// status file does not exist in this layer
			continue
		}
		if file, err := os.Open(statusFile); err == nil {
			// make sure it gets closed
			defer file.Close()

			// create a new scanner and read the file line by line
			scanner := bufio.NewScanner(file)
			var currPackage string
			for scanner.Scan() {
				currPackage = parseLine(scanner.Text(), currPackage, packages)
			}
		} else {
			return packages, err
		}
	}
	return packages, nil
}

func parseLine(text string, currPackage string, packages map[string]PackageInfo) string {
	line := strings.Split(text, ": ")
	if len(line) == 2 {
		key := line[0]
		value := line[1]

		switch key {
		case "Package":
			return value
		case "Version":
			if packages[currPackage].Version != "" {
				glog.Warningln("Multiple versions of same package detected.  Diffing such multi-versioning not yet supported.")
				return currPackage
			}
			tempPackage := PackageInfo{Version: value}
			packages[currPackage] = tempPackage
			return currPackage

		case "Installed-Size":
			tempPackage := PackageInfo{Size: value}
			packages[currPackage] = tempPackage
			return currPackage
		default:
			return currPackage
		}
	}
	return currPackage
}

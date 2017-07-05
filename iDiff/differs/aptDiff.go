package differs

import (
	"bufio"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/golang/glog"
)

func output(diff utils.PackageDiff) error {
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
func AptDiff(d1file, d2file string) (string, error) {
	d1, err := utils.GetDirectory(d1file)
	if err != nil {
		glog.Errorf("Error reading directory structure from file %s: %s\n", d1file, err)
		return "", err
	}
	d2, err := utils.GetDirectory(d2file)
	if err != nil {
		glog.Errorf("Error reading directory structure from file %s: %s\n", d2file, err)
		return "", err
	}

	dirPath1 := d1.Root
	dirPath2 := d2.Root
	pack1, err := getPackages(dirPath1)
	if err != nil {
		return "", err
	}
	pack2, err := getPackages(dirPath2)
	if err != nil {
		return "", err
	}

	diff := utils.DiffMaps(pack1, pack2)
	diff.Image1 = dirPath1
	diff.Image2 = dirPath2
	diff.OutputDiff("")
	return "", nil
}

func getPackages(path string) (map[string]utils.PackageInfo, error) {
	packages := make(map[string]utils.PackageInfo)
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

func parseLine(text string, currPackage string, packages map[string]utils.PackageInfo) string {
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
			modifiedValue := strings.Replace(value, "+", " ", 1)
			tempPackage := utils.PackageInfo{Version: modifiedValue}
			packages[currPackage] = tempPackage
			return currPackage

		case "Installed-Size":
			tempPackage := utils.PackageInfo{Size: value}
			packages[currPackage] = tempPackage
			return currPackage
		default:
			return currPackage
		}
	}
	return currPackage
}

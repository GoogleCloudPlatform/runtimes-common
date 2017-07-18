package differs

import (
	"bufio"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/golang/glog"
)

// AptDiff compares the packages installed by apt-get.
func AptDiff(img1, img2 string, json bool, eng bool) (string, error) {
	pack1, err := getPackages(img1)
	if err != nil {
		return "", err
	}
	pack2, err := getPackages(img2)
	if err != nil {
		return "", err
	}

	diff := utils.GetMapDiff(pack1, pack2)
	diff.Image1 = img1
	diff.Image2 = img2

	if json {
		return utils.JSONify(diff)
	}
	utils.Output(diff)
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
			currPackageInfo, ok := packages[currPackage]
			if !ok {
				currPackageInfo = utils.PackageInfo{}
			}
			currPackageInfo.Version = modifiedValue
			packages[currPackage] = currPackageInfo
			return currPackage

		case "Installed-Size":
			currPackageInfo, ok := packages[currPackage]
			if !ok {
				currPackageInfo = utils.PackageInfo{}
			}
			currPackageInfo.Size = value
			packages[currPackage] = currPackageInfo
			return currPackage
		default:
			return currPackage
		}
	}
	return currPackage
}

package differs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/runtimes-common/iDiff/utils"
)

// NodeDiff compares the packages installed by apt-get.
func NodeDiff(d1file, d2file string) (string, error) {
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
	pack1, err := getNodePackages(dirPath1)
	if err != nil {
		return "", err
	}
	pack2, err := getNodePackages(dirPath2)
	if err != nil {
		return "", err
	}

	diff := utils.DiffMaps(pack1, pack2)
	diff.Image1 = dirPath1
	diff.Image2 = dirPath2
	output(diff)
	return "", nil

}

func buildNodePaths(path string) []string {
	globalPaths, _ := utils.BuildLayerTargets(path, "layer/node_modules")
	localPaths, _ := utils.BuildLayerTargets(path, "layer/usr/local/lib/node_modules")
	return append(globalPaths, localPaths...)
}

func getPackageSize(path string) int64 {
	packagePath := strings.TrimSuffix(path, "package.json")
	packageStat, _ := os.Stat(packagePath)
	return packageStat.Size()
}

func getNodePackages(path string) (map[string]utils.PackageInfo, error) {
	packages := make(map[string]utils.PackageInfo)

	layerStems := buildNodePaths(path)

	for _, modulesDir := range layerStems {

		packageJSONs, _ := utils.BuildLayerTargets(modulesDir, "package.json")
		for _, currPackage := range packageJSONs {
			packageJSON, _ := readPackageJSON(currPackage)
			var currInfo utils.PackageInfo
			currInfo.Version = packageJSON.Version
			currInfo.Size = string(getPackageSize(path))
			packages[packageJSON.Name] = currInfo
		}
	}
	return packages, nil
}

type nodePackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func readPackageJSON(path string) (nodePackage, error) {
	var currPackage nodePackage
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return currPackage, err
	}
	err = json.Unmarshal(jsonBytes, &currPackage)
	if err != nil {
		return currPackage, err
	}
	return currPackage, err
}

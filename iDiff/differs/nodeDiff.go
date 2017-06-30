package differs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/golang/glog"
)

// NodeDiff compares the packages installed by apt-get.
// TODO: Move this code to a place so that it isn't repeated within each specific differ.
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
		glog.Errorf("Error reading packages from directory %s: %s\n", dirPath1, err)
		return "", err
	}
	pack2, err := getNodePackages(dirPath2)
	if err != nil {
		glog.Errorf("Error reading packages from directory %s: %s\n", dirPath2, err)
		return "", err
	}

	diff := utils.DiffMaps(pack1, pack2)
	diff.Image1 = dirPath1
	diff.Image2 = dirPath2
	output(diff)
	return "", nil
}

func buildNodePaths(path string) ([]string, error) {
	globalPaths, err := utils.BuildLayerTargets(path, "layer/node_modules")
	if err != nil {
		return []string{}, err
	}
	localPaths, err := utils.BuildLayerTargets(path, "layer/usr/local/lib/node_modules")
	if err != nil {
		return []string{}, err
	}
	return append(globalPaths, localPaths...), nil
}

func getPackageSize(path string) (int64, error) {
	packagePath := strings.TrimSuffix(path, "package.json")
	packageStat, err := os.Stat(packagePath)
	if err != nil {
		return 0, err
	}
	return packageStat.Size(), nil
}

func getNodePackages(path string) (map[string]utils.PackageInfo, error) {
	packages := make(map[string]utils.PackageInfo)

	layerStems, err := buildNodePaths(path)
	if err != nil {
		glog.Warningf("Error building JSON paths at %s: %s\n", path, err)
		return packages, err
	}

	for _, modulesDir := range layerStems {
		packageJSONs, _ := utils.BuildLayerTargets(modulesDir, "package.json")
		for _, currPackage := range packageJSONs {
			if _, err := os.Stat(currPackage); err != nil {
				// package.json file does not exist at this target path
				continue
			}
			packageJSON, _ := readPackageJSON(currPackage)
			if err != nil {
				glog.Warningf("Error reading package JSON at %s: %s\n", currPackage, err)
				return packages, err
			}
			var currInfo utils.PackageInfo
			currInfo.Version = packageJSON.Version
			size, _ := getPackageSize(currPackage)
			if err != nil {
				glog.Warningf("Error getting package size at %s: %s\n", currPackage, err)
				return packages, err
			}
			currInfo.Size = strconv.FormatInt(size, 10)
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

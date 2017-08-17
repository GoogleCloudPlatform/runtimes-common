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

type NodeDiffer struct {
}

// NodeDiff compares the packages installed by apt-get.
func (d NodeDiffer) Diff(image1, image2 utils.Image) (utils.DiffResult, error) {
	diff, err := multiVersionDiff(image1, image2, d)
	return diff, err
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

func (d NodeDiffer) getPackages(path string) (map[string]map[string]utils.PackageInfo, error) {
	packages := make(map[string]map[string]utils.PackageInfo)

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
			packageJSON, err := readPackageJSON(currPackage)
			if err != nil {
				glog.Warningf("Error reading package JSON at %s: %s\n", currPackage, err)
				return packages, err
			}
			// Build PackageInfo for this package occurence
			var currInfo utils.PackageInfo
			currInfo.Version = packageJSON.Version
			packagePath := strings.TrimSuffix(currPackage, "package.json")
			size, err := utils.GetDirectorySize(packagePath)
			if err != nil {
				glog.Warningf("Error getting package size at %s: %s\n", currPackage, err)
				return packages, err
			}
			currInfo.Size = strconv.FormatInt(size, 10)

			// Check if other package version already recorded
			if _, ok := packages[packageJSON.Name]; !ok {
				// package not yet seen
				infoMap := make(map[string]utils.PackageInfo)
				infoMap[currPackage] = currInfo
				packages[packageJSON.Name] = infoMap
				continue
			}
			packages[packageJSON.Name][currPackage] = currInfo

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

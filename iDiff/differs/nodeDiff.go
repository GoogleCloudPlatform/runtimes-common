package differs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/golang/glog"
)

type NodeDiffer struct {
}

// NodeDiff compares the packages installed by apt-get.
func (d NodeDiffer) Diff(image1, image2 utils.Image) (utils.DiffResult, error) {
	img1 := image1.FSPath
	img2 := image2.FSPath

	pack1, err := getNodePackages(img1)
	if err != nil {
		glog.Errorf("Error reading packages from directory %s: %s\n", img1, err)
		return &utils.MultiVersionPackageDiffResult{}, err
	}
	pack2, err := getNodePackages(img2)
	if err != nil {
		glog.Errorf("Error reading packages from directory %s: %s\n", img2, err)
		return &utils.MultiVersionPackageDiffResult{}, err
	}

	diff := utils.GetMultiVersionMapDiff(pack1, pack2, img1, img2)
	diff.DiffType = "Node Diff"
	return &diff, nil
}

func buildNodePaths(path string) ([]string, error) {
	globalPaths := filepath.Join(path, "node_modules")
	localPath := filepath.Join(path, "usr/local/lib/node_modules")
	return []string{globalPaths, localPath}, nil
}

func getPackageSize(path string) (int64, error) {
	packagePath := strings.TrimSuffix(path, "package.json")
	packageStat, err := os.Stat(packagePath)
	if err != nil {
		return 0, err
	}
	return packageStat.Size(), nil
}

type packageLock struct {
	PackageMap map[string]packageObj `json:"dependencies"`
}

type packageObj struct {
	Version string `json:"version"`
}

func readPackages(path string) (map[string]utils.PackageInfo, error) {
	packages := make(map[string]utils.PackageInfo)
	packageFile, err := os.Open(path)
	if err != nil {
		return packages, err
	}
	jsonParser := json.NewDecoder(packageFile)
	var packagesStruct packageLock
	if err = jsonParser.Decode(&packagesStruct); err != nil {
		return packages, err
	}

	for pack, obj := range packagesStruct.PackageMap {
		// pack := packagesStruct.PackageMap[key]
		newPack := utils.PackageInfo{Version: obj.Version}
		packages[pack] = newPack
	}
	return packages, nil
}

func getNodePackages(path string) (map[string]map[string]utils.PackageInfo, error) {
	packages := make(map[string]map[string]utils.PackageInfo)
	if _, err := os.Stat(path); err != nil {
		// path provided invalid
		return packages, err
	}
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
			// Build PackageInfo for this package occurence
			var currInfo utils.PackageInfo
			currInfo.Version = packageJSON.Version
			size, _ := getPackageSize(currPackage)
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

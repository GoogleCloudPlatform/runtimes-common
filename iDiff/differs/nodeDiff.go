package differs

import (
	"encoding/json"
	"io/ioutil"
)

// NodeDiff compares the packages installed by apt-get.
// func NodeDiff(img1, img2 string) (string, error) {
// 	pack1, err := getPackages(img1)
// 	if err != nil {
// 		return "", err
// 	}
// 	pack2, err := getPackages(img2)
// 	if err != nil {
// 		return "", err
// 	}

// 	diff1, diff2 := diffMaps(pack1, pack2)

// }

func buildNodePaths(path string) []string {
	"layer/node_modules"
	"layer/usr/local/lib"
}

func getPackages(path string) (map[string]utils.PackageInfo, error) {
	packages := make(map[string]utils.PackageInfo)

	layerStems := buildNodePaths(path)

	for _, modulesDir := range layerStems {

		packageJSONs := utils.BuildLayerTargets(modulesDir, "package.json")
			for _, currPackage := range packageJSONs {
				packages[packageJSON.Name] := utils.PackageInfo{Version:packageJSON.Version}
			}
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

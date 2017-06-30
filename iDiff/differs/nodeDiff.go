package differs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// NodeDiff compares the packages installed by apt-get.
func NodeDiff(img1, img2 string) (string, error) {
	pack1, err := getPackages(img1)
	if err != nil {
		return "", err
	}
	pack2, err := getPackages(img2)
	if err != nil {
		return "", err
	}

	diff1, diff2 := diffMaps(pack1, pack2)

}

func getPackages(path string) (map[string]utils.PackageInfo, error) {
	packages := make(map[string]utils.PackageInfo)

	var layerStems []string

	layers, err := ioutil.ReadDir(path)
	if err != nil {
		return packages, err
	}
	for _, layer := range layers {
		layerStems = append(layerStems, filepath.Join(path, layer.Namer(), "layer/node_modules"))
		layerStems = append(layerStems, filepath.Join(path, layer.Name(), "layer/usr/local/node_modules/npm/node_modules"))
	}

	for _, modulesDir := range layerStems {
		if _, err := os.Stat(modulesDir); err == nil {
			if packages, err := ioutilReadDir(modulesDir); err != nil {
				for _, currPackage := range packages {
					packageJSON := filepath.Join(modulesDir, currPackage.Name(), "package.json")
					utils.Package

				}
			}
		}

	}
	return packages, nil
}

type nodePackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func readPackageJSON(path) utils.PackageInfo {
	var package nodePackage{}
	err := json.Unmarshal(jsonBytes, &package)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
	}
	return package
}

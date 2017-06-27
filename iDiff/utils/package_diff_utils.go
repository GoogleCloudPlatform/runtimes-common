package utils

import (
	"io/ioutil"
	"path/filepath"
)

// BuildLayerTargets creates a string slice of the changed layers with the target path concatenated.
func BuildLayerTargets(path, target string) ([]string, error) {
	var layerStems []string
	layers, err := ioutil.ReadDir(path)
	if err != nil {
		return layerStems, err
	}
	if err != nil {
		return layerStems, err
	}
	for _, layer := range layers {
		layerStems = append(layerStems, filepath.Join(path, layer.Name(), target))
	}
	return layerStems, nil
}

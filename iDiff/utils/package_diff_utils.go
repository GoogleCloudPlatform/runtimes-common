package utils

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
)

// PackageDiff stores the difference information between two images.
type PackageDiff struct {
	Image1    string
	Packages1 map[string]PackageInfo
	Image2    string
	Packages2 map[string]PackageInfo
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

// DiffMaps determines the differences between maps of package names to PackageInfo structs
// The return struct includes a list of packages only in the first map, a list of packages only in
// the second map, and a list of packages which differed only in their PackageInfo (version, size, etc.)
func DiffMaps(map1, map2 map[string]PackageInfo) PackageDiff {
	diff1 := map[string]PackageInfo{}
	diff2 := map[string]PackageInfo{}
	infoDiff := []Info{}
	for key1, value1 := range map1 {
		value2, ok := map2[key1]
		if !ok {
			diff1[key1] = value1
		} else if !reflect.DeepEqual(value2, value1) {
			infoDiff = append(infoDiff, Info{key1, value1, value2})
			delete(map2, key1)
		} else {
			delete(map2, key1)
		}
	}
	for key2, value2 := range map2 {
		diff2[key2] = value2
	}
	diff := PackageDiff{Packages1: diff1, Packages2: diff2, InfoDiff: infoDiff}
	return diff
}

func (pi PackageInfo) string() string {
	return pi.Version
}

// BuildLayerTargets creates a string slice of the layers found at path with the target concatenated.
func BuildLayerTargets(path, target string) ([]string, error) {
	layerStems := []string{}
	layers, err := ioutil.ReadDir(path)
	if err != nil {
		return layerStems, err
	}
	for _, layer := range layers {
		if layer.IsDir() {
			layerStems = append(layerStems, filepath.Join(path, layer.Name(), target))
		}
	}
	return layerStems, nil
}

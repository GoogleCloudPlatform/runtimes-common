package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"

	"github.com/golang/glog"
)

type MultiVersionPackageDiff struct {
	Image1    string
	Packages1 map[string]map[string]PackageInfo
	Image2    string
	Packages2 map[string]map[string]PackageInfo
	InfoDiff  []Info
}

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
	// Path    string
}

func multiVersionDiff(infoDiff []Info, key1 string, value1, value2 map[string]PackageInfo) {
	return
}

func checkPackageMapType(map1, map2 interface{}) (reflect.Type, bool, error) {
	// check types and determine multi-version package maps or not
	map1Kind := reflect.ValueOf(map1)
	map2Kind := reflect.ValueOf(map2)
	if map1Kind.Kind() != reflect.Map && map2Kind.Kind() != reflect.Map {
		return nil, false, fmt.Errorf("Package maps were not maps.  Instead were: %s and %s", map1Kind.Kind(), map2Kind.Kind())
	}
	mapType := reflect.TypeOf(map1)
	if mapType != reflect.TypeOf(map2) {
		return nil, false, fmt.Errorf("Package maps were of different types")
	}
	multiVersion := false
	if mapType == reflect.TypeOf(map[string]map[string]PackageInfo{}) {
		multiVersion = true
	}
	return mapType, multiVersion, nil
}

// DiffMaps determines the differences between maps of package names to PackageInfo structs
// The return struct includes a list of packages only in the first map, a list of packages only in
// the second map, and a list of packages which differed only in their PackageInfo (version, size, etc.)
func DiffMaps(map1, map2 interface{}) interface{} {
	mapType, multiV, err := checkPackageMapType(map1, map2)
	if err != nil {
		glog.Error(err)
	}

	map1Value := reflect.ValueOf(map1)
	map2Value := reflect.ValueOf(map2)

	diff1 := reflect.MakeMap(mapType)
	diff2 := reflect.MakeMap(mapType)
	infoDiff := []Info{}
	for _, key1 := range map1Value.MapKeys() {
		value1 := map1Value.MapIndex(key1)
		value2 := map2Value.MapIndex(key1)
		if !value2.IsValid() { //reflect.New(reflect.TypeOf(value2)) {
			diff1.SetMapIndex(key1, value1)
		} else if !reflect.DeepEqual(value2, value1) {
			if multiV {
				multiVersionDiff(infoDiff, key1.String(),
					value1.Interface().(map[string]PackageInfo), value2.Interface().(map[string]PackageInfo))
			} else {
				infoDiff = append(infoDiff, Info{key1.String(), value1.Interface().(PackageInfo),
					value2.Interface().(PackageInfo)})
				map2Value.SetMapIndex(key1, reflect.Value{})
				// delete(map2, key1)
			}
		} else {
			map2Value.SetMapIndex(key1, reflect.Value{})
		}
	}
	for _, key2 := range map2Value.MapKeys() {
		value2 := map2Value.MapIndex(key2)
		diff2.SetMapIndex(key2, value2)
	}
	if multiV {
		return MultiVersionPackageDiff{Packages1: diff1.Interface().(map[string]map[string]PackageInfo),
			Packages2: diff2.Interface().(map[string]map[string]PackageInfo), InfoDiff: infoDiff}
	}
	return PackageDiff{Packages1: diff1.Interface().(map[string]PackageInfo),
		Packages2: diff2.Interface().(map[string]PackageInfo), InfoDiff: infoDiff}
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

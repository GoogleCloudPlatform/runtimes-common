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
	InfoDiff  []MultiVersionInfo
}

type MultiVersionInfo struct {
	Package string
	Info1   []PackageInfo
	Info2   []PackageInfo
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
	// Layer   string
}

func multiVersionDiff(infoDiff []MultiVersionInfo, key1 string, value1,
	value2 map[string]PackageInfo) []MultiVersionInfo {
	diff := GetMapDiff(value1, value2)
	fmt.Println(diff)
	// diffVal := reflect.ValueOf(diff)
	// packDiff := diffVal.Interface().(PackageDiff)
	packageVersions1 := []PackageInfo{}
	packageVersions2 := []PackageInfo{}
	for _, val := range diff.Packages1 {
		packageVersions1 = append(packageVersions1, val)
	}
	for _, val2 := range diff.Packages2 {
		packageVersions2 = append(packageVersions2, val2)
	}
	for _, val3 := range diff.InfoDiff {

		if !reflect.DeepEqual(val3.Info1, val3.Info2) {

			packageVersions1 = append(packageVersions1, val3.Info1)
			packageVersions2 = append(packageVersions2, val3.Info2)
		}
	}

	infoDiff = append(infoDiff, MultiVersionInfo{key1, packageVersions1, packageVersions2})
	return infoDiff
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

// GetMapDiff determines the differences between maps of package names to PackageInfo structs
// This getter supports only single version packages.
func GetMapDiff(map1, map2 map[string]PackageInfo) PackageDiff {
	diff := diffMaps(map1, map2)
	diffVal := reflect.ValueOf(diff)
	packDiff := diffVal.Interface().(PackageDiff)
	return packDiff
}

// GetMultiVersionMapDiff determines the differences between two image package maps with multi-version packages
// This getter supports multi version packages.
func GetMultiVersionMapDiff(map1, map2 map[string]map[string]PackageInfo) MultiVersionPackageDiff {
	diff := diffMaps(map1, map2)
	diffVal := reflect.ValueOf(diff)
	packDiff := diffVal.Interface().(MultiVersionPackageDiff)
	return packDiff
}

// DiffMaps determines the differences between maps of package names to PackageInfo structs
// The return struct includes a list of packages only in the first map, a list of packages only in
// the second map, and a list of packages which differed only in their PackageInfo (version, size, etc.)
func diffMaps(map1, map2 interface{}) interface{} {
	mapType, multiV, err := checkPackageMapType(map1, map2)
	if err != nil {
		glog.Error(err)
	}

	map1Value := reflect.ValueOf(map1)
	map2Value := reflect.ValueOf(map2)

	diff1 := reflect.MakeMap(mapType)
	diff2 := reflect.MakeMap(mapType)
	infoDiff := []Info{}
	multiInfoDiff := []MultiVersionInfo{}

	for _, key1 := range map1Value.MapKeys() {
		value1 := map1Value.MapIndex(key1)
		value2 := map2Value.MapIndex(key1)
		if !value2.IsValid() {
			diff1.SetMapIndex(key1, value1)
		} else if !reflect.DeepEqual(value2.Interface(), value1.Interface()) {
			if multiV {
				multiInfoDiff = multiVersionDiff(multiInfoDiff, key1.String(),
					value1.Interface().(map[string]PackageInfo), value2.Interface().(map[string]PackageInfo))
			} else {
				infoDiff = append(infoDiff, Info{key1.String(), value1.Interface().(PackageInfo),
					value2.Interface().(PackageInfo)})
			}
			map2Value.SetMapIndex(key1, reflect.Value{})
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
			Packages2: diff2.Interface().(map[string]map[string]PackageInfo), InfoDiff: multiInfoDiff}
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

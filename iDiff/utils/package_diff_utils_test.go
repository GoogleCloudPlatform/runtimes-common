package utils

import (
	"reflect"
	"sort"
	"testing"
)

type ByPackage []Info

func (a ByPackage) Len() int {
	return len(a)
}

func (a ByPackage) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByPackage) Less(i, j int) bool {
	return a[i].Package < a[j].Package
}

func TestDiffMaps(t *testing.T) {
	testCases := []struct {
		descrip  string
		map1     interface{}
		map2     interface{}
		expected interface{}
	}{
		{
			descrip: "Missing Packages.",
			map1: map[string]PackageInfo{
				"pac1": {"1.0", "40"},
				"pac3": {"3.0", "60"}},
			map2: map[string]PackageInfo{
				"pac4": {"4.0", "70"},
				"pac5": {"5.0", "80"}},
			expected: PackageDiff{
				Packages1: map[string]PackageInfo{
					"pac1": {"1.0", "40"},
					"pac3": {"3.0", "60"}},
				Packages2: map[string]PackageInfo{
					"pac4": {"4.0", "70"},
					"pac5": {"5.0", "80"}},
				InfoDiff: []Info{}},
		},
		{
			descrip: "Different Versions and Sizes.",
			map1: map[string]PackageInfo{
				"pac2": {"2.0", "50"},
				"pac3": {"3.0", "60"}},
			map2: map[string]PackageInfo{
				"pac2": {"2.0", "45"},
				"pac3": {"4.0", "60"}},
			expected: PackageDiff{
				Packages1: map[string]PackageInfo{},
				Packages2: map[string]PackageInfo{},
				InfoDiff: []Info{
					{"pac2", PackageInfo{"2.0", "50"}, PackageInfo{"2.0", "45"}},
					{"pac3", PackageInfo{"3.0", "60"}, PackageInfo{"4.0", "60"}}},
			},
		},
		{
			descrip: "Identical packages, versions, and sizes",
			map1: map[string]PackageInfo{
				"pac1": {"1.0", "40"},
				"pac2": {"2.0", "50"},
				"pac3": {"3.0", "60"}},
			map2: map[string]PackageInfo{
				"pac1": {"1.0", "40"},
				"pac2": {"2.0", "50"},
				"pac3": {"3.0", "60"}},
			expected: PackageDiff{
				Packages1: map[string]PackageInfo{},
				Packages2: map[string]PackageInfo{},
				InfoDiff:  []Info{}},
		},
		// {
		// 	descrip: "MultiVersion Packages",
		// 	map1: map[string]map[string]PackageInfo{
		// 		"pac1": {"layer1/layer/node_modules/pac1": {"1.0", "40"}},
		// 		"pac2": {"layer1/layer/usr/local/lib/node_modules/pac2": {"2.0", "50"}},
		// 		"pac3": {"layer2/layer/usr/local/lib/node_modules/pac2": {"3.0", "50"}}},
		// 	map2: map[string]map[string]PackageInfo{
		// 		"pac1": {"layer1/layer/node_modules/pac1": {"2.0", "40"}},
		// 		"pac2": {"layer1/layer/usr/local/lib/node_modules/pac2": {"4.0", "50"}},
		// 		"pac3": {"layer2/layer/usr/local/lib/node_modules/pac2": {"3.0", "50"}}},
		// 	expected: MultiVersionPackageDiff{
		// 		Packages1: map[string]map[string]PackageInfo{},
		// 		Packages2: map[string]map[string]PackageInfo{},
		// 		InfoDiff:  []Info{}},
		// },
	}
	for _, test := range testCases {
		diff := DiffMaps(test.map1, test.map2)
		switch test.expected
		sort.Sort(ByPackage(test.expected.InfoDiff))
		sort.Sort(ByPackage(diff.InfoDiff))
		if !reflect.DeepEqual(test.expected, diff) {
			t.Errorf("Expected Diff to be: %s but got: %s", test.expected, diff)
		}
	}
}

func TestCheckPackageMapType(t *testing.T) {
	testCases := []struct {
		descrip       string
		map1          interface{}
		map2          interface{}
		expectedType  reflect.Type
		expectedMulti bool
		err           bool
	}{
		{
			descrip: "Map arguments not maps",
			map1:    "not a map",
			map2:    "not a map either",
			err:     true,
		},
		{
			descrip: "Map arguments not same type",
			map1:    map[string]int{},
			map2:    map[int]string{},
			err:     true,
		},
		{
			descrip:      "Single Version Package Maps",
			map1:         map[string]PackageInfo{},
			map2:         map[string]PackageInfo{},
			expectedType: reflect.TypeOf(map[string]PackageInfo{}),
		},
		{
			descrip:       "MultiVersion Package Maps",
			map1:          map[string]map[string]PackageInfo{},
			map2:          map[string]map[string]PackageInfo{},
			expectedType:  reflect.TypeOf(map[string]map[string]PackageInfo{}),
			expectedMulti: true,
		},
	}
	for _, test := range testCases {
		actualType, actualMulti, err := checkPackageMapType(test.map1, test.map2)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Error("Expected error but got none.")
		}
		if actualType != test.expectedType {
			t.Errorf("Expected type: %s but got: %s", test.expectedType, actualType)
		}
		if actualMulti != test.expectedMulti {
			t.Errorf("Expected multi: %t but got %t", test.expectedMulti, actualMulti)
		}
	}
}
func TestBuildLayerTargets(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		target   string
		expected []string
		err      bool
	}{
		{
			descrip:  "Filter out non directories",
			path:     "testTars/la-croix1-actual",
			target:   "123",
			expected: []string{},
		},
		{
			descrip:  "Error on bad directory path",
			path:     "test_files/notReal",
			target:   "123",
			expected: []string{},
			err:      true,
		},
		{
			descrip:  "Filter out non-directories and get directories",
			path:     "testTars/la-croix3-full",
			target:   "123",
			expected: []string{"testTars/la-croix3-full/nest/123", "testTars/la-croix3-full/nested-dir/123"},
		},
	}
	for _, test := range testCases {
		layers, err := BuildLayerTargets(test.path, test.target)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Errorf("Expected error but got none: %s", err)
		}
		sort.Strings(test.expected)
		sort.Strings(layers)
		if !reflect.DeepEqual(test.expected, layers) {
			t.Errorf("Expected: %s, but got: %s.", test.expected, layers)
		}
	}
}

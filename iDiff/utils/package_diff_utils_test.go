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
		map1     map[string]PackageInfo
		map2     map[string]PackageInfo
		expected PackageDiff
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
	}
	for _, test := range testCases {
		diff := DiffMaps(test.map1, test.map2)
		sort.Sort(ByPackage(test.expected.InfoDiff))
		sort.Sort(ByPackage(diff.InfoDiff))
		if !reflect.DeepEqual(test.expected, diff) {
			t.Errorf("Expected packages only in map1 to be: %s but got: %s", test.expected, diff)
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

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
				"pac1": PackageInfo{"1.0", "40"},
				"pac3": PackageInfo{"3.0", "60"}},
			map2: map[string]PackageInfo{
				"pac4": PackageInfo{"4.0", "70"},
				"pac5": PackageInfo{"5.0", "80"}},
			expected: PackageDiff{
				Packages1: []string{"pac1:1.0", "pac3:3.0"},
				Packages2: []string{"pac4:4.0", "pac5:5.0"},
				InfoDiff:  []Info{}},
		},
		{
			descrip: "Different Versions and Sizes.",
			map1: map[string]PackageInfo{
				"pac2": PackageInfo{"2.0", "50"},
				"pac3": PackageInfo{"3.0", "60"}},
			map2: map[string]PackageInfo{
				"pac2": PackageInfo{"2.0", "45"},
				"pac3": PackageInfo{"4.0", "60"}},
			expected: PackageDiff{
				Packages1: []string{},
				Packages2: []string{},
				InfoDiff: []Info{
					Info{"pac2", PackageInfo{"2.0", "50"}, PackageInfo{"2.0", "45"}},
					Info{"pac3", PackageInfo{"3.0", "60"}, PackageInfo{"4.0", "60"}}},
			},
		},
		{
			descrip: "Identical packages, versions, and sizes",
			map1: map[string]PackageInfo{
				"pac1": PackageInfo{"1.0", "40"},
				"pac2": PackageInfo{"2.0", "50"},
				"pac3": PackageInfo{"3.0", "60"}},
			map2: map[string]PackageInfo{
				"pac1": PackageInfo{"1.0", "40"},
				"pac2": PackageInfo{"2.0", "50"},
				"pac3": PackageInfo{"3.0", "60"}},
			expected: PackageDiff{
				Packages1: []string{},
				Packages2: []string{},
				InfoDiff:  []Info{}},
		},
	}
	for _, test := range testCases {
		diff := DiffMaps(test.map1, test.map2)
		sort.Strings(test.expected.Packages1)
		sort.Strings(diff.Packages1)
		sort.Strings(test.expected.Packages2)
		sort.Strings(diff.Packages2)
		sort.Sort(ByPackage(test.expected.InfoDiff))
		sort.Sort(ByPackage(diff.InfoDiff))
		if !reflect.DeepEqual(test.expected, diff) {
			t.Errorf("Expected packages only in map1 to be: %s but got: %s", test.expected, diff)
		}
	}
}

func TestBuildLayerTargets(t *testing.T) {
	path := "test_files/dir1"
	target := "123"
	expected := []string{"test_files/dir1/file1/123", "test_files/dir1/file2/123", "test_files/dir1/file3/123"}
	layers, err := BuildLayerTargets(path, target)
	if err != nil {
		t.Errorf("Unexpected error occured: %s", err)
	}
	sort.Strings(expected)
	sort.Strings(layers)
	if !reflect.DeepEqual(expected, layers) {
		t.Errorf("Expected: %s, but got: %s.", expected, layers)
	}
}

func TestBuildLayerTargetsFailure(t *testing.T) {
	path := "test_files/notReal"
	target := "123"
	expected := []string{}
	layers, err := BuildLayerTargets(path, target)
	if err == nil {
		t.Errorf("Expected error but none occurred")
	}
	if !reflect.DeepEqual(expected, layers) {
		t.Errorf("Expected: %s, but got: %s.", expected, layers)
	}
}

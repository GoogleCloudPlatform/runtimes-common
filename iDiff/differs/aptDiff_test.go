package differs

import (
	"reflect"
	"sort"
	"testing"
)

type ByPackage []Info

func (a ByPackage) Len() int           { return len(a) }
func (a ByPackage) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPackage) Less(i, j int) bool { return a[i].Package < a[j].Package }
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
				InfoDiff:  []Info{},
			},
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
				InfoDiff: []Info{Info{"pac2", PackageInfo{"2.0", "50"}, PackageInfo{"2.0", "45"}},
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
				InfoDiff:  []Info{},
			},
		},
	}

	for _, test := range testCases {
		diff := diffMaps(test.map1, test.map2)
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

func TestParseLine(t *testing.T) {
	testCases := []struct {
		descrip     string
		line        string
		packages    map[string]PackageInfo
		currPackage string
		expPackage  string
		expected    map[string]PackageInfo
	}{
		{
			descrip:    "Not applicable line",
			line:       "Garbage: garbage info",
			packages:   map[string]PackageInfo{},
			expPackage: "",
			expected:   map[string]PackageInfo{},
		},
		{
			descrip:     "Package line",
			line:        "Package: La-Croix",
			currPackage: "Tea",
			expPackage:  "La-Croix",
			packages:    map[string]PackageInfo{},
			expected:    map[string]PackageInfo{},
		},
		{
			descrip:     "Version line",
			line:        "Version: Lime",
			packages:    map[string]PackageInfo{},
			currPackage: "La-Croix",
			expPackage:  "La-Croix",
			expected:    map[string]PackageInfo{"La-Croix": PackageInfo{Version: "Lime"}},
		},
		{
			descrip:     "Version line",
			line:        "Version: Lime",
			packages:    map[string]PackageInfo{},
			currPackage: "La-Croix",
			expPackage:  "La-Croix",
			expected:    map[string]PackageInfo{"La-Croix": PackageInfo{Version: "Lime"}},
                },
		{
			descrip:     "Size line",
			line:        "Installed-Size: 12floz",
			packages:    map[string]PackageInfo{},
			currPackage: "La-Croix",
			expPackage:  "La-Croix",
			expected:    map[string]PackageInfo{"La-Croix": PackageInfo{Size: "12floz"}},
		},
	}

	for _, test := range testCases {
		currPackage := parseLine(test.line, test.currPackage, test.packages)
		if currPackage != test.expPackage {
			t.Errorf("Expected current package to be: %s, but got: %s.", test.expPackage, currPackage)
		}
		if !reflect.DeepEqual(test.packages, test.expected) {
			t.Errorf("Expected: %s but got: %s", test.expected, test.packages)
		}
	}
}

func TestGetPackages(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		expected map[string]PackageInfo
	}{
		{
			descrip:  "no packages",
			path:     "testDirs/noPackages",
			expected: map[string]PackageInfo{},
		},
		{
			descrip: "all packages in one layer",
			path:    "testDirs/packageOne",
			expected: map[string]PackageInfo{
				"pac1": PackageInfo{Version: "1.0"},
				"pac2": PackageInfo{Version: "2.0"},
				"pac3": PackageInfo{Version: "3.0"}},
		},
		{
			descrip: "many packages in different layers",
			path:    "testDirs/packageMany",
			expected: map[string]PackageInfo{
				"pac1": PackageInfo{Version: "1.0"},
				"pac2": PackageInfo{Version: "2.0"},
				"pac3": PackageInfo{Version: "3.0"},
				"pac4": PackageInfo{Version: "4.0"},
				"pac5": PackageInfo{Version: "5.0"}},
		},
	}
	for _, test := range testCases {
		packages, err := getPackages(test.path)
		if err != nil {
			t.Errorf("Got unexpected error: %s", err)
		}
		if !reflect.DeepEqual(packages, test.expected) {
			t.Errorf("Expected: %s but got: %s", test.expected, packages)
		}
	}
}

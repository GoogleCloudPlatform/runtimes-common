package differs

import (
	"reflect"
	"testing"

	"github.com/runtimes-common/iDiff/utils"
)

func TestParseLine(t *testing.T) {
	testCases := []struct {
		descrip     string
		line        string
		packages    map[string]utils.PackageInfo
		currPackage string
		expPackage  string
		expected    map[string]utils.PackageInfo
	}{
		{
			descrip:    "Not applicable line",
			line:       "Garbage: garbage info",
			packages:   map[string]utils.PackageInfo{},
			expPackage: "",
			expected:   map[string]utils.PackageInfo{},
		},
		{
			descrip:     "Package line",
			line:        "Package: La-Croix",
			currPackage: "Tea",
			expPackage:  "La-Croix",
			packages:    map[string]utils.PackageInfo{},
			expected:    map[string]utils.PackageInfo{},
		},
		{
			descrip:     "Version line",
			line:        "Version: Lime",
			packages:    map[string]utils.PackageInfo{},
			currPackage: "La-Croix",
			expPackage:  "La-Croix",
			expected:    map[string]utils.PackageInfo{"La-Croix": utils.PackageInfo{Version: "Lime"}},
		},
		{
			descrip:     "Version line",
			line:        "Version: Lime",
			packages:    map[string]utils.PackageInfo{},
			currPackage: "La-Croix",
			expPackage:  "La-Croix",
			expected:    map[string]utils.PackageInfo{"La-Croix": utils.PackageInfo{Version: "Lime"}},
		},
		{
			descrip:     "Size line",
			line:        "Installed-Size: 12floz",
			packages:    map[string]utils.PackageInfo{},
			currPackage: "La-Croix",
			expPackage:  "La-Croix",
			expected:    map[string]utils.PackageInfo{"La-Croix": utils.PackageInfo{Size: "12floz"}},
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
		expected map[string]utils.PackageInfo
	}{
		{
			descrip:  "no packages",
			path:     "testDirs/noPackages",
			expected: map[string]utils.PackageInfo{},
		},
		{
			descrip: "all packages in one layer",
			path:    "testDirs/packageOne",
			expected: map[string]utils.PackageInfo{
				"pac1": utils.PackageInfo{Version: "1.0"},
				"pac2": utils.PackageInfo{Version: "2.0"},
				"pac3": utils.PackageInfo{Version: "3.0"}},
		},
		{
			descrip: "many packages in different layers",
			path:    "testDirs/packageMany",
			expected: map[string]utils.PackageInfo{
				"pac1": utils.PackageInfo{Version: "1.0"},
				"pac2": utils.PackageInfo{Version: "2.0"},
				"pac3": utils.PackageInfo{Version: "3.0"},
				"pac4": utils.PackageInfo{Version: "4.0"},
				"pac5": utils.PackageInfo{Version: "5.0"}},
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

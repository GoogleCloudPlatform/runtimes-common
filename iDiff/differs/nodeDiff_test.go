package differs

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

func TestGetNodePackages(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		expected map[string]map[string]utils.PackageInfo
		err      bool
	}{
		{
			descrip:  "no directory",
			path:     "testDirs/notThere",
			expected: map[string]map[string]utils.PackageInfo{},
			err:      true,
		},
		{
			descrip:  "no packages",
			path:     "testDirs/noPackages",
			expected: map[string]map[string]utils.PackageInfo{},
		},
		{
			descrip: "all packages in one layer",
			path:    "testDirs/packageOne",
			expected: map[string]map[string]utils.PackageInfo{
				"pac1": {"testDirs/packageOne/layer1/layer/node_modules/pac1/package.json": {Version: "1.0", Size: "41"}},
				"pac2": {"testDirs/packageOne/layer1/layer/usr/local/lib/node_modules/pac2/package.json": {Version: "2.0", Size: "41"}},
				"pac3": {"testDirs/packageOne/layer1/layer/node_modules/pac3/package.json": {Version: "3.0", Size: "41"}}},
		},
		{
			descrip: "many packages in different layers",
			path:    "testDirs/packageMany",
			expected: map[string]map[string]utils.PackageInfo{
				"pac1": {"testDirs/packageMany/layer1/layer/node_modules/pac1/package.json": {Version: "1.0", Size: "41"}},
				"pac2": {"testDirs/packageMany/layer1/layer/usr/local/lib/node_modules/pac2/package.json": {Version: "2.0", Size: "41"}},
				"pac3": {"testDirs/packageMany/layer2/layer/node_modules/pac3/package.json": {Version: "3.0", Size: "41"}},
				"pac4": {"testDirs/packageMany/layer2/layer/node_modules/pac4/package.json": {Version: "4.0", Size: "41"}},
				"pac5": {"testDirs/packageMany/layer2/layer/node_modules/pac5/package.json": {Version: "5.0", Size: "41"}}},
		},
		{
			descrip: "Multi version packages",
			path:    "testDirs/packageMulti",
			expected: map[string]map[string]utils.PackageInfo{
				"pac1": {"testDirs/packageMulti/layer1/layer/node_modules/pac1/package.json": {Version: "1.0", Size: "41"}},
				"pac2": {"testDirs/packageMulti/layer1/layer/usr/local/lib/node_modules/pac2/package.json": {Version: "2.0", Size: "41"},
					"testDirs/packageMulti/layer2/layer/usr/local/lib/node_modules/pac2/package.json": {Version: "3.0", Size: "41"}}},
		},
	}

	for _, test := range testCases {
		d := NodeDiffer{}
		packages, err := d.getPackages(test.path)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Errorf("Expected error but got none.")
		}
		if !reflect.DeepEqual(packages, test.expected) {
			t.Errorf("Expected: %s but got: %s", test.expected, packages)
		}
	}
}
func TestReadPackageJSON(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		expected nodePackage
		err      bool
	}{
		{
			descrip: "Error on non-existent file",
			path:    "testDirs/not_real.json",
			err:     true,
		},
		{
			descrip:  "Parse JSON with exact fields",
			path:     "testDirs/exact.json",
			expected: nodePackage{"La-croix", "Lime"},
		},
		{
			descrip:  "Parse JSON with additional fields",
			path:     "testDirs/extra.json",
			expected: nodePackage{"La-croix", "Lime"},
		},
	}
	for _, test := range testCases {
		actual, err := readPackageJSON(test.path)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Error("Expected errorbut got none.")
		}
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Expected: %s but got: %s", test.expected, actual)
		}
	}
}

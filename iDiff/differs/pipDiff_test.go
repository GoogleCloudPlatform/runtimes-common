package differs

import (
	"reflect"
	"testing"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

func TestGetPythonVersion(t *testing.T) {
	testCases := []struct {
		layerPath       string
		expectedVersion string
		expectedSuccess bool
	}{
		{
			layerPath:       "testDirs/pipTests/pythonVersionTests/noLibLayer",
			expectedVersion: "",
			expectedSuccess: false,
		},
		{
			layerPath:       "testDirs/pipTests/pythonVersionTests/noPythonLayer",
			expectedVersion: "",
			expectedSuccess: false,
		},
		{
			layerPath:       "testDirs/pipTests/pythonVersionTests/version2.7Layer",
			expectedVersion: "python2.7",
			expectedSuccess: true,
		},
		{
			layerPath:       "testDirs/pipTests/pythonVersionTests/version3.6Layer",
			expectedVersion: "python3.6",
			expectedSuccess: true,
		},
	}
	for _, test := range testCases {
		version, success := getPythonVersion(test.layerPath)
		if success != test.expectedSuccess {
			if test.expectedSuccess {
				t.Error("Expected success finding version but got none")
			} else {
				t.Errorf("Expected failure finding version but found one: %s", version)
			}
		} else if version != test.expectedVersion {
			t.Errorf("Expected: %s.  Got: %s", test.expectedVersion, version)
		}
	}
}

func TestGetPythonPackages(t *testing.T) {
	testCases := []struct {
		path             string
		expectedPackages map[string]utils.PackageInfo
	}{
		{
			path:             "testDirs/pipTests/noPackagesTest",
			expectedPackages: map[string]utils.PackageInfo{},
		},
		{
			path: "testDirs/pipTests/packagesManyLayers",
			expectedPackages: map[string]utils.PackageInfo{
				"packageone":   {Version: "3.6.9", Size: "0"},
				"packagetwo":   {Version: "4.6.2", Size: "0"},
				"packagethree": {Version: "2.4.5", Size: "0"},
				"packagefour":  {Version: "2.4.6", Size: "0"},
			},
		},
		{
			path: "testDirs/pipTests/packagesOneLayer",
			expectedPackages: map[string]utils.PackageInfo{
				"packageone": {Version: "3.6.9", Size: "0"},
				"packagetwo": {Version: "4.6.2", Size: "0"},
			},
		},
	}
	for _, test := range testCases {
		d := PipDiffer{}
		packages, _ := d.getPackages(test.path)
		if !reflect.DeepEqual(packages, test.expectedPackages) {
			t.Errorf("Expected: %s but got: %s", test.expectedPackages, packages)
		}
	}
}

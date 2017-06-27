package differs

import (
	"reflect"
	"sort"
	"testing"
)

func TestDiffMaps(t *testing.T) {
	map1 := map[string]string{"pac1": "1.0", "pac2": "2.0", "pac3": "3.0"}
	map2 := map[string]string{"pac1": "1.0", "pac2": "2.0", "pac3": "4.0",
		"pac4": "4.0", "pac5": "5.0"}
	expected := PackageDiff{
		Packages1:   []string{},
		Packages2:   []string{"pac4:4.0", "pac5:5.0"},
		VersionDiff: []VDiff{VDiff{"pac3", "3.0", "4.0"}},
	}
	diff := diffMaps(map1, map2)
	sort.Strings(expected.Packages1)
	sort.Strings(diff.Packages1)
	sort.Strings(expected.Packages2)
	sort.Strings(diff.Packages2)
	if !reflect.DeepEqual(expected, diff) {
		t.Errorf("Expected: %s but got: %s", expected, diff)
	}
}

func TestGetPackages(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		expected map[string]string
	}{
		{
			descrip:  "no packages",
			path:     "testDirs/noPackages",
			expected: map[string]string{},
		},
		{
			descrip:  "all packages in one layer",
			path:     "testDirs/packageOne",
			expected: map[string]string{"pac1": "1.0", "pac2": "2.0", "pac3": "3.0"},
		},
		{
			descrip: "many packages in different layers",
			path:    "testDirs/packageMany",
			expected: map[string]string{"pac1": "1.0", "pac2": "2.0", "pac3": "3.0",
				"pac4": "4.0", "pac5": "5.0"},
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

package differs

import (
	"reflect"
	"testing"
)

func TestDiffMaps(t *testing.T) {
	map1 := map[string]string{"pac1": "1.0", "pac2": "2.0", "pac3": "3.0"}
	map2 := map[string]string{"pac1": "1.0", "pac2": "2.0", "pac3": "3.0",
		"pac4": "4.0", "pac5": "5.0"}
	expected1 := []string{}
	expected2 := []string{"pac4:4.0", "pac5:5.0"}
	diff1, diff2 := diffMaps(map1, map2)
	if !reflect.DeepEqual(expected1, diff1) {
		t.Errorf("Expected: %s but got: %s", expected1, diff1)
	}
	if !reflect.DeepEqual(expected2, diff2) {
		t.Errorf("Expected: %s but got: %s", expected2, diff2)
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

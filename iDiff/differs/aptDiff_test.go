package differs

import (
	"reflect"
	"testing"
)

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

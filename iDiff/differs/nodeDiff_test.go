package differs

import (
	"reflect"
	"testing"
)

func TestReadPackageJSON(t *testing.T) {
	testCases := []struct {
		descrip  string
		path     string
		expected nodePackage
		err      bool
	}{
		{
			descrip: "Error on non-existent file",
			path:    "testDirs/nodeTestFiles/not_real.json",
			err:     true,
		},
		{
			descrip:  "Parse JSON with exact fields",
			path:     "testDirs/nodeTestFiles/exact.json",
			expected: nodePackage{"La-croix", "Lime"},
		},
		{
			descrip:  "Parse JSON with additional fields",
			path:     "testDirs/nodeTestFiles/extra.json",
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

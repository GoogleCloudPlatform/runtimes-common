package utils

import (
	"reflect"
	"sort"
	"testing"
)

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

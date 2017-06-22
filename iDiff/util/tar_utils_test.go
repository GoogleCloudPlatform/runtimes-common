package tarUtil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestUnTar(t *testing.T) {
	testCases := []struct {
		descrip  string
		tarPath  string
		target   string
		expected string
		err      bool
	}{
		{
			descrip:  "Tar with files",
			tarPath:  "testTars/la-croix1.tar",
			target:   "testTars/la-croix1",
			expected: "testTars/la-croix1-actual",
		},
		{
			descrip:  "Tar with folders with files",
			tarPath:  "testTars/la-croix2.tar",
			target:   "testTars/la-croix2",
			expected: "testTars/la-croix2-actual",
		},
		{
			descrip:  "Tar with folders with files and a tar file",
			tarPath:  "testTars/la-croix3.tar",
			target:   "testTars/la-croix3",
			expected: "testTars/la-croix3-actual",
		},
	}
	for _, test := range testCases {
		err := UnTar(test.tarPath, test.target)
		if err != nil && !test.err {
			t.Errorf("Got unexpected error: %s", err)
		}
		if err == nil && test.err {
			t.Errorf("Expected error but got none: %s", err)
		}
		if !dirEquals(test.expected, test.target) || !dirEquals(test.target, test.expected) {
			t.Errorf("Directory created not correct structure.")
		}
		os.RemoveAll(test.target)
	}

}

func TestIsTar(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{input: "/testTar/la-croix1.tar", expected: true},
		{input: "/testTar/la-croix1-actual", expected: false},
	}
	for _, test := range testCases {
		actual := isTar(test.input)
		if test.expected != actual {
			t.Errorf("Expected: %t but got: %t", test.expected, actual)
		}
	}
}

func TestExtractTar(t *testing.T) {
	tarPath := "testTars/la-croix3.tar"
	target := "testTars/la-croix3"
	expected := "testTars/la-croix3-full"
	err := ExtractTar(tarPath)
	if err != nil {
		t.Errorf("Got unexpected error: %s", err)
	}
	if !dirEquals(expected, target) || !dirEquals(target, expected) {
		t.Errorf("Directory created not correct structure.")
	}
	os.RemoveAll(target)

}

func dirEquals(actual string, path string) bool {
	files1, _ := ioutil.ReadDir(actual)

	for _, file := range files1 {
		newActualPath := filepath.Join(actual, file.Name())
		newExpectedPath := filepath.Join(path, file.Name())
		fstat, ok := os.Stat(newExpectedPath)
		if ok != nil {
			return false
		}

		if file.IsDir() && !dirEquals(newActualPath, newExpectedPath) {
			return false
		}

		if fstat.Name() != file.Name() {
			return false
		}
		if fstat.Size() != file.Size() {
			return false
		}
		if filepath.Ext(file.Name()) == ".tar" {
			continue
		}

		content1, _ := ioutil.ReadFile(newActualPath)
		content2, _ := ioutil.ReadFile(newExpectedPath)

		if 0 != bytes.Compare(content1, content2) {
			return false
		}
	}
	return true
}

func TestDirToJSON(t *testing.T) {
	path := "testTars/la-croix3-full"
	target := "testTars/la-croix3-full.json"
	expected := "testTars/la-croix3-actual.json"
	err := DirToJSON(path, target)
	if err != nil {
		t.Errorf("Error converting struture to JSON")
	}

	var actualJSON Dir
	var expectedJSON Dir
	content1, _ := ioutil.ReadFile(target)
	content2, _ := ioutil.ReadFile(expected)

	json.Unmarshal(content1, &actualJSON)
	json.Unmarshal(content2, &expectedJSON)

	if !reflect.DeepEqual(actualJSON, expectedJSON) {
		t.Errorf("JSON was incorrect")
	}
	os.Remove(target)
}

package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"
)

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ", ")
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

var schemaVersions map[string]interface{} = map[string]interface{}{
	"1.0.0": new(StructureTestv0),
	"2.0.0": new(StructureTestv1),
}

type SchemaVersion struct {
	SchemaVersion string
}

type Unmarshaller func([]byte, interface{}) error

func combineTests(tests *StructureTest, tmpTests *StructureTest) {
	tests.CommandTests = append(tests.CommandTests, tmpTests.CommandTests...)
	tests.FileExistenceTests = append(tests.FileExistenceTests, tmpTests.FileExistenceTests...)
	tests.FileContentTests = append(tests.FileContentTests, tmpTests.FileContentTests...)
}

func parseFile(tests *StructureTest, configFile string) error {
	var tmpTests StructureTest
	testContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var unmarshal Unmarshaller
	var versionHolder SchemaVersion

	switch {
	case strings.HasSuffix(fp, ".json"):
		unmarshal = json.Unmarshal
	case strings.HasSuffix(fp, ".yaml"):
		unmarshal = yaml.Unmarshal
	default:
		return nil, errors.New("Please provide valid JSON or YAML config file.")
	}

	if err := unmarshal(testContents, &versionHolder); err != nil {
		return nil, err
	}

	version := versionHolder.SchemaVersion
	if version == "" {
		return nil, errors.New("Please provide JSON schema version.")
	} else {
		st := schemaVersions[version]
		if st == nil {
			return nil, errors.New("Unsupported schema version: " + version)
		}
		unmarshal(testContents, &st)
		tests, ok := st.(StructureTest) //type assertion
		if !ok {
			return nil, errors.New("Error encountered when type casting Structure Test interface!")
		}
		combineTests(tests, &tmpTests)
	}
	return nil
}

func Parse(configFiles []string, tests *StructureTest) error {
	for _, file := range configFiles {
		if err := parseFile(tests, file); err != nil {
			return err
		}
	}
}

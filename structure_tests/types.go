package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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
}

type SchemaVersion struct {
	SchemaVersion string
}

type Unmarshaller func([]byte, interface{}) error

func Parse(fp string) (StructureTest, error) {
	testContents, err := ioutil.ReadFile(fp)
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
		unmarshal(testContents, st)
		tests, ok := st.(StructureTest) //type assertion
		if !ok {
			return nil, errors.New("Error encountered when type casting Structure Test interface!")
		}
		return tests, nil
	}
}

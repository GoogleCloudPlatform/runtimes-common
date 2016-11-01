package main

import (
	"strings"
	"testing"
)

type StructureTest interface {
	RunAll(t *testing.T)
}

type CommandTest interface {
	GetName() string
	GetCommand() string
	GetFlags() []string
	GetExpectedOutput() []string
	GetExcludedOutput() []string
	GetExpectedError() []string
	GetExcludedError() []string // excluded error from running command
}

type FileExistenceTest interface {
	GetName() string      // name of test
	GetPath() string      // file to check existence of
	GetIsDirectory() bool // whether or not the path points to a directory
	GetShouldExist() bool // whether or not the file should exist
}

type FileContentTest interface {
	GetName() string               // name of test
	GetPath() string               // file to check existence of
	GetExpectedContents() []string // list of expected contents of file
	GetExcludedContents() []string // list of excluded contents of file
}

func validateCommandTest(t *testing.T, tt CommandTest) {
	if tt.GetName() == "" {
		t.Fatalf("Please provide a valid name for every test!")
	}
	if tt.GetCommand() == "" {
		t.Fatalf("Please provide a valid command to run for test %s", tt.GetName())
	}
	t.Logf("COMMAND TEST: %s", tt.GetName())
}

func validateFileExistenceTest(t *testing.T, tt FileExistenceTest) {
	if tt.GetName() == "" {
		t.Fatalf("Please provide a valid name for every test!")
	}
	if tt.GetPath() == "" {
		t.Fatalf("Please provide a valid file path for test %s", tt.GetName())
	}
	t.Logf("FILE EXISTENCE TEST: %s", tt.GetName())
}

func validateFileContentTest(t *testing.T, tt FileContentTest) {
	if tt.GetName() == "" {
		t.Fatalf("Please provide a valid name for every test!")
	}
	if tt.GetPath() == "" {
		t.Fatalf("Please provide a valid file path for test %s", tt.GetName())
	}
	t.Logf("FILE CONTENT TEST: %s", tt.GetName())
}

/*****************************************************************************/
/*                           MISC TYPE DEFINITIONS                           */
/*****************************************************************************/

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

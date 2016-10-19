package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"testing"
)

func TestRunCommand(t *testing.T) {
	for _, tt := range tests.GetCommandTests() {
		validateCommandTest(t, tt)
		var cmd *exec.Cmd
		if tt.GetFlags() != nil && len(tt.GetFlags()) > 0 {
			cmd = exec.Command(tt.GetCommand(), tt.GetFlags()...)
		} else {
			cmd = exec.Command(tt.GetCommand())
		}
		t.Logf("Executing: %s", cmd.Args)
		var outbuf, errbuf bytes.Buffer

		cmd.Stdout = &outbuf
		cmd.Stderr = &errbuf

		if err := cmd.Run(); err != nil {
			// The test might be designed to run a command that exits with an error.
			t.Logf("Error running command: %s. Continuing.", err)
		}

		stdout := outbuf.String()
		if stdout != "" {
			t.Logf("stdout: %s", stdout)
		}
		stderr := errbuf.String()
		if stderr != "" {
			t.Logf("stderr: %s", stderr)
		}

		for _, errStr := range tt.GetExpectedError() {
			errMsg := fmt.Sprintf("Expected string '%s' not found in error!", errStr)
			compileAndRunRegex(errStr, stderr, t, errMsg, true)
		}
		for _, errStr := range tt.GetExcludedError() {
			errMsg := fmt.Sprintf("Excluded string '%s' found in error!", errStr)
			compileAndRunRegex(errStr, stderr, t, errMsg, false)
		}

		for _, outStr := range tt.GetExpectedOutput() {
			errMsg := fmt.Sprintf("Expected string '%s' not found in output!", outStr)
			compileAndRunRegex(outStr, stdout, t, errMsg, true)
		}
		for _, outStr := range tt.GetExcludedError() {
			errMsg := fmt.Sprintf("Excluded string '%s' found in output!", outStr)
			compileAndRunRegex(outStr, stdout, t, errMsg, false)
		}
	}
}

func TestFileExists(t *testing.T) {
	for _, tt := range tests.GetFileExistenceTests() {
		validateFileExistenceTest(t, tt)
		var err error
		if tt.GetIsDirectory() {
			_, err = ioutil.ReadDir(tt.GetPath())
		} else {
			_, err = ioutil.ReadFile(tt.GetPath())
		}
		if tt.GetShouldExist() && err != nil {
			t.Errorf("File %s should exist but does not!", tt.GetPath())
		} else if !tt.GetShouldExist() && err == nil {
			t.Errorf("File %s should not exist but does!", tt.GetPath())
		}
	}
}

func TestFileContents(t *testing.T) {
	for _, tt := range tests.GetFileContentTests() {
		validateFileContentTest(t, tt)
		actualContents, err := ioutil.ReadFile(tt.GetPath())
		if err != nil {
			t.Errorf("Failed to open %s. Error: %s", tt.GetPath(), err)
		}

		contents := string(actualContents[:])

		var errMessage string
		for _, s := range tt.GetExpectedContents() {
			errMessage = "Expected string " + s + " not found in file contents!"
			compileAndRunRegex(s, contents, t, errMessage, true)
		}
		for _, s := range tt.GetExcludedContents() {
			errMessage = "Excluded string " + s + " found in file contents!"
			compileAndRunRegex(s, contents, t, errMessage, false)
		}
	}
}

func compileAndRunRegex(regex string, base string, t *testing.T, err string, shouldMatch bool) {
	r, rErr := regexp.Compile(regex)
	if rErr != nil {
		t.Errorf("Error compiling regex %s : %s", regex, rErr.Error())
		return
	}
	if shouldMatch != r.MatchString(base) {
		t.Errorf(err)
	}
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

var configFiles arrayFlags
var tests StructureTest

func TestMain(m *testing.M) {
	flag.Var(&configFiles, "config", "path to the .yaml file containing test definitions.")
	flag.Parse()

	if len(configFiles) == 0 {
		configFiles = append(configFiles, "/workspace/structure_test.json")
	}

	var err error
	for _, file := range configFiles {
		if tests, err = Parse(file); err != nil {
			log.Fatalf("Error parsing config file: %s", err)
		}
		log.Printf("Running tests for file %s", file)
		if exit := m.Run(); exit != 0 {
			os.Exit(exit)
		}
	}
}

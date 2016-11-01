package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"
)

type StructureTestv0 struct {
	CommandTests       []CommandTestv0
	FileExistenceTests []FileExistenceTestv0
	FileContentTests   []FileContentTestv0
}

func (st StructureTestv0) RunAll(t *testing.T) {
	st.RunCommandTests(t)
	st.RunFileExistenceTests(t)
	st.RunFileContentTests(t)
}

func (st StructureTestv0) RunCommandTests(t *testing.T) {
	for _, tt := range st.CommandTests {
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

func (st StructureTestv0) RunFileExistenceTests(t *testing.T) {
	for _, tt := range st.FileExistenceTests {
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

func (st StructureTestv0) RunFileContentTests(t *testing.T) {
	for _, tt := range st.FileContentTests {
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

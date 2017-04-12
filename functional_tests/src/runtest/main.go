// Copyright 2017 Google Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"syscall"
)

const FAILED = "FAILED"
const PASSED = "PASSED"

var verbose = flag.Bool("verbose", false, "Verbose logging")

func main() {
	testSpec := flag.String("test_spec", "", "Path to a yaml or json file containing the test spec")
	flag.Parse()

	if *testSpec == "" {
		log.Fatal("--test_spec must be specified")
	}

	suite := LoadSuite(*testSpec)
	doSetup(suite)
	report := doRunTests(suite)
	defer report()
	doTeardown(suite)
}

func info(text string, arg ...interface{}) {
	log.Printf(text, arg...)
}

func runCommand(name string, args ...string) (err error, stdout string, stderr string) {
	if *verbose {
		info("Running command: %v", append([]string{name}, args...))
	}
	cmd := exec.Command(name, args...)
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	cmd.Stdout = &stdoutBuffer
	cmd.Stderr = &stderrBuffer
	err = cmd.Run()
	stdout = stdoutBuffer.String()
	stderr = stderrBuffer.String()
	if *verbose {
		commandOutput("STDOUT", stdout)
		commandOutput("STDERR", stderr)
	}
	return
}

func commandOutput(name string, content string) {
	if len(content) <= 0 {
		return
	}

	info("%s>>%s<<%s", name, content, name)
}

func doSetup(suite Suite) {
	info(">>> Setting up...")
	for _, setup := range suite.Setup {
		err, _, _ := runCommand(setup.Command[:1][0], setup.Command[1:]...)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func doTeardown(suite Suite) {
	info(">>> Tearing down...")
	for _, teardown := range suite.Teardown {
		err, _, _ := runCommand(teardown.Command[:1][0], teardown.Command[1:]...)
		if err != nil {
			info(" > Warning: Teardown command failed: %s", err)
		}
	}
}

func doRunTests(suite Suite) func() {
	info(">>> Testing...")
	results := make(map[int]string)
	passing := make(map[int]bool)
	for index, test := range suite.Tests {
		doOneTest(index, test, suite, results, passing)
	}

	report := func() {
		allPassing := true
		for index := range suite.Tests {
			allPassing = allPassing && passing[index]
		}
		if allPassing {
			info(">>> Summary: %s", PASSED)
		} else {
			info(">>> Summary: %s", FAILED)
		}
		for index := range suite.Tests {
			info(" > %s", results[index])
		}
	}
	return report
}

func doOneTest(index int, test Test, suite Suite, results map[int]string, passing map[int]bool) {
	var name string
	var msg string
	var result string

	if len(test.Name) > 0 {
		name = fmt.Sprintf("test-%d (%s)", index, test.Name)
	} else {
		name = fmt.Sprintf("test-%d", index)
	}
	info(" > %s", name)

	recordResult := func(b *string) {
		if *b == PASSED {
			results[index] = fmt.Sprintf("%s: %s", name, PASSED)
			passing[index] = true
		} else {
			results[index] = fmt.Sprintf("%s: %s", name, FAILED)
			passing[index] = false
		}
		info(" %s", *b)
	}
	defer recordResult(&result)

	args := append([]string{"exec", suite.Target}, test.Command...)
	err, stdout, stderr := runCommand("docker", args...)

	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				result = fmt.Sprintf("%s: Exit status %d", FAILED, status.ExitStatus())
			}
		} else {
			result = fmt.Sprintf("%s: Encountered error: %v", FAILED, err)
		}
		if len(stdout) > 0 {
			result = fmt.Sprintf("%s\nSTDOUT>>>%s<<<STDOUT", result, stdout)
		}
		if len(stderr) > 0 {
			result = fmt.Sprintf("%s\nSTDERR>>>%s<<<STDERR", result, stderr)
		}
		return
	}

	msg = DoStringAssert(stdout, test.Expect.Stdout)
	if len(msg) > 0 {
		result = fmt.Sprintf("%s: stdout assertion failure\n%s", FAILED, msg)
		return
	}
	msg = DoStringAssert(stderr, test.Expect.Stderr)
	if len(msg) > 0 {
		result = fmt.Sprintf("%s: stderr assertion failure\n%s", FAILED, msg)
		return
	}

	result = PASSED
}

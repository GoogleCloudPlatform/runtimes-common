// Copyright 2016 Google Inc. All rights reserved.

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
	"fmt"
	"log"
	"os"
	"testing"
)

const (
	RED    = "\033[0;31m"
	GREEN  = "\033[0;32m"
	YELLOW = "\033[1;33m"
	CYAN   = "\033[0;36m"
	BLUE   = "\033[0;34m"
	PURPLE = "\033[0;35m"
	NC     = "\033[0m" // No Color

	LOG_TEMPLATE     = "LOG: %s"
	INFO_TEMPLATE    = YELLOW + "INFO: %s" + NC
	HEADER_TEMPLATE  = GREEN + "%s" + NC
	SPECIAL_TEMPLATE = BLUE + "%s" + NC
	ERROR_TEMPLATE   = "\n" + RED + "ERROR: %s" + NC + "\n"
	FATAL_TEMPLATE   = "\n" + PURPLE + "FATAL: %s" + NC + "\n"
)

// ANSI Color Escape Codes
// Black        0;30     Dark Gray     1;30
// Red          0;31     Light Red     1;31
// Green        0;32     Light Green   1;32
// Brown/Orange 0;33     Yellow        1;33
// Blue         0;34     Light Blue    1;34
// Purple       0;35     Light Purple  1;35
// Cyan         0;36     Light Cyan    1;36
// Light Gray   0;37     White         1;37

var Log = log.New(os.Stdout,
	"",
	log.Ldate|log.Ltime)

var Info = log.New(os.Stdout,
	"", 0)

var Error = log.New(os.Stderr,
	"",
	log.Ldate|log.Ltime)

func _Log(text string, args ...interface{}) {
	Log.Println(fmt.Sprintf(text, args...))
}

func _Info(text string, args ...interface{}) {
	Info.Println(fmt.Sprintf(INFO_TEMPLATE, fmt.Sprintf(text, args...)))
}

func _Special(text string, args ...interface{}) {
	Info.Println(fmt.Sprintf(SPECIAL_TEMPLATE, fmt.Sprintf(text, args...)))
}

func _Header(text string, args ...interface{}) string {
	// returns formatted string to be passed to t.Run()
	return fmt.Sprintf(HEADER_TEMPLATE, fmt.Sprintf(text, args...))
}

func _Error(t *testing.T, text string, args ...interface{}) {
	Error.Println(fmt.Sprintf(ERROR_TEMPLATE, fmt.Sprintf(text, args...)))
	t.Fail()
}

func _Fatal(t *testing.T, text string, args ...interface{}) {
	Error.Println(fmt.Sprintf(FATAL_TEMPLATE, fmt.Sprintf(text, args...)))
	t.FailNow()
}

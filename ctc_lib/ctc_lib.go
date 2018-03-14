/*
Copyright 2018 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ctc_lib

// This file declares all the package level globals

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

var exitOnError = true
var Version string
var emptyTemplate = "{{.}}"
var defaultLogLevel = "info"
var Log *log.Logger

func SetExitOnError(value bool) {
	exitOnError = value
}

func GetExitOnError() bool {
	return exitOnError
}

func CommandExit(err error) {
	if err != nil && exitOnError {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Error(err)
}

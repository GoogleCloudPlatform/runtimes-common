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
	"bytes"
	"os"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/config"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/constants"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/flags"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/logging"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var exitOnError = true
var Version string
var ConfigFile string
var ConfigType = constants.ConfigType

var UpdateCheck bool

var Log *log.Logger

var toolName string

func SetExitOnError(value bool) {
	exitOnError = value
}

func GetExitOnError() bool {
	return exitOnError
}

func CommandExit(err error) {
	if err != nil {
		logging.Out.Error(err)
		if exitOnError {
			os.Exit(1)
		}
	}
}

func readDefaultConfig() {
	viper.SetConfigType(config.DefaultConfigType)
	viper.ReadConfig(bytes.NewBuffer(config.DefaultConfig))
}

func initConfig() {
	if ConfigFile == "" {
		logging.Out.Debugf(`No Config provided. Using Tools Defaults.
You can override it via ctc_lib.ConfigFile pkg variable`)
		readDefaultConfig()
		return
	}

	viper.SetConfigFile(ConfigFile)
	viper.SetConfigType(ConfigType)

	err := viper.ReadInConfig()
	if err != nil {
		logging.Out.Warningf("Error reading config file at %s: %s. Using Defaults", ConfigFile, err)
		readDefaultConfig()
	}
}

func initLogging() {
	Log = logging.NewLogger(viper.GetString("logDir"), toolName,
		flags.Verbosity.Level, flags.EnableColors)
	Log.SetLevel(flags.Verbosity.Level)
	Log.AddHook(logging.NewFatalHook(exitOnError))
}

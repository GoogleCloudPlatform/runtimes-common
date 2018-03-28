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

package logging

import (
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func NewLogger(path string, level log.Level, enableColors bool) *log.Logger {
	if level == log.DebugLevel {
		// Log to File when verbosity=debug
		writer, _ := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
			rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
		)
		return &log.Logger{
			Out:       writer,
			Formatter: new(log.JSONFormatter),
			Hooks:     make(log.LevelHooks),
			Level:     log.DebugLevel,
		}
	}
	return &log.Logger{
		Out:       os.Stderr,
		Formatter: NewCTCLogFormatter(enableColors),
		Hooks:     make(log.LevelHooks),
		Level:     level,
	}
}

// Define Explicit StdOut Loggers which can be used to always print to StdOut.
// This can also add other functionalilty like colored output etc.
var Out = log.New()

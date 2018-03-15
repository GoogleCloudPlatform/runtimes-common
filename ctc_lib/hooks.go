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

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

// Hook which exits when Log.Panic and Log.Fatal is Called
type FatalHook struct {
	writer io.Writer
}

func NewFatalHook(writer io.Writer) *FatalHook {
	return &FatalHook{
		writer: writer,
	}
}

func (hook *FatalHook) Fire(entry *log.Entry) error {
	switch entry.Level {
	case log.PanicLevel:
		CommandExit(fmt.Errorf(entry.Message))
	case log.FatalLevel:
		CommandExit(fmt.Errorf(entry.Message))
		log.Warn("Avoid using Log.Fatal. Consider Log.Panic instead to exit gracefully")
	}
	return nil
}

func (hook *FatalHook) Levels() []log.Level {
	return []log.Level{log.PanicLevel, log.FatalLevel}
}

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
	"bytes"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 36
	gray    = 37
)

// This Files Defines a LogFormatter With Colors. This Format used is
//<Level>: <Message>

// Implements Interface Logrus.Formatter
// https://github.com/sirupsen/logrus/blob/master/formatter.go
type CTCLogFormatter struct {
	EnableColors bool
}

func NewCTCLogFormatter(enableColors bool) *CTCLogFormatter {
	return &CTCLogFormatter{
		EnableColors: enableColors,
	}
}

func (f *CTCLogFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	levelColor := f.getColor(entry)
	levelText := strings.ToUpper(entry.Level.String())
	if entry.Message != "" {
		fmt.Fprintf(b, "\x1b[%dm%s:\x1b[0m %-44s ", levelColor, levelText, entry.Message)
	}
	for _, k := range keys {
		v := entry.Data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *CTCLogFormatter) getColor(entry *log.Entry) int {
	if !f.EnableColors {
		return nocolor
	}
	switch entry.Level {
	case log.DebugLevel:
		return gray
	case log.WarnLevel:
		return yellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		return red
	default:
		return green
	}
}

func (f *CTCLogFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *CTCLogFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}
	b.WriteString(stringVal)
}

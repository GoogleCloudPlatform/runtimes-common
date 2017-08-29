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
	"strings"
	"testing"
)

type StructureTest interface {
	RunAll(t *testing.T) int
}

type arrayFlags []string

func (a *arrayFlags) String() string {
	return strings.Join(*a, ", ")
}

func (a *arrayFlags) Set(value string) error {
	*a = append(*a, value)
	return nil
}

var schemaVersions map[string]VersionHolder = map[string]VersionHolder{
	"1.0.0": new(VersionHolderv000),
	"1.1.0": new(VersionHolderv100),
}

type SchemaVersion struct {
	SchemaVersion string
}

type Unmarshaller func([]byte, interface{}) error

type VersionHolder interface {
	New() StructureTest
}

type VersionHolderv000 struct{}
type VersionHolderv100 struct{}

func (v VersionHolderv000) New() StructureTest {
	return new(StructureTestv000)
}

func (v VersionHolderv100) New() StructureTest {
	return new(StructureTestv100)
}

type EnvVar struct {
	Key   string
	Value string
}

type Command []string

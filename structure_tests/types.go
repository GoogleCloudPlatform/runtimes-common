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

	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/types/v1"
	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/types/v2"
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

var schemaVersions map[string]func() StructureTest = map[string]func() StructureTest{
	"1.0.0": func() StructureTest { return new(v1.StructureTest) },
	"2.0.0": func() StructureTest { return new(v2.StructureTest) },
}

type SchemaVersion struct {
	SchemaVersion string
}

type Unmarshaller func([]byte, interface{}) error

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

package testutil

import (
	"strings"
)

var TestTUFConfig = `
gcsProjectID: testProjectID
kmsProjectID: testKmsProjectID
kmlLocation: global
kmsKeyringID: testKeyRing
cryptoKey: testKey
gcsBucketID: testBucket
`

func IsErrorEqualOrContains(err error, subErr error) bool {
	if err == nil && subErr == nil {
		return true // Return true Both of them are nil
	} else if err == nil || subErr == nil {
		return false // Return false if either of them are nil
	} else if strings.Contains(err.Error(), subErr.Error()) {
		return true // Return true if Messages are equal
	}
	return false // Return false
}

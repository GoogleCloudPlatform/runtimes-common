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

package drivers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	pkgutil "github.com/GoogleCloudPlatform/container-diff/pkg/util"
	"github.com/GoogleCloudPlatform/runtimes-common/structure_tests/types/unversioned"
)

type TarDriver struct {
	Image pkgutil.Image
}

// TODO(nkubala): need to call 'pkgutil.CleanupImage(image)' somewhere
func NewTarDriver(imageName string) Driver {
	// container-diff will infer from the source the correct prepper to use

	// ip := pkgutil.DaemonPrepper{
	// 	Source: imageName,
	// }

	ip := pkgutil.ImagePrepper{
		Source: imageName,
	}
	image, err := ip.GetImage()
	if err != nil {
		panic(err)
	}
	return &TarDriver{
		Image: image,
	}
}

func (d *TarDriver) Setup(t *testing.T, envVars []unversioned.EnvVar, fullCommand []unversioned.Command) {
	// this driver is unable to process commands, so this is a noop.
	return
}

func (d *TarDriver) ProcessCommand(t *testing.T, envVars []unversioned.EnvVar, fullCommand []string) (string, string, int) {
	// this driver is unable to process commands, so this a noop.
	return "", "", 0
}

func (d *TarDriver) StatFile(t *testing.T, path string) (os.FileInfo, error) {
	return os.Stat(filepath.Join(d.Image.FSPath, path))
}

func (d *TarDriver) ReadFile(t *testing.T, path string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(d.Image.FSPath, path))
}

func (d *TarDriver) ReadDir(t *testing.T, path string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(filepath.Join(d.Image.FSPath, path))
}

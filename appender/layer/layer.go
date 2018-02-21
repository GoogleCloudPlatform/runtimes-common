/*
Copyright 2018 Google LLC
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

package layer

import (
	"github.com/containers/image/copy"
	"github.com/containers/image/signature"

	"github.com/GoogleCloudPlatform/container-diff/pkg/image"

	"github.com/containers/image/docker"
)

func AppendLayer(base, dest string, layer []byte) error {
	baseRef, err := docker.ParseReference("//" + base)
	if err != nil {
		return err
	}

	ms, err := image.NewMutableSource(baseRef)
	if err != nil {
		return err
	}

	if err := ms.AppendLayer(layer); err != nil {
		return err
	}

	dstRef, err := docker.ParseReference("//" + dest)
	if err != nil {
		return err
	}

	proxyRef := &image.ProxyReference{
		Src: ms,
	}

	pc, err := signature.NewPolicyContext(&signature.Policy{
		Default: signature.PolicyRequirements{signature.NewPRInsecureAcceptAnything()},
	})
	if err != nil {
		return err
	}

	return copy.Image(pc, dstRef, proxyRef, nil)
}

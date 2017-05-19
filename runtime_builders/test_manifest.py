# Copyright 2017 Google Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import pytest

import builder_util
import yaml


def test_simple_manifest():
    _load_build_verify('test_manifests/simple.yaml')


def test_circular_manifest():
    with pytest.raises(SystemExit):
        _load_build_verify('test_manifests/circular.yaml')


def test_broken_manifest():
    with pytest.raises(SystemExit):
        _load_build_verify('test_manifests/broken_link.yaml')


def _load_build_verify(manifest_file):
    with open(manifest_file) as f:
        manifest = yaml.load(f)
    graph = builder_util._build_manifest_graph(manifest)

    builder_util._verify_manifest_graph(graph)

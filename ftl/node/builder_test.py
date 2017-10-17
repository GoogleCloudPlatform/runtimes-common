# Copyright 2017 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import unittest
import tempfile

from containerregistry.client.v2_2 import docker_image

from ftl.common import context
from ftl.common import test_util
from ftl.node import builder


class TarDockerImage():

    def __init__(self, config_path, tarball_path):
        self._config = open(config_path, 'r').read()
        # TODO(aaron-prindle) use fast image format instead of tarball
        self._docker_image = docker_image.FromDisk(self._config,
                                                   zip([], []), tarball_path)

    def GetConfig(self):
        return self._config

    def GetDockerImage(self):
        return self._docker_image


class BuilderTestCase():
    def __init__(self, builder_fxn, ctx, cash, base_image):
        self._ctx = ctx
        self._builder = builder_fxn(ctx)
        self._cash = cash
        self._base_image = base_image

    def CreatePackageBase(self):
        self._base_image.GetDockerImage().__enter__()
        return self._builder.CreatePackageBase(
            self._base_image.GetDockerImage(), self._cash)

    def GetCacheEntries(self):
        return len(self._cash._registry._registry)

    def GetCacheMap(self):
        return self._cash._registry._registry

    def GetCacheEntryByStringKey(self, key):
        # cast key to docker_tag
        return self._cash._registry.getImage(key)


class NodeTest(unittest.TestCase):
    def setup(self, builder, ctx, cash):
        pass

    def test_create_package_base(self):
        current_dir = os.path.dirname(__file__)
        test_case = BuilderTestCase(
            builder.Node,
            context.Workspace(os.path.join(current_dir, "testdata/node_app")),
            test_util.MockHybridRegistry('fake.gcr.io/google-appengine',
                                         tempfile.mkdtemp()),
            TarDockerImage(
                os.path.join(current_dir, "testdata/base_image/config_file"),
                os.path.join(
                    current_dir,
                    "testdata/base_image/distroless-nodejs-latest.tar.gz")), )
        test_case.CreatePackageBase()
        # check that image was added to the cache
        self.assertEqual(1, test_case.GetCacheEntries())
        for k in test_case.GetCacheMap():
            self.assertEqual(
                str(k).startswith("fake.gcr.io/google-appengine/node-package"),
                True)
        test_case.CreatePackageBase()
        # check that image was added to the cache
        self.assertEqual(1, len(test_case.GetCacheMap()))
        for k in test_case.GetCacheMap():
            self.assertEqual(
                str(k).startswith("fake.gcr.io/google-appengine/node-package"),
                True)


if __name__ == '__main__':
    unittest.main()

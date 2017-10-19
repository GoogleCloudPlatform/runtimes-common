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

import json
import os
import unittest
import shutil
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

    @classmethod
    def setUpClass(cls):
        current_dir = os.path.dirname(__file__)
        cls.base_image = TarDockerImage(
            os.path.join(current_dir, "testdata/base_image/config_file"),
            os.path.join(
                current_dir,
                "testdata/base_image/distroless-nodejs-latest.tar.gz"))
        cls.ctx = context.Workspace(
            os.path.join(current_dir, "testdata/node_app"))

    def setUp(self):
        self._tmpdir = tempfile.mkdtemp()
        self.cache = test_util.MockHybridRegistry(
            'fake.gcr.io/google-appengine',
            self._tmpdir)

    def tearDown(self):
        shutil.rmtree(self._tmpdir)

    def test_create_package_base(self):
        test_case = BuilderTestCase(
            builder.Node,
            self.ctx,
            self.cache,
            self.base_image)

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

    def test_overrides(self):
        b = builder.Node(self.ctx)
        app_base = b.CreatePackageBase(self.base_image.GetDockerImage(),
                                       self.cache)
        cfg = json.loads(app_base.config_file())
        self.assertEqual(cfg['config']['Entrypoint'],
                         ['sh', '-c', "'node server.js'"])


class ParseEntrypointTest(unittest.TestCase):

    def add_shell(self, args):
        return ['sh', '-c', "'%s'" % args]

    def test_no_scripts(self):
        self.assertEqual(
            builder.parse_entrypoint({}),
            self.add_shell('node server.js'))

    def test_prestart(self):
        self.assertEqual(
            builder.parse_entrypoint(
                {'scripts': {'prestart': 'foo'}}),
            self.add_shell('foo && node server.js')
        )

    def test_start(self):
        self.assertEqual(
            builder.parse_entrypoint(
                {'scripts': {'start': 'foo'}}),
            self.add_shell('foo')
        )

    def test_both(self):
        self.assertEqual(
            builder.parse_entrypoint(
                {'scripts': {'prestart': 'foo', 'start': 'baz'}}),
            self.add_shell('foo && baz')
        )


if __name__ == '__main__':
    unittest.main()

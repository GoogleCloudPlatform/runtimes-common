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
import mock
import os
import unittest
import shutil
import tempfile

from containerregistry.client.v2_2 import docker_image

from ftl.common import context
from ftl.common import test_util
from ftl.node import builder

_PACKAGE_JSON = json.loads("""
{
  "name": "hello-world-express",
  "description": "Hello World test app",
  "version": "0.0.1",
  "private": true,
  "dependencies": {
    "express": "3.x"
  }
}
""")
_PACKAGE_JSON_TEXT = json.dumps(_PACKAGE_JSON)

_APP = """
var express = require('express');
var app = express();

// Routes
app.get('/', function(req, res) {
  res.send('Hello World!');
});

// Listen
var port = process.env.PORT || 3000;
app.listen(port);
console.log('Listening on localhost:'+ port);
"""


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

        # Mock out the calls to NPM for speed.
        self._builder._gen_package_tar = mock.Mock()
        self._builder._gen_package_tar.return_value = ('layer', 'sha')

        self._cash = cash
        self._base_image = base_image

    def CreatePackageBase(self):
        with self._base_image.GetDockerImage():
            return self._builder.CreatePackageBase(
                self._base_image.GetDockerImage(),
                self._cash)

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

    def setUp(self):
        self._tmpdir = tempfile.mkdtemp()
        self.cache = test_util.MockHybridRegistry(
            'fake.gcr.io/google-appengine',
            self._tmpdir)
        self.ctx = context.Memory()
        self.ctx.AddFile("app.js", _APP)
        self.test_case = BuilderTestCase(
                builder.Node,
                self.ctx,
                self.cache,
                self.base_image)

    def tearDown(self):
        shutil.rmtree(self._tmpdir)

    def test_create_package_base_cache(self):
        self.ctx.AddFile('package.json', _PACKAGE_JSON_TEXT)

        self.test_case.CreatePackageBase()
        # check that image was added to the cache
        self.assertEqual(1, self.test_case.GetCacheEntries())
        for k in self.test_case.GetCacheMap():
            self.assertEqual(
                str(k).startswith("fake.gcr.io/google-appengine/node-package"),
                True)

        self.test_case.CreatePackageBase()
        # check that image was added to the cache
        self.assertEqual(1, len(self.test_case.GetCacheMap()))
        for k in self.test_case.GetCacheMap():
            self.assertEqual(
                str(k).startswith("fake.gcr.io/google-appengine/node-package"),
                True)

    def test_create_package_base_entrypoint(self):
        pj = _PACKAGE_JSON.copy()
        pj['scripts'] = {
            'start': 'foo bar'
        }
        self.ctx.AddFile('package.json', json.dumps(pj))

        base = self.test_case.CreatePackageBase()
        self.assertEqual(_entrypoint(base), ['sh', '-c', 'foo bar'])

    def test_create_package_base_no_entrypoint(self):
        self.ctx.AddFile('package.json', _PACKAGE_JSON_TEXT)

        base = self.test_case.CreatePackageBase()
        self.assertEqual(_entrypoint(base), ['sh', '-c', 'node server.js'])

    def test_create_package_base_prestart(self):
        pj = _PACKAGE_JSON.copy()
        pj['scripts'] = {
            'prestart': 'foo bar',
            'start': 'baz'
        }
        self.ctx.AddFile('package.json', json.dumps(pj))

        base = self.test_case.CreatePackageBase()
        self.assertEqual(_entrypoint(base), ['sh', '-c', 'foo bar && baz'])

    def test_create_package_base_no_descriptor(self):
        self.assertFalse(self.ctx.Contains('package.json'))
        self.assertFalse(self.ctx.Contains('package-lock.json'))
        base = self.test_case.CreatePackageBase()
        self.assertEqual(_entrypoint(base), ['sh', '-c', 'node server.js'])


def _entrypoint(image):
    cfg = json.loads(image.config_file())
    return cfg['config']['Entrypoint']


if __name__ == '__main__':
    unittest.main()

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


class NodeTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        current_dir = os.path.dirname(__file__)
        cls.base_image = test_util.TarDockerImage(
            os.path.join(current_dir, "testdata/base_image/config_file"),
            os.path.join(
                current_dir,
                "testdata/base_image/distroless-nodejs-latest.tar.gz"))

    def setUp(self):
        self._tmpdir = tempfile.mkdtemp()
        self.cache = test_util.MockHybridRegistry(
            'fake.gcr.io/google-appengine', self._tmpdir)
        self.ctx = context.Memory()
        self.ctx.AddFile("app.js", _APP)
        self.test_case = test_util.BuilderTestCase(builder.Node, self.ctx,
                                                   self.cache, self.base_image)

    def tearDown(self):
        shutil.rmtree(self._tmpdir)

    def test_create_package_base_image(self):
        self.assertIsInstance(self.test_case.CreatePackageBase(),
                              docker_image.DockerImage)

    def test_create_package_base_entrypoint(self):
        pj = _PACKAGE_JSON.copy()
        pj['scripts'] = {'start': 'foo bar'}
        self.ctx.AddFile('package.json', json.dumps(pj))

        base = self.test_case.CreatePackageBase()
        self.assertEqual(_entrypoint(base), ['sh', '-c', 'foo bar'])

    def test_create_package_base_no_entrypoint(self):
        self.ctx.AddFile('package.json', _PACKAGE_JSON_TEXT)

        base = self.test_case.CreatePackageBase()
        self.assertEqual(_entrypoint(base), ['sh', '-c', 'node server.js'])

    def test_create_package_base_prestart(self):
        pj = _PACKAGE_JSON.copy()
        pj['scripts'] = {'prestart': 'foo bar', 'start': 'baz'}
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

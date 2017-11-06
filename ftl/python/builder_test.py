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
import datetime

from containerregistry.client.v2_2 import docker_image

from ftl.common import context
from ftl.common import test_util
from ftl.python import builder

_REQUIREMENTS_TXT = """
Flask==0.7.2
"""

_APP = """
import os
from flask import Flask
app = Flask(__name__)


@app.route("/")
def hello():
    return "Hello from Python!"


if __name__ == "__main__":
    port = int(os.environ.get("PORT", 5000))
    app.run(host='0.0.0.0', port=port)
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
                self._base_image.GetDockerImage(), self._cash)

    def GetCacheEntries(self):
        return len(self._cash._registry._registry)

    def GetCacheMap(self):
        return self._cash._registry._registry

    def GetCacheEntryByStringKey(self, key):
        # cast key to docker_tag
        return self._cash._registry.getImage(key)


class PythonTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        current_dir = os.path.dirname(__file__)
        cls.base_image = TarDockerImage(
            os.path.join(current_dir, "testdata/base_image/config_file"),
            os.path.join(
                current_dir,
                "testdata/base_image/distroless-python2.7-latest.tar.gz"))

    def setUp(self):
        self._tmpdir = tempfile.mkdtemp()
        self.cache = test_util.MockHybridRegistry(
            'fake.gcr.io/google-appengine', self._tmpdir)
        self.ctx = context.Memory()
        self.ctx.AddFile("app.py", _APP)
        self.ctx.AddFile('requirements.txt', _REQUIREMENTS_TXT)
        self.test_case = BuilderTestCase(builder.Python, self.ctx, self.cache,
                                         self.base_image)

    def tearDown(self):
        shutil.rmtree(self._tmpdir)

    def test_create_package_base_cache(self):
        self.test_case.CreatePackageBase()
        # check that image was added to the cache
        self.assertEqual(1, self.test_case.GetCacheEntries())
        for k in self.test_case.GetCacheMap():
            self.assertEqual(
                str(k).startswith(
                    "fake.gcr.io/google-appengine/python-requirements-cache"),
                True)

        self.test_case.CreatePackageBase()
        # check that image was added to the cache
        self.assertEqual(1, len(self.test_case.GetCacheMap()))
        for k in self.test_case.GetCacheMap():
            self.assertEqual(
                str(k).startswith(
                    "fake.gcr.io/google-appengine/python-requirements-cache"),
                True)

    def test_create_package_base_ttl_written(self):
        base = self.test_case.CreatePackageBase()
        self.assertNotEqual(_creation_time(base), "1970-01-01T00:00:00Z")
        last_created = _timestamp_to_time(_creation_time(base))
        now = datetime.datetime.now()
        self.assertTrue(last_created > now - datetime.timedelta(days=2))

    # TODO(aaron-prindle) add test to check expired/unexpired logic for TTL


def _creation_time(image):
    cfg = json.loads(image.config_file())
    print cfg
    return cfg['created']


def _timestamp_to_time(dt_str):
    return datetime.datetime.strptime(dt_str, "%Y-%m-%d")


if __name__ == '__main__':
    unittest.main()

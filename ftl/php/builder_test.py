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
from ftl.php import builder


_COMPOSER_JSON = json.loads("""
{
  "name": "hello-world",
  "require": {
    "php": ">=5.5",
    "silex/silex": "^1.3"
  },
  "require-dev": {
    "behat/mink": "^1.7",
    "behat/mink-goutte-driver": "^1.2",
    "phpunit/phpunit": "~4",
    "symfony/browser-kit": "^3.0",
    "symfony/http-kernel": "^3.0",
    "google/cloud-tools": "^0.6"
  }
}
""")

_COMPOSER_JSON_TEXT = json.dumps(_COMPOSER_JSON)

_APP = """
require_once __DIR__ . '/../vendor/autoload.php';
$app = new Silex\Application();

$app->get('/', function () {
    return 'Hello World';
});
$app->get('/goodbye', function () {
    return 'Goodbye World';
});

// @codeCoverageIgnoreStart
if (PHP_SAPI != 'cli') {
    $app->run();
}
// @codeCoverageIgnoreEnd

return $app;
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

        # Mock out the calls to Composer for speed.
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


class PHPTest(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        current_dir = os.path.dirname(__file__)
        cls.base_image = TarDockerImage(
            os.path.join(current_dir, "testdata/base_image/config_file"),
            os.path.join(
                current_dir,
                "testdata/base_image/distroless-base-latest.tar.gz"))

    def setUp(self):
        self._tmpdir = tempfile.mkdtemp()
        self.cache = test_util.MockHybridRegistry(
            'fake.gcr.io/google-appengine',
            self._tmpdir)
        self.ctx = context.Memory()
        self.ctx.AddFile("app.php", _APP)
        self.test_case = BuilderTestCase(
                builder.PHP,
                self.ctx,
                self.cache,
                self.base_image)

    def tearDown(self):
        shutil.rmtree(self._tmpdir)

    def test_create_package_base_cache(self):
        self.ctx.AddFile('composer.json', _COMPOSER_JSON_TEXT)

        self.test_case.CreatePackageBase()
        # check that image was added to the cache
        self.assertEqual(1, self.test_case.GetCacheEntries())
        for k in self.test_case.GetCacheMap():
            self.assertEqual(
                str(k).startswith("fake.gcr.io/google-appengine/php-composer"),
                True)

        self.test_case.CreatePackageBase()
        # check that image was added to the cache
        self.assertEqual(1, len(self.test_case.GetCacheMap()))
        for k in self.test_case.GetCacheMap():
            self.assertEqual(
                str(k).startswith("fake.gcr.io/google-appengine/php-composer"),
                True)


if __name__ == '__main__':
    unittest.main()

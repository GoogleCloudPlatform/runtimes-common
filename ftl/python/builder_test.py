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
import datetime

from containerregistry.client.v2_2 import docker_image

from ftl.common import context
from ftl.common import test_util
from ftl.python import builder

_REQUIREMENTS_TXT = """
Flask==0.12.0
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


class PythonTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        current_dir = os.path.dirname(__file__)
        cls.base_image = test_util.TarDockerImage(
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
        self.test_case = test_util.BuilderTestCase(builder.Python, self.ctx,
                                                   self.cache, self.base_image)

    def tearDown(self):
        shutil.rmtree(self._tmpdir)

    def test_create_package_base_image(self):
        # check that image was added to the cache
        self.assertIsInstance(self.test_case.CreatePackageBase(),
                              docker_image.DockerImage)

    def test_create_package_base_ttl_written(self):
        base = self.test_case.CreatePackageBase()
        self.assertNotEqual(_creation_time(base), "1970-01-01T00:00:00Z")
        last_created = _timestamp_to_time(_creation_time(base))
        now = datetime.datetime.now()
        self.assertTrue(last_created > now - datetime.timedelta(days=2))

    # TODO(aaron-prindle) add test to check expired/unexpired logic for TTL


def _creation_time(image):
    cfg = json.loads(image.config_file())
    return cfg.get('created')


def _timestamp_to_time(dt_str):
    dt = dt_str.rstrip("Z")
    return datetime.datetime.strptime(dt, "%Y-%m-%dT%H:%M:%S")


if __name__ == '__main__':
    unittest.main()

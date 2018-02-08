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

import unittest
import datetime
import mock
import json

from ftl.common import context
from ftl.common import ftl_util
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
    @mock.patch('containerregistry.client.v2_2.docker_image.FromRegistry')
    def setUp(self, mock_from):
        mock_from.return_value.__enter__.return_value = None
        self.ctx = context.Memory()
        self.ctx.AddFile("app.py", _APP)
        args = mock.Mock()
        args.name = 'gcr.io/test/test:latest'
        args.base = 'gcr.io/google-appengine/python:latest'
        args.python_version = 'python2.7'
        self.builder = builder.Python(self.ctx, args, "")
        self.builder._pip_install = mock.Mock()

        # Mock out the calls to package managers for speed.
        self.builder.PackageLayer._gen_package_tar = mock.Mock()
        self.builder.PackageLayer._gen_package_tar.return_value = ('layer',
                                                                   'sha')

    def test_build_interpreter_layer_ttl_written(self):
        interpreter = self.builder.InterpreterLayer(
            self.builder._venv_dir, self.builder._args.python_version)
        interpreter._setup_venv = mock.Mock()
        interpreter.BuildLayer()
        overrides = ftl_util.CfgDctToOverrides(
            json.loads(interpreter.GetImage().config_file()))

        self.assertNotEqual(overrides.creation_time, "1970-01-01T00:00:00Z")
        last_created = ftl_util.timestamp_to_time(overrides.creation_time)
        now = datetime.datetime.now()
        self.assertTrue(last_created > now - datetime.timedelta(days=2))


if __name__ == '__main__':
    unittest.main()

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
import tarfile
import cStringIO
import mock
import httplib2
import tempfile
import hashlib

from containerregistry.client.v2_2 import docker_image

from ftl.common import context
from ftl.common import cache
from ftl.python import builder

class PythonTest(unittest.TestCase):

    def test_create_package_base_uncached(self):
        b = builder.Python(None)
        current_dir = os.path.dirname(__file__)
        base_config_path = os.path.join(current_dir, "testdata/base_image/config_file")
        with open(base_config_path, 'r') as reader:
            base_config = reader.read()
        ctx = context.Workspace(os.path.join(current_dir, "testdata/python_app"))
        with builder.Python(ctx) as b:
            fast_dir = "testdata/base_image/fast"
            # TODO(aaron-prindle) issues traversing fast_dir w/ bazel 'data', fix list creations
            sha_list = []
            for i in range(10):
                sha_list.append(os.path.join(fast_dir, "00"+str(i)+".sha256"))
            tar_list = []
            for i in range(10):
                sha_list.append(os.path.join(fast_dir, "00"+str(i)+".tar.gz"))
            with docker_image.FromDisk(base_config, zip(sha_list, tar_list)) as base_image:
                tmp = tempfile.mkdtemp()
                cash = cache.MockHybridRegistry('fake.gcr.io/google-appengine', tmp)
                with b.CreatePackageBase(base_image, cash) as python_app_layer:
                    # check that the python_app_layer image was added to the cache
                    self.assertEqual(1, len(cash.GetMap()))
                    for k in cash.GetMap():
                        self.assertEqual(str(k).startswith("fake.gcr.io/google-appengine/python-requirements-cache"),
                        True)
                cash.ResetCacheMiss()
                with b.CreatePackageBase(base_image, cash) as python_app_layer:
                    # check that the python_app_layer image was added to the cache
                    self.assertEqual(1, len(cash.GetMap()))
                    for k in cash.GetMap():
                        self.assertEqual(str(k).startswith("fake.gcr.io/google-appengine/python-requirements-cache"),
                        True)
                    # check that no additional cache misses occur on rebuild
                    self.assertEqual(0, cash.GetCacheMiss())

if __name__ == '__main__':
    unittest.main()

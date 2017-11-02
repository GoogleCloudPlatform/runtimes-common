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
import tempfile
import os
from mock import patch
from testing.lib import mock_registry_test_base
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import save
import main


class NodeTest(mock_registry_test_base.MockRegistryTestBase):
    def setUp(self):
        super(NodeTest, self).setUp()

    @patch('containerregistry.client.v2_2.append.Layer')
    def test_main(self, append_layer_mock):
        with docker_image.FromTarball('ftl/node_builder_base.tar') as img:
            self.registry.setImage('fake.gcr.io/base/image:initial', img)

        with docker_image.FromRegistry('fake.gcr.io/base/image:initial'
                                       ) as img:
            self.registry.setImage('fake.gcr.io/base/image:appended', img)

        self.AssertPushed(self.registry, 'fake.gcr.io/base/image:initial')
        self.AssertPushed(self.registry, 'fake.gcr.io/base/image:appended')

        append_layer_mock.return_value = self.registry.getImage(
            'fake.gcr.io/base/image:appended')

        args = ["--base=fake.gcr.io/base/image:initial",
                "--name=fake.gcr.io/base/image:latest", "--directory= "]

        main.main(args)

        self.AssertPushed(self.registry, 'fake.gcr.io/base/image:latest')
    
    @patch('containerregistry.client.v2_2.save.tarball')
    @patch('tempfile.mkdtemp')
    @patch('containerregistry.client.v2_2.append.Layer')
    def test_main_local(self, append_layer_mock, tmpdir_mock, save_tarball_mock):
        with docker_image.FromTarball('ftl/node_builder_base.tar') as img:
            self.registry.setImage('fake.gcr.io/base/image:initial', img)

        with docker_image.FromRegistry('fake.gcr.io/base/image:initial'
                                       ) as img:
            self.registry.setImage('fake.gcr.io/base/image:appended', img)

        self.AssertPushed(self.registry, 'fake.gcr.io/base/image:initial')
        self.AssertPushed(self.registry, 'fake.gcr.io/base/image:appended')

        append_layer_mock.return_value = self.registry.getImage(
            'fake.gcr.io/base/image:appended')
        tmpdir_mock.return_value = os.getcwd()
        save_tarball_mock.return_value = self.registry.getImage('fake.gcr.io/base/image:initial')

        args = ["--base=fake.gcr.io/base/image:initial",
                "--name=fake.gcr.io/base/image:latest", "--directory= ", "--local"]

        main.main(args)

        self.AssertNotPushed(self.registry, 'fake.gcr.io/base/image:latest')
        image_tar_path = os.path.join(os.getcwd(), 'image.tar')
        self.assertTrue(os.path.isfile(image_tar_path))
        os.remove(image_tar_path)


if __name__ == '__main__':
    unittest.main()

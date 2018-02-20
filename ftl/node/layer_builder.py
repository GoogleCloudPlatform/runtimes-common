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
"""This package implements the Node package layer builder."""


import os
import subprocess
import tempfile
import json
import datetime

from ftl.common import constants
from ftl.common import ftl_util
from ftl.common import single_layer_image
from ftl.common import tar_to_dockerimage


class LayerBuilder(single_layer_image.CacheableLayerBuilder):
    def __init__(self, ctx=None, descriptor_files=None, pkg_descriptor=None,
                 destination_path=constants.DEFAULT_DESTINATION_PATH):
        super(LayerBuilder, self).__init__()
        self._ctx = ctx
        self._descriptor_files = descriptor_files
        self._pkg_descriptor = pkg_descriptor
        self._destination_path = destination_path

    def GetCacheKeyRaw(self):
        return ftl_util.descriptor_parser(self._descriptor_files,
                                          self._ctx)

    def BuildLayer(self):
        """Override."""
        blob, u_blob = self._gen_npm_install_tar(self._pkg_descriptor,
                                                 self._destination_path)
        self._img = tar_to_dockerimage.FromFSImage(
            [blob], [u_blob], self._generate_overrides())

    def _gen_npm_install_tar(self, pkg_descriptor, destination_path):
        # Create temp directory to write package descriptor to
        pkg_dir = tempfile.mkdtemp()
        app_dir = os.path.join(pkg_dir, destination_path.strip("/"))
        os.makedirs(app_dir)

        # Copy out the relevant package descriptors to a tempdir.
        ftl_util.descriptor_copy(self._ctx, self._descriptor_files,
                                 app_dir)

        self._check_gcp_build(
            json.loads(self._ctx.GetFile(constants.PACKAGE_JSON)), app_dir)
        subprocess.check_call(
            ['rm', '-rf',
                os.path.join(app_dir, 'node_modules')])
        with ftl_util.Timing("npm_install"):
            if pkg_descriptor is None:
                subprocess.check_call(
                    ['npm', 'install', '--production'], cwd=app_dir)
            else:
                subprocess.check_call(
                    ['npm', 'install', '--production', pkg_descriptor],
                    cwd=app_dir)

        return ftl_util.zip_dir_to_layer_sha(pkg_dir)

    def _generate_overrides(self):
        overrides_dct = {
            'created': str(datetime.date.today()) + "T00:00:00Z",
        }
        return overrides_dct

    def _check_gcp_build(self, package_json, app_dir):
        scripts = package_json.get('scripts', {})
        gcp_build = scripts.get('gcp-build')

        if not gcp_build:
            return

        env = os.environ.copy()
        env["NODE_ENV"] = "development"
        subprocess.check_call(['npm', 'install'], cwd=app_dir, env=env)
        subprocess.check_call(
            ['npm', 'run-script', 'gcp-build'], cwd=app_dir, env=env)

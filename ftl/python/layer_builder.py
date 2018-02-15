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
"""This package implements the Python package layer builder."""

import datetime
import os
import subprocess

from ftl.common import ftl_util
from ftl.common import single_layer_image
from ftl.common import tar_to_dockerimage


class PackageLayerBuilder(single_layer_image.CacheableLayerBuilder):
    def __init__(self, ctx, descriptor_files, pkg_dir, dep_img_lyr):
        super(PackageLayerBuilder, self).__init__()
        self._ctx = ctx
        self._pkg_dir = pkg_dir
        self._descriptor_files = descriptor_files
        self._dep_img_lyr = dep_img_lyr

    def GetCacheKeyRaw(self):
        descriptor_contents = ftl_util.descriptor_parser(
            self._descriptor_files, self._ctx)
        return "%s %s" % (descriptor_contents,
                          self._dep_img_lyr.GetCacheKeyRaw())

    def BuildLayer(self):
        blob, u_blob = ftl_util.zip_dir_to_layer_sha(self._pkg_dir)
        self._img = tar_to_dockerimage.FromFSImage(
            [blob], [u_blob], _generate_overrides(False))


class InterpreterLayerBuilder(single_layer_image.CacheableLayerBuilder):
    def __init__(self, venv_dir, python_version):
        super(InterpreterLayerBuilder, self).__init__()
        self._venv_dir = venv_dir
        self._python_version = python_version

    def GetCacheKeyRaw(self):
        return self._python_version

    def BuildLayer(self):
        self._setup_venv(self._python_version)
        blob, u_blob = ftl_util.zip_dir_to_layer_sha(
            os.path.abspath(os.path.join(self._venv_dir, os.pardir)))
        self._img = tar_to_dockerimage.FromFSImage(
            [blob], [u_blob], _generate_overrides(True))

    def _setup_venv(self, python_version):
        with ftl_util.Timing("create_virtualenv"):
            subprocess.check_call([
                'virtualenv', '--no-download', self._venv_dir, '-p',
                python_version
            ])


def _generate_overrides(set_path):
    env = {
        "VIRTUAL_ENV": "/env",
    }
    if set_path:
        env['PATH'] = '/env/bin:$PATH'
    overrides_dct = {
        "created": str(datetime.date.today()) + "T00:00:00Z",
        "env": env
    }
    return overrides_dct

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
import logging

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
        self._img = tar_to_dockerimage.FromFSImage([blob], [u_blob],
                                                   _generate_overrides(False))


class InterpreterLayerBuilder(single_layer_image.CacheableLayerBuilder):
    def __init__(self, venv_dir, python_cmd, venv_cmd):
        super(InterpreterLayerBuilder, self).__init__()
        self._venv_dir = venv_dir
        self._python_cmd = python_cmd
        self._venv_cmd = venv_cmd

    def GetCacheKeyRaw(self):
        return "%s %s" % (self._python_cmd, self._venv_cmd)

    def BuildLayer(self):
        self._setup_venv()
        blob, u_blob = ftl_util.zip_dir_to_layer_sha(
            os.path.abspath(os.path.join(self._venv_dir, os.pardir)))
        self._img = tar_to_dockerimage.FromFSImage([blob], [u_blob],
                                                   _generate_overrides(True))

    def _setup_venv(self):
        with ftl_util.Timing("create_virtualenv"):
            venv_cmd_args = list(self._venv_cmd)
            venv_cmd_args.extend([
                '--no-download',
                self._venv_dir,
                '-p',
            ])
            venv_cmd_args.extend(self._python_cmd)
            proc_pipe = subprocess.Popen(
                venv_cmd_args,
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
            stdout, stderr = proc_pipe.communicate()
            logging.info("`virtualenv` stdout:\n%s" % stdout)
            if stderr:
                logging.error(
                    "`virtualenv` had error output:\n%s" % stderr)
            if proc_pipe.returncode:
                raise Exception("error: `virtualenv` returned code: %d" %
                                proc_pipe.returncode)

            subprocess.check_call(venv_cmd_args)


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

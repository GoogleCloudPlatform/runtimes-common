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
"""This package defines the interface for orchestrating image builds."""

import os
import subprocess
import tempfile
import logging
import datetime

from containerregistry.client.v2_2 import append
from containerregistry.transform.v2_2 import metadata

from ftl.common import builder
from ftl.common import ftl_util

_PYTHON_NAMESPACE = 'python-requirements-cache'
_REQUIREMENTS_TXT = 'requirements.txt'
_VENV_DIR = 'env'
_TMP_APP = 'app'
_WHEEL_DIR = 'wheel'


class Python(builder.JustApp):
    def __init__(self, ctx):
        self.descriptor_files = [_REQUIREMENTS_TXT]
        self.namespace = _PYTHON_NAMESPACE
        self.ctx = ctx
        self._tmp_app = self._gen_tmp_dir(_TMP_APP)
        self._venv_dir = self._gen_tmp_dir(_VENV_DIR)
        self._wheel_dir = self._gen_tmp_dir(_WHEEL_DIR)
        super(Python, self).__init__(ctx)

    def __enter__(self):
        """Override."""
        return self

    def _generate_overrides(self, set_path):
        env = {
            "VIRTUAL_ENV": "/env",
        }
        if set_path:
            env['PATH'] = '/env/bin:$PATH'
        return metadata.Overrides(
            creation_time=str(datetime.date.today()) + "T00:00:00Z", env=env)

    def CreatePackageBase(self, base, python_version='python2.7'):
        """Override."""
        package_base = base

        self._setup_app_dir(self._tmp_app)
        self._setup_venv(python_version)
        layer, sha = ftl_util.zip_dir_to_layer_sha(
            os.path.abspath(os.path.join(self._venv_dir, os.pardir)))
        package_base = append.Layer(
            package_base,
            layer,
            diff_id=sha,
            overrides=self._generate_overrides(True))

        self._pip_install()
        whls = self._resolve_whls()
        pkg_dirs = [self._whl_to_fslayer(whl) for whl in whls]
        logging.info("pkg_dirs" + str(pkg_dirs))
        for pkg_dir in pkg_dirs:
            layer, sha = ftl_util.zip_dir_to_layer_sha(pkg_dir)
            logging.info('Generated layer with sha: %s', sha)
            package_base = append.Layer(
                package_base,
                layer,
                diff_id=sha,
                overrides=self._generate_overrides(False))
        return package_base

    def _gen_dirs(self, dirs):
        tmp_dir = tempfile.mkdtemp()
        dir_map = {}
        for dir in dirs:
            dir_name = os.path.join(tmp_dir, dir)
            dir_map[dir] = dir_name
            os.mkdir(dir_name)
        return dir_map

    def _gen_tmp_dir(self, dirr):
        tmp_dir = tempfile.mkdtemp()
        dir_name = os.path.join(tmp_dir, dirr)
        os.mkdir(dir_name)
        return dir_name

    def _gen_pip_env(self):
        pip_env = os.environ.copy()
        # bazel adds its own PYTHONPATH to the env
        # which must be removed for the pip calls to work properly
        del pip_env['PYTHONPATH']
        pip_env['VIRTUAL_ENV'] = self._venv_dir
        pip_env['PATH'] = self._venv_dir + "/bin" + ":" + os.environ['PATH']
        return pip_env

    def _setup_app_dir(self, app_dir):
        # Copy out the relevant package descriptors to a tempdir.
        for f in self.descriptor_files:
            if self._ctx.Contains(f):
                with open(os.path.join(app_dir, f), 'w') as w:
                    w.write(self._ctx.GetFile(f))

    def _setup_venv(self, python_version):
        with ftl_util.Timing("create_virtualenv"):
            subprocess.check_call(
                [
                    'virtualenv', '--no-download', self._venv_dir, '-p',
                    python_version
                ],
                cwd=self._tmp_app)

    def _pip_install(self):
        with ftl_util.Timing("pip_install_wheels"):
            subprocess.check_call(
                [
                    'pip', 'wheel', '-w', self._wheel_dir, '-r',
                    'requirements.txt'
                ],
                cwd=self._tmp_app,
                env=self._gen_pip_env())

    def _resolve_whls(self):
        return [
            os.path.join(self._wheel_dir, f)
            for f in os.listdir(self._wheel_dir)
        ]

    def _whl_to_fslayer(self, whl):
        tmp_dir = tempfile.mkdtemp()
        pkg_dir = os.path.join(tmp_dir, 'env')
        os.makedirs(pkg_dir)
        subprocess.check_call(
            ['pip', 'install', '--prefix', pkg_dir, whl],
            env=self._gen_pip_env())
        return tmp_dir


def From(ctx):
    return Python(ctx)

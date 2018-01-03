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
import datetime

from ftl.common import builder
from ftl.common import ftl_util
from ftl.common import single_layer_image
from ftl.common import tar_to_dockerimage

_VENV_DIR = 'env'
_WHEEL_DIR = 'wheel'
_THREADS = 32
_REQUIREMENTS_TXT = 'requirements.txt'
_PYTHON_NAMESPACE = 'python-requirements-cache'


def _generate_overrides(set_path):
    env = {
        "VIRTUAL_ENV": "/env",
    }
    if set_path:
        env['PATH'] = '/env/bin:$PATH'
    overrides_dct = {
        "creation_time": str(datetime.date.today()) + "T00:00:00Z",
        "env": env
    }
    return overrides_dct


class Python(builder.RuntimeBase):
    def __init__(self, ctx, args, cache_version_str):
        super(Python, self).__init__(ctx, _PYTHON_NAMESPACE, args,
                                     cache_version_str, [_REQUIREMENTS_TXT])
        self._venv_dir = ftl_util.gen_tmp_dir(_VENV_DIR)
        self._wheel_dir = ftl_util.gen_tmp_dir(_WHEEL_DIR)

    def Build(self):
        lyr_imgs = []
        lyr_imgs.append(self._base)
        if ftl_util.has_pkg_descriptor(self._descriptor_files, self._ctx):

            interpreter = self.InterpreterLayer(self._venv_dir,
                                                self._args.python_version)
            cached_int_img = self._cash.GetAndCheckTTL(
                self._base, self._namespace, interpreter.GetCacheKey())
            if cached_int_img is not None:
                interpreter.SetImage(cached_int_img)
            else:
                interpreter.BuildLayer()
                self._cash.Store(self._base, self._namespace,
                                 interpreter.GetCacheKey(),
                                 interpreter.GetImage())
            lyr_imgs.append(interpreter)

            pkg_descriptor = ftl_util.descriptor_parser(
                self._descriptor_files, self._ctx)
            self._pip_install(pkg_descriptor)

            whls = self._resolve_whls()
            pkg_dirs = [self._whl_to_fslayer(whl) for whl in whls]

            for whl_pkg_dir in pkg_dirs:
                pkg = self.PackageLayer(self._ctx, self._descriptor_files,
                                        whl_pkg_dir, interpreter)
                cached_pkg_img = self._cash.GetAndCheckTTL(
                    self._base, self._namespace, pkg.GetCacheKey())
                if cached_pkg_img is not None:
                    pkg.SetImage(cached_pkg_img)
                else:
                    pkg.BuildLayer()
                    self._cash.Store(self._base, self._namespace,
                                     pkg.GetCacheKey(), pkg.GetImage())
                lyr_imgs.append(pkg)

        app = self.AppLayer(self._ctx)
        app.BuildLayer()
        lyr_imgs.append(app)
        ftl_image = self.AppendLayersIntoImage(lyr_imgs)
        self.StoreImage(ftl_image)

    class InterpreterLayer(single_layer_image.CacheableLayer):
        def __init__(self, venv_dir, python_version):
            super(Python.InterpreterLayer, self).__init__()
            self._venv_dir = venv_dir
            self._python_version = python_version

        def GetCacheKeyRaw(self):
            return self._python_version

        def BuildLayer(self):
            self._setup_venv(self._python_version)
            lyr, sha = ftl_util.zip_dir_to_layer_sha(
                os.path.abspath(os.path.join(self._venv_dir, os.pardir)))
            self._img = tar_to_dockerimage.FromFSImage(
                lyr, _generate_overrides(True))

        def _setup_venv(self, python_version):
            with ftl_util.Timing("create_virtualenv"):
                subprocess.check_call([
                    'virtualenv', '--no-download', self._venv_dir, '-p',
                    python_version
                ])

    def _pip_install(self, pkg_txt):
        with ftl_util.Timing("pip_install_wheels"):
            args = ['pip', 'wheel', '-w', self._wheel_dir, '-r', "/dev/stdin"]

            pipe1 = subprocess.Popen(
                args,
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                env=self._gen_pip_env(), )
            pipe1.communicate(input=pkg_txt)[0]

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

    def _gen_pip_env(self):
        pip_env = os.environ.copy()
        # bazel adds its own PYTHONPATH to the env
        # which must be removed for the pip calls to work properly
        del pip_env['PYTHONPATH']
        pip_env['VIRTUAL_ENV'] = self._venv_dir
        pip_env['PATH'] = self._venv_dir + "/bin" + ":" + os.environ['PATH']
        return pip_env

    class PackageLayer(single_layer_image.CacheableLayer):
        def __init__(self, ctx, descriptor_files, pkg_dir, dep_img_lyr):
            super(Python.PackageLayer, self).__init__()
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
            lyr, sha = ftl_util.zip_dir_to_layer_sha(self._pkg_dir)
            self._img = tar_to_dockerimage.FromFSImage(
                lyr, _generate_overrides(False))

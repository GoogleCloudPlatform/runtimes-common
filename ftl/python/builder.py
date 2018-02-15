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

from ftl.common import builder
from ftl.common import ftl_util
from ftl.common import layer_builder as base_builder
from ftl.python import layer_builder as package_builder

_VENV_DIR = 'env'
_WHEEL_DIR = 'wheel'
_THREADS = 32
_REQUIREMENTS_TXT = 'requirements.txt'
_PYTHON_NAMESPACE = 'python-requirements-cache'


class Python(builder.RuntimeBase):
    def __init__(self, ctx, args, cache_version_str):
        super(Python, self).__init__(ctx, _PYTHON_NAMESPACE, args,
                                     cache_version_str, [_REQUIREMENTS_TXT])
        self._venv_dir = ftl_util.gen_tmp_dir(_VENV_DIR)
        self._wheel_dir = ftl_util.gen_tmp_dir(_WHEEL_DIR)

    def Build(self):
        lyr_imgs = []
        lyr_imgs.append(self._base_image)
        if ftl_util.has_pkg_descriptor(self._descriptor_files, self._ctx):
            interpreter_builder = package_builder.InterpreterLayerBuilder(
                self._venv_dir,
                self._args.python_version)
            cached_int_img = None
            if self._args.cache:
                with ftl_util.Timing("checking cached int layer"):
                    key = interpreter_builder.GetCacheKey()
                    cached_int_img = self._cache.Get(key)
            if cached_int_img is not None:
                interpreter_builder.SetImage(cached_int_img)
            else:
                with ftl_util.Timing("building int layer"):
                    interpreter_builder.BuildLayer()
                if self._args.cache:
                    with ftl_util.Timing("uploading int layer"):
                        self._cache.Set(interpreter_builder.GetCacheKey(),
                                        interpreter_builder.GetImage())
            lyr_imgs.append(interpreter_builder)

            pkg_descriptor = ftl_util.descriptor_parser(
                self._descriptor_files, self._ctx)

            with ftl_util.Timing("installing pip packages"):
                self._pip_install(pkg_descriptor)

            with ftl_util.Timing("resolving whl paths"):
                whls = self._resolve_whls()
                pkg_dirs = [self._whl_to_fslayer(whl) for whl in whls]

            for whl_pkg_dir in pkg_dirs:
                layer_builder = package_builder.PackageLayerBuilder(
                    self._ctx, self._descriptor_files,
                    whl_pkg_dir, interpreter_builder)
                cached_pkg_img = None
                if self._args.cache:
                    with ftl_util.Timing("checking cached pkg layer"):
                        key = layer_builder.GetCacheKey()
                        cached_pkg_img = self._cache.Get(key)
                if cached_pkg_img is not None:
                    layer_builder.SetImage(cached_pkg_img)
                else:
                    with ftl_util.Timing("building pkg layer"):
                        layer_builder.BuildLayer()
                    if self._args.cache:
                        with ftl_util.Timing("uploading pkg layer"):
                            self._cache.Set(layer_builder.GetCacheKey(),
                                            layer_builder.GetImage())
                lyr_imgs.append(layer_builder)

        app = base_builder.AppLayerBuilder(
            ctx=self._ctx,
            destination_path=self._args.destination_path,
            entrypoint=self._args.entrypoint,
            exposed_ports=self._args.exposed_ports)
        with ftl_util.Timing("builder app layer"):
            app.BuildLayer()
        lyr_imgs.append(app)
        with ftl_util.Timing("stitching lyrs into final image"):
            ftl_image = self.AppendLayersIntoImage(lyr_imgs)
        with ftl_util.Timing("uploading final image"):
            self.StoreImage(ftl_image)

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

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

import logging
import os
import tempfile
import subprocess
import concurrent.futures

from ftl.common import constants
from ftl.common import ftl_util
from ftl.common import ftl_error
from ftl.common import single_layer_image
from ftl.common import tar_to_dockerimage

from ftl.python import python_util


class PackageLayerBuilder(single_layer_image.CacheableLayerBuilder):
    def __init__(self,
                 ctx=None,
                 descriptor_files=None,
                 pkg_dir=None,
                 dep_img_lyr=None,
                 cache_key_version=None,
                 cache=None):
        super(PackageLayerBuilder, self).__init__()
        self._ctx = ctx
        self._pkg_dir = pkg_dir
        self._descriptor_files = descriptor_files
        self._dep_img_lyr = dep_img_lyr
        self._cache_key_version = cache_key_version
        self._cache = cache

    def GetCacheKeyRaw(self):
        cache_key = ""
        return "%s %s" % (cache_key, self._cache_key_version)

    def BuildLayer(self):
        with ftl_util.Timing('building_python_pkg_layer'):
            self._build_layer()
        if self._cache:
            with ftl_util.Timing('uploading_python_pkg_layer'):
                self._cache.Set(self.GetCacheKey(), self.GetImage())

    def _build_layer(self):
        blob, u_blob = ftl_util.zip_dir_to_layer_sha(self._pkg_dir, "")
        overrides = ftl_util.generate_overrides(False)
        self._img = tar_to_dockerimage.FromFSImage([blob], [u_blob], overrides)

    def _log_cache_result(self, hit):
        if hit:
            cache_str = constants.PHASE_1_CACHE_HIT
        else:
            cache_str = constants.PHASE_1_CACHE_MISS
        logging.info(
            cache_str.format(
                key_version=constants.CACHE_KEY_VERSION,
                language='PYTHON (package)',
                key=self.GetCacheKey()))


class RequirementsLayerBuilder(single_layer_image.CacheableLayerBuilder):
    def __init__(self,
                 ctx=None,
                 descriptor_files=None,
                 directory=None,
                 pkg_dir=None,
                 dep_img_lyr=None,
                 cache_key_version=None,
                 wheel_dir=constants.WHEEL_DIR,
                 virtualenv_dir=constants.VIRTUALENV_DIR,
                 python_cmd=[constants.PYTHON_DEFAULT_CMD],
                 pip_cmd=[constants.PIP_DEFAULT_CMD],
                 virtualenv_cmd=[constants.VIRTUALENV_DEFAULT_CMD],
                 venv_cmd=[constants.VENV_DEFAULT_CMD],
                 cache=None):
        super(RequirementsLayerBuilder, self).__init__()
        self._ctx = ctx
        self._pkg_dir = pkg_dir
        self._wheel_dir = wheel_dir
        self._virtualenv_dir = virtualenv_dir
        self._python_cmd = python_cmd
        self._pip_cmd = pip_cmd
        self._virtualenv_cmd = virtualenv_cmd
        self._venv_cmd = venv_cmd
        self._descriptor_files = descriptor_files
        self._directory = directory
        self._dep_img_lyr = dep_img_lyr
        self._cache_key_version = cache_key_version
        self._cache = cache

    def GetCacheKeyRaw(self):
        descriptor_contents = ftl_util.descriptor_parser(
            self._descriptor_files, self._ctx)
        cache_key = '%s %s' % (descriptor_contents,
                               self._dep_img_lyr.GetCacheKeyRaw())
        return "%s %s" % (cache_key, self._cache_key_version)

    def BuildLayer(self):
        cached_img = None
        if self._cache:
            with ftl_util.Timing('checking_cached_requirements.txt_layer'):
                key = self.GetCacheKey()
                cached_img = self._cache.Get(key)
                self._log_cache_result(False if cached_img is None else True)
        if cached_img:
            self.SetImage(cached_img)
        else:
            python_util.setup_virtualenv(self._virtualenv_dir,
                                         self._virtualenv_cmd,
                                         self._python_cmd,
                                         self._venv_cmd)

            pkg_descriptor = ftl_util.descriptor_parser(
                self._descriptor_files, self._ctx)
            self._pip_download_wheels(pkg_descriptor)
            whls = self._resolve_whls()
            pkg_dirs = [self._whl_to_fslayer(whl) for whl in whls]

            req_txt_imgs = []
            with ftl_util.Timing('uploading_all_package_layers'):
                with concurrent.futures.ThreadPoolExecutor(
                        max_workers=constants.THREADS) as executor:
                    future_to_params = {
                        executor.submit(self._build_pkg, whl_pkg_dir,
                                        req_txt_imgs): whl_pkg_dir
                        for whl_pkg_dir in pkg_dirs
                    }
                    for future in concurrent.futures.as_completed(
                            future_to_params):
                        future.result()

            req_txt_image = ftl_util.AppendLayersIntoImage(req_txt_imgs)

            if req_txt_image:
                self.SetImage(req_txt_image)

                if self._cache:
                    with ftl_util.Timing('uploading_requirements.txt_pkg_lyr'):
                        self._cache.Set(self.GetCacheKey(), self.GetImage())

    def _build_pkg(self, whl_pkg_dir, req_txt_imgs):
        layer_builder = PackageLayerBuilder(
            ctx=self._ctx,
            descriptor_files=self._descriptor_files,
            pkg_dir=whl_pkg_dir,
            dep_img_lyr=self._dep_img_lyr,
            cache_key_version=self._cache_key_version,
            cache=self._cache)
        layer_builder.BuildLayer()
        req_txt_imgs.append(layer_builder.GetImage())

    def _resolve_whls(self):
        with ftl_util.Timing('resolving_whl_paths'):
            return [
                os.path.join(self._wheel_dir, f)
                for f in os.listdir(self._wheel_dir)
            ]

    def _whl_to_fslayer(self, whl):
        tmp_dir = tempfile.mkdtemp()
        pkg_dir = os.path.join(tmp_dir, self._virtualenv_dir.lstrip('/'))
        os.makedirs(pkg_dir)

        pip_cmd_args = list(self._pip_cmd)
        pip_cmd_args.extend(['install', '--no-deps', '--prefix', pkg_dir, whl])
        pip_cmd_args.extend(constants.PIP_OPTIONS)
        ftl_util.run_command('pip_install_from_wheels', pip_cmd_args, None,
                             self._gen_pip_env())
        return tmp_dir

    def _pip_download_wheels(self, pkg_txt):
        ftl_util.run_command(
            'pip_install_wheel',
            ['pip', 'install', 'wheel'],
            cmd_cwd=self._directory,
            cmd_env=self._gen_pip_env(),
            cmd_input=pkg_txt,
            err_type=ftl_error.FTLErrors.USER())

        pip_cmd_args = list(self._pip_cmd)
        pip_cmd_args.extend(
            ['wheel', '-w', self._wheel_dir, '-r', 'requirements.txt'])
        pip_cmd_args.extend(constants.PIP_OPTIONS)
        ftl_util.run_command(
            'pip_download_wheels',
            pip_cmd_args,
            cmd_cwd=self._directory,
            cmd_env=self._gen_pip_env(),
            cmd_input=pkg_txt,
            err_type=ftl_error.FTLErrors.USER())

    def _gen_pip_env(self):
        pip_env = os.environ.copy()
        # bazel adds its own PYTHONPATH to the env
        # which must be removed for the pip calls to work properly
        pip_env.pop('PYTHONPATH', None)
        pip_env['VIRTUAL_ENV'] = self._virtualenv_dir
        pip_env['PATH'] = self._virtualenv_dir + '/bin' + ':' + os.environ[
            'PATH']
        return pip_env

    def _log_cache_result(self, hit):
        if hit:
            cache_str = constants.PHASE_1_CACHE_HIT
        else:
            cache_str = constants.PHASE_1_CACHE_MISS
        logging.info(
            cache_str.format(
                key_version=constants.CACHE_KEY_VERSION,
                language='PYTHON (requirements)',
                key=self.GetCacheKey()))


class PipfileLayerBuilder(RequirementsLayerBuilder):
    def __init__(self,
                 ctx=None,
                 descriptor_files=None,
                 directory=None,
                 pkg_descriptor=None,
                 pkg_dir=None,
                 dep_img_lyr=None,
                 cache_key_version=None,
                 wheel_dir=constants.WHEEL_DIR,
                 virtualenv_dir=constants.VIRTUALENV_DIR,
                 python_cmd=[constants.PYTHON_DEFAULT_CMD],
                 pip_cmd=[constants.PIP_DEFAULT_CMD],
                 virtualenv_cmd=[constants.VIRTUALENV_DEFAULT_CMD],
                 cache=None):
        super(PipfileLayerBuilder, self).__init__()
        self._ctx = ctx
        self._pkg_dir = pkg_dir
        self._wheel_dir = wheel_dir
        self._virtualenv_dir = virtualenv_dir
        self._python_cmd = python_cmd
        self._pip_cmd = pip_cmd
        self._virtualenv_cmd = virtualenv_cmd
        self._descriptor_files = descriptor_files
        self._directory = directory
        self._dep_img_lyr = dep_img_lyr
        self._cache_key_version = cache_key_version
        self._cache = cache
        self._pkg_descriptor = pkg_descriptor

    def GetCacheKeyRaw(self):
        cache_key = "%s %s %s" % (self._pkg_descriptor[0],
                                  self._pkg_descriptor[1],
                                  self._dep_img_lyr.GetCacheKeyRaw())
        return "%s %s" % (cache_key, self._cache_key_version)

    def _log_cache_result(self, hit):
        if hit:
            cache_str = constants.PHASE_2_CACHE_HIT
        else:
            cache_str = constants.PHASE_2_CACHE_MISS
        logging.info(
            cache_str.format(
                key_version=constants.CACHE_KEY_VERSION,
                language='PYTHON',
                package_name=self._pkg_descriptor[0],
                package_version=self._pkg_descriptor[1],
                key=self.GetCacheKey()))

    def BuildLayer(self):
        cached_img = None
        if self._cache:
            with ftl_util.Timing('checking_cached_pipfile_pkg_layer'):
                key = self.GetCacheKey()
                cached_img = self._cache.Get(key)
                self._log_cache_result(False if cached_img is None else True)
        if cached_img:
            self.SetImage(cached_img)
        else:
            self._pip_download_wheels(' '.join(self._pkg_descriptor))
            whls = self._resolve_whls()
            if len(whls) != 1:
                raise Exception("expected one whl for one installed pkg")
            pkg_dir = self._whl_to_fslayer(whls[0])
            blob, u_blob = ftl_util.zip_dir_to_layer_sha(pkg_dir, "")
            overrides = ftl_util.generate_overrides(False)
            self._img = tar_to_dockerimage.FromFSImage([blob], [u_blob],
                                                       overrides)
            if self._cache:
                with ftl_util.Timing('uploading_pipfile_pkg_layer'):
                    self._cache.Set(self.GetCacheKey(), self.GetImage())

    def _pip_download_wheels(self, pkg_txt):
        pip_cmd_args = list(self._pip_cmd)
        pip_cmd_args.extend(
            ['wheel', '-w', self._wheel_dir, '-r', '/dev/stdin'])
        pip_cmd_args.extend(['--no-deps'])
        pip_cmd_args.extend(constants.PIP_OPTIONS)
        ftl_util.run_command(
            'pip_download_wheels',
            pip_cmd_args,
            cmd_cwd=self._directory,
            cmd_env=self._gen_pip_env(),
            cmd_input=pkg_txt,
            err_type=ftl_error.FTLErrors.USER())


class InterpreterLayerBuilder(single_layer_image.CacheableLayerBuilder):
    def __init__(self,
                 virtualenv_dir=constants.VIRTUALENV_DIR,
                 python_cmd=[constants.PYTHON_DEFAULT_CMD],
                 virtualenv_cmd=[constants.VIRTUALENV_DEFAULT_CMD],
                 venv_cmd=[constants.VENV_DEFAULT_CMD],
                 cache_key_version=None,
                 cache=None):
        super(InterpreterLayerBuilder, self).__init__()
        self._virtualenv_dir = virtualenv_dir
        self._python_cmd = python_cmd
        self._virtualenv_cmd = virtualenv_cmd
        self._venv_cmd = venv_cmd
        self._cache_key_version = cache_key_version
        self._cache = cache

    def GetCacheKeyRaw(self):
        cache_key = '%s %s %s' % (self._python_version(), self._virtualenv_cmd,
                                  self._virtualenv_dir)
        return "%s %s" % (cache_key, self._cache_key_version)

    def _python_version(self):
        with ftl_util.Timing('check python version'):
            python_version_cmd = list(self._python_cmd)
            python_version_cmd.append('--version')
            logging.info("`python version` full cmd:\n%s" %
                         " ".join(python_version_cmd))
            proc_pipe = subprocess.Popen(
                python_version_cmd,
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
            )
            stdout, stderr = proc_pipe.communicate()
            logging.info("`python version` stderr:\n%s" % stderr)
            if proc_pipe.returncode:
                raise Exception("error: `python version` returned code: %d" %
                                proc_pipe.returncode)
            #  up until Python 3.4 the version info gets written to stderr
            return stdout if len(stdout) >= len(stderr) else stderr

    def BuildLayer(self):
        cached_img = None
        if self._cache:
            with ftl_util.Timing('checking_cached_interpreter_layer'):
                key = self.GetCacheKey()
                cached_img = self._cache.Get(key)
                self._log_cache_result(False if cached_img is None else True)
        if cached_img:
            self.SetImage(cached_img)
        else:
            with ftl_util.Timing('building_interpreter_layer'):
                self._build_layer()
            if self._cache:
                with ftl_util.Timing('uploading_interpreter_layer'):
                    self._cache.Set(self.GetCacheKey(), self.GetImage())

    def _build_layer(self):
        python_util.setup_virtualenv(self._virtualenv_dir,
                                     self._virtualenv_cmd,
                                     self._python_cmd,
                                     self._venv_cmd)

        blob, u_blob = ftl_util.zip_dir_to_layer_sha(self._virtualenv_dir,
                                                     self._virtualenv_dir)

        overrides = ftl_util.generate_overrides(True, self._virtualenv_dir)
        self._img = tar_to_dockerimage.FromFSImage([blob], [u_blob], overrides)

    def _log_cache_result(self, hit):
        if hit:
            cache_str = constants.PHASE_1_CACHE_HIT
        else:
            cache_str = constants.PHASE_1_CACHE_MISS
        logging.info(
            cache_str.format(
                key_version=constants.CACHE_KEY_VERSION,
                language='PYTHON (interpreter)',
                key=self.GetCacheKey()))

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

import logging
import os
import json

from ftl.common import constants
from ftl.common import ftl_util
from ftl.common import ftl_error
from ftl.common import single_layer_image
from ftl.common import tar_to_dockerimage


class LayerBuilder(single_layer_image.CacheableLayerBuilder):
    def __init__(self,
                 ctx=None,
                 descriptor_files=None,
                 pkg_descriptor=None,
                 directory=None,
                 destination_path=constants.DEFAULT_DESTINATION_PATH,
                 should_use_yarn=None,
                 cache_key_version=None,
                 cache=None):
        super(LayerBuilder, self).__init__()
        self._ctx = ctx
        self._descriptor_files = descriptor_files
        self._pkg_descriptor = pkg_descriptor
        self._directory = directory
        self._destination_path = destination_path
        self._should_use_yarn = should_use_yarn
        self._cache_key_version = cache_key_version
        self._cache = cache

    def GetCacheKeyRaw(self):
        all_descriptor_contents = ftl_util.all_descriptor_contents(
            self._descriptor_files, self._ctx)
        cache_key = '%s %s' % (all_descriptor_contents, self._destination_path)
        return "%s %s" % (cache_key, self._cache_key_version)

    def BuildLayer(self):
        """Override."""
        cached_img = None
        is_gcp_build = False
        if self._ctx and self._ctx.Contains(constants.PACKAGE_JSON):
            is_gcp_build = self._is_gcp_build(
                json.loads(self._ctx.GetFile(constants.PACKAGE_JSON)))

        if is_gcp_build:
            if self._should_use_yarn:
                self._gcp_build(self._directory, 'yarn', 'run')
                self._cleanup_build_layer()
            else:
                self._gcp_build(self._directory, 'npm', 'run-script')
                self._cleanup_build_layer()

        key = self.GetCacheKey()
        if self._cache:
            with ftl_util.Timing('checking_cached_packages_json_layer'):
                cached_img = self._cache.Get(key)
                self._log_cache_result(False if cached_img is None else True,
                                       key)
        if cached_img:
            self.SetImage(cached_img)
        else:
            with ftl_util.Timing('building_packages_json_layer'):
                self._build_layer()
                self._cleanup_build_layer()
            if self._cache:
                with ftl_util.Timing('uploading_packages_json_layer'):
                    self._cache.Set(key, self.GetImage())

    def _build_layer(self):
        if self._should_use_yarn:
            blob, u_blob = self._gen_yarn_install_tar(self._directory)
        else:
            blob, u_blob = self._gen_npm_install_tar(self._directory)
        self._img = tar_to_dockerimage.FromFSImage([blob], [u_blob],
                                                   ftl_util.generate_overrides(
                                                       False))

    def _cleanup_build_layer(self):
        if self._directory:
            modules_dir = os.path.join(self._directory, "node_modules")
            rm_cmd = ['rm', '-rf', modules_dir]
            ftl_util.run_command('rm_node_modules', rm_cmd)

    def _gen_yarn_install_tar(self, app_dir):
        yarn_install_cmd = ['yarn', 'install', '--production']
        ftl_util.run_command(
            'yarn_install',
            yarn_install_cmd,
            cmd_cwd=app_dir,
            err_type=ftl_error.FTLErrors.USER())

        module_destination = os.path.join(self._destination_path,
                                          'node_modules')
        modules_dir = os.path.join(self._directory, "node_modules")
        return ftl_util.zip_dir_to_layer_sha(modules_dir, module_destination)

    def _gen_npm_install_tar(self, app_dir):
        npm_install_cmd = ['npm', 'install', '--production']
        ftl_util.run_command(
            'npm_install',
            npm_install_cmd,
            cmd_cwd=app_dir,
            err_type=ftl_error.FTLErrors.USER())

        module_destination = os.path.join(self._destination_path,
                                          'node_modules')
        modules_dir = os.path.join(self._directory, "node_modules")
        return ftl_util.zip_dir_to_layer_sha(modules_dir, module_destination)

    def _is_gcp_build(self, package_json):
        scripts = package_json.get('scripts', {})
        if scripts.get('gcp-build'):
            return True
        return False

    def _gcp_build(self, app_dir, install_bin, run_cmd):
        env = os.environ.copy()
        env["NODE_ENV"] = "development"
        install_cmd = [install_bin, 'install']
        ftl_util.run_command(
            '%s_install' % install_bin,
            install_cmd,
            app_dir,
            env,
            err_type=ftl_error.FTLErrors.USER())

        npm_run_script_cmd = [install_bin, run_cmd, 'gcp-build']
        ftl_util.run_command(
            '%s_%s_gcp_build' % (install_bin, run_cmd),
            npm_run_script_cmd,
            app_dir,
            env,
            err_type=ftl_error.FTLErrors.USER())

    def _log_cache_result(self, hit, key):
        if self._pkg_descriptor:
            if hit:
                cache_str = constants.PHASE_2_CACHE_HIT
            else:
                cache_str = constants.PHASE_2_CACHE_MISS
            logging.info(
                cache_str.format(
                    key_version=constants.CACHE_KEY_VERSION,
                    language='NODE',
                    package_name=self._pkg_descriptor[0],
                    package_version=self._pkg_descriptor[1],
                    key=key))
        else:
            if hit:
                cache_str = constants.PHASE_1_CACHE_HIT
            else:
                cache_str = constants.PHASE_1_CACHE_MISS
            logging.info(
                cache_str.format(
                    key_version=constants.CACHE_KEY_VERSION,
                    language='NODE',
                    key=key))

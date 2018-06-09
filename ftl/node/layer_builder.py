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
import datetime

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
                 cache=None):
        super(LayerBuilder, self).__init__()
        self._ctx = ctx
        self._descriptor_files = descriptor_files
        self._pkg_descriptor = pkg_descriptor
        self._directory = directory
        self._destination_path = destination_path
        self._cache = cache

    def GetCacheKeyRaw(self):
        all_descriptor_contents = ftl_util.all_descriptor_contents(
            self._descriptor_files, self._ctx)
        return '%s %s' % (all_descriptor_contents, self._destination_path)

    def BuildLayer(self):
        """Override."""
        cached_img = None
        if self._cache:
            with ftl_util.Timing('checking_cached_packages_json_layer'):
                key = self.GetCacheKey()
                cached_img = self._cache.Get(key)
                self._log_cache_result(False if cached_img is None else True)
        if cached_img:
            self.SetImage(cached_img)
        else:
            with ftl_util.Timing('building_packages_json_layer'):
                self._build_layer()
            if self._cache:
                with ftl_util.Timing('uploading_packages_json_layer'):
                    self._cache.Set(self.GetCacheKey(), self.GetImage())

    def _build_layer(self):
        blob, u_blob = self._gen_npm_install_tar(self._pkg_descriptor,
                                                 self._directory)
        self._img = tar_to_dockerimage.FromFSImage([blob], [u_blob],
                                                   self._generate_overrides())

    def _gen_npm_install_tar(self, pkg_descriptor, app_dir):
        if self._ctx:
            if self._ctx.Contains(constants.PACKAGE_JSON):
                self._check_gcp_build(
                    json.loads(self._ctx.GetFile(constants.PACKAGE_JSON)),
                    app_dir)
        rm_cmd = ['rm', '-rf', os.path.join(app_dir, 'node_modules')]
        ftl_util.run_command('rm_node_modules', rm_cmd)

        if not pkg_descriptor:
            npm_install_cmd = ['npm', 'install', '--production']
            ftl_util.run_command(
                'npm_install',
                npm_install_cmd,
                cmd_cwd=app_dir,
                err_type=ftl_error.FTLErrors.USER())
        else:
            npm_install_cmd = [
                'npm', 'install', '--production', pkg_descriptor
            ]
            ftl_util.run_command(
                'npm_install',
                npm_install_cmd,
                cmd_cwd=app_dir,
                err_type=ftl_error.FTLErrors.USER())

        return ftl_util.zip_dir_to_layer_sha(app_dir,
                                             self._destination_path)

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
        npm_install_cmd = ['npm', 'install']
        ftl_util.run_command(
            'npm_install',
            npm_install_cmd,
            app_dir,
            env,
            err_type=ftl_error.FTLErrors.USER())

        npm_run_script_cmd = ['npm', 'run-script', 'gcp-build']
        ftl_util.run_command(
            'npm_run_script_gcp_build',
            npm_run_script_cmd,
            app_dir,
            env,
            err_type=ftl_error.FTLErrors.USER())

    def _log_cache_result(self, hit):
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
                    key=self.GetCacheKey()))
        else:
            if hit:
                cache_str = constants.PHASE_1_CACHE_HIT
            else:
                cache_str = constants.PHASE_1_CACHE_MISS
            logging.info(
                cache_str.format(
                    key_version=constants.CACHE_KEY_VERSION,
                    language='NODE',
                    key=self.GetCacheKey()))

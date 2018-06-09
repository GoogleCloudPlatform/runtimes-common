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
"""This package implements the PHP package layer builder."""

import logging
import os
import datetime

from ftl.common import constants
from ftl.common import ftl_util
from ftl.common import ftl_error
from ftl.common import single_layer_image
from ftl.common import tar_to_dockerimage

from ftl.php import php_util


class PhaseOneLayerBuilder(single_layer_image.CacheableLayerBuilder):
    def __init__(self,
                 ctx=None,
                 descriptor_files=None,
                 destination_path=constants.DEFAULT_DESTINATION_PATH,
                 directory=None,
                 cache=None):
        super(PhaseOneLayerBuilder, self).__init__()
        self._ctx = ctx
        self._descriptor_files = descriptor_files
        self._destination_path = destination_path
        self._directory = directory
        self._cache = cache

    def GetCacheKeyRaw(self):
        return "%s %s" % (
            ftl_util.descriptor_parser(self._descriptor_files, self._ctx),
            self._destination_path)

    def BuildLayer(self):
        """Override."""
        cached_img = None
        if self._cache:
            with ftl_util.Timing('checking_cached_composer_json_layer'):
                key = self.GetCacheKey()
                cached_img = self._cache.Get(key)
                self._log_cache_result(False if cached_img is None else True)
        if cached_img:
            self.SetImage(cached_img)
        else:
            with ftl_util.Timing('building_composer_json_layer'):
                self._build_layer()
            if self._cache:
                with ftl_util.Timing('uploading_composer_json_layer'):
                    self._cache.Set(self.GetCacheKey(), self.GetImage())

    def _build_layer(self):
        blob, u_blob = self._gen_composer_install_tar(self._directory,
                                                      self._destination_path)
        overrides_dct = {'created': str(datetime.date.today()) + 'T00:00:00Z'}
        self._img = tar_to_dockerimage.FromFSImage([blob], [u_blob],
                                                   overrides_dct)

    def _gen_composer_install_tar(self, app_dir, destination_path):
        vendor_dir = os.path.join(app_dir, 'vendor')
        rm_cmd = ['rm', '-rf', vendor_dir]
        ftl_util.run_command('rm_vendor_dir', rm_cmd)

        composer_install_cmd = [
            'composer', 'install', '--no-dev', '--no-progress', '--no-suggest',
            '--no-interaction'
        ]
        ftl_util.run_command(
            'composer_install',
            composer_install_cmd,
            cmd_cwd=app_dir,
            cmd_env=php_util.gen_composer_env(),
            err_type=ftl_error.FTLErrors.USER())

        vendor_destination = os.path.join(destination_path, 'vendor')
        return ftl_util.zip_dir_to_layer_sha(vendor_dir, vendor_destination)

    def _log_cache_result(self, hit):
        if hit:
            cache_str = constants.PHASE_1_CACHE_HIT
        else:
            cache_str = constants.PHASE_1_CACHE_MISS
        logging.info(
            cache_str.format(
                key_version=constants.CACHE_KEY_VERSION,
                language='PHP',
                key=self.GetCacheKey()))

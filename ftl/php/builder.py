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

from ftl.common import builder
from ftl.common import ftl_util
from ftl.common import single_layer_image
from ftl.common import tar_to_dockerimage

_PHP_NAMESPACE = 'php-package-lock-cache'
_COMPOSER_LOCK = 'composer.lock'
_COMPOSER_JSON = 'composer.json'


class PHP(builder.RuntimeBase):
    def __init__(self, ctx, args, cache_version_str):
        super(PHP, self).__init__(ctx, _PHP_NAMESPACE, args, cache_version_str,
                                  [_COMPOSER_LOCK, _COMPOSER_JSON])

    def Build(self):
        lyr_imgs = []
        lyr_imgs.append(self._base)
        if ftl_util.has_pkg_descriptor(self._descriptor_files, self._ctx):
            pkg = self.PackageLayer(self._ctx, self._descriptor_files, None,
                                    self._args.destination_path)
            cached_pkg_img = self._cash.GetAndCheckTTL(self._base,
                                                       self._namespace,
                                                       pkg.GetCacheKey())
            if cached_pkg_img is not None:
                pkg.SetImage(cached_pkg_img)
            else:
                pkg.BuildLayer()
                self._cash.Store(self._base, self._namespace,
                                 pkg.GetCacheKey(), pkg.GetImage())
            lyr_imgs.append(pkg)

        app = self.AppLayer(self._ctx, self._args.destination_path)
        app.BuildLayer()
        lyr_imgs.append(app)
        ftl_image = self.AppendLayersIntoImage(lyr_imgs)
        self.StoreImage(ftl_image)

    class PackageLayer(single_layer_image.CacheLayer):
        def __init__(self, ctx, descriptor_files, pkg_descriptor,
                     destination_path):
            super(PHP.PackageLayer, self).__init__()
            self._ctx = ctx
            self._descriptor_files = descriptor_files
            self._pkg_descriptor = pkg_descriptor
            self._destination_path = destination_path

        def GetCacheKeyRaw(self):
            return ftl_util.descriptor_parser(self._descriptor_files,
                                              self._ctx)

        def BuildLayer(self):
            """Override."""
            lyr, sha = self._gen_composer_install_tar(self._pkg_descriptor,
                                                      self._destination_path)
            logging.info('Generated layer with sha: %s', sha)
            self._img = tar_to_dockerimage.FromFSImage(lyr, {
                "creation_time":
                str(datetime.date.today()) + "T00:00:00Z"
            })

        def _gen_composer_install_tar(self, pkg_descriptor, destination_path):
            # Create temp directory to write package descriptor to
            pkg_dir = tempfile.mkdtemp()
            app_dir = os.path.join(pkg_dir, destination_path.strip("/"))
            os.makedirs(app_dir)

            # Copy out the relevant package descriptors to a tempdir.
            ftl_util.descriptor_copy(self._ctx, self._descriptor_files,
                                     app_dir)

            subprocess.check_call(
                ['rm', '-rf', os.path.join(app_dir, 'vendor')])

            with ftl_util.Timing("composer_install"):
                if pkg_descriptor is None:
                    subprocess.check_call(
                        ['composer', 'install', '--no-dev', '--no-scripts'],
                        cwd=app_dir)
                else:
                    subprocess.check_call(
                        [
                            'composer', 'install', '--no-dev', '--no-scripts',
                            pkg_descriptor
                        ],
                        cwd=app_dir)
            return ftl_util.zip_dir_to_layer_sha(pkg_dir)

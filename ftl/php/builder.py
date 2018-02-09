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
import json

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

    def _parse_composer_pkgs(self):
        descriptor_contents = ftl_util.descriptor_parser(
            self._descriptor_files, self._ctx)
        composer_json = json.loads(descriptor_contents)
        pkgs = []
        for k, v in composer_json['require'].iteritems():
            pkgs.append((k, v))
        return pkgs

    def Build(self):
        lyr_imgs = []
        lyr_imgs.append(self._base_image)
        if ftl_util.has_pkg_descriptor(self._descriptor_files, self._ctx):
            pkgs = self._parse_composer_pkgs()
            # if there are 42 or more packages, revert to using phase 1
            if len(pkgs) > 41:
                pkgs = [None]
            for pkg_txt in pkgs:
                pkg = self.PackageLayer(self._ctx, self._descriptor_files,
                                        pkg_txt, self._args.destination_path,
                                        self._args.entrypoint)
                cached_pkg_img = None
                if self._args.cache:
                    with ftl_util.Timing("checking cached pkg layer"):
                        cached_pkg_img = self._cache.Get(pkg.GetCacheKey())
                if cached_pkg_img is not None:
                    pkg.SetImage(cached_pkg_img)
                else:
                    with ftl_util.Timing("building pkg layer"):
                        pkg.BuildLayer()
                        # keep track of mappings for new cache entries only
                        mapping = pkg.GetCacheMapping()                        
                        self._cache_mappings[mapping[0]] = mapping[1]
                    if self._args.cache:
                        with ftl_util.Timing("uploading pkg layer"):
                            self._cache.Set(pkg.GetCacheKey(), pkg.GetImage())
                lyr_imgs.append(pkg)

        app = self.AppLayer(self._ctx, self._args.destination_path)
        with ftl_util.Timing("builder app layer"):
            app.BuildLayer()
        lyr_imgs.append(app)
        with ftl_util.Timing("stitching layers into final image"):
            ftl_image = self.AppendLayersIntoImage(lyr_imgs)
        with ftl_util.Timing("uploading final image"):
            self.StoreImage(ftl_image)

    class PackageLayer(single_layer_image.CacheableLayer):
        def __init__(self, ctx, descriptor_files, pkg_descriptor,
                     destination_path, entrypoint):
            super(PHP.PackageLayer, self).__init__()
            self._ctx = ctx
            self._descriptor_files = descriptor_files
            self._pkg_descriptor = pkg_descriptor
            self._destination_path = destination_path
            self._entrypoint = entrypoint

        def GetCacheMapping(self):
            return (self.GetCacheKeyRaw(), self.GetCacheKey())

        def GetCacheKeyRaw(self):
            if self._pkg_descriptor is not None:
                # phase 2 cache key
                return self._pkg_descriptor[0] + ' ' + self._pkg_descriptor[1]
            # phase 1 cache key
            return ftl_util.descriptor_parser(self._descriptor_files,
                                              self._ctx)

        def BuildLayer(self):
            """Override."""
            blob, u_blob = self._gen_composer_install_tar(
                self._pkg_descriptor, self._destination_path)
            overrides_dct = {
                    'created': str(datetime.date.today()) + "T00:00:00Z"
                }
            if self._entrypoint:
                overrides_dct['Entrypoint'] = self._entrypoint
            self._img = tar_to_dockerimage.FromFSImage(
                [blob], [u_blob], overrides_dct)

        def _gen_composer_install_tar(self, pkg_descriptor, destination_path):
            # Create temp directory to write package descriptor to
            pkg_dir = tempfile.mkdtemp()
            app_dir = os.path.join(pkg_dir, destination_path.strip("/"))
            os.makedirs(app_dir)

            # Copy out the relevant package descriptors to a tempdir.
            if pkg_descriptor is None:
                # phase 1 copy whole descriptor
                ftl_util.descriptor_copy(self._ctx, self._descriptor_files,
                                         app_dir)

            subprocess.check_call(
                ['rm', '-rf', os.path.join(app_dir, 'vendor')])

            with ftl_util.Timing("composer_install"):
                if pkg_descriptor is None:
                    # phase 1 install entire descriptor
                    subprocess.check_call(
                        ['composer', 'install', '--no-dev', '--no-scripts'],
                        cwd=app_dir)
                else:
                    pkg, version = pkg_descriptor
                    subprocess.check_call(
                        ['composer', 'require',
                         str(pkg),
                         str(version)],
                        cwd=app_dir)
            return ftl_util.zip_dir_to_layer_sha(pkg_dir)

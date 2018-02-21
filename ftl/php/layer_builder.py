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

import os
import subprocess
import tempfile
import datetime

from ftl.common import constants
from ftl.common import ftl_util
from ftl.common import single_layer_image
from ftl.common import tar_to_dockerimage


class LayerBuilder(single_layer_image.CacheableLayerBuilder):
    def __init__(self,
                 ctx=None,
                 descriptor_files=None,
                 pkg_descriptor=None,
                 destination_path=constants.DEFAULT_DESTINATION_PATH):
        super(LayerBuilder, self).__init__()
        self._ctx = ctx
        self._descriptor_files = descriptor_files
        self._pkg_descriptor = pkg_descriptor
        self._destination_path = destination_path

    def GetCacheKeyRaw(self):
        if self._pkg_descriptor is not None:
            # phase 2 cache key
            return "%s %s %s" % (self._pkg_descriptor[0],
                                 self._pkg_descriptor[1],
                                 self._destination_path)
        # phase 1 cache key
        return "%s %s" % (
            ftl_util.descriptor_parser(self._descriptor_files, self._ctx),
            self._destination_path)

    def BuildLayer(self):
        """Override."""
        blob, u_blob = self._gen_composer_install_tar(self._pkg_descriptor,
                                                      self._destination_path)
        overrides_dct = {'created': str(datetime.date.today()) + "T00:00:00Z"}
        self._img = tar_to_dockerimage.FromFSImage([blob], [u_blob],
                                                   overrides_dct)

    def _gen_composer_install_tar(self, pkg_descriptor, destination_path):
        # Create temp directory to write package descriptor to
        pkg_dir = tempfile.mkdtemp()
        app_dir = os.path.join(pkg_dir, destination_path.strip("/"))
        print 'app_dir: %s' % app_dir
        os.makedirs(app_dir)

        # Copy out the relevant package descriptors to a tempdir.
        if pkg_descriptor is None:
            # phase 1 copy whole descriptor
            ftl_util.descriptor_copy(self._ctx, self._descriptor_files,
                                     app_dir)

        subprocess.check_call(['rm', '-rf', os.path.join(app_dir, 'vendor')])

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
                     str(pkg), str(version)],
                    cwd=app_dir)
        return ftl_util.zip_dir_to_layer_sha(pkg_dir)

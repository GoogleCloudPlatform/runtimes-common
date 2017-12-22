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

_PHP_NAMESPACE = 'php-composer-lock-cache'
_COMPOSER_LOCK = 'composer.lock'
_COMPOSER_JSON = 'composer.json'


class PHP(builder.JustApp):
    def __init__(self, ctx):
        self.descriptor_files = [_COMPOSER_LOCK, _COMPOSER_JSON]
        self.namespace = _PHP_NAMESPACE
        super(PHP, self).__init__(ctx)

    def __enter__(self):
        """Override."""
        return self

    def _generate_overrides(self):
        return metadata.Overrides(
            creation_time=str(datetime.date.today()) + "T00:00:00Z")

    def CreatePackageBase(self, base, destination_path="srv"):
        """Override."""
        overrides = self._generate_overrides()
        layer, sha = self._gen_package_tar(destination_path)
        logging.info('Generated layer with sha: %s', sha)

        with append.Layer(
                base, layer, diff_id=sha, overrides=overrides) as dep_image:
            return dep_image

    def _gen_package_tar(self, destination_path):
        # Create temp directory to write package descriptor to
        pkg_dir = tempfile.mkdtemp()
        app_dir = os.path.join(pkg_dir, destination_path.strip("/"))
        os.makedirs(app_dir)

        # Copy out the relevant package descriptors to a tempdir.
        for f in self.descriptor_files:
            if self._ctx.Contains(f):
                with open(os.path.join(app_dir, f), 'w') as w:
                    w.write(self._ctx.GetFile(f))

        subprocess.check_call(['rm', '-rf', os.path.join(app_dir, 'vendor')])

        with ftl_util.Timing("composer_install"):
            subprocess.check_call(
                ['composer', 'install', '--no-dev', '--no-scripts'],
                cwd=app_dir)

        return ftl_util.zip_dir_to_layer_sha(pkg_dir)


def From(ctx):
    return PHP(ctx)

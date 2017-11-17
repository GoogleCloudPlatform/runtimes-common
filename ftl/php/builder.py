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

import hashlib
import os
import subprocess
import tempfile
import logging
import datetime

from containerregistry.client.v2_2 import append
from containerregistry.transform.v2_2 import metadata

from ftl.common import builder

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

    def CreatePackageBase(self, base_image, cache,
                          destination_path="/workspace/app"):
        """Override."""
        overrides = self._generate_overrides()
        layer, sha = self._gen_package_tar(destination_path)
        logging.info('Generated layer with sha: %s', sha)

        with append.Layer(
                base_image, layer, diff_id=sha,
                overrides=overrides) as dep_image:
            return dep_image

    def _gen_package_tar(self, destination_path):
        # Create temp directory to write package descriptor to
        tmp = tempfile.mkdtemp()
        app_dir = os.path.join(tmp, destination_path.strip("/"))
        os.makedirs(app_dir)

        # Copy out the relevant package descriptors to a tempdir.
        for f in self.descriptor_files:
            if self._ctx.Contains(f):
                with open(os.path.join(app_dir, f), 'w') as w:
                    w.write(self._ctx.GetFile(f))

        tar_path = tempfile.mktemp()
        logging.info('Starting composer install ...')
        subprocess.check_call(
            ['composer', 'install', '--no-dev', '--no-scripts'], cwd=app_dir)
        logging.info('Finished composer install.')

        logging.info('Starting to tar composer packages...')
        subprocess.check_call(['tar', '-C', tmp, '-cf', tar_path, '.'])
        logging.info('Finished generating tarfile for composer packages.')

        # We need the sha of the unzipped and zipped tarball.
        # So for performance, tar, sha, zip, sha.
        # We use gzip for performance instead of python's zip.
        sha = 'sha256:' + hashlib.sha256(open(tar_path).read()).hexdigest()

        logging.info('Starting to gzip composer package tarfile...')
        subprocess.check_call(['gzip', tar_path])
        logging.info('Finished generating gzip composer package tarfile.')
        return open(os.path.join(tmp, tar_path + '.gz'), 'rb').read(), sha


def From(ctx):
    return PHP(ctx)

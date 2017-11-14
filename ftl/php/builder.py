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

from containerregistry.client.v2_2 import append

from ftl.common import builder

_PHP_NAMESPACE = 'php-composer-lock-cache'
_COMPOSER_LOCK = 'composer.lock'
_COMPOSER_JSON = 'composer.json'


class PHP(builder.JustApp):
    def __init__(self, ctx):
        self._overrides = None
        super(PHP, self).__init__(ctx)

    def __enter__(self):
        """Override."""
        return self

    def CreatePackageBase(self, base_image, cache, use_cache=True):
        """Override."""

        descriptor = None
        for f in [_COMPOSER_LOCK, _COMPOSER_JSON]:
            if self._ctx.Contains(f):
                descriptor = f
                descriptor_contents = self._ctx.GetFile(f)
                break

        if not descriptor:
            logging.info('No package descriptor found. No packages installed.')

            return append.Layer(base_image, tar_gz=None)

        checksum = hashlib.sha256(descriptor_contents).hexdigest()
        if use_cache:
            hit = cache.Get(base_image, _PHP_NAMESPACE, checksum)
            if hit:
                logging.info('Found cached dependency layer for %s' % checksum)
                return hit
            else:
                logging.info('No cached dependency layer for %s' % checksum)
        else:
            logging.info('Skipping checking cache for dependency layer %s'
                         % checksum)

        layer, sha = self._gen_package_tar()

        with append.Layer(
          base_image, layer, diff_id=sha) as dep_image:
            if use_cache:
                logging.info('Storing layer %s in cache.', sha)
                cache.Store(base_image, _PHP_NAMESPACE, checksum, dep_image)
            else:
                logging.info('Skipping storing layer %s in cache.', sha)
            return dep_image

    def _gen_package_tar(self):
        # Create temp directory to write package descriptor to

        tmp = tempfile.mkdtemp()
        app_dir = os.path.join(tmp, 'app')
        os.mkdir(app_dir)

        # Copy out the relevant package descriptors to a tempdir.
        for f in [_COMPOSER_LOCK, _COMPOSER_JSON]:
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

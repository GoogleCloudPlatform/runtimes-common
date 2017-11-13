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
import json
import datetime

from containerregistry.client.v2_2 import append
from containerregistry.transform.v2_2 import metadata

from ftl.common import builder

_PYTHON_NAMESPACE = 'python-requirements-cache'
_REQUIREMENTS_TXT = 'requirements.txt'
_DEFAULT_TTL_WEEKS = 1


class Python(builder.JustApp):
    def __init__(self, ctx):
        self._overrides = None
        super(Python, self).__init__(ctx)

    def __enter__(self):
        """Override."""
        return self

    def CreatePackageBase(self, base_image, cache, use_cache=True):
        """Override."""
        # Figure out if we need to override entrypoint.
        # Save the overrides for later to avoid adding an extra layer.
        descriptor = None
        if self._ctx.Contains(_REQUIREMENTS_TXT):
            descriptor = _REQUIREMENTS_TXT
            descriptor_contents = self._ctx.GetFile(_REQUIREMENTS_TXT)

        if not descriptor:
            logging.info('No requirements.txt found.  No packages installed.')
            return base_image
        overrides = metadata.Overrides(
            creation_time=str(datetime.date.today()) + "T00:00:00Z")
        checksum = hashlib.sha256(descriptor_contents).hexdigest()
        if use_cache:
            hit = cache.Get(base_image, _PYTHON_NAMESPACE, checksum)
            if hit:
                logging.info('Found cached dependency layer for %s' % checksum)
                last_created = _timestamp_to_time(_creation_time(hit))
                now = datetime.datetime.now()
                if last_created > now - datetime.timedelta(
                        seconds=_DEFAULT_TTL_WEEKS):
                    return hit
                else:
                    logging.info('TTL expired for cached image, rebuilding %s'
                                 % checksum)
            else:
                logging.info('No cached dependency layer for %s' % checksum)
        else:
            logging.info(
                'Skipping checking cache for dependency layer %s' % checksum)

        layer, sha = self._gen_package_tar(descriptor, descriptor_contents)

        with append.Layer(
                base_image, layer, diff_id=sha,
                overrides=overrides) as dep_image:
            if use_cache:
                logging.info('Storing layer %s in cache.', sha)
                cache.Store(base_image, _PYTHON_NAMESPACE, checksum, dep_image)
            else:
                logging.info('Skipping storing layer %s in cache.', sha)
            return dep_image

    def _gen_package_tar(self, descriptor, descriptor_contents):
        tmp_app = tempfile.mkdtemp()
        tmp_venv = tempfile.mkdtemp()

        tmp_app = os.path.join(tmp_app, 'app')
        venv_dir = os.path.join(tmp_venv, 'env')
        os.makedirs(tmp_app)
        os.makedirs(venv_dir)

        # Copy out the relevant package descriptors to a tempdir.
        with open(os.path.join(tmp_app, descriptor), 'w') as f:
            f.write(descriptor_contents)

        tar_path = tempfile.mktemp()
        logging.info('Starting venv creation ...')

        # TODO(aaron-prindle) add support for different python versions
        subprocess.check_call(
            ['virtualenv', '--no-download', venv_dir, '-p', 'python3.6'],
            cwd=tmp_app)
        os.environ['VIRTUAL_ENV'] = venv_dir
        os.environ['PATH'] = venv_dir + "/bin" + ":" + os.environ['PATH']
        # bazel adds its own PYTHONPATH to the env
        # which must be removed for the pip calls to work properly
        my_env = os.environ.copy()
        my_env.pop('PYTHONPATH', None)

        subprocess.check_call(
            ['pip', 'install', '-r', 'requirements.txt'],
            cwd=tmp_app,
            env=my_env)
        logging.info('Finished pip install.')

        logging.info('Starting to tar pip packages...')
        subprocess.check_call(['tar', '-C', tmp_venv, '-cf', tar_path, '.'])
        logging.info('Finished generating tarfile for pip packages...')

        # We need the sha of the unzipped and zipped tarball.
        # So for performance, tar, sha, zip, sha.
        # We use gzip for performance instead of python's zip.
        sha = 'sha256:' + hashlib.sha256(open(tar_path).read()).hexdigest()

        logging.info('Starting to gzip pip package tarfile...')
        subprocess.check_call(['gzip', tar_path])
        logging.info('Finished generating gzip pip package tarfile.')
        return open(os.path.join(tmp_venv, tar_path + '.gz'), 'rb').read(), sha


def From(ctx):
    return Python(ctx)


def _creation_time(image):
    logging.info(image.config_file())
    cfg = json.loads(image.config_file())
    return cfg.get('created')


def _timestamp_to_time(dt_str):
    dt = dt_str.rstrip("Z")
    return datetime.datetime.strptime(dt, "%Y-%m-%dT%H:%M:%S")

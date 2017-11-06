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

# _DEFAULT_ENTRYPOINT = 'node server.js'


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
        for f in [_REQUIREMENTS_TXT]:
            if self._ctx.Contains(f):
                descriptor = f
                descriptor_contents = self._ctx.GetFile(f)
                break

        if not descriptor:
            logging.info('No package descriptor found. No packages installed.')
        overrides = metadata.Overrides(
            creation_time=datetime.datetime.now().isoformat())
        checksum = hashlib.sha256(descriptor_contents).hexdigest()
        if use_cache:
            hit = cache.Get(base_image, _PYTHON_NAMESPACE, checksum)
            if hit:
                logging.info('Found cached dependency layer for %s' % checksum)
                # TODO(aaron-prindle) check that cached dep layer is in TTL
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
            logging.info(dep_image.config_file())
            logging.info(_creation_time(dep_image))
            return dep_image

    def _gen_package_tar(self, descriptor, descriptor_contents):
        # python2.7 packages installed w/ pip default to /usr/lib/python2.7
        # /usr/lib/python2.7/dist-packages
        tmp = tempfile.mkdtemp()
        # TODO(aaron-prindle) make the python version detected from the base image
        app_dir = os.path.join(tmp, 'usr', 'lib', 'python2.7', 'dist-packages')
        os.makedirs(app_dir)

        # Copy out the relevant package descriptors to a tempdir.
        with open(os.path.join(app_dir, descriptor), 'w') as f:
            f.write(descriptor_contents)

        tar_path = tempfile.mktemp()
        #TODO(aaron-prindle) verify check_gcp_build substitute not needed
        subprocess.check_call(
            ['pip', 'install', '-r', 'requirements.txt', '--target', app_dir],
            cwd=app_dir)
        subprocess.check_call(['tar', '-C', tmp, '-cf', tar_path, '.'])

        # We need the sha of the unzipped and zipped tarball.
        # So for performance, tar, sha, zip, sha.
        # We use gzip for performance instead of python's zip.
        sha = 'sha256:' + hashlib.sha256(open(tar_path).read()).hexdigest()
        subprocess.check_call(['gzip', tar_path])
        return open(os.path.join(tmp, tar_path + '.gz'), 'rb').read(), sha


# def check_gcp_build(package_json, app_dir):
#     scripts = package_json.get('scripts', {})
#     gcp_build = scripts.get('gcp-build')

#     if not gcp_build:
#         return

#     env = os.environ.copy()
#     env["PYTHON_ENV"] = "development"
#     subprocess.check_call(['npm', 'install'], cwd=app_dir, env=env)
#     subprocess.check_call(['npm', 'run-script', 'gcp-build'],
#                           cwd=app_dir, env=env)


def From(ctx):
    return Python(ctx)


def _creation_time(image):
    cfg = json.loads(image.config_file())
    return cfg['created']


def _timestamp_to_time(dt_str):
    dt, _, us = dt_str.partition(".")
    dt = datetime.datetime.strptime(dt, "%Y-%m-%dT%H:%M:%S")
    us = int(us.rstrip("Z"), 10)
    return dt + datetime.timedelta(microseconds=us)

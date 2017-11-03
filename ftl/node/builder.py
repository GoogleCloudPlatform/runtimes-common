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

from containerregistry.client.v2_2 import append
from containerregistry.transform.v2_2 import metadata

from ftl.common import builder

_NODE_NAMESPACE = 'node-package-lock-cache'
_PACKAGE_LOCK = 'package-lock.json'
_PACKAGE_JSON = 'package.json'

_DEFAULT_ENTRYPOINT = 'node server.js'


class Node(builder.JustApp):
    def __init__(self, ctx):
        self._overrides = None
        super(Node, self).__init__(ctx)

    def __enter__(self):
        """Override."""
        return self

    def CreatePackageBase(self, base_image, cache, use_cache=True):
        """Override."""
        # Figure out if we need to override entrypoint.
        # Save the overrides for later to avoid adding an extra layer.
        pj_contents = {}
        if self._ctx.Contains(_PACKAGE_JSON):
            pj_contents = json.loads(self._ctx.GetFile(_PACKAGE_JSON))
        entrypoint = parse_entrypoint(pj_contents)
        overrides = metadata.Overrides(entrypoint=entrypoint)

        descriptor = None
        for f in [_PACKAGE_LOCK, _PACKAGE_JSON]:
            if self._ctx.Contains(f):
                descriptor = f
                descriptor_contents = self._ctx.GetFile(f)
                break

        if not descriptor:
            logging.info('No package descriptor found. No packages installed.')

            # Add the overrides now.
            return append.Layer(
                base_image, tar_gz=None, overrides=overrides)

        checksum = hashlib.sha256(descriptor_contents).hexdigest()
        if use_cache:
            hit = cache.Get(base_image, _NODE_NAMESPACE, checksum)
            if hit:
                logging.info('Found cached dependency layer for %s' % checksum)
                return hit
            else:
                logging.info('No cached dependency layer for %s' % checksum)
        else:
            logging.info('Skipping checking cache for dependency layer %s'
                         % checksum)

        layer, sha = self._gen_package_tar(descriptor, descriptor_contents)

        with append.Layer(
          base_image, layer, diff_id=sha, overrides=overrides) as dep_image:
            if use_cache:
                logging.info('Storing layer %s in cache.', sha)
                cache.Store(base_image, _NODE_NAMESPACE, checksum, dep_image)
            else:
                logging.info('Skipping storing layer %s in cache.', sha)
            return dep_image

    def _gen_package_tar(self, descriptor, descriptor_contents):
        # We want the node_modules directory rooted at /app/node_modules in
        # the final image.
        # So we build a hierarchy like:
        # /$tmp/app/node_modules
        # And use the -C flag to tar to root the tarball at /$tmp.
        tmp = tempfile.mkdtemp()
        app_dir = os.path.join(tmp, 'app')
        os.mkdir(app_dir)

        # Copy out the relevant package descriptors to a tempdir.
        with open(os.path.join(app_dir, descriptor), 'w') as f:
            f.write(descriptor_contents)

        tar_path = tempfile.mktemp()
        check_gcp_build(json.loads(self._ctx.GetFile(_PACKAGE_JSON)), app_dir)
        subprocess.check_call(['rm', '-rf',
                              os.path.join(app_dir, 'node_modules')])
        subprocess.check_call(['npm', 'install', '--production', '--no-cache'],
                              cwd=app_dir)
        subprocess.check_call(['tar', '-C', tmp, '-cf', tar_path, '.'])

        # We need the sha of the unzipped and zipped tarball.
        # So for performance, tar, sha, zip, sha.
        # We use gzip for performance instead of python's zip.
        sha = 'sha256:' + hashlib.sha256(open(tar_path).read()).hexdigest()
        subprocess.check_call(['gzip', tar_path])
        return open(os.path.join(tmp, tar_path + '.gz'), 'rb').read(), sha


def check_gcp_build(package_json, app_dir):
    scripts = package_json.get('scripts', {})
    gcp_build = scripts.get('gcp-build')

    if not gcp_build:
        return

    env = os.environ.copy()
    env["NODE_ENV"] = "development"
    subprocess.check_call(['npm', 'install'], cwd=app_dir, env=env)
    subprocess.check_call(['npm', 'run-script', 'gcp-build'],
                          cwd=app_dir, env=env)


def From(ctx):
    return Node(ctx)


def parse_entrypoint(package_json):
    entrypoint = []

    scripts = package_json.get('scripts', {})
    start = scripts.get('start', _DEFAULT_ENTRYPOINT)
    prestart = scripts.get('prestart')

    if prestart:
        entrypoint = '%s && %s' % (prestart, start)
    else:
        entrypoint = start
    return ['sh', '-c', entrypoint]

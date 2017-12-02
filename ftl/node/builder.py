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
import json
import datetime

from containerregistry.client.v2_2 import append
from containerregistry.transform.v2_2 import metadata

from ftl.common import builder
from ftl.common import ftl_util

_NODE_NAMESPACE = 'node-package-lock-cache'
_PACKAGE_LOCK = 'package-lock.json'
_PACKAGE_JSON = 'package.json'

_DEFAULT_ENTRYPOINT = 'node server.js'


class Node(builder.JustApp):
    def __init__(self, ctx):
        self.descriptor_files = [_PACKAGE_LOCK, _PACKAGE_JSON]
        self.namespace = _NODE_NAMESPACE
        super(Node, self).__init__(ctx)

    def __enter__(self):
        """Override."""
        return self

    def _generate_overrides(self):
        pj_contents = {}
        if self._ctx.Contains(_PACKAGE_JSON):
            pj_contents = json.loads(self._ctx.GetFile(_PACKAGE_JSON))
        entrypoint = parse_entrypoint(pj_contents)
        return metadata.Overrides(
            creation_time=str(datetime.date.today()) + "T00:00:00Z",
            entrypoint=entrypoint)

    def CreatePackageBase(self, base, destination_path="/app"):
        """Override."""
        overrides = self._generate_overrides()

        layer, sha = self._gen_package_tar(destination_path)
        logging.info('Generated layer with sha: %s', sha)

        with append.Layer(
                base, layer, diff_id=sha, overrides=overrides) as dep_image:
            return dep_image

    def _gen_package_tar(self, destination_path):
        # We want the node_modules directory rooted at:
        # $destination_path/node_modules in the final image.
        # So we build a hierarchy like:
        # /$tmp/$destination_path/node_modules
        # And use the -C flag to tar to root the tarball at /$tmp.
        pkg_dir = tempfile.mkdtemp()
        app_dir = os.path.join(pkg_dir, destination_path.strip("/"))
        os.makedirs(app_dir)

        # Copy out the relevant package descriptors to a tempdir.
        for f in [_PACKAGE_LOCK, _PACKAGE_JSON]:
            if self._ctx.Contains(f):
                with open(os.path.join(app_dir, f), 'w') as w:
                    w.write(self._ctx.GetFile(f))

        check_gcp_build(json.loads(self._ctx.GetFile(_PACKAGE_JSON)), app_dir)
        subprocess.check_call(
            ['rm', '-rf', os.path.join(app_dir, 'node_modules')])
        with ftl_util.Timing("npm_install"):
            subprocess.check_call(
                ['npm', 'install', '--production'], cwd=app_dir)

        return ftl_util.folder_to_layer_sha(pkg_dir, "npm")


def check_gcp_build(package_json, app_dir):
    scripts = package_json.get('scripts', {})
    gcp_build = scripts.get('gcp-build')

    if not gcp_build:
        return

    env = os.environ.copy()
    env["NODE_ENV"] = "development"
    subprocess.check_call(['npm', 'install'], cwd=app_dir, env=env)
    subprocess.check_call(
        ['npm', 'run-script', 'gcp-build'], cwd=app_dir, env=env)


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

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
import errno


from containerregistry.client.v2_2 import append
from ftl.common import builder


_PYTHON_NAMESPACE = 'python-requirements-cache'
_REQUIREMENTS_TXT = 'requirements.txt'


class Python(builder.JustApp):

    def __init__(self, ctx):
        super(Python, self).__init__(ctx)

    def __enter__(self):
        """Override."""
        return self

    def CreatePackageBase(self, base_image, cache):
        """Override."""
        if self._ctx.Contains(_REQUIREMENTS_TXT):
            print('Using %s as a package descriptor.' % _REQUIREMENTS_TXT)
            descriptor = self._ctx.GetFile(_REQUIREMENTS_TXT)
        else:
            print('No package descriptor found. Not installing anything.')
            return base_image

        checksum = hashlib.sha256(descriptor).hexdigest()
        hit = cache.Get(base_image, _PYTHON_NAMESPACE, checksum)
        if hit:
            print('Found cached dependency layer for %s' % checksum)
            return hit
        else:
            print('No cached dependency layer for %s' % checksum)

        # We want the node_modules directory rooted at /app/node_modules in
        # the final image.
        # So we build a hierarchy like:
        # /$tmpdir/app/node_modules
        # And use the -C flag to tar to root the tarball at /$tmpdir.
        tmpdir = tempfile.mkdtemp()
        # TODO(aaron-prindle) make the python$VERSION detected from the base image
        app_dir = os.path.join(tmpdir, 'usr', 'lib', 'python2.7', 'dist-packages')
        os.makedirs(app_dir)

        # Copy out the relevant package descriptors to a tempdir.
        if self._ctx.Contains(_REQUIREMENTS_TXT):
            with open(os.path.join(app_dir, _REQUIREMENTS_TXT), 'w') as f:
                f.write(self._ctx.GetFile(_REQUIREMENTS_TXT))

        tar_path = tempfile.mktemp()
        subprocess.check_call(['pip', 'install',
        '-r', 'requirements.txt',
        '--target', app_dir], cwd=app_dir)
        subprocess.check_call([
            'tar',
            '-C', tmpdir,
            '-cf', tar_path,
            '.'
        ])

        # We need the sha of the unzipped and zipped tarball.
        # So for performance, tar, sha, zip, sha.
        # We use gzip for performance instead of python's zip.
        sha = 'sha256:' + hashlib.sha256(open(tar_path).read()).hexdigest()
        subprocess.check_call(['gzip', tar_path])
        layer = open(os.path.join(tmpdir, tar_path + '.gz'), 'rb').read()

        with append.Layer(base_image, layer, diff_id=sha) as dep_image:
            print('Storing layer %s in cache.', sha)
            cache.Store(base_image, _PYTHON_NAMESPACE, checksum, dep_image)
            return dep_image


def From(ctx):
    return Python(ctx)

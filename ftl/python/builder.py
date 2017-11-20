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

_PYTHON_NAMESPACE = 'python-requirements-cache'
_REQUIREMENTS_TXT = 'requirements.txt'


class Python(builder.JustApp):
    def __init__(self, ctx):
        self.descriptor_files = [_REQUIREMENTS_TXT]
        self.namespace = _PYTHON_NAMESPACE
        self.ctx = ctx
        super(Python, self).__init__(ctx)

    def __enter__(self):
        """Override."""
        return self

    def _generate_overrides(self):
        return metadata.Overrides(
            creation_time=str(datetime.date.today()) + "T00:00:00Z")

    def CreatePackageBase(self, base_image):
        """Override."""
        descriptor = None
        for f in self.descriptor_files:
            if self.ctx.Contains(f):
                descriptor = f
                descriptor_contents = self.ctx.GetFile(descriptor)
                logging.info(descriptor_contents)
                for pkg_txt in descriptor_contents.splitlines():
                    if len(pkg_txt) <= 1 or pkg_txt[0] == '#':
                        continue
                    logging.info("passing: " + pkg_txt + "to CreatePackageBaseHelper")
                    base_image = self.CreatePackageBaseHelper(base_image,
                    pkg_txt)

        if not descriptor:
            logging.info('No package descriptor found. No packages installed.')
            return None

        return base_image

    def CreatePackageBaseHelper(self, base_image, pkg_txt=None):
        """Override."""
        overrides = self._generate_overrides()

        layer, sha = self._gen_package_tar(pkg_txt)
        logging.info('Generated layer with sha: %s', sha)

        with append.Layer(
                base_image, layer, diff_id=sha,
                overrides=overrides) as dep_image:
            return dep_image

    def _gen_package_tar(self, pkg_txt):
        tmp_app = tempfile.mkdtemp()
        tmp_venv = tempfile.mkdtemp()

        tmp_app = os.path.join(tmp_app, 'app')
        venv_dir = os.path.join(tmp_venv, 'env')
        wheel_dir = os.path.join(tmp_venv, 'env', 'wheel')
        os.makedirs(tmp_app)
        os.makedirs(wheel_dir)

        # Copy out the relevant package descriptors to a tempdir.
        for f in self.descriptor_files:
            # if self._ctx.Contains(f):
            with open(os.path.join(tmp_app, f), 'w') as w:
                w.write(self._ctx.GetFile(f))

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

        args = ['pip', 'wheel',
            '-w', wheel_dir,
            '-r', "/dev/stdin"]

        logging.info("pip requirements input: " + pkg_txt)
        pipe1 = subprocess.Popen(args,
                                stdin=subprocess.PIPE,
                                stdout=subprocess.PIPE,
                                stderr=subprocess.PIPE,
                                env=my_env,
                                cwd=tmp_app)
        result = pipe1.communicate(input=pkg_txt)[0]
        logging.info("pip output: " + str(result))
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

        self._cleanup()
        return open(os.path.join(tmp_venv, tar_path + '.gz'), 'rb').read(), sha

    def _cleanup(self):
        subprocess.check_call(['rm', '-rf', '/env'])

def From(ctx):
    return Python(ctx)

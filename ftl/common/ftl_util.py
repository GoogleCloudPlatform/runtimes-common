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
"""This package defines helpful utilities for FTL ."""
import os
import time
import logging
import subprocess
import hashlib
import tempfile


class Timing(object):
    def __init__(self, descriptor):
        self.descriptor = descriptor

    def __enter__(self):
        self.start = time.time()
        return self

    def __exit__(self, unused_type, unused_value, unused_traceback):
        end = time.time()
        logging.info('%s took %d seconds', self.descriptor, end - self.start)


def folder_to_layer_sha(pkg_dir, pkg_mngr_name):
    tar_path = tempfile.mktemp()
    with Timing("tar_%s_package" % pkg_mngr_name):
        subprocess.check_call(['tar', '-C', pkg_dir, '-cf', tar_path, '.'])

    # We need the sha of the unzipped and zipped tarball.
    # So for performance, tar, sha, zip, sha.
    # We use gzip for performance instead of python's zip.
    sha = 'sha256:' + hashlib.sha256(open(tar_path).read()).hexdigest()

    with Timing("gzip_%s_tar" % pkg_mngr_name):
        subprocess.check_call(['gzip', tar_path, '-1'])
    return open(os.path.join(pkg_dir, tar_path + '.gz'), 'rb').read(), sha

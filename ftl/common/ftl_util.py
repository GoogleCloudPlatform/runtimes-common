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
import datetime
import json

_OVERRIDES = ["creation_time", "entrypoint", "env"]


def CfgDctToOverrides(config_dct):
    output = {}
    for k, v in config_dct.iteritems():
        if k in _OVERRIDES:
            output[k] = v
    return output


class Timing(object):
    def __init__(self, descriptor):
        self.descriptor = descriptor

    def __enter__(self):
        self.start = time.time()
        return self

    def __exit__(self, unused_type, unused_value, unused_traceback):
        end = time.time()
        logging.info('%s took %d seconds', self.descriptor, end - self.start)


def zip_dir_to_layer_sha(pkg_dir):
    tar_path = tempfile.mktemp()
    with Timing("tar_runtime_package"):
        subprocess.check_call(['tar', '-C', pkg_dir, '-cf', tar_path, '.'])

    # We need the sha of the unzipped and zipped tarball.
    # So for performance, tar, sha, zip, sha.
    # We use gzip for performance instead of python's zip.
    sha = 'sha256:' + hashlib.sha256(open(tar_path).read()).hexdigest()

    with Timing("gzip_runtime_tar"):
        subprocess.check_call(['gzip', tar_path, '-1'])
    return open(os.path.join(pkg_dir, tar_path + '.gz'), 'rb').read(), sha


def has_pkg_descriptor(descriptor_files, ctx):
    for f in descriptor_files:
        if ctx.Contains(f):
            return True
    return False


def descriptor_parser(descriptor_files, ctx):
    descriptor = None
    for f in descriptor_files:
        if ctx.Contains(f):
            descriptor = f
            descriptor_contents = ctx.GetFile(descriptor)
            break
    if not descriptor:
        logging.info('No package descriptor found. No packages installed.')
        return None
    return descriptor_contents


def descriptor_copy(ctx, descriptor_files, app_dir):
    for f in descriptor_files:
        if ctx.Contains(f):
            with open(os.path.join(app_dir, f), 'w') as w:
                w.write(ctx.GetFile(f))


def gen_tmp_dir(dirr):
    tmp_dir = tempfile.mkdtemp()
    dir_name = os.path.join(tmp_dir, dirr)
    os.mkdir(dir_name)
    return dir_name


def creation_time(image):
    logging.info(image.config_file())
    cfg = json.loads(image.config_file())
    return cfg.get('created')


def timestamp_to_time(dt_str):
    dt = dt_str.rstrip("Z")
    return datetime.datetime.strptime(dt, "%Y-%m-%dT%H:%M:%S")

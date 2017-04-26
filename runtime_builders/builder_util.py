#!/usr/bin/python

# Copyright 2017 Google Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import logging
import os
from ruamel import yaml
import subprocess
import tempfile


RUNTIME_BUCKET = 'runtime-builders'
RUNTIME_BUCKET_PREFIX = 'gs://{0}/'.format(RUNTIME_BUCKET)
MANIFEST_FILE = RUNTIME_BUCKET_PREFIX + 'runtimes.yaml'

SCHEMA_VERSION = 1
SCHEMA_LINE = 'schema_version: {0}\n'.format(SCHEMA_VERSION)


def write_to_gcs(gcs_path, file_contents):
    try:
        logging.info(gcs_path)
        fd, f_name = tempfile.mkstemp(text=True)
        os.write(fd, file_contents)

        command = ['gsutil', 'cp', f_name, gcs_path]
        try:
            output = subprocess.check_output(command)
        except subprocess.CalledProcessError as e:
            logging.error('Error encountered when writing to GCS!: {0}'
                          .format(output))
            logging.error(e)
    finally:
        os.remove(f_name)


def get_file_from_gcs(gcs_file, temp_file):
    command = ['gsutil', 'cp', gcs_file, temp_file]
    try:
        subprocess.check_output(command, stderr=subprocess.STDOUT)
        return True
    except subprocess.CalledProcessError as e:
        logging.error('Error when retrieving file from GCS! {0}'
                      .format(e.output))
        return False


def write_manifest_file(manifest):
    manifest_contents = yaml.round_trip_dump(manifest,
                                             default_flow_style=False)
    write_to_gcs(MANIFEST_FILE, SCHEMA_LINE + manifest_contents)


def load_manifest_file():
    try:
        _, tmp = tempfile.mkstemp(text=True)
        command = ['gsutil', 'cp', MANIFEST_FILE, tmp]
        subprocess.check_output(command, stderr=subprocess.STDOUT)
        with open(tmp) as f:
            return yaml.safe_load(f)
    except subprocess.CalledProcessError:
        logging.info('Manifest file not found in GCS: creating new one.')
    finally:
        os.remove(tmp)

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

import argparse
import filecmp
import glob
import json
import logging
import os
import subprocess
import sys
import tempfile

GCS_FILE_PREFIX = 'gs://'


def main():
    logging.getLogger().setLevel(logging.INFO)
    parser = argparse.ArgumentParser()
    parser.add_argument('--directory', '-d',
                        help='directory containing all builder config files',
                        required=True)
    args = parser.parse_args()

    return _verify(args.directory)


def _verify(directory):
    failures = 0

    try:
        for config_file in glob.glob(os.path.join(directory, '*.json')):
            with open(config_file, 'r') as f:
                config = json.load(f)
                for release in config['releases']:
                    staging_path = release['path']
                    for tag in release['tags']:
                        failures += _verify_files(staging_path, tag)
        return failures
    except ValueError as ve:
        logging.error('Error when parsing JSON! Check file formatting. \n{0}'
                      .format(ve))
    except KeyError as ke:
        logging.error('Config file is missing required field! \n{0}'
                      .format(ke))


def _verify_files(staging_path, tagged_path):
    bucket_name = staging_path.replace(GCS_FILE_PREFIX, '').split('/')[0]
    if bucket_name not in tagged_path:
        logging.error('Buckets do not match!')
        logging.error('{0} || {1}'.format(staging_path, tagged_path))
        return 1

    try:
        tmp1 = os.path.join(tempfile._get_default_tempdir(),
                            next(tempfile._get_candidate_names()))

        tmp2 = os.path.join(tempfile._get_default_tempdir(),
                            next(tempfile._get_candidate_names()))

        print tmp1
        print tmp2

        _get_file_from_gcs(staging_path, tmp1)
        _get_file_from_gcs(tagged_path, tmp2)

        if not filecmp.cmp(tmp1, tmp2):
            logging.error('Files {0} and {1} do not match!'
                          .format(staging_path, tagged_path))
            return 1
        return 0
    finally:
        os.remove(tmp1)
        os.remove(tmp2)


def _get_file_from_gcs(gcs_file, temp_file):
    command = ['gsutil', 'cp', gcs_file, temp_file]
    try:
        output = subprocess.check_output(command)
    except subprocess.CalledProcessError as e:
        logging.error('Error when retrieving file from GCS! {0}'
                      .format(output))
        logging.error(e)


if __name__ == '__main__':
    sys.exit(main())

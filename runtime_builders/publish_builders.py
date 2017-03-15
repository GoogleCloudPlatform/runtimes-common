#!/usr/bin/python

# Copyright 2016 Google Inc. All rights reserved.

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
import json
import logging
import os
import sys

import builder_util


def main():
    logging.getLogger().setLevel(logging.INFO)

    parser = argparse.ArgumentParser()
    parser.add_argument('--directory', '-d',
                        help='directory containing builder config files',
                        required=True)
    args = parser.parse_args()

    try:
        for filename in os.listdir(args.directory):
            filepath = os.path.join(args.directory, filename)
            if filepath.endswith('.json'):
                with open(filepath, 'r') as f:
                    config = json.load(f)
                    project_name = config['project']
                    latest = config['latest']
                    prefix = builder_util.RUNTIME_BUCKET_PREFIX
                    if not latest.startswith(prefix):
                        logging.error('Please provide fully qualified '
                                      'path to config file in GCS!')
                        logging.error('Path should start with \'{0}\''
                                      ''.format(prefix))
                        sys.exit(1)
                    parts = os.path.splitext(latest)
                    if parts[1] != '.yaml':
                        logging.error('Please provide yaml config file to '
                                      'publish as latest!')
                        sys.exit(1)
                    full_prefix = prefix + project_name + '-'
                    latest_file = parts[0][len(full_prefix):]
                    logging.info(latest_file)
                    _write_version_file(project_name, latest_file)
    except ValueError as ve:
        logging.error('Error when parsing JSON! Check file formatting. \n{0}'
                      .format(ve))
    except KeyError as ke:
        logging.error('Config file is missing required field! \n{0}'
                      .format(ke))


def _write_version_file(project_name, latest_version):
    file_name = '{0}.version'.format(project_name)
    full_path = builder_util.RUNTIME_BUCKET_PREFIX + file_name

    logging.info(full_path)

    builder_util.write_to_gcs(full_path, latest_version)


if __name__ == '__main__':
    sys.exit(main())

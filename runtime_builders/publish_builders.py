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
import subprocess
import sys


def main():
    logging.getLogger().setLevel(logging.INFO)

    parser = argparse.ArgumentParser()
    parser.add_argument('--directory', '-d',
                        help='directory containing builder config files',
                        required=True)
    args = parser.parse_args()

    failures = 0
    try:
        for filename in os.listdir(args.directory):
            filepath = os.path.join(args.directory, filename)
            if filepath.endswith('.json'):
                with open(filepath, 'r') as f:
                    config = json.load(f)
                    for builder in config['releases']:
                        staged_builder = builder['path']
                        for tag in builder['tags']:
                            failures += _copy(staged_builder, tag)
    except ValueError as ve:
        logging.error('Error when parsing JSON! Check file formatting. \n{0}'
                      .format(ve))
    except KeyError as ke:
        logging.error('Config file is missing required field! \n{0}'
                      .format(ke))
    return failures


def _copy(builder, tag):
    logging.info('Copying builder {0} to: {1}'.format(builder, tag))
    try:
        output = subprocess.check_output(['gsutil', 'cp', builder, tag])
    except subprocess.CalledProcessError as e:
        logging.error(e)
        return 1
    logging.debug(output)
    return 0


if __name__ == '__main__':
    sys.exit(main())

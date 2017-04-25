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

    manifest = builder_util.load_manifest_file()
    if manifest is None:
        manifest = {}
    logging.info(manifest)

    try:
        for filename in os.listdir(args.directory):
            filepath = os.path.join(args.directory, filename)
            if filepath.endswith('.json'):
                with open(filepath, 'r') as f:
                    config = json.load(f)
                prefix = builder_util.RUNTIME_BUCKET_PREFIX
                latest = config['latest']
                if not latest.startswith(prefix):
                    logging.error('Please provide fully qualified path to '
                                  'config file in GCS!')
                    logging.error('Path should start with \'{0}\''
                                  ''.format(prefix))
                    sys.exit(1)
                _format_and_write(config['project'], latest, prefix, manifest)
        builder_util.write_manifest_file(manifest)
    except ValueError as ve:
        logging.error('Error when parsing JSON! Check file formatting. \n{0}'
                      .format(ve))
    except KeyError as ke:
        logging.error('Config file is missing required field! \n{0}'
                      .format(ke))


def _format_and_write(project_name, latest, prefix, manifest):
    parts = os.path.splitext(latest)
    if parts[1] != '.yaml':
        logging.error('Please provide yaml config file to publish as latest!')
        sys.exit(1)
    try:
        project_prefix = prefix + project_name + '-'
        latest_file = parts[0][len(project_prefix):]
        full_latest_file = latest[len(prefix):]
        if project_name not in manifest.keys():
            manifest[project_name] = {}
            manifest[project_name]['target'] = {}
            manifest[project_name]['target']['file'] = None
            # TODO: add logic here for 'file' vs 'runtime' when
            # we add support for versioning
        m_project = manifest[project_name]
        prev_builder = m_project['target']['file']
        if prev_builder is not None and prev_builder != full_latest_file:
            logging.warn('Overwriting old {0} builder: {1}'
                         .format(project_name, prev_builder))
        m_project['target']['file'] = full_latest_file
        _write_version_file(project_name, latest_file)
    except KeyError as ke:
        logging.error('FATAL: Formatting issue encountered in manifest. '
                      'Exiting. \n{0}'.format(ke))
        sys.exit(1)


def _write_version_file(project_name, latest_version):
    file_name = '{0}.version'.format(project_name)
    full_path = builder_util.RUNTIME_BUCKET_PREFIX + file_name
    builder_util.write_to_gcs(full_path, latest_version)


if __name__ == '__main__':
    sys.exit(main())

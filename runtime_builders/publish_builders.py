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
from ruamel import yaml
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
    logging.debug(yaml.round_trip_dump(manifest, indent=4))

    try:
        for filename in os.listdir(args.directory):
            filepath = os.path.join(args.directory, filename)
            if filepath.endswith('.json'):
                with open(filepath, 'r') as f:
                    config = json.load(f)
            elif filepath.endswith('.yaml'):
                with open(filepath, 'r') as f:
                    config = yaml.round_trip_load(f)
            else:
                continue
            _parse_and_write(config, manifest)
        builder_util.verify_and_write_manifest(manifest)
    except ValueError as ve:
        logging.error('Error when parsing JSON! Check file formatting. \n{0}'
                      .format(ve))
    except KeyError as ke:
        logging.error('Config file is missing required field! \n{0}'
                      .format(ke))


def _parse_and_write(config, manifest):
    try:
        project_name = config['project']
        builders = config['builders']
        aliases = config['aliases']
        for builder in builders:
            _process_builder(builder, manifest)
        for alias in aliases:
            _process_alias(alias, manifest)
    except KeyError as ke:
        logging.error('Fatal error encountered when parsing config: {0}'
                      .format(ke))
        sys.exit(1)

    # TODO: remove once we deprecate old <runtime>.version file
    latest = config['latest']
    _publish_latest(latest, project_name)


def _publish_latest(latest, project_name):
    parts = os.path.splitext(latest)
    project_prefix = builder_util.RUNTIME_BUCKET_PREFIX + project_name + '-'
    latest_file = parts[0][len(project_prefix):]
    file_name = '{0}.version'.format(project_name)
    full_path = builder_util.RUNTIME_BUCKET_PREFIX + file_name
    builder_util.write_to_gcs(full_path, latest_file)


def _process_builder(builder, manifest):
    prefix = builder_util.RUNTIME_BUCKET_PREFIX
    latest = builder['latest']
    if not latest.startswith(prefix):
        logging.error('Please provide fully qualified path to '
                      'config file in GCS!')
        logging.error('Path \'{0}\' should start with \'{1}\''
                      ''.format(latest, prefix))
        sys.exit(1)
    parts = os.path.splitext(latest)
    if parts[1] != '.yaml':
        logging.error('Please provide yaml config file to publish as latest!')
        sys.exit(1)

    full_latest_file = latest[len(prefix):]
    _process_entry(builder['name'], full_latest_file, 'file', manifest)


def _process_alias(alias, manifest):
    _process_entry(alias['name'], alias['alias'], 'runtime', manifest)


def _process_entry(entry_key, entry_value, manifest_key, manifest):
    try:
        if entry_key not in manifest.keys():
            manifest[entry_key] = {}
            manifest[entry_key]['target'] = {}
            manifest[entry_key]['target'][manifest_key] = None
        m_project = manifest[entry_key]
        prev_entry = m_project['target'][manifest_key]
        if prev_entry is not None and prev_entry != entry_value:
            logging.warn('Overwriting old {0} {1} entry: {2}'
                         .format(entry_key, manifest_key, prev_entry))
        manifest[entry_key]['target'][manifest_key] = entry_value
    except KeyError as ke:
        logging.error('FATAL: Formatting issue encountered in manifest. '
                      'Exiting. \n{0}'.format(ke))
        sys.exit(1)


if __name__ == '__main__':
    sys.exit(main())

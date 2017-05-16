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
import sys
import tempfile


RUNTIME_BUCKET = 'runtime-builders'
RUNTIME_BUCKET_PREFIX = 'gs://{0}/'.format(RUNTIME_BUCKET)
MANIFEST_FILE = RUNTIME_BUCKET_PREFIX + 'runtimes.yaml'

SCHEMA_VERSION = 1


def copy_to_gcs(file_path, gcs_path):
    command = ['gsutil', 'cp', file_path, gcs_path]
    try:
        output = subprocess.check_output(command)
        logging.debug(output)
    except subprocess.CalledProcessError:
        logging.error('Error encountered when writing to GCS!')
    except Exception as e:
        logging.error('Fatal error encountered when shelling command {0}'
                      .format(command))
        logging.error(e)


def write_to_gcs(gcs_path, file_contents):
    try:
        logging.info(gcs_path)
        fd, f_name = tempfile.mkstemp(text=True)
        os.write(fd, file_contents)

        copy_to_gcs(f_name, gcs_path)
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


def verify_manifest(manifest):
    """Verify that the provided runtime manifest is valid before publishing.

    Aliases are provided for runtime 'names' that can be included in users'
    application configuration files: this method ensures that all the aliases
    can resolve to actual builder files.

    All builders and aliases are turned into nodes in a graph, which is then
    traversed to be sure that all nodes lead down to a builder node.

    Example formatting of the manifest, showing both an 'alias' and
    an actual builder file:

    runtimes:
      java:
        target:
          runtime: java-openjdk
      java-openjdk:
        target:
          file: gs://runtimes/java-openjdk-1234.yaml
        deprecation:
          message: "openjdk is deprecated."
    """
    try:
        node_graph = {}
        for key, val in manifest.get('runtimes').iteritems():
            target = val.get('target', {})
            if not target:
                deprecation = val.get('deprecation', {})
                if not deprecation:
                    logging.error('No target or deprecation specified for '
                                  'runtime: {0}'.format(key))
                    sys.exit(1)
                continue
            child = None
            isBuilder = 'file' in target.keys()
            if not isBuilder:
                child = target['runtime']
            node = node_graph.get(key, {})
            if not node:
                node_graph[key] = Node(key, isBuilder, child)
        for _, node in node_graph.items():
            child = node
            while True:
                if not child.child:
                    break
                elif child.child not in node_graph.keys():
                    logging.error('Non-existent alias provided for {0}: {1}'
                                  .format(child.name, child.child))
                    sys.exit(1)
                child = node_graph[child.child]
            if not child.isBuilder:
                logging.error('No terminating builder for alias {0}'
                              .format(node.name))
                sys.exit(1)
    except KeyError as ke:
        logging.error('Error encountered when verifying manifest:', ke)
        sys.exit(1)


def load_manifest_file():
    try:
        _, tmp = tempfile.mkstemp(text=True)
        command = ['gsutil', 'cp', MANIFEST_FILE, tmp]
        subprocess.check_output(command, stderr=subprocess.STDOUT)
        with open(tmp) as f:
            return yaml.round_trip_load(f)
    except subprocess.CalledProcessError:
        logging.info('Manifest file not found in GCS: creating new one.')
        return {'schema_version': SCHEMA_VERSION}
    finally:
        os.remove(tmp)


class Node:
    def __init__(self, name, isBuilder, child):
        self.name = name
        self.isBuilder = isBuilder
        self.child = child

    def __repr__(self):
        return '{0}: {1}|{2}'.format(self.name, self.isBuilder, self.child)

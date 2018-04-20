# Copyright 2018 Google Inc. All Rights Reserved.
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

import argparse
import httplib2
import json
import logging
import requests
import subprocess
import sys

from containerregistry.client import docker_name
from containerregistry.client import docker_creds
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_session
from containerregistry.transport import transport_pool

from ftl.common import ftl_util
from ftl.common import constants

from ftl.php import layer_builder as php_builder
from ftl.python import layer_builder as python_builder

PHP = 'php'
PYTHON = 'python'
LANGUAGES = [PHP, PYTHON]
LANGUAGE_CACHES = {
    PHP: constants.PHP_CACHE_NAMESPACE,
    PYTHON: constants.PYTHON_CACHE_NAMESPACE
}

MAPPING_BUCKET = 'ftl-global-cache'
MAPPING_FILE = '{language}-mapping.json'
LOCAL_MAPPING_FILE = '/workspace/' + MAPPING_FILE


def main():
    logging.getLogger().setLevel(logging.INFO)

    parser = argparse.ArgumentParser()
    parser.add_argument(
        '--packages',
        action='store',
        dest='packages',
        nargs='*',
        required=True,
        type=str,
        help='')
    parser.add_argument(
        '--language',
        '-l',
        action='store',
        dest='language',
        required=True,
        help='',
        choices=LANGUAGES)
    args = parser.parse_args()

    runner = CacheRunner(args.packages, args.language)
    runner.populate_cache()


class CacheRunner(object):
    def __init__(self, packages, language):
        self._packages = packages
        self._language = language

        _cache = LANGUAGE_CACHES[language]
        self._cache_name = constants.GLOBAL_CACHE_REGISTRY + '/' + _cache

        self._reg = docker_name.Registry('gcr.io', strict=False)
        self._creds = docker_creds.DefaultKeychain.Resolve(self._reg)
        self._transport = transport_pool.Http(httplib2.Http, size=2)
        self._cache = docker_name.Tag(self._cache_name, strict=False)

        # retrieve mappings when initializing runner
        self._mappings = self.read_mappings()
        logging.info('existing mapping: %s' % self._mappings)

    def _tag(self, tag):
        return docker_name.Tag(self._cache_name + ':' + tag)

    def populate_cache(self):
        existing_entries = self.retrieve_cache_entries()

        # Determine which images exist in the cache that should not be there,
        # and remove them
        self.remove_old_entries(existing_entries)

        # Populate cache with new entries
        self.populate_cache_entries(existing_entries)

        # Finally, write back the new key mapping to the filesystem
        # to be copied to GCS
        self.write_mapping_to_workspace()

    def read_mappings(self):
        """read cache_key -> package tuple mappings from GCS config file
        return map of key to package,which we'll use as lookup
        when pushing images"""

        r = requests.get(
            'https://www.googleapis.com/storage/v1'
            '/b/{bucket}/o/{file}?alt=media'.format(
                bucket=MAPPING_BUCKET,
                file=MAPPING_FILE.format(language=self._language)))

        if not r.ok:
            logging.error('Error retrieving mapping: %s' % r.text)
        try:
            return json.loads(r.text)
        except ValueError:
            # no mapping found: create new one
            return {}

    def retrieve_cache_entries(self):
        # returns all images stored in the cache currently
        with docker_image.FromRegistry(self._cache, self._creds,
                                       self._transport) as session:
            entries = set(tag.rstrip() for tag in session.tags() if tag)
            logging.info('existing entries in cache: %s' % entries)
            return entries

    def remove_old_entries(self, existing_entries):
        # for each existing entry in the mapping,
        # if it isn't in the package list, remove it
        for entry in existing_entries:
            entry_info = self._mappings.get(entry, '')
            if entry_info and entry_info not in self._packages:
                logging.info(
                    'removing entry {0} from cache'.format(entry_info))
                self._remove_entry(entry)

    def _remove_entry(self, entry):
        # delete entry from mapping and cache
        docker_session.Delete(self._tag(entry), self._creds, self._transport)
        del self._mappings[entry]

    def populate_cache_entries(self, existing_entries):
        # for each package, either verify it is already in the cache,
        # or build the image and push it to the cache
        for package in self._packages:
            if package:
                try:
                    name = None
                    version = None
                    if self._language == PHP:
                        name, version = package.split(':')
                    elif self._language == PYTHON:
                        name, version = package.split('==')
                    if name not in existing_entries:
                        # builder._pip_install() expects the double equals
                        # on the version
                        self._build_image_and_push(name, '==' + version)
                except ValueError:
                    logging.error(
                        'Encountered malformed package: {0}'.format(package))

    def _build_image_and_push(self, package_name, package_version):
        logging.info('building package {name}, version {version}'.format(
            name=package_name, version=package_version))
        image = None
        builder = None
        if self._language == PHP:
            builder = php_builder.PhaseTwoLayerBuilder(
                pkg_descriptor=(package_name, package_version))
        elif self._language == PYTHON:
            interpreter_builder = python_builder.InterpreterLayerBuilder()
            interpreter_builder.BuildLayer()
            builder = python_builder.PipfileLayerBuilder(
                pkg_descriptor=(package_name, package_version),
                wheel_dir=self._setup_pip_and_wheel(),
                dep_img_lyr=interpreter_builder)
        if not builder:
            logging.error('Could not find builder for language {0}'.format(
                self._language))
            sys.exit(1)
        builder.BuildLayer()
        # since we only have one layer, just use the builder's image
        # itself as the final image
        image = builder._img

        # TODO(nkubala): we should refactor AppendLayersIntoImage to not
        # have to set a base image
        # image = ftl_util.AppendLayersIntoImage([builder])

        key = builder.GetCacheKey()
        tag = self._tag(key)

        with docker_session.Push(
                tag, self._creds, self._transport, threads=2) as session:
            session.upload(image)
        self._mappings['%s:%s' % (package_name, package_version)] = key

    def write_mapping_to_workspace(self):
        with open(MAPPING_FILE.format(language=self._language), 'w') as f:
            json.dump(self._mappings, f)

    def _setup_pip_and_wheel(self):
        cmd = [constants.PIP_DEFAULT_CMD]
        cmd.extend(['install', '--upgrade', 'pip'])
        subprocess.check_call(cmd)

        cmd = [constants.PIP_DEFAULT_CMD]
        cmd.extend(['install', 'wheel'])
        subprocess.check_call(cmd)

        return ftl_util.gen_tmp_dir(constants.WHEEL_DIR)


if __name__ == '__main__':
    main()

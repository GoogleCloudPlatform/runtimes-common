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
"""This package defines the utilities for testing FTL objects."""

import hashlib
import logging
import os

from ftl.common import cache
from testing.lib import mock_registry

from containerregistry.client import docker_name
from containerregistry.client.v2_2 import docker_image


class MockHybridRegistry(cache.Base):
    """MockHybridRegistry is a cache implementation that stores layers in
    memory and can get docker tarballs from the local file system.

    It stores layers under a 'namespace', with a tag derived from the layer
    checksum. For example: gcr.io/$repo/$namespace:$checksum
    """

    def __init__(self, repo, directory):
        super(MockHybridRegistry, self).__init__()
        self._directory = directory
        self._repo = repo
        self._registry = mock_registry.MockRegistry()
        self._cache_miss = 0

    def _tag(self, base_image, namespace, checksum):
        fingerprint = '%s %s' % (base_image.digest(), checksum)
        return docker_name.Tag('{base}/{namespace}:{tag}'.format(
            base=str(self._repo),
            namespace=namespace,
            tag=hashlib.sha256(fingerprint).hexdigest()))

    def Get(self, base_image, namespace, checksum):
        entry = self._tag(base_image, namespace, checksum)
        if self._registry.existsImage(entry):
            return self._registry.getImage(entry)

        tarball = os.path.join(self._directory, str(entry))
        if os.path.isfile(tarball):
            logging.info('Found cached base image: %s.' % entry)
            self._registry.setImage(entry, docker_image.FromTarball(tarball))
            return self._registry.getImage(entry)
        logging.info('No cached base image found for entry: %s.' % entry)
        self._cache_miss += 1
        return None

    def Store(self, base_image, namespace, checksum, value):
        entry = self._tag(base_image, namespace, checksum)
        self._registry.setImage(entry, value)

    def StoreTarImage(self, namespace, checksum, tarpath, config_text):
        with docker_image.FromDisk(config_text, zip([], []),
                                   tarpath) as base_image:
            entry = self._tag(base_image, namespace, checksum)
            self._registry.setImage(entry, base_image)

    def GetRegistry(self):
        return self._registry.getRegistry()

    def GetCacheMiss(self):
        return self._cache_miss

    def ResetCacheMiss(self):
        self._cache_miss = 0

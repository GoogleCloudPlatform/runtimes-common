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
"""This package defines the interface for caching objects."""

import os
import abc
import hashlib
import logging

from containerregistry.client import docker_name
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_session
from testing.lib import mock_registry


class Base(object):
    """Base is an abstract base class representing a layer cache.

    It provides methods used by builders for accessing and storing layers.
    """

    # __enter__ and __exit__ allow use as a context manager.
    @abc.abstractmethod
    def __enter__(self):
        """Initialize the context."""

    def __exit__(self, unused_type, unused_value, unused_traceback):
        """Cleanup the context."""

    @abc.abstractmethod
    def Get(self, base_image, namespace, checksum):
        """Lookup a cached image.
        Args:
          base_image: the docker_image.Image on which things are based.
          namespace: a namespace for this cache.
          checksum: the checksum of the package descriptor atop our base.
        Returns:
          the docker_image.Image of the cache hit, or None.
        """

    @abc.abstractmethod
    def Store(self, base_image, namespace, checksum, value):
        """Lookup a cached image.
        Args:
          base_image: the docker_image.Image on which things are based.
          namespace: a namespace for this cache.
          checksum: the checksum of the package descriptor atop our base.
          value: the docker_image.Image to store into the cache.
        """


class Registry(Base):
    """Registry is a cache implementation that stores layers in a registry.

    It stores layers under a 'namespace', with a tag derived from the layer
    checksum. For example: gcr.io/$repo/$namespace:$checksum
    """

    def __init__(self, repo, creds, transport, threads=1, mount=None):
        super(Registry, self).__init__()
        self._repo = repo
        self._creds = creds
        self._transport = transport
        self._threads = threads
        self._mount = mount or []

    def _tag(self, base_image, namespace, checksum):
        fingerprint = '%s %s' % (base_image.digest(), checksum)
        return docker_name.Tag('{base}/{namespace}:{tag}'.format(
            base=str(self._repo),
            namespace=namespace,
            tag=hashlib.sha256(fingerprint).hexdigest()))

    def Get(self, base_image, namespace, checksum):
        entry = self._tag(base_image, namespace, checksum)
        with docker_image.FromRegistry(entry, self._creds,
                                       self._transport) as img:
            if img.exists():
                logging.info('Found cached base image: %s.' % entry)
                return img
            logging.info('No cached base image found for entry: %s.' % entry)
        return None

    def Store(self, base_image, namespace, checksum, value):
        entry = self._tag(base_image, namespace, checksum)
        with docker_session.Push(
                entry,
                self._creds,
                self._transport,
                threads=self._threads,
                mount=self._mount) as session:
            session.upload(value)


class MockHybridRegistry(Base):
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

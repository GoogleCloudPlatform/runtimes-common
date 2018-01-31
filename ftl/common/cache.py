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

import abc
import hashlib
import logging
import datetime

from containerregistry.client import docker_name
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_session

from ftl.common import ftl_util

_DEFAULT_TTL_WEEKS = 1


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
    def Get(self, base_image, namespace, cache_key):
        """Lookup a cached image.
        Args:
          base_image: the docker_image.Image on which things are based.
          namespace: a namespace for this cache.
          cache_key: the cache_key of the package descriptor atop our base.
        Returns:
          the docker_image.Image of the cache hit, or None.
        """

    @abc.abstractmethod
    def Set(self, base_image, namespace, cache_key, value):
        """Set an entry in the cache.
        Args:
          base_image: the docker_image.Image on which things are based.
          namespace: a namespace for this cache.
          cache_key: the cache_key of the package descriptor atop our base.
          value: the docker_image.Image to store into the cache.
        """


class Registry(Base):
    """Registry is a cache implementation that stores layers in a registry.

    It stores layers under a 'namespace', with a tag derived from the layer
    cache_key. For example: gcr.io/$repo/$namespace:$cache_key
    """

    def __init__(self,
                 repo,
                 creds,
                 transport,
                 cache_version=None,
                 threads=1,
                 mount=None):
        super(Registry, self).__init__()
        self._repo = repo
        self._creds = creds
        self._transport = transport
        self._cache_version = cache_version
        self._threads = threads
        self._mount = mount or []

    def _tag(self, base_image, namespace, cache_key):
        fingerprint = '%s %s' % (base_image.digest(), cache_key)
        if self._cache_version:
            fingerprint += ' ' + self._cache_version
        return docker_name.Tag('{base}/{namespace}:{tag}'.format(
            base=str(self._repo),
            namespace=namespace,
            tag=hashlib.sha256(fingerprint).hexdigest()))

    def Get(self, base_image, namespace, cache_key):
        """Attempt to retrieve value from cache."""
        logging.debug("Checking cache for base %s, namespace %s, cache_key %s",
                      base_image, namespace, cache_key)
        hit = self.getEntry(base_image, namespace, cache_key)
        if hit:
            logging.info('Found cached dependency layer for %s' % cache_key)
            if self.checkTTL(hit):
                return hit
            else:
                logging.info('TTL expired for cached image, rebuilding %s'
                             % cache_key)
        else:
            logging.info('No cached dependency layer for %s' % cache_key)

    def getEntry(self, base_image, namespace, cache_key):
        """Retrieve value from cache."""
        entry = self._tag(base_image, namespace, cache_key)
        logging.debug("Checking cache for entry %s", entry)
        with docker_image.FromRegistry(entry, self._creds,
                                       self._transport) as img:
            if img.exists():
                logging.info('Found cached base image: %s.' % entry)
                return img
            logging.info('No cached base image found for entry: %s.' % entry)

    def checkTTL(self, entry):
        """Check TTL of cache entry.
        Return whether or not the entry is expired."""
        last_created = ftl_util.timestamp_to_time(
                ftl_util.creation_time(entry))
        now = datetime.datetime.now()
        return last_created > now - datetime.timedelta(
                weeks=_DEFAULT_TTL_WEEKS)

    def Set(self, base_image, namespace, cache_key, value):
        entry = self._tag(base_image, namespace, cache_key)
        with docker_session.Push(
                entry,
                self._creds,
                self._transport,
                threads=self._threads,
                mount=self._mount) as session:
            session.upload(value)

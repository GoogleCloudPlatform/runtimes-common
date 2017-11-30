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
"""This package defines the shared cli args for ftl binaries."""

import abc

from containerregistry.client.v2_2 import docker_digest

_DEFAULT_TTL_WEEKS = 1


class BaseLayer(object):
    """BaseLayer is an abstract base class representing a container layer.

    It provides methods for generating a dependency layer and an application
    layer.
    """

    __metaclass__ = abc.ABCMeta  # For enforcing that methods are overriden.

    def __init__(self):
        self._img = None

    def GetImage(self):
        if self._img is None:
            raise Exception('error: layer image was not built yet so \
                             image cannot be accessed')
        return self._img

    def SetImage(self, img):
        self._img = img

    @abc.abstractmethod
    def BuildLayer(self):
        """Synthesizes the application layer from the context.
        Returns:
          a raw string of the layer's .tar.gz
        """


class CacheLayer(BaseLayer):
    @abc.abstractmethod
    def GetCacheKeyRaw(self):
        """Synthesizes the application layer from the context.
        Returns:
          a raw string of the layer's .tar.gz
        """

    def BuildLayer(self):
        pass

    def GetCacheKey(self):
        return docker_digest.SHA256(self.GetCacheKeyRaw())

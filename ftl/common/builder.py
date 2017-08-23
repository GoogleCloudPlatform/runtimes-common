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


import abc
import cStringIO
import gzip
import hashlib
import os
import tarfile

from containerregistry.client.v2_2 import append


import context


class Base(object):

  __metaclass__ = abc.ABCMeta  # For enforcing that methods are overriden.

  def __init__(self, ctx):
    self._ctx = ctx

  @abc.abstractmethod
  def CreatePackageBase(self, base_image, cache):
    """Create an image exists with the packages on this base.
    Args:
      base_image: docker_name.Tag, the base image atop which we install pkgs.
      cache: cache.Base, a cache into which artifacts may be read/written.
    Returns:
      a v2_2.docker_image.DockerImage of the above.
    """

  @abc.abstractmethod
  def BuildAppLayer(self):
    """Synthesizes the application layer from the context.
    Returns:
      a raw string of the layer's .tar.gz
    """

  # __enter__ and __exit__ allow use as a context manager.
  @abc.abstractmethod
  def __enter__(self):
    """Initialize the builder."""

  def __exit__(self, unused_type, unused_value, unused_traceback):
    """Cleanup after the builder."""
    pass


class JustApp(Base):

  def __init__(self, ctx):
    super(JustApp, self).__init__(ctx)

  def __enter__(self):
    """Override."""
    return self

  def CreatePackageBase(self, base_image, cache):
    """Override."""
    # JustApp doesn't install anything, it just appends
    # the application layer, so return the base image as
    # our package base.
    return base_image

  def BuildAppLayer(self):
    """Override."""
    buf = cStringIO.StringIO()
    with tarfile.open(fileobj=buf, mode='w') as out:
      for name in self._ctx.ListFiles():
        content = self._ctx.GetFile(name)
        info = tarfile.TarInfo(os.path.join('app', name))
        info.size = len(content)
        out.addfile(info, fileobj=cStringIO.StringIO(content))
    
    tar = buf.getvalue()
    sha = 'sha256:' + hashlib.sha256(tar).hexdigest()

    gz = cStringIO.StringIO()
    with gzip.GzipFile(fileobj=gz, mode='w', compresslevel=1) as f:
      f.write(tar)
    return gz.getvalue(), sha

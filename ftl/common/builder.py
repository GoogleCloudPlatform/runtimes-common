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
import hashlib
import os
import tarfile
import logging
import subprocess


class Base(object):
    """Base is an abstract base class representing a container builder.

    It provides methods for generating a dependency layer and an application
    layer.
    """

    __metaclass__ = abc.ABCMeta  # For enforcing that methods are overriden.

    def __init__(self, ctx):
        self._ctx = ctx

    @abc.abstractmethod
    def BuildAppLayer(self):
        """Synthesizes the application layer from the context.
        Returns:
          a raw string of the layer's .tar.gz
        """

    def __enter__(self):
        """Initialize the builder."""

    def __exit__(self, unused_type, unused_value, unused_traceback):
        """Cleanup after the builder."""


class JustApp(Base):
    """JustApp is an implementation of a builder that only generates an
    application layer.
    """

    def __init__(self, ctx):
        super(JustApp, self).__init__(ctx)

    def CreatePackageBase(self, base_image, cache):
        """Override."""
        # JustApp doesn't install anything, it just appends
        # the application layer, so return the base image as
        # our package base.
        return base_image

    def BuildAppLayer(self):
        """Override."""
        buf = cStringIO.StringIO()
        logging.info('Starting to generate app layer tarfile from context...')
        with tarfile.open(fileobj=buf, mode='w') as out:
            for name in self._ctx.ListFiles():
                content = self._ctx.GetFile(name)
                info = tarfile.TarInfo(os.path.join('app', name))
                info.size = len(content)
                out.addfile(info, fileobj=cStringIO.StringIO(content))
        logging.info('Finished generating app layer tarfile from context.')

        tar = buf.getvalue()
        sha = 'sha256:' + hashlib.sha256(tar).hexdigest()

        logging.info('Starting to gzip app layer tarfile...')
        gzip_process = subprocess.Popen(['gzip', '-f'],
                                        stdout=subprocess.PIPE,
                                        stdin=subprocess.PIPE,
                                        stderr=subprocess.PIPE)
        gz = gzip_process.communicate(input=tar)[0]
        logging.info('Finished gzipping tarfile.')
        return gz, sha

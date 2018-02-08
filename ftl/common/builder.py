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
import os
import tarfile
import logging
import subprocess
import httplib2
import json
import datetime

from containerregistry.client import docker_creds
from containerregistry.client import docker_name
from containerregistry.client.v2_2 import append
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_session
from containerregistry.client.v2_2 import save
from containerregistry.transport import transport_pool

from ftl.common import cache
from ftl.common import ftl_util
from ftl.common import single_layer_image
from ftl.common import tar_to_dockerimage

_THREADS = 32


class Base(object):
    """Base is an abstract base class representing a container builder.
    It provides methods for generating runtime layers and an application
    layer.
    """

    __metaclass__ = abc.ABCMeta  # For enforcing that methods are overriden.

    def __init__(self, ctx):
        self._ctx = ctx

    @abc.abstractmethod
    def Build(self):
        """Build method encapsulates all layer building and image creation.
        """


class JustApp(Base):
    """JustApp is an implementation of a builder that has logic to build an
    application layer.
    """

    def __init__(self, ctx):
        super(JustApp, self).__init__(ctx)

    def Build(self):
        """Override."""
        # this can't be abstract as it is instantiated in tests
        return

    class AppLayer(single_layer_image.BaseLayer):
        def __init__(self, ctx, destination_path='srv', entrypoint=None,
                     exposed_ports=None):
            self._ctx = ctx
            self._destination_path = destination_path
            self._entrypoint = entrypoint
            self._exposed_ports = exposed_ports

        def GetCacheKeyRaw(self):
            return None

        def BuildLayer(self):
            """Override."""
            buf = cStringIO.StringIO()
            logging.info('Starting to generate app layer \
                tarfile from context...')
            with tarfile.open(fileobj=buf, mode='w') as out:
                for name in self._ctx.ListFiles():
                    content = self._ctx.GetFile(name)
                    info = tarfile.TarInfo(
                        os.path.join(self._destination_path.strip("/"), name))
                    info.size = len(content)
                    out.addfile(info, fileobj=cStringIO.StringIO(content))
            logging.info('Finished generating app layer tarfile from context.')

            tar = buf.getvalue()

            logging.info('Starting to gzip app layer tarfile...')
            gzip_process = subprocess.Popen(
                ['gzip', '-f'],
                stdout=subprocess.PIPE,
                stdin=subprocess.PIPE,
                stderr=subprocess.PIPE)
            gz = gzip_process.communicate(input=tar)[0]
            overrides_dct = {
                    'created': str(datetime.date.today()) + "T00:00:00Z"
                }
            if self._entrypoint:
                overrides_dct['Entrypoint'] = self._entrypoint
            if self._exposed_ports:
                overrides_dct['ExposedPorts'] = self._exposed_ports
            logging.info('Finished gzipping tarfile.')
            self._img = tar_to_dockerimage.FromFSImage([gz], [tar],
                                                       overrides_dct)


class RuntimeBase(JustApp):
    """RuntimeBase is an abstract base class representing a container builder
    for runtime applications with dependencies.

    It provides methods for generating appending layers and caching images
    """

    def __init__(self, ctx, namespace, args, cache_version_str,
                 descriptor_files):
        super(RuntimeBase, self).__init__(ctx)
        self._namespace = namespace
        if args.entrypoint:
            args.entrypoint = args.entrypoint.split(" ")
        if args.exposed_ports:
            args.exposed_ports = args.exposed_ports.split(" ")
        self._args = args
        self._base_name = docker_name.Tag(self._args.base)
        self._base_creds = docker_creds.DefaultKeychain.Resolve(
            self._base_name)
        self._target_image = docker_name.Tag(self._args.name)
        self._target_creds = docker_creds.DefaultKeychain.Resolve(
            self._target_image)
        self._transport = transport_pool.Http(httplib2.Http, size=_THREADS)
        self._base_image = docker_image.FromRegistry(
            self._base_name, self._base_creds, self._transport)
        self._base_image.__enter__()
        self._cache = cache.Registry(
            repo=self._target_image.as_repository(),
            namespace=self._namespace,
            base_image=self._base_image,
            creds=self._target_creds,
            transport=self._transport,
            cache_version=cache_version_str,
            threads=_THREADS,
            mount=[self._base_name],
            use_global=self._args.global_cache)
        self._descriptor_files = descriptor_files

    def Build(self):
        return

    def AppendLayersIntoImage(self, lyr_imgs):
        for i, lyr_img in enumerate(lyr_imgs):
            if i == 0:
                result_image = lyr_img
                continue
            img = lyr_img.GetImage()
            diff_ids = img.diff_ids()
            for diff_id in diff_ids:
                lyr = img.blob(img._diff_id_to_digest(diff_id))
                overrides = ftl_util.CfgDctToOverrides(
                    json.loads(img.config_file()))

                result_image = append.Layer(
                    result_image, lyr, diff_id=diff_id, overrides=overrides)
        return result_image

    def StoreImage(self, result_image):
        if self._args.output_path:
            with ftl_util.Timing("saving_tarball_image"):
                with tarfile.open(
                        name=self._args.output_path, mode='w') as tar:
                    save.tarball(self._target_image, result_image, tar)
                logging.info("{0} tarball located at {1}".format(
                    str(self._target_image), self._args.output_path))
            return
        if self._args.upload:
            with ftl_util.Timing("pushing_image_to_docker_registry"):
                with docker_session.Push(
                        self._target_image,
                        self._target_creds,
                        self._transport,
                        threads=_THREADS,
                        mount=[self._base_name]) as session:
                    logging.info('Pushing final image...')
                    session.upload(result_image)
                return

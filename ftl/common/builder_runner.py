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
"""A binary for constructing images from a source context."""

import tarfile
import json
import datetime
from containerregistry.client import docker_creds
from containerregistry.client import docker_name
from containerregistry.client.v2_2 import append
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_session
from containerregistry.client.v2_2 import save
from containerregistry.transport import transport_pool

import httplib2
import logging
import hashlib

from ftl.common import cache
from ftl.common import context

_THREADS = 32
_DEFAULT_TTL_WEEKS = 1


class BuilderRunner():
    def __init__(self, args, builder, cache_version_str):
        self.args = args
        self.transport = transport_pool.Http(httplib2.Http, size=_THREADS)
        self.base_name = docker_name.Tag(args.base)
        self.base_creds = docker_creds.DefaultKeychain.Resolve(self.base_name)
        self.target_image = docker_name.Tag(args.name)
        self.target_creds = docker_creds.DefaultKeychain.Resolve(
            self.target_image)
        self.ctx = context.Workspace(args.directory)
        self.cash = cache.Registry(
            self.target_image.as_repository(),
            self.target_creds,
            self.transport,
            cache_version=cache_version_str,
            threads=_THREADS,
            mount=[self.base_name])
        self.builder = builder.From(self.ctx)

    def GetCacheKey(self, descriptor_files):
        descriptor = None
        for f in descriptor_files:
            if self.ctx.Contains(f):
                descriptor = f
                descriptor_contents = self.ctx.GetFile(descriptor)
                break
        if not descriptor:
            logging.info('No package descriptor found. No packages installed.')
            return None
        return hashlib.sha256(descriptor_contents).hexdigest()

    def GetCachedDepsImage(self, checksum):
        if not checksum:
            # TODO(aaron-prindle) verify this makes sense to use None
            # as sentinel for no descriptor and to handle this here
            return self.args.base

        hit = self.cash.Get(self.args.base, self.builder.namespace, checksum)
        if hit:
            logging.info('Found cached dependency layer for %s' % checksum)
            last_created = _timestamp_to_time(_creation_time(hit))
            now = datetime.datetime.now()
            if last_created > now - datetime.timedelta(
                    seconds=_DEFAULT_TTL_WEEKS):
                return hit
            else:
                logging.info(
                    'TTL expired for cached image, rebuilding %s' % checksum)
        else:
            logging.info('No cached dependency layer for %s' % checksum)
        return None

    def StoreDepsImage(self, dep_image, checksum):
        if self.args.cache:
            logging.info('Storing layer cash.')
            self.cash.Store(self.args.base, self.builder.namespace, checksum,
                            dep_image)
        else:
            logging.info('Skipping storing layer cash.')

    def GenerateFTLImage(self):
        with docker_image.FromRegistry(self.base_name, self.base_creds,
                                       self.transport) as self.args.base:

            # Create (or pull from cache) the base image with the
            # package descriptor installation overlaid.
            logging.info('Generating dependency layer...')
            checksum = self.GetCacheKey(self.builder.descriptor_files)
            deps_image = self.GetCachedDepsImage(checksum)
            if not deps_image:
                # TODO(aaron-prindle) make this better, prob pass args to bldr
                if self.args.destination_path:
                    deps_image = self.builder.CreatePackageBase(
                        self.args.base,
                        self.args.destination_path)
                else:
                    deps_image = self.builder.CreatePackageBase(
                        self.args.base)
                self.StoreDepsImage(deps_image, checksum)
            # Construct the application layer from the context.
            logging.info('Generating app layer...')
            app_layer, diff_id = self.builder.BuildAppLayer()
            with append.Layer(
                    deps_image, app_layer, diff_id=diff_id) as app_image:
                if self.args.output_path:
                    with tarfile.open(
                            name=self.args.output_path, mode='w') as tar:
                        save.tarball(self.target_image, app_image, tar)
                    logging.info("{0} tarball located at {1}".format(
                        str(self.target_image), self.args.output_path))
                    return
                with docker_session.Push(
                        self.target_image,
                        self.target_creds,
                        self.transport,
                        threads=_THREADS,
                        mount=[self.base_name]) as session:
                    logging.info('Pushing final image...')
                    session.upload(app_image)


def _creation_time(image):
    logging.info(image.config_file())
    cfg = json.loads(image.config_file())
    return cfg.get('created')


def _timestamp_to_time(dt_str):
    dt = dt_str.rstrip("Z")
    return datetime.datetime.strptime(dt, "%Y-%m-%dT%H:%M:%S")

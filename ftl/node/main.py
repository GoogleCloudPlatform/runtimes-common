
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

import argparse

from containerregistry.client import docker_creds
from containerregistry.client import docker_name
from containerregistry.client.v2_2 import append
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_session
from containerregistry.transport import transport_pool

import httplib2

from ftl.common import cache
from ftl.common import context

from ftl.node import builder


_THREADS = 32


parser = argparse.ArgumentParser(
    description='Construct images from source.')

parser.add_argument('--base', action='store',
                    help=('The name of the docker base image.'))

parser.add_argument('--name', action='store',
                    help=('The name of the docker image to push.'))

parser.add_argument('--directory', action='store',
                    help='The path where the application data sits.')


def main():
    args = parser.parse_args()

    transport = transport_pool.Http(httplib2.Http, size=_THREADS)

    # TODO(mattmoor): Support digest base images.
    base_name = docker_name.Tag(args.base)
    base_creds = docker_creds.DefaultKeychain.Resolve(base_name)

    target_image = docker_name.Tag(args.name)
    target_creds = docker_creds.DefaultKeychain.Resolve(target_image)

    with context.Workspace(args.directory) as ctx:
        with cache.Registry(
          target_image.as_repository(), target_creds, transport,
          threads=_THREADS, mount=[base_name]) as cash:
            with builder.From(ctx) as bldr:
                with docker_image.FromRegistry(
                  base_name, base_creds, transport) as base_image:

                    # Create (or pull from cache) the base image with the
                    # package descriptor installation overlaid.
                    print('Generating dependency layer...')
                    with bldr.CreatePackageBase(base_image, cash) as deps:
                        # Construct the application layer from the context.
                        print('Generating app layer...')
                        app_layer, diff_id = bldr.BuildAppLayer()
                        with append.Layer(
                          deps,
                          app_layer,
                          diff_id=diff_id) as app_image:
                            with docker_session.Push(
                              target_image, target_creds, transport,
                              threads=_THREADS, mount=[base_name]) as session:
                                print('Pushing final image...')
                                session.upload(app_image)


if __name__ == '__main__':
    main()

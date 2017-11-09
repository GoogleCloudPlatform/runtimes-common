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
import sys
import tarfile
from containerregistry.client import docker_creds
from containerregistry.client import docker_name
from containerregistry.client.v2_2 import append
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_session
from containerregistry.client.v2_2 import save
from containerregistry.tools import patched
from containerregistry.transport import transport_pool


import httplib2
import logging

from ftl.common import cache
from ftl.common import context

from ftl.node import builder

_THREADS = 32
_LEVEL_MAP = {
    "NOTSET": logging.NOTSET,
    "DEBUG": logging.DEBUG,
    "INFO": logging.INFO,
    "WARNING": logging.WARNING,
    "ERROR": logging.ERROR,
    "CRITICAL": logging.CRITICAL,
}

parser = argparse.ArgumentParser(
    description='Construct node images from source.')

parser.add_argument(
    '--base',
    action='store',
    required=True,
    help=('The name of the docker base image.'))

parser.add_argument(
    '--name',
    required=True,
    action='store',
    help=('The name of the docker image to push.'))

parser.add_argument(
    '--directory',
    required=True,
    action='store',
    help='The path where the application data sits.')

parser.add_argument(
    '--no-cache',
    dest='cache',
    action='store_false',
    help='Do not use cache during build.')

parser.add_argument(
    '--cache',
    dest='cache',
    default=True,
    action='store_true',
    help='Use cache during build (default).')

parser.add_argument(
    '--output-path',
    dest='output_path',
    action='store',
    help='Store final image as local tarball at output path \
          instead of pushing to registry')

parser.add_argument(
    "-v",
    "--verbosity",
    default="NOTSET",
    nargs="?",
    action='store',
    choices=_LEVEL_MAP.keys())


def main(args):
    args = parser.parse_args(args)
    logging.getLogger().setLevel(_LEVEL_MAP[args.verbosity])
    logging.basicConfig(
        format='%(asctime)s.%(msecs)03d %(levelname)-8s %(message)s',
        datefmt='%Y-%m-%d,%H:%M:%S')
    transport = transport_pool.Http(httplib2.Http, size=_THREADS)

    # TODO(mattmoor): Support digest base images.
    base_name = docker_name.Tag(args.base)
    base_creds = docker_creds.DefaultKeychain.Resolve(base_name)

    target_image = docker_name.Tag(args.name)
    target_creds = docker_creds.DefaultKeychain.Resolve(target_image)

    ctx = context.Workspace(args.directory)
    cash = cache.Registry(
        target_image.as_repository(),
        target_creds,
        transport,
        threads=_THREADS,
        mount=[base_name])
    bldr = builder.From(ctx)
    with docker_image.FromRegistry(base_name, base_creds,
                                   transport) as base_image:

        # Create (or pull from cache) the base image with the
        # package descriptor installation overlaid.
        logging.info('Generating dependency layer...')
        with bldr.CreatePackageBase(base_image, cash, args.cache) as deps:
            # Construct the application layer from the context.
            logging.info('Generating app layer...')
            app_layer, diff_id = bldr.BuildAppLayer()
            with append.Layer(deps, app_layer, diff_id=diff_id) as app_image:
                if args.output_path:
                    with tarfile.open(name=args.output_path, mode='w') as tar:
                        save.tarball(target_image, app_image, tar)
                    logging.info("{0} tarball located at {1}".format(
                                 str(target_image), args.output_path))
                    return
                with docker_session.Push(
                        target_image,
                        target_creds,
                        transport,
                        threads=_THREADS,
                        mount=[base_name]) as session:
                    logging.info('Pushing final image...')
                    session.upload(app_image)


if __name__ == '__main__':
    with patched.Httplib2():
        main(sys.argv[1:])

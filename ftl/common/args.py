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

from ftl.common import logger
import argparse


def base_parser():
    parser = argparse.ArgumentParser()

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
        choices=logger.LEVEL_MAP.keys())
    return parser


node_flgs = ['destination_path']
php_flgs = ['destination_path']


def extra_args(parser, opt_list):
    opt_dict = {
        'destination_path': [
            '--destination', {
                "dest":
                'destination_path',
                "action":
                'store',
                "default":
                None,
                "help":
                'The base path that the node_modules will be installed in the \
                        final image'
            }
        ],
    }
    for opt in opt_list:
        arg_vars = opt_dict[opt]
        parser.add_argument(arg_vars[0], **arg_vars[1])

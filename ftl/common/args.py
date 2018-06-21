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

import argparse
import os

from ftl.common import constants
from ftl.common import logger


def version_parser():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        '--version', action='version', version=constants.FTL_VERSION)
    return parser


def base_parser():
    parser = argparse.ArgumentParser()
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument(
        '--base', action='store', help=('The name of the docker base image.'))
    group.add_argument(
        '--tar_base_image_path',
        dest='tar_base_image_path',
        action='store',
        default=None,
        help='The tar path for the base dockerimage for the FTL build')
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
        '--additional-directory',
        action='store',
        help='Additional \
        path whose contents gets appended to the final image.')
    parser.add_argument(
        '--cache-repository',
        action='store',
        required=False,
        help=('The name of the repository to use as the root for the cache.'))
    parser.add_argument(
        '--no-cache',
        dest='cache',
        action='store_false',
        help='Do not check cache during build.')
    parser.add_argument(
        '--cache',
        dest='cache',
        default=True,
        action='store_true',
        help='Check cache during build (default).')
    parser.add_argument(
        '--global-cache',
        dest='global_cache',
        default=False,
        action='store_true',
        help='Use global cache')
    parser.add_argument(
        '--no-upload',
        dest='upload',
        action='store_false',
        help='Do not upload to cache during build.')
    parser.add_argument(
        '--upload',
        dest='upload',
        default=True,
        action='store_true',
        help='Upload to cache during build (default).')
    parser.add_argument(
        '--output-path',
        dest='output_path',
        action='store',
        help='Store final image as local tarball at output path \
            instead of pushing to registry')
    parser.add_argument(
        "-v",
        "--verbosity",
        default=constants.DEFAULT_LOG_LEVEL,
        nargs="?",
        action='store',
        choices=logger.LEVEL_MAP.keys())
    parser.add_argument(
        '--destination',
        dest='destination_path',
        action='store',
        default=constants.DEFAULT_DESTINATION_PATH,
        help='The base path that the app and dependency files will be \
        installed in the final image')
    parser.add_argument(
        '--entrypoint',
        dest='entrypoint',
        action='store',
        default=constants.DEFAULT_ENTRYPOINT,
        help='The entrypoint for the dockerimage')
    parser.add_argument(
        '--sh-c-prefix',
        dest='sh_c_prefix',
        action='store_true',
        default=False,
        help='If sh -c should be prepended to the entrypoint')
    parser.add_argument(
        '--succeed-on-error',
        dest='fail_on_error',
        action='store_false',
        default=True,
        help='If sh -c should be prepended to the entrypoint')
    parser.add_argument(
        '--exposed-ports',
        dest='exposed_ports',
        action='store',
        default=None,
        help='The port to expose for the dockerimage')
    parser.add_argument(
        '--builder-output-path',
        dest='builder_output_path',
        action='store',
        default=(os.environ.get(constants.BUILDER_OUTPUT)
                 if os.environ.get(constants.BUILDER_OUTPUT) else None),
        help='The path to store FTL logs')
    return parser


node_flgs = []
php_flgs = []
python_flgs = ['python_cmd', 'pip_cmd', 'venv_cmd', 'venv_dir']


def extra_args(parser, opt_list):
    opt_dict = {
        'python_cmd': [
            '--python-cmd', {
                "dest": 'python_cmd',
                "action": 'store',
                "default": constants.PYTHON_DEFAULT_CMD,
                "help": 'The python command to be run (ex: python2.7)'
            }
        ],
        'pip_cmd': [
            '--pip-cmd', {
                "dest": 'pip_cmd',
                "action": 'store',
                "default": constants.PIP_DEFAULT_CMD,
                "help": 'The pip command to be run (ex: pip)'
            }
        ],
        'venv_cmd': [
            '--virtualenv-cmd', {
                "dest": 'venv_cmd',
                "action": 'store',
                "default": constants.VENV_DEFAULT_CMD,
                "help": 'The virtualenv command to be run (ex: virtualenv)'
            }
        ],
        'venv_dir': [
            '--virtualenv-dir', {
                "dest": 'venv_dir',
                "action": 'store',
                "default": constants.VENV_DIR,
                "help": 'The virtualenv command to be run (ex: virtualenv)'
            }
        ],
    }
    for opt in opt_list:
        arg_vars = opt_dict[opt]
        parser.add_argument(arg_vars[0], **arg_vars[1])

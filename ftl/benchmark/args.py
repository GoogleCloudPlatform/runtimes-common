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


def base_parser():
    parser = argparse.ArgumentParser()

    parser.add_argument(
        '--base',
        action='store',
        help=('The name of the docker base image.'))

    parser.add_argument(
        '--name',
        action='store',
        help=('The name of the docker image to push.'))

    parser.add_argument(
        '--directory',
        action='store',
        help='The path where the application data sits.')

    parser.add_argument(
        '--iterations',
        action='store',
        type=int,
        default=5,
        help='Number of times to build the image')

    parser.add_argument(
        '--description',
        action='store',
        help=('Description of the app being benchmarked.'))

    return parser

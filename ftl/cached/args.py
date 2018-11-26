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
from ftl.common import args


def str2bool(v):
    if v.lower() in ('yes', 'true', 't', 'y', '1'):
        return True
    elif v.lower() in ('no', 'false', 'f', 'n', '0'):
        return False
    else:
        raise argparse.ArgumentTypeError('Boolean value expected.')


def base_parser():
    parser = args.base_parser()

    parser.add_argument(
        '--label-1',
        dest='label_1',
        action='store',
        default='original',
        help='image label for original uploaded image')

    parser.add_argument(
        '--label-2',
        dest='label_2',
        action='store',
        default='reupload',
        help='image label for reuploaded image')

    parser.add_argument(
        '--dir-1',
        dest='dir_1',
        action='store',
        required=True,
        help='image label for original uploaded image')

    parser.add_argument(
        '--dir-2',
        dest='dir_2',
        action='store',
        required=True,
        help='image label for reuploaded image')

    parser.add_argument(
        '--layer-offset',
        dest='layer_offset',
        type=int,
        action='store',
        help='the number of expected differing layers')
    parser.add_argument(
        '--should-cache',
        dest='should_cache',
        type=str2bool,
        default=False,
        action='store',
        help='if the test is expecting a cache hit')

    return parser

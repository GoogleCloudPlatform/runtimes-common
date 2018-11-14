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
"""This package defines helpful utilities for FTL ."""
import os

from ftl.common import ftl_util


def setup_virtualenv(virtualenv_dir, virtualenv_cmd, python_cmd):
    if os.path.isdir(virtualenv_dir):
        return
    virtualenv_cmd_args = list(virtualenv_cmd)
    virtualenv_cmd_args.extend([
        '--no-download',
        virtualenv_dir,
        '-p',
    ])
    virtualenv_cmd_args.extend(list(python_cmd))
    ftl_util.run_command(
        'create_virtualenv',
        virtualenv_cmd_args,
        cmd_cwd="/")

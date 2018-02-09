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

import os
import sys
import argparse

from containerregistry.tools import patched

from ftl.common import args
from ftl.common import logger
from ftl.common import context
from ftl.common import ftl_util

from ftl.php import builder as php_builder

parser = args.base_parser()
php_parser = argparse.ArgumentParser(
    add_help=False,
    parents=[parser],
    description='Construct php images from source.')
args.extra_args(php_parser, args.php_flgs)

# Version string used to bust caches.
_PHP_CACHE_VERSION = 'v1'


def main(cli_args):
    builder_args = php_parser.parse_args(cli_args)
    logger.setup_logging(builder_args)
    with ftl_util.Timing("full build"):
        with ftl_util.Timing("builder initialization"):
            php_ftl = php_builder.PHP(
                context.Workspace(builder_args.directory), builder_args,
                _PHP_CACHE_VERSION)
        with ftl_util.Timing("build process for FTL image"):
            php_ftl.Build()
    newpid = os.fork()
    if newpid == 0:
        with ftl_util.Timing("GCS mapping upload"):
            php_ftl.SaveCacheMappings()


if __name__ == '__main__':
    with patched.Httplib2():
        main(sys.argv[1:])

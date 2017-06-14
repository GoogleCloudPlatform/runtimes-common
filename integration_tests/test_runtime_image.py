#!/usr/bin/python

# Copyright 2017 Google Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import logging
import sys

from testsuite import deploy_app
from testsuite import test_util


def _main(appdir):
    try:
        logging.debug('Testing runtime image.')
        version = deploy_app.deploy_app_without_image(appdir)
        application_url = test_util.retrieve_url_for_version(version)
        output, status_code = test_util.get(application_url)
        if status_code:
            raise Exception('Error pinging application!')
    except Exception as e:
        logging.debug('{0}'.format(e))
        sys.exit(1)


if __name__ == '__main__':
    sys.exit(_main(sys.argv[1]))

#!/usr/bin/python

# Copyright 2017 Google Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the 'License');
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an 'AS IS' BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import json
import logging
import unittest
from retrying import retry

import test_util


class TestCustom(unittest.TestCase):

    def __init__(self, url, methodName='runTest'):
        self._base_url = url
        self._url = url + test_util.CUSTOM_ENDPOINT
        unittest.TestCase.__init__(self)

    def runTest(self):
        logging.debug('Retrieving list of custom test endpoints.')
        output, status_code = test_util.get(self._url)
        self.assertEquals(status_code, 0,
                          'Cannot connect to sample application!')

        logging.debug('output: {0}'.format(output))

        for endpoint_info in json.loads(output):
            endpoint = endpoint_info[0]
            timeout = endpoint_info[1]

            full_endpoint = self._base_url + endpoint
            logging.info('making get request to {0}'.format(full_endpoint))
            response, code = test_util.get(full_endpoint, timeout=timeout)

            logging.debug(response)




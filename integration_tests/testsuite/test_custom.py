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
import urlparse

import test_util


class TestCustom(unittest.TestCase):

    def __init__(self, url, methodName='runTest'):
        self._base_url = url
        self._url = urlparse.urljoin(url, test_util.CUSTOM_ENDPOINT)
        unittest.TestCase.__init__(self)

    def runTest(self):
        logging.debug('Retrieving list of custom test endpoints.')
        output, status_code = test_util.get(self._url)
        self.assertEquals(status_code, 0,
                          'Cannot connect to sample application!')

        test_num = 0
        if not output:
            logging.debug('No custom tests specified.')
        else:
            for test_info in json.loads(output):
                test_num += 1
                name = test_info.get('name', 'test_{0}'.format(test_num))
                path = test_info.get('path')
                if path is None:
                    logging.warn('Test \'%s\' has no path specified! '
                                 'Skipping...', name)
                    continue

                timeout = test_info.get('timeout', 500)

                test_endpoint = urlparse.urljoin(self._base_url, path)
                logging.info('Running custom test: %s', name)
                response, _ = test_util.get(test_endpoint, timeout=timeout)

                logging.debug(response)

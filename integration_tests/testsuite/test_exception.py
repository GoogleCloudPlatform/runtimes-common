#!/usr/bin/python

# Copyright 2016 Google Inc. All rights reserved.

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
import unittest

import test_util


class TestException(unittest.TestCase):

    def __init__(self, url, methodName='runTest'):
        self._url = url + test_util.EXCEPTION_ENDPOINT
        unittest.TestCase.__init__(self)

    def runTest(self):
        payload = test_util._generate_exception_payload()
        response_code = test_util._post(self._url, payload)
        if response_code != 0:
            return self.fail('Error encountered inside sample application!')
        logging.info('Token {0} written to Stackdriver '
                     'Error Reporting'.format(payload.get('token')))

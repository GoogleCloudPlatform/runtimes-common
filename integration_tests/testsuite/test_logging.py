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
import unittest
import urlparse
from retrying import retry

import google.cloud.logging

import test_util


class TestLogging(unittest.TestCase):

    def __init__(self, url, methodName='runTest'):
        self._url = urlparse.urljoin(url, test_util.LOGGING_ENDPOINT)
        unittest.TestCase.__init__(self)

    def runTest(self):
        logging.debug('Posting to endpoint: {0}'.format(self._url))

        payload = test_util.generate_logging_payload()
        if test_util.post(self._url, payload) != 0:
            return self.fail('Error encountered inside sample application!')

        client = google.cloud.logging.Client()
        log_name = payload.get('log_name')
        token = payload.get('token')

        logging.info('log name is {0}, '
                     'token is {1}'.format(log_name, token))

        project_id = test_util._project_id()
        FILTER = 'logName = projects/{0}/logs/' \
                 'appengine.googleapis.com%2Fstdout'.format(project_id)

        self.assertTrue(self._read_log(client, log_name, token, FILTER),
                        'Log entry not found for posted token!')

    @retry(wait_fixed=4000, stop_max_attempt_number=8)
    def _read_log(self, client, log_name, token, filter):
        for entry in client.list_entries(filter_=filter):
                if token in entry.payload:
                    logging.info('Token {0} found in '
                                 'Stackdriver logs!'.format(token))
                    return True
        raise Exception('Log entry not found for posted token!')

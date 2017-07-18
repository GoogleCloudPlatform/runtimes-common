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
from retrying import retry
import unittest
import urlparse

import google.cloud.monitoring

import test_util


class TestMonitoring(unittest.TestCase):

    def __init__(self, url, methodName='runTest'):
        self._url = urlparse.urljoin(url, test_util.MONITORING_ENDPOINT)
        super(TestMonitoring, self).__init__()

    def runTest(self):
        payload = test_util.generate_metrics_payload()
        _, response_code = test_util.post(self._url, payload,
                                          test_util.METRIC_TIMEOUT)
        self.assertEquals(response_code, 0,
                          'Error encountered inside sample application!')

        client = google.cloud.monitoring.Client()

        self.assertTrue(self._read_metric(payload.get('name'),
                                          payload.get('token'), client),
                        'Token not found in Stackdriver monitoring!')

    @retry(wait_fixed=8000, stop_max_attempt_number=10)
    def _read_metric(self, name, target, client):
        query = client.query(name, minutes=2)
        if self._query_is_empty(query):
            raise Exception('Metric read retries exceeded!')

        for timeseries in query:
            for point in timeseries.points:
                if point.value == target:
                    logging.info('Token {0} found in Stackdriver '
                                 'metrics'.format(target))
                    return True
                print(point.value)
        return False

    def _query_is_empty(self, query):
        if query is None:
            logging.info('query is none')
            return True
        # query is a generator, so sum over it to get the length
        query_length = sum(1 for timeseries in query)
        if query_length == 0:
            logging.info('query is empty')
            return True
        return False

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
import time
from retrying import retry

import google.cloud.monitoring

import test_util


def _test_monitoring(base_url):
    url = base_url + test_util.MONITORING_ENDPOINT

    payload = test_util._generate_metrics_payload()
    if test_util._post(url, payload, test_util.METRIC_TIMEOUT) != 0:
        return test_util._fail('Error encountered inside test application!')

    # wait for metric to propagate
    time.sleep(test_util.METRIC_PROPAGATION_TIME)

    try:
        client = google.cloud.monitoring.Client()

        if not _read_metric(payload.get('name'),
                            payload.get('token'), client):
            return test_util._fail('Token not found in Stackdriver ' +
                                   'monitoring!')
        return 0
    except Exception as e:
        return test_util._fail(e)


@retry(wait_exponential_multiplier=1000,
       stop_max_attempt_number=8,
       wait_exponential_max=8000)
def _read_metric(name, target, client):
    query = client.query(name, minutes=2)
    if _query_is_empty(query):
        raise Exception('Metric read retries exceeded!')

    for timeseries in query:
        for point in timeseries.points:
            if point.value == target:
                logging.info('Token {0} found in Stackdriver '
                             'metrics'.format(target))
                return True
            print point.value
    return False


def _query_is_empty(query):
    if query is None:
        logging.info('query is none')
        return True
    # query is a generator, so sum over it to get the length
    query_length = sum(1 for timeseries in query)
    if query_length == 0:
        logging.info('query is empty')
        return True
    return False

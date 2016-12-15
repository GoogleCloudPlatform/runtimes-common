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

import json
import logging
import requests
import time
import unittest
from retrying import retry

import google.cloud.monitoring

import test_util

def _test_monitoring(base_url):
  logging.info('testing monitoring')
  url = base_url + test_util.MONITORING_ENDPOINT

  payload = test_util._generate_metrics_payload()

  try:
    headers = {'Content-Type': 'application/json'}
    response = requests.post(url, json.dumps(payload), timeout=test_util.METRIC_TIMEOUT, headers=headers)
    test_util._check_response(response, 'error when posting metric request!')
  except requests.exceptions.Timeout:
    logging.error('Timeout when posting metric data!')

  time.sleep(test_util.METRIC_PROPAGATION_TIME) # wait for metric to propagate

  try:
    client = google.cloud.monitoring.Client()

    for descriptor in client.list_resource_descriptors():
      print descriptor.type

    if not _read_metric(payload.get('name'), payload.get('token'), client):
      logging.error('Token not found in Stackdriver monitoring!')
      return False
    return True

  except Exception as e:
    logging.error(e)

if __name__ == '__main__':
  unittest.main()


@retry(wait_exponential_multiplier=1000, stop_max_attempt_number=8, wait_exponential_max=8000)
def _read_metric(name, target, client):
  query = client.query(name, minutes=2)
  if _query_is_empty(query):
    raise Exception('Metric read retries exceeded!')

  for timeseries in query:
    for point in timeseries.points:
      logging.info(point)
      if point.value == target:
        logging.info('Token {0} found in Stackdriver metric'.format(target))
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

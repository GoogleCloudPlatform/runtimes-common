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
    query = client.query(payload.get('name'), minutes=5)
    for timeseries in query:
      for point in timeseries.points:
        logging.debug(point)
        if point.value == payload.get('token'):
          logging.info('Token {0} found in Stackdriver metric'.format(payload.get('token')))
          return True
        print point.value

    logging.error('Token not found in Stackdriver monitoring!')
    return False

    for descriptor in client.list_resource_descriptors():
      print descriptor.type
  except Exception as e:
    logging.error(e)

if __name__ == '__main__':
  unittest.main()

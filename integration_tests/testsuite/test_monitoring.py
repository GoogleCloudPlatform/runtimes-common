#!/usr/bin/python

import json
import logging
import requests
import time

import test_util

from google.cloud import monitoring as gcloud_monitoring


def _test_monitoring(base_url):
  logging.info("testing monitoring")
  url = base_url + test_util.MONITORING_ENDPOINT

  payload = test_util._generate_metrics_payload()

  try:
    headers = {'Content-Type': 'application/json'}
    response = requests.post(url, json.dumps(payload), timeout=test_util.METRIC_TIMEOUT, headers=headers)
    test_util._check_response(response, "error when posting metric request!")
  except requests.exceptions.Timeout:
    logging.error("timeout when posting metric data!")

  time.sleep(test_util.METRIC_PROPAGATION_TIME) # wait for metric to propagate

  try:
    client = gcloud_monitoring.Client()
    query = client.query(payload.get('name'), minutes=5)
    for timeseries in query:
      for point in timeseries.points:
        logging.debug(point)
        if point.value == payload.get('token'):
          logging.info("token {0} found in stackdriver metric".format(payload.get('token')))
          return True
        print point.value

    logging.error("token not found in stackdriver monitoring!")
    return False

    for descriptor in client.list_resource_descriptors():
      print descriptor.type
  except Exception as e:
    logging.error(e)

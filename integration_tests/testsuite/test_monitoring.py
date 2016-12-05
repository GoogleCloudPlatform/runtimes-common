#!/usr/bin/python

import json
import logging
import requests

import test_util

from google.cloud import monitoring as gcloud_monitoring


def _test_monitoring(base_url):
	logging.info("testing monitoring")
	url = base_url + test_util.MONITORING_ENDPOINT

	payload = test_util._generate_metrics_payload()

	try:
		headers = {'Content-Type': 'application/json'}
		response = requests.post(url, json.dumps(payload), timeout=5, headers=headers)
	except requests.exceptions.Timeout:
		logging.error("timeout when posting metric data!")

	if response.status_code - 200 >= 100: # 2xx
		logging.error("error when posting metric request: exit code: {0}, text: {1}".format(response.status_code, response.text))

	try:
		client = gcloud_monitoring.Client()
		for descriptor in client.list_resource_descriptors():
			print descriptor.type
	except Exception as e:
		logging.error(e)

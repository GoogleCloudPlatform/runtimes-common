#!/usr/bin/python

import json
import logging
import requests

import test_util

from google.cloud import logging as gcloud_logging


def _test_logging(base_url):
	logging.info("testing logging")
	url = base_url + test_util.LOGGING_ENDPOINT
	logging.debug("posting to endpoint: {0}".format(url))

	# TODO (nkubala): possibly handle multiple log destinations depending on environments
	payload = test_util._generate_logging_payload()
	logging.debug("data: {0}".format(payload))

	try:
		headers = {'Content-Type': 'application/json'}
		response = requests.post(url, json.dumps(payload), timeout=5, headers=headers)
		test_util._check_response(response, "error when posting log request!")
	except requests.exceptions.Timeout:
		logging.error("timeout when posting log data!")

	try:
		client = gcloud_logging.Client(credentials=GoogleCredentials.get_application_default())
		log_name = payload.get('log_name')
		logging.info("log name is {0}".format(log_name))

		FILTER = 'logName = projects/nick-cloudbuild/logs/appengine.googleapis.com%2Fstdout'
		for entry in client.list_entries(filter_=FILTER):
			logging.info("entry is {0}".format(entry))
			logging.info("entry.payload is {0}".format(entry.payload))
			if entry.payload == token:
				return
		logging.error("log entry not found for posted token!")
	except Exception as e:
		logging.error("error encountered when retrieving log entries!")
		logging.error("exception type: {0}".format(type(e)))
		logging.error(e)

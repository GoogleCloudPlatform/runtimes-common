#!/usr/bin/python

import binascii
from oauth2client.client import GoogleCredentials
import os
import logging
import random
import requests
import string
import subprocess
import sys

from google.cloud import logging as gcloud_logging

LOGNAME_LENGTH = 20

ROOT_ENDPOINT = "/"
ROOT_EXPECTED_OUTPUT = "Hello World!"

LOGGING_ENDPOINT = "/log"


def _test_app(base_url):
	logging.info("starting app test with base url {0}".format(base_url))
	
	_test_root(base_url)
	_test_logging(base_url)


def _test_root(base_url):
	url = base_url + ROOT_ENDPOINT
	logging.debug("hitting endpoint: {0}".format(url))
	response = requests.get(url)
	output = response.content
	if response.status_code != 200:
		logging.error("error when making get request: {0}".format(output))
	logging.info("output is: {0}".format(output))
	if output != ROOT_EXPECTED_OUTPUT:
		# TODO: best way to handle error?
		logging.error("unexpected output: expected {0}, received {1}".format(ROOT_EXPECTED_OUTPUT, output))
	# assert output == v


def _test_logging(base_url):
	logging.info("testing logging")
	url = base_url + LOGGING_ENDPOINT
	logging.debug("posting to endpoint: {0}".format(url))

	payload = _generate_payload()
	logging.debug("data: {0}".format(payload))

	try:
		headers = {'Content-Type': 'application/json'}
		response = requests.post(url, payload, timeout=5, headers=headers)
	except requests.exceptions.Timeout:
		logging.error("timeout when posting log data!")

	if response.status_code != 200:
		logging.error("error when posting log request: exit code: {0}, text: {1}".format(response.status_code, response.text))

	try:
		client = gcloud_logging.Client(credentials=GoogleCredentials.get_application_default())
		# client = gcloud_logging.Client()
		log_name = payload.get('log_name')
		logging.info("log name is {0}".format(log_name))

		gcloud_logger = client.logger(log_name)
		for entry in gcloud_logger.list_entries():
			logging.info("entry is {0}".format(entry))
			logging.info("entry.payload is {0}".format(entry.payload))
			if entry.payload == token:
				return
		logging.error("log entry not found for posted token!")
	except Exception as e:
		logging.error("error encountered when retrieving log entries!")
		logging.error("exception type: {0}".format(type(e)))
		logging.error(e)


def _generate_logname_and_token():
	logname = ''.join(random.choice(string.ascii_uppercase + string.ascii_lowercase) for i in range(LOGNAME_LENGTH))
	token = binascii.b2a_hex(os.urandom(16))
	return logname, token


def _generate_payload():
  	logname, token = _generate_logname_and_token()
  	data = {'log_name':logname, 'token':token}
  	return data


if __name__ == "__main__":
	sys.exit(main())

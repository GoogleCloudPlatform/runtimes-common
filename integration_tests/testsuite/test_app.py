#!/usr/bin/python

import binascii
import os
import logging
import random
import string
import subprocess
import sys

from google.cloud import logging as gcloudlogging

LOGNAME_LENGTH = 20

ROOT_ENDPOINT = "/"
ROOT_EXPECTED_OUTPUT = "Hello World!"

LOGGING_ENDPOINT = "/log?logname={logname}&token={token}"

def _test_app(base_url):
	logging.info("starting app test with base url {0}".format(base_url))
	
	_test_root(base_url)
	# _test_logging(base_url)

def _test_root(base_url):
	url = base_url + ROOT_ENDPOINT
	logging.debug("hitting endpoint: {0}".format(url))
	output = subprocess.check_output(['curl', '--silent', '-L', url])
	logging.info("output is: {0}".format(output))
	if output != ROOT_EXPECTED_OUTPUT:
		# TODO: best way to handle error?
		logging.error("unexpected output")
	# assert output == v

def _test_logging(base_url):
	logging.info("testing logging")
	logname, token = _generate_logname_and_token()
	url = base_url + LOGGING_ENDPOINT.format(logname=logname, token=token)
	logging.debug("posting to endpoint: {0}".format(url))
	output = subprocess.check_output(
		['curl', '-X', 'POST', '--silent', '-L', url])

	client = gcloudlogging.Client()
	gcloud_logger = client.logger(logname)
	for entry in gcloud_logger.list_entries():
		if entry.payload == token:
			return

	logging.error("log entry not found for posted token!")

def _generate_logname_and_token():
	logname = ''.join(random.choice(string.ascii_uppercase + string.ascii_lowercase) for i in range(LOGNAME_LENGTH))
	token = binascii.b2a_hex(os.urandom(16))
	return logname, token

if __name__ == "__main__":
	sys.exit(main())

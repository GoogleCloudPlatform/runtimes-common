#!/usr/bin/python

import binascii
import os
import logging
import subprocess
import sys

ROOT_ENDPOINT = "/"
ROOT_EXPECTED_OUTPUT = "Hello World!"

LOGGING_ENDPOINT = "/log?token={token}"

def _test_app(base_url):
	logging.info("starting app test with base url {0}".format(base_url))
	
	_test_root(base_url)
	_test_logging(base_url)

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
	token = _generate_token()
	url = base_url + LOGGING_ENDPOINT.format(token=token)
	logging.debug("posting to endpoint: {0}".format(url))
	output = subprocess.check_output(
		['curl', '-X', 'POST', '--silent', '-L', url])

def _generate_token():
	return binascii.b2a_hex(os.urandom(16))

if __name__ == "__main__":
	sys.exit(main())

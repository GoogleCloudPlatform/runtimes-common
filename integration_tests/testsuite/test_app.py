#!/usr/bin/python

import binascii
import json
from oauth2client.client import GoogleCredentials
import os
import logging
import random
import requests
import string
import subprocess
import sys

from Config import *


def _test_app(base_url):
	logging.info("starting app test with base url {0}".format(base_url))
	
	_test_root(base_url)
	# _test_logging(base_url)
	_test_monitoring(base_url)


def _test_root(base_url):
	url = base_url + ROOT_ENDPOINT
	logging.debug("hitting endpoint: {0}".format(url))
	response = requests.get(url)
	output = response.content
	if response.status_code != 200:
		logging.error("error when making get request: {0}".format(output))
	logging.info("output is: {0}".format(output))
	if output != ROOT_EXPECTED_OUTPUT:
		# TODO (nkubala): best way to handle error?
		# should probably raise "FailedTestException" that is caught by the driver
		logging.error("unexpected output: expected {0}, received {1}".format(ROOT_EXPECTED_OUTPUT, output))


if __name__ == "__main__":
	sys.exit(main())

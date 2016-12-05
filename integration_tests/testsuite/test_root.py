#!/usr/bin/python

import logging
import requests

import test_util


def _test_root(base_url):
	url = base_url + test_util.ROOT_ENDPOINT
	logging.debug("hitting endpoint: {0}".format(url))
	response = requests.get(url)
	output = response.content
	if response.status_code != 200:
		logging.error("error when making get request: {0}".format(output))
	logging.info("output is: {0}".format(output))
	if output != test_util.ROOT_EXPECTED_OUTPUT:
		# TODO (nkubala): best way to handle error?
		# should probably raise "FailedTestException" that is caught by the driver
		logging.error("unexpected output: expected {0}, received {1}".format(test_util.ROOT_EXPECTED_OUTPUT, output))

#!/usr/bin/python

import logging
import requests

import test_util


def _test_root(base_url):
	url = base_url + test_util.ROOT_ENDPOINT
	logging.debug("hitting endpoint: {0}".format(url))
	response = requests.get(url)
	util._check_response(response, "error when making get request!")
	output = response.content
	logging.info("output is: {0}".format(output))
	if output != test_util.ROOT_EXPECTED_OUTPUT:
		# TODO (nkubala): best way to handle error?
		# should probably raise "FailedTestException" that is caught by the driver
		logging.error("unexpected output: expected {0}, received {1}".format(test_util.ROOT_EXPECTED_OUTPUT, output))

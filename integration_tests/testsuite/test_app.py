#!/usr/bin/python

import logging
import subprocess
import sys

TEST_ENDPOINTS = {
	"/": "Hello World!"
}

def _test_app(base_url):
	logging.info("starting app test with base url {0}".format(base_url))
	for k, v in TEST_ENDPOINTS.iteritems():
		url = base_url + k
		logging.debug("hitting endpoint: {0}".format(url))
		output = subprocess.check_output(['curl', '--silent', '-L', url])
		logging.info("output is: {0}".format(output))
		if output != v:
			# TODO: best way to handle error?
			logging.error("unexpected output")
		# assert output == v

if __name__ == "__main__":
	sys.exit(main())

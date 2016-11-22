#!/usr/bin/python

import logging
import subprocess
import sys

logging.getLogger().setLevel(logging.DEBUG)

BASE_URL = "https://nick-cloudbuild.appspot.com"

TEST_ENDPOINTS = {
	"/": "ruby app"
}

def main():
	hit_endpoints_and_check_output()


def hit_endpoints_and_check_output():
	for k, v in TEST_ENDPOINTS.iteritems():
		url = BASE_URL + k
		output = subprocess.check_output(['curl', '--silent', '-L', url])
		logging.info("output is: {0}".format(output))
		if output != v:
			# TODO: best way to handle error?
			logging.error("unexpected output")
		# assert output == v

if __name__ == "__main__":
	sys.exit(main())

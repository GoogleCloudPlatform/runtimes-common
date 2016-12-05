#!/usr/bin/python

import binascii
import os
import random
import string

LOGNAME_LENGTH = 20

ROOT_ENDPOINT = "/"
ROOT_EXPECTED_OUTPUT = "Hello World!"

LOGGING_ENDPOINT = "/logging"
MONITORING_ENDPOINT= "/monitoring"

def _generate_name_and_token():
	# TODO (nkubala): log directly to stdout since we're in GAE flex???
	# TODO (nkubala): possibly handle multiple log destinations depending on environments
	name = ''.join(random.choice(string.ascii_uppercase + string.ascii_lowercase) for i in range(LOGNAME_LENGTH))
	# name = 'stdout'
	token = binascii.b2a_hex(os.urandom(16))
	return name, token


def _generate_logging_payload():
  	logname, token = _generate_name_and_token()
  	data = {'log_name':logname, 'token':token}
  	return data


def _generate_metrics_payload():
	name, token = _generate_name_and_token()
	data = {'name':name, 'token':token}
	return data

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

METRIC_PREFIX = "custom.googleapis.com/{0}"


def _generate_name():
	# TODO (nkubala): log directly to stdout since we're in GAE flex???
	name = ''.join(random.choice(string.ascii_uppercase + string.ascii_lowercase) for i in range(LOGNAME_LENGTH))
	# name = 'stdout'
	return name


def _generate_hex_token():
	return binascii.b2a_hex(os.urandom(16))


def _generate_int64_token():
	return random.randint(-(2 ** 31), (2 ** 31)-1)


def _generate_logging_payload():
	data = {'log_name':_generate_name(), 'token':_generate_hex_token()}
	return data


def _generate_metrics_payload():
	data = {'name':METRIC_PREFIX.format(_generate_name()), 'token':_generate_int64_token()}
	return data

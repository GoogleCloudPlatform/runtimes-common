#!/usr/bin/python

# Copyright 2016 Google Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import binascii
import logging
import os
import random
import string

LOGNAME_LENGTH = 20

ROOT_ENDPOINT = "/"
ROOT_EXPECTED_OUTPUT = "Hello World!"

LOGGING_ENDPOINT = "/logging"
MONITORING_ENDPOINT= "/monitoring"

METRIC_PREFIX = "custom.googleapis.com/{0}"
METRIC_TIMEOUT = 60 # seconds
METRIC_PROPAGATION_TIME = 45 # seconds


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


def _check_response(response, error_message):
  if response.status_code - 200 >= 100: # 2xx
    logging.error("{0} exit code: {1}, text: {2}".format(error_message, response.status_code, response.text))

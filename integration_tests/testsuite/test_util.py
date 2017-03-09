#!/usr/bin/python

# Copyright 2017 Google Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the 'License');
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an 'AS IS' BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import binascii
from enum import Enum
import json
import logging
import os
import random
import requests
import string
import subprocess

requests.packages.urllib3.disable_warnings()

LOGNAME_LENGTH = 16

DEFAULT_TIMEOUT = 30  # seconds

ROOT_ENDPOINT = '/'
ROOT_EXPECTED_OUTPUT = 'Hello World!'

STANDARD_LOGGING_ENDPOINT = '/logging_standard'
CUSTOM_LOGGING_ENDPOINT = '/logging_custom'
MONITORING_ENDPOINT = '/monitoring'
EXCEPTION_ENDPOINT = '/exception'

METRIC_PREFIX = 'custom.googleapis.com/{0}'
METRIC_TIMEOUT = 60  # seconds


class Severity(Enum):
    DEBUG = 100
    INFO = 200
    WARNING = 400
    ERROR = 500
    CRITICAL = 600


def _generate_name():
    name = ''.join(random.choice(string.ascii_uppercase +
                   string.ascii_lowercase) for i in range(LOGNAME_LENGTH))
    return name


def _generate_hex_token():
    return binascii.b2a_hex(os.urandom(16))


def _generate_int64_token():
    return random.randint(-(2 ** 31), (2 ** 31)-1)


def generate_logging_payloads():
    payloads = []
    for l in list(Severity):
        payloads.append({
            'log_name': _generate_name(),
            'token': _generate_hex_token(),
            'level': l.name
            })
    return payloads


def generate_metrics_payload():
    data = {'name': METRIC_PREFIX.format(_generate_name()),
            'token': _generate_int64_token()}
    return data


def generate_exception_payload():
    data = {'token': _generate_int64_token()}
    return data


def _get(url, timeout=DEFAULT_TIMEOUT):
    logging.info('making get request to url {0}'.format(url))
    try:
        response = requests.get(url)
        return _check_response(response,
                               'error when making get ' +
                               'request! url: {0}'
                               .format(url))
    except Exception as e:
        logging.error('Error encountered when making get request!')
        logging.error(e)
        return None, 1


def _post(url, payload, timeout=DEFAULT_TIMEOUT):
    try:
        headers = {'Content-Type': 'application/json'}
        response = requests.post(url,
                                 json.dumps(payload),
                                 timeout=timeout,
                                 headers=headers)
        return _check_response(response, 'error when posting request! url: {0}'
                               .format(url))
    except requests.exceptions.Timeout:
        logging.error('POST to {0} timed out after {1} seconds!'
                      .format(url, timeout))
        return 'ERROR', 1


def _check_response(response, error_message):
    if response.status_code - 200 >= 100:  # 2xx
        logging.error('{0} exit code: {1}, text: {2}'
                      .format(error_message,
                              response.status_code,
                              response.text))
        return response.text, 1
    return response.text, 0


def _project_id():
    try:
        cmd = ['gcloud', 'config', 'list', '--format=json']
        entries = json.loads(subprocess.check_output(cmd))
        return entries.get('core').get('project')
    except Exception as e:
        logging.error('Error encountered when retrieving project id!')
        logging.error(e)


def get_default_url():
    return 'https://{0}.appspot.com'.format(_project_id())

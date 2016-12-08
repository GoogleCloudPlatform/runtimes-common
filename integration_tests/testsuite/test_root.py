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

import logging
import requests

import test_util


def _test_root(base_url):
  url = base_url + test_util.ROOT_ENDPOINT
  logging.debug("hitting endpoint: {0}".format(url))
  response = requests.get(url)
  test_util._check_response(response, "error when making get request!")
  output = response.content
  logging.info("output is: {0}".format(output))
  if output != test_util.ROOT_EXPECTED_OUTPUT:
    # TODO (nkubala): best way to handle error?
    # should probably raise "FailedTestException" that is caught by the driver
    logging.error("unexpected output: expected {0}, received {1}".format(test_util.ROOT_EXPECTED_OUTPUT, output))

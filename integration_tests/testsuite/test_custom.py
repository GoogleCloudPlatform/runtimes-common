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

import re
import json
import logging
import requests
import unittest
import urlparse

import test_util


class TestCustom(unittest.TestCase):
    """ This TestCase fetch a configuration from the endpoint '/custom' describing
        a series of tests, then run each of them and report their results.

        In the case where a test of the series fail, this TestCase will be
        considered as failed.

        The specification for the custom tests can be found at:
        https://github.com/GoogleCloudPlatform/runtimes-common/tree/master/integration_tests#custom-tests
    """

    def __init__(self, url, methodName='runTest'):
        self._base_url = url
        self._url = urlparse.urljoin(url, test_util.CUSTOM_ENDPOINT)
        unittest.TestCase.__init__(self)

    def runTest(self):
        """ Retrieve the configuration for the custom tests and launch the tests.

        :return: None.
        """
        logging.debug('Retrieving list of custom test endpoints.')
        output, status_code = test_util.get(self._url)
        self.assertEquals(status_code, 0,
                          'Cannot connect to sample application!')

        test_num = 0
        if not output:
            logging.debug('No custom tests specified.')
        else:
            for specification in json.loads(output):
                test_num += 1
                self._runTestForSpecification(specification, test_num)

    def _runTestForSpecification(self, specification, test_num):
        """ Given the specification for a test execute the steps and
            validate the result.

        :param specification: Dictionary containing the specification.
        :param test_num: Identifier of the test.
        :return: None.
        """
        name = specification.get('name', 'test_{0}'.format(test_num))
        timeout = specification.get('timeout', 500)
        path = specification.get('path')
        steps = specification.get('steps')
        validation = specification.get('validation')

        logging.info('Running custom test: %s', name)

        if path is not None:
            if steps is not None or validation is not None:
                logging.warn('When the field path is specified, the fields '
                             'validation and steps should not be present')
                return

            # Run the old test
            test_endpoint = urlparse.urljoin(self._base_url, path)
            response, _ = test_util.get(test_endpoint, timeout=timeout)
            logging.debug(response)
            return

        context = {'name': name}
        step_num = 0

        for step in steps:
            self._runStep(context, step, step_num)

        logging.debug("context : %s", json.dumps(context,
                                                 sort_keys=True,
                                                 indent=4,
                                                 separators=(',', ': ')))

        self._validate(context, validation)

    def _runStep(self, context, step, step_num):
        """ Use the provided step's configuration to send a request to the
            specified path and store the result into the context.

        :param context: A dictionary containing the context for the test.
        :param step: A dictionary containing the configuration of the step,
               this include:
                 name (optional): name of the step.
                 configuration (optional):
                    method: 'GET' or 'POST'.
                    headers: Dictionary containing the headers of the request.
                    content: Payload attached to the request.
                 path: Url of the request
        :param step_num: Index of the step.
        :return: None.
        """
        step_name = step.get('name', 'step_{0}'.format(step_num))
        configuration = step.get('configuration', dict())
        path = step.get('path')

        logging.info("Running step {0} of test {1}".format(
            step_name,
            context.get('name')
        ))

        response = requests.request(method=configuration.get('method', 'GET'),
                                    url=path,
                                    headers=configuration.get('headers'),
                                    data=configuration.get('content'))

        if 'application/json' in response.headers.get("Content-Type"):
            content = response.json()
        else:
            content = response.text

        context[step_name] = {
            'request': {
                'configuration': configuration,
                'path': path
            },
            'response': {
                'headers': dict(response.headers),
                'status': response.status_code,
                'content': content
            }
        }

    def _validate(self, context, specification):
        """ Compare the specification with the context and assert that every key
            present in the specification is also present in the context, and
            that the value associated to that key in the context match the
            regular expression specified in the specification.

        :param context: Dictionary containing for each step the request and
               the response.
        :param specification: Dictionary with the following fields:
               match: List of object containing:
                 key: Path in the context e.g step.response.headers.property .
                 pattern: Regular expression to be compared with the value
                          present at the path `key` in the context.
        :return: None.
        """

        match = specification.get('match')
        for test in match:
            key = test.get('key')
            value = self._evaluate_substitution(context, key)
            pattern = test.get('pattern')
            self.assertIsNotNone(re.search(pattern, value),
                                 "The value `{0}` for the key `{1}` "
                                 "do not match the pattern `{2}`"
                                 .format(value, key, pattern))

    def _evaluate_substitution(self, context, path):
        """ Search for the path `path` in the context and return the associated
            value.

            If the path is not valid the test is considered failed.

        :param context: A dictionary in which the key will be searched.
        :param path: A list of keys separated by dots, representing a path
                     in the context.
        :return: The value present in the context at the path `path`.
        """
        for key in path.split('.'):
            context = context.get(key)
            self.assertIsNotNone(context, "An error occurred during the "
                                          "substitution: the key {0} of path "
                                          "{1} is not present in the context"
                                          .format(key, path))
        return context

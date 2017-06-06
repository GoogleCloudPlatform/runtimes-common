#!/usr/bin/python

# Copyright 2017 Google Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import argparse
import logging
import sys
import unittest

import deploy_app
import test_custom
import test_exception
import test_logging_standard
import test_logging_custom
import test_monitoring
import test_root


def _main():
    logging.getLogger().setLevel(logging.INFO)

    parser = argparse.ArgumentParser()
    parser.add_argument('--image', '-i',
                        help='Newly-constructed base ' +
                        'image to build sample app on')
    parser.add_argument('--directory', '-d',
                        help='Root directory of sample app')
    parser.add_argument('--no-deploy',
                        action='store_false',
                        dest='deploy',
                        help='Flag to skip deployment of app ' +
                        '(must provide app URL)')
    parser.add_argument('--skip-standard-logging-tests',
                        action='store_false',
                        dest='standard_logging',
                        help='Flag to skip standard logging tests')
    parser.add_argument('--skip-custom-logging-tests',
                        action='store_false',
                        dest='custom_logging',
                        help='Flag to skip custom logging tests')
    parser.add_argument('--skip-monitoring-tests',
                        action='store_false',
                        dest='monitoring',
                        help='Flag to skip monitoring tests')
    parser.add_argument('--skip-exception-tests',
                        action='store_false',
                        dest='exception',
                        help='Flag to skip error reporting tests')
    parser.add_argument('--skip-custom-tests',
                        action='store_false',
                        dest='custom',
                        help='Flag to skip custom integration tests')
    parser.add_argument('--url', '-u',
                        help='URL where deployed app is ' +
                        'exposed (if applicable)')
    parser.add_argument('--verbose', '-v', action='store_true')
    args = parser.parse_args()

    if args.verbose:
        logging.getLogger().setLevel(logging.DEBUG)

    deploy_url = ''
    application_url = ''

    if args.deploy:
        if args.image is None:
            logging.error('Please specify base image name.')
            sys.exit(1)

        if args.directory is None:
            logging.error('Please specify at least one application to deploy.')
            sys.exit(1)

        logging.debug('Deploying app!')
        version, deploy_url = deploy_app.deploy_app(args.image, args.directory)

    application_url = args.url or deploy_url

    code = _test_app(application_url, args)

    if args.deploy:
        deploy_app.stop_app(version)

    return code


def _test_app(base_url, args):
    logging.info('Starting app test with base url {0}'.format(base_url))

    suite = unittest.TestSuite()

    suite.addTest(test_root.TestRoot(base_url))

    if args.standard_logging:
        suite.addTest(test_logging_standard.TestStandardLogging(base_url))

    if args.custom_logging:
        suite.addTest(test_logging_custom.TestCustomLogging(base_url))

    if args.monitoring:
        suite.addTest(test_monitoring.TestMonitoring(base_url))

    if args.exception:
        suite.addTest(test_exception.TestException(base_url))

    if args.custom:
        suite.addTest(test_custom.TestCustom(base_url))

    return not unittest.TextTestRunner(verbosity=4).run(suite).wasSuccessful()


if __name__ == '__main__':
    sys.exit(_main())

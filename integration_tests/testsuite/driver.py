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

import argparse
import logging
import sys

from args import Args
import deploy_app
import test_exception
import test_logging
import test_monitoring
import test_root


def _main():
    logging.getLogger().setLevel(logging.DEBUG)

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
    parser.add_argument('--no-logging',
                        action='store_false',
                        dest='logging',
                        help='Flag to skip logging tests')
    parser.add_argument('--no-monitoring',
                        action='store_false',
                        dest='monitoring',
                        help='Flag to skip monitoring tests')
    parser.add_argument('--no-exception',
                        action='store_false',
                        dest='exception',
                        help='Flag to skip error reporting tests')
    parser.add_argument('--url', '-u',
                        help='URL where deployed app is ' +
                        'exposed (if applicable)')
    args = Args(parser.parse_args())

    # this is only necessary until container builder's build robot is granted
    # the full scopes to do a deploy, and only when being run through
    # a container build

    # deploy_app._authenticate(args.directory)

    if args.deploy:
        if args.image is None:
            logging.error('Please specify base image name.')
            sys.exit(1)

        if args.directory is None:
            logging.error('Please specify at least one application to deploy.')
            sys.exit(1)

        logging.debug('Deploying app!')
        deploy_url = deploy_app._deploy_app(args.image, args.directory)
        if deploy_url is not '' and not deploy_url.startswith('https://'):
            deploy_url = 'https://' + deploy_url
        if args.url is None or args.url == '':
            args.url = deploy_url

    return _test_app(args)


def _test_app(args):
    base_url = args.url
    logging.info('Starting app test with base url {0}'.format(base_url))
    error_count = 0

    error_count += test_root._test_root(base_url)
    if args.logging:
        logging.info('Testing app logging')
        error_count += test_logging._test_logging(base_url)

    if args.monitoring:
        logging.info('Testing app monitoring')
        error_count += test_monitoring._test_monitoring(base_url)

    if args.exception:
        logging.info('Testing error reporting')
        error_count += test_exception._test_exception(base_url)

    return error_count


if __name__ == '__main__':
    sys.exit(_main())

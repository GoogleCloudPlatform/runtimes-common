#!/usr/bin/python

import argparse
import logging
import sys

import deploy_app
import test_root
import test_util
import test_monitoring
import test_logging

# TODO (nkubala): make this a configurable parameter from cloudbuild
# required to be paired with '--no-deploy'
DEFAULT_URL = "https://nick-cloudbuild.appspot.com"

def _main():
  logging.getLogger().setLevel(logging.DEBUG)

  parser = argparse.ArgumentParser()
  parser.add_argument('--image', '-i', help='Newly-constructed base image to build sample app on')

  parser.add_argument('--directory', '-d', help='Root directory of sample app')
  # TODO (nkubala): potentially add support for multiple app directories to deploy
  # parser.add_argument('--directory', '-d', help='Root directory of sample app', action='append')
  parser.add_argument('--no-deploy', action='store_true', help='Flag to skip deployment of app (must provide app URL)')
  parser.add_argument('--url', '-u', help='URL where deployed app is exposed (if applicable)', default=DEFAULT_URL)
  args = parser.parse_args()

  deploy_app._authenticate(args.directory)

  if not args.no_deploy:
    if args.image is None:
      logging.error("Please specify base image name!")
      sys.exit(1)

    if args.directory is None:
      logging.error("Please specify at least one application to deploy!")
      sys.exit(1)

    logging.debug("deploying app!")
    deploy_app._deploy_app(args.image, args.directory)

  _test_app(args.url)


def _test_app(base_url):
  # TODO (nkubala): check output from each test to log individul failures
  logging.info("starting app test with base url {0}".format(base_url))
  test_root._test_root(base_url)
  # test_logging._test_logging(base_url)
  test_monitoring._test_monitoring(base_url)


if __name__ == '__main__':
  sys.exit(_main())

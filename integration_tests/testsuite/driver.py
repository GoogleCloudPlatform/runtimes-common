#!/usr/bin/python

import argparse
import logging
import sys

import deploy_app
import test_app

DEFAULT_URL = "https://nick-cloudbuild.appspot.com"

def _main():
	logging.getLogger().setLevel(logging.DEBUG)

	parser = argparse.ArgumentParser()
	parser.add_argument('--image', '-i', help='Newly-constructed base image to build sample app on')

	parser.add_argument('--directory', '-d', help='Root directory of sample app')
	# TODO: potentially add support for multiple app directories to deploy
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

	test_app._test_app(args.url)

if __name__ == '__main__':
	sys.exit(_main())

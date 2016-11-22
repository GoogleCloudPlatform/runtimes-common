#!/usr/bin/python

import argparse
import logging
import os
import subprocess
import sys
import time

from shutil import copy

PROJECT_ID = "nick-cloudbuild"
DEPLOY_DELAY_SECONDS=20

def cleanup(appdir):
	try:
		os.remove(os.path.join(appdir, "Dockerfile"))
	except:
		pass


def _deploy_app(image, appdir):
	try:
		os.chdir(appdir)
		current_dir = os.path.realpath('.')
		try:
			copy("/app.yaml", current_dir)
		except:
			logging.error("error copying app.yaml from root dir!")
			sys.exit(1)

		try:
			os.remove("Dockerfile")
		except:
			pass

		# substitute vars in Dockerfile (equivalent of envsubst)
		with open("Dockerfile.in", 'r') as fin:
			with open("Dockerfile", 'a+') as fout:
				for line in fin:
					fout.write(line.replace('${STAGING_IMAGE}', image))
			fout.close()
		fin.close()

		subprocess.call(['gcloud', 'auth', 'activate-service-account', '--key-file=/auth.json'])
		subprocess.call(['gcloud', 'app', 'deploy', '--stop-previous-version', '--verbosity=debug'])

		print 'waiting {0} seconds for app to deploy'.format(DEPLOY_DELAY_SECONDS)
		for i in range(0,DEPLOY_DELAY_SECONDS):
			# sys.stdout.write('.')
			# sys.stdout.flush()
			time.sleep(1)
		print

	finally:
		cleanup(appdir)

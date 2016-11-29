#!/usr/bin/python

import argparse
import logging
import os
import subprocess
import sys
import time

from shutil import copy

PROJECT_ID = "nick-cloudbuild"
DEPLOY_DELAY_SECONDS = 20

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

		auth_command = ['gcloud', 'auth', 'activate-service-account', '--key-file=/auth.json']
		auth_proc = subprocess.Popen(auth_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

		output, error = auth_proc.communicate()
		if auth_proc.returncode != 0:
			sys.exit("Error encountered when authenticating. Full log: \n\n" + output)

		deploy_command = ['gcloud', 'app', 'deploy', '--stop-previous-version', '--verbosity=debug']
		deploy_proc = subprocess.Popen(deploy_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

		output, error = deploy_proc.communicate()
		if deploy_proc.returncode != 0:
			sys.exit("Error encountered when deploying app. Full log: \n\n" + output)

		print 'waiting {0} seconds for app to deploy'.format(DEPLOY_DELAY_SECONDS)
		for i in range(0, DEPLOY_DELAY_SECONDS):
			time.sleep(1)
		print

	finally:
		cleanup(appdir)

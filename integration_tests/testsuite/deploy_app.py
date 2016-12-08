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
import os
import subprocess
import sys
import time

from shutil import copy

# TODO (nkubala): make this configurable param from caller
PROJECT_ID = "nick-cloudbuild"
DEPLOY_DELAY_SECONDS = 30 # time to give GAE to start app after deploy


def cleanup(appdir):
  try:
    os.remove(os.path.join(appdir, "Dockerfile"))
  except:
    pass


def _authenticate(appdir):
  logging.debug("authentcating service account credentials...")
  auth_command = ['gcloud', 'auth', 'activate-service-account', '--key-file=/auth.json']
  subprocess.call(auth_command)

  p = subprocess.check_output(['gcloud', 'auth', 'list'])
  logging.info(p)

  # TODO (nkubala): make this work so we can handle failed authentication
  # auth_command = ['gcloud', 'auth', 'activate-service-account', '--key-file=/auth.json']
  # auth_proc = subprocess.Popen(auth_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

  # output, error = auth_proc.communicate()
  # if auth_proc.returncode != 0:
  #   sys.exit("Error encountered when authenticating. Full log: \n\n" + output)

  # authenticate the same service account INSIDE the same application
  # so it can write logs to this test driver's project
  try:
    owd = os.getcwd()
    os.chdir(appdir)
    current_dir = os.path.realpath('.')
    copy("/auth.json", current_dir)
  except:
    logging.error("error copying auth.json from root dir!")
    sys.exit(1)
  finally:
    os.chdir(owd)


def _deploy_app(image, appdir):
  try:
    owd = os.getcwd()
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


    deploy_command = ['gcloud', 'app', 'deploy', '--stop-previous-version', '--verbosity=debug']
    subprocess.call(deploy_command)

    # TODO (nkubala): make this work so we can handle failed deploys
    # deploy_proc = subprocess.Popen(deploy_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

    # output, error = deploy_proc.communicate()
    # if deploy_proc.returncode != 0:
    #   sys.exit("Error encountered when deploying app. Full log: \n\n" + output)

    print 'waiting {0} seconds for app to deploy...'.format(DEPLOY_DELAY_SECONDS)
    for i in range(0, DEPLOY_DELAY_SECONDS):
      time.sleep(1)
    print

  finally:
    cleanup(appdir)
    os.chdir(owd)

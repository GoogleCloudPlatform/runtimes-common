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

import json
import logging
import os
import subprocess
import sys
from retrying import retry

import test_util

DEPLOYED_VERSION = None


def _cleanup(appdir):
    try:
        os.remove(os.path.join(appdir, 'Dockerfile'))
    except Exception:
        pass


def deploy_app(image, appdir):
    try:
        # change to app directory (and remember original directory)
        owd = os.getcwd()
        os.chdir(appdir)

        # substitute vars in Dockerfile (equivalent of envsubst)
        with open('Dockerfile.in', 'r') as fin:
            with open('Dockerfile', 'w') as fout:
                for line in fin:
                    fout.write(line.replace('${STAGING_IMAGE}', image))
            fout.close()
        fin.close()

        global DEPLOYED_VERSION
        DEPLOYED_VERSION = test_util.generate_version()

        # TODO: once sdk driver is published, use it here
        deploy_command = ['gcloud', 'app', 'deploy', '--version',
                          DEPLOYED_VERSION, '--verbosity=debug']

        deploy_proc = subprocess.Popen(deploy_command,
                                       stdout=subprocess.PIPE,
                                       stdin=subprocess.PIPE)

        output, _ = deploy_proc.communicate()
        if deploy_proc.returncode != 0:
            sys.exit('Error encountered when deploying app. ' +
                     'Full log: \n\n' + (output or ''))

        return _retrieve_url()

    finally:
        _cleanup(appdir)
        os.chdir(owd)


def stop_app():
    if not DEPLOYED_VERSION:
        logging.error('App was never deployed!')
    logging.debug('Stopping application at version %s', DEPLOYED_VERSION)
    stop_command = ['gcloud', 'app', 'versions', 'stop', DEPLOYED_VERSION]
    stop_proc = subprocess.Popen(stop_command,
                                 stdout=subprocess.PIPE,
                                 stdin=subprocess.PIPE)

    output, _ = stop_proc.communicate()
    if stop_proc.returncode != 0:
        sys.exit('Error encountered when stopping version! ' +
                 'Full log: \n\n' + (output or ''))
    delete_command = ['gcloud', 'app', 'services', 'delete', 'default',
                      '--version', DEPLOYED_VERSION]

    delete_proc = subprocess.Popen(delete_command,
                                   stdout=subprocess.PIPE,
                                   stdin=subprocess.PIPE)

    output, _ = delete_proc.communicate()
    if stop_proc.returncode != 0:
        sys.exit('Error encountered when deleting version! ' +
                 'Full log: \n\n' + (output or ''))


@retry(wait_fixed=10000, stop_max_attempt_number=4)
def _retrieve_url():
    try:
        # retrieve url of deployed app for test driver
        url_command = ['gcloud', 'app', 'versions', 'describe',
                       DEPLOYED_VERSION, '--service',
                       'default', '--format=json']
        app_dict = json.loads(subprocess.check_output(url_command))
        return app_dict.get('versionUrl')
    except (subprocess.CalledProcessError, ValueError, KeyError):
        logging.warn('Error encountered when retrieving app URL!')
        return None
    raise Exception('Unable to contact deployed application!')

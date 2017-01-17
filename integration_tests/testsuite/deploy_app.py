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

import json
import logging
import os
import subprocess
import sys
import shutil
import time

DEPLOY_DELAY_SECONDS = 30  # time to give GAE to start app after deploy


def cleanup(appdir):
    try:
        os.remove(os.path.join(appdir, 'Dockerfile'))
    except:
        pass


def _deploy_app(image, appdir):
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

        # TODO: use sdk driver here
        deploy_command = ['gcloud', 'app', 'deploy',
                          '--stop-previous-version', '--verbosity=debug']

        deploy_proc = subprocess.Popen(deploy_command,
                                       stdout=subprocess.PIPE,
                                       stderr=subprocess.STDOUT)

        output, error = deploy_proc.communicate()
        if deploy_proc.returncode != 0:
            sys.exit('Error encountered when deploying app. ' +
                     'Full log: \n\n' + output)

        print 'waiting {0} seconds for ' \
              'app to deploy...'.format(DEPLOY_DELAY_SECONDS)

        time.sleep(DEPLOY_DELAY_SECONDS)

        try:
            # retrieve url of deployed app for test driver
            url_command = ['gcloud', 'app', 'describe', '--format=json']
            app_dict = json.loads(subprocess.check_output(url_command))
            hostname = app_dict.get('defaultHostname')
            if hostname is None:
                return ''
            return hostname.encode('ascii', 'ignore')
        except Exception:
            print 'Error encountered when retrieving app URL!'
            print 'Defaulting to provided URL parameter.'
            return ''

    finally:
        cleanup(appdir)
        os.chdir(owd)

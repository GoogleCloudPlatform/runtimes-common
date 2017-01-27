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
import os
from retrying import retry
import subprocess
import sys


def _cleanup(appdir):
    try:
        os.remove(os.path.join(appdir, 'Dockerfile'))
    except:
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

        # TODO: use sdk driver here
        deploy_command = ['gcloud', 'app', 'deploy',
                          '--stop-previous-version', '--verbosity=debug']

        deploy_proc = subprocess.Popen(deploy_command,
                                       stdout=subprocess.PIPE,
                                       stdin=subprocess.PIPE)

        output, error = deploy_proc.communicate()
        if deploy_proc.returncode != 0:
            sys.exit('Error encountered when deploying app. ' +
                     'Full log: \n\n' + (output or ''))

        return _retrieve_url()

    finally:
        _cleanup(appdir)
        os.chdir(owd)


@retry(wait_fixed=10000, stop_max_attempt_number=4)
def _retrieve_url():
    try:
        # retrieve url of deployed app for test driver
        url_command = ['gcloud', 'app', 'describe', '--format=json']
        app_dict = json.loads(subprocess.check_output(url_command))
        hostname = app_dict.get('defaultHostname')
        return hostname.encode('ascii', 'ignore')
    except (subprocess.CalledProcessError, ValueError, KeyError):
        print('Error encountered when retrieving app URL!')
        print('Defaulting to provided URL parameter.')
        return ''
    raise Exception('Unable to contact deployed application!')

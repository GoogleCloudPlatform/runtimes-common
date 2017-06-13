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

import logging
import os
import subprocess
import sys

import deploy_app
import test_util


def _deploy_app(appdir):
    try:
        # change to app directory (and remember original directory)
        owd = os.getcwd()
        os.chdir(appdir)

        deployed_version = test_util.generate_version()

        # TODO: once sdk driver is published, use it here
        deploy_command = ['gcloud', 'app', 'deploy',
                          '--version', deployed_version, '-q']

        subprocess.check_output(deploy_command)

        return deployed_version
    except subprocess.CalledProcessError as cpe:
        logging.error('Error encountered when deploying application! %s',
                      cpe.output)
        sys.exit(1)

    finally:
        deploy_app._cleanup(appdir)
        os.chdir(owd)


def _main(appdir):
    try:
        logging.debug('Testing runtime image.')
        version = _deploy_app(appdir)
        application_url = test_util.retrieve_url_for_version(version)
        output, status_code = test_util.get(application_url)
        if status_code:
            raise Exception('Error pinging application!')
    except Exception as e:
        logging.debug('{0}'.format(e))
        sys.exit(1)


if __name__ == '__main__':
    sys.exit(_main(sys.argv[1]))

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
import test_util


def _cleanup(appdir):
    try:
        os.remove(os.path.join(appdir, 'Dockerfile'))
    except Exception:
        pass


def _set_app_image(image):
    # substitute vars in Dockerfile (equivalent of envsubst)
    with open('Dockerfile.in', 'r') as fin:
        with open('Dockerfile', 'w') as fout:
            for line in fin:
                fout.write(line.replace('${STAGING_IMAGE}', image))
        fout.close()
    fin.close()


def deploy_app(image, appdir):
    try:
        # change to app directory (and remember original directory)
        owd = os.getcwd()
        os.chdir(appdir)

        # fills in image field in templated Dockerfile if image is specified
        if image:
            _set_app_image(image)

        deployed_version = test_util.generate_version()

        # TODO: once sdk driver is published, use it here
        deploy_command = ['gcloud', 'app', 'deploy', '--no-promote'
                          '--version', deployed_version, '-q']

        subprocess.check_output(deploy_command)

        return deployed_version
    except subprocess.CalledProcessError as cpe:
        logging.error('Error encountered when deploying application! %s',
                      cpe.output)
        sys.exit(1)

    finally:
        _cleanup(appdir)
        os.chdir(owd)


def deploy_app_without_image(appdir):
    return deploy_app(None, appdir)


def stop_app(deployed_version):
    logging.debug('Removing application version %s', deployed_version)
    try:
        delete_command = ['gcloud', 'app', 'services', 'delete', 'default',
                          '--version', deployed_version, '-q']

        subprocess.check_output(delete_command)
    except subprocess.CalledProcessError as cpe:
        logging.error('Error encountered when deleting app version! %s',
                      cpe.output)

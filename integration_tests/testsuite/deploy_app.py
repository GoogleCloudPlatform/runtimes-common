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


def _set_base_image(image):
    # substitute vars in Dockerfile (equivalent of envsubst)
    with open('Dockerfile.in', 'r') as fin:
        with open('Dockerfile', 'w') as fout:
            for line in fin:
                fout.write(line.replace('${STAGING_IMAGE}', image))


def _set_builder_image(builder):
    with open('test.yaml.in', 'r') as fin:
        with open('test.yaml', 'w') as fout:
            for line in fin:
                fout.write(line.replace('${STAGING_BUILDER_IMAGE}', builder))


def deploy_app(base_image, builder_image, appdir, yaml):
    try:
        if yaml:
            # convert yaml to absolute path before changing directory
            yaml = os.path.abspath(yaml)

        # change to app directory (and remember original directory)
        owd = os.getcwd()
        os.chdir(appdir)

        # fills in image field in templated Dockerfile and/or builder yaml
        if base_image:
            _set_base_image(base_image)
        if builder_image:
            _set_builder_image(builder_image)

        deployed_version = test_util.generate_version()

        # TODO: once sdk driver is published, use it here
        deploy_command = ['gcloud', 'app', 'deploy', '--no-promote',
                          '--version', deployed_version, '-q']
        if yaml:
            logging.info(yaml)
            deploy_command.append(yaml)

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
    return deploy_app(None, None, appdir, None)


def stop_app(deployed_version):
    logging.debug('Removing application version %s', deployed_version)
    try:
        delete_command = ['gcloud', 'app', 'services', 'delete', 'default',
                          '--version', deployed_version, '-q']

        subprocess.check_output(delete_command)
    except subprocess.CalledProcessError as cpe:
        logging.error('Error encountered when deleting app version! %s',
                      cpe.output)

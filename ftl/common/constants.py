# Copyright 2017 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

DEFAULT_LOG_LEVEL = 'NOTSET'

DEFAULT_DESTINATION_PATH = 'srv'
DEFAULT_ENTRYPOINT = None

# docker transport thread config
THREADS = 32

# ftl version
FTL_VERSION = "v0.3.1"

# cache constants
DEFAULT_TTL_WEEKS = 1

# Google Cloud Builder env options
BUILDER_OUTPUT = 'BUILDER_OUTPUT'
BUILDER_OUTPUT_FILE = 'output'

# Google Cloud Builder Args
GLOBAL_CACHE_REGISTRY = 'gcr.io/ftl-global-cache'

# php constants
PHP_CACHE_NAMESPACE = 'php-cache'
COMPOSER_LOCK = 'composer.lock'
COMPOSER_JSON = 'composer.json'

# node constants
NODE_CACHE_NAMESPACE = 'node-cache'
PACKAGE_LOCK = 'package-lock.json'
PACKAGE_JSON = 'package.json'
NODE_DEFAULT_ENTRYPOINT = 'node server.js'
NPMRC = '.npmrc'

# python constants
PIPFILE_LOCK = 'Pipfile.lock'
PIPFILE = 'Pipfile'
REQUIREMENTS_TXT = 'requirements.txt'
PYTHON_CACHE_NAMESPACE = 'python-cache'
VENV_DIR = '/env'
WHEEL_DIR = 'wheel'
PIP_DEFAULT_CMD = 'pip'
PYTHON_DEFAULT_CMD = 'python2.7'
VENV_DEFAULT_CMD = 'virtualenv'
PIP_OPTIONS = ['--disable-pip-version-check']

# logging constants
PHASE_1_CACHE_STR = '{key_version}:{language}->{key}'
PHASE_2_CACHE_STR = '{key_version}:{language}:{package_name}:' \
            '{package_version}->{key}'
CACHE_HIT = '[CACHE][HIT] '
CACHE_MISS = '[CACHE][MISS] '

PHASE_1_CACHE_HIT = CACHE_HIT + PHASE_1_CACHE_STR
PHASE_2_CACHE_HIT = CACHE_HIT + PHASE_2_CACHE_STR
PHASE_1_CACHE_MISS = CACHE_MISS + PHASE_1_CACHE_STR
PHASE_2_CACHE_MISS = CACHE_MISS + PHASE_2_CACHE_STR

CACHE_KEY_VERSION = 'v1'

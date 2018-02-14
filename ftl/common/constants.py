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

DEFAULT_LOG_LEVEL = "NOTSET"

DEFAULT_DESTINATION_PATH = 'srv'
DEFAULT_ENTRYPOINT = None

# cache constants
DEFAULT_TTL_WEEKS = 1
GLOBAL_CACHE_REGISTRY = 'gcr.io/ftl-global-cache'

# node constants
NODE_NAMESPACE = 'node-package-lock-cache'
PACKAGE_LOCK = 'package-lock.json'
PACKAGE_JSON = 'package.json'
NODE_DEFAULT_ENTRYPOINT = 'node server.js'

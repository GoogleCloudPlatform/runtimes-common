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


class Args():

    def __init__(self, args):
        self._image = args.image
        self._directory = args.directory
        self._deploy = args.deploy
        self._logging = args.logging
        self._monitoring = args.monitoring
        self._exception = args.exception
        self._url = args.url

    @property
    def image(self):
        return self._image

    @property
    def directory(self):
        return self._directory

    @property
    def deploy(self):
        return self._deploy

    @property
    def logging(self):
        return self._logging

    @property
    def monitoring(self):
        return self._monitoring

    @property
    def exception(self):
        return self._exception

    @property
    def url(self):
        return self._url

    @url.setter
    def url(self, value):
        self._url = value

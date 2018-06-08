# Copyright 2018 Google Inc. All Rights Reserved.
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

import os
import logging

from ftl.common import constants


class FTLErrors():
    @classmethod
    def USER(self):
        return "USER"

    @classmethod
    def INTERNAL(self):
        return "INTERNAL"


class UserError(Exception):
    def __init__(self, message):
        super(UserError, self).__init__(message)


class InternalError(Exception):
    def __init__(self, message):
        super(InternalError, self).__init__(message)


def UserErrorHandler(err, path, fail_on_error):
    logging.error(err)
    if path:
        with open(os.path.join(path, constants.BUILDER_OUTPUT_FILE), "w") as f:
            f.write("USER ERROR:\n%s" % str(err))
    if  fail_on_error:
        exit(1)
    else:
        exit(0)


def InternalErrorHandler(err, path):
    logging.error(err)
    if path:
        with open(os.path.join(path, constants.BUILDER_OUTPUT_FILE), "w") as f:
            f.write("INTERNAL ERROR:\n%s" % str(err))
    if  fail_on_error:
        exit(1)
    else:
        exit(0)

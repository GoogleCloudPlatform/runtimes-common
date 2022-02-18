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
import json
import logging
import hashlib

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


def genErrorId(s):
    return hashlib.sha256(s).hexdigest().upper()[:8]


def UserErrorHandler(err, path, fail_on_error):
    logging.error(err)
    if path:
        resp = {
            "error": {
                "errorType": constants.FTL_ERROR_TYPE,
                "canonicalCode": constants.FTL_USER_ERROR,
                "errorId": genErrorId(str(err)),
                "errorMessage": str(err)
            }
        }
        with open(os.path.join(path, constants.BUILDER_OUTPUT_FILE), "w") as f:
            f.write(json.dumps(resp))
    if fail_on_error:
        exit(1)
    else:
        exit(0)


def InternalErrorHandler(err, path, fail_on_error):
    logging.error(err)
    if path:
        resp = {
            "error": {
                "errorType": constants.FTL_ERROR_TYPE,
                "canonicalCode": constants.FTL_INTERNAL_ERROR,
                "errorId": genErrorId(str(err)),
                "errorMessage": str(err)
            }
        }
        with open(os.path.join(path, constants.BUILDER_OUTPUT_FILE), "w") as f:
            f.write(json.dumps(resp))
    if fail_on_error:
        exit(1)
    else:
        exit(0)

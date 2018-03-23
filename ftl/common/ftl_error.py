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
import shutil

from ftl.common import constants


class UserError(Exception):
    def __init__(self, message):
        super(UserError, self).__init__(message)
        logging.error(message)


class InternalError(Exception):
    def __init__(self, message):
        super(InternalError, self).__init__(message)
        logging.error(message)


def UserErrorHandler(u_err, log_path):
    if log_path:
        with open(os.path.join(log_path, constants.FTL_USER_LOG), "w") as f:
            f.write(str(u_err))


def InternalErrorHandler(log_path):
    if log_path:
        shutil.copyfile(
            os.path.join(log_path, constants.FTL_FULL_LOG),
            os.path.join(log_path, constants.FTL_INTERNAL_LOG))

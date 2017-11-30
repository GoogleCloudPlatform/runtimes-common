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
"""This package defines helpful utilities for FTL ."""
import time
import logging


class Timing(object):
    def __init__(self, descriptor):
        self.descriptor = descriptor

    def __enter__(self):
        self.start = time.time()
        return self

    def __exit__(self, unused_type, unused_value, unused_traceback):
        end = time.time()
        logging.info('%s took %d seconds', self.descriptor, end - self.start)

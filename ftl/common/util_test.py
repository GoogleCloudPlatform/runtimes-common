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
"""Unit tests for ftl_util"""

import unittest
import constants
import StringIO
import logging
import mock

import ftl_util
import logger


class UtilTest(unittest.TestCase):
    def setUp(self):
        args = mock.Mock()
        args.verbosity = "DEBUG"
        args.log_path = None

        logger.setup_logging(args)
        defaultLogger = logging.getLogger()

        self.log_capture_string = StringIO.StringIO()
        ch = logging.StreamHandler(self.log_capture_string)
        ch.setLevel(logging.DEBUG)

        defaultLogger.addHandler(ch)

    def test_parse_phase1_entry(self):
        language = 'PYTHON'
        key = 'test key'
        phase_1_entry = constants.PHASE_1_CACHE_HIT.format(
            key_version=constants.CACHE_KEY_VERSION,
            language=language,
            key=key
        )

        logging.info(phase_1_entry)
        entry = self.log_capture_string.getvalue()
        log_pieces = ftl_util.parseCacheLogEntry(entry)
        print log_pieces

        self.assertEqual(log_pieces['key'], key)
        self.assertEqual(log_pieces['language'], language)

    def test_parse_phase2_entry(self):
        language = 'PYTHON'
        package = 'flask'
        version = '==0.12.0'
        key = 'test key'
        phase_2_entry = constants.PHASE_2_CACHE_HIT.format(
            key_version=constants.CACHE_KEY_VERSION,
            language=language,
            package_name=package,
            package_version=version,
            key=key
        )

        logging.info(phase_2_entry)
        entry = self.log_capture_string.getvalue()
        log_pieces = ftl_util.parseCacheLogEntry(entry)
        print log_pieces

        self.assertEqual(log_pieces['key'], key)
        self.assertEqual(log_pieces['language'], language)
        self.assertEqual(log_pieces['package'], package)
        self.assertEqual(log_pieces['version'], version)

    def test_parse_bad_entry(self):
        logging.info('a malfomed, non-cache log entry')
        entry = self.log_capture_string.getvalue()
        log_pieces = ftl_util.parseCacheLogEntry(entry)

        self.assertEqual(log_pieces, None)


if __name__ == '__main__':
    unittest.main()

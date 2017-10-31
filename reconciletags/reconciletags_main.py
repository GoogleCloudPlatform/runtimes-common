# Copyright 2017 Google Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
Reads json files mapping docker digests to tags and reconciles them.
"""

import argparse
import json
import logging
import os
import shutil
import unittest
from containerregistry.tools import patched
from reconciletags import tag_reconciler
from reconciletags import data_integrity_test
from reconciletags import config_integrity_test


def create_temp_dir(files):
    dir = "../config/tag/"
    os.makedirs(dir)

    for file in files:
        if os.path.isfile(file):
            shutil.copy(file, dir)
        else:
            raise AssertionError("{0} is not a valid file".format(file))


def delete_temp_dir():
    shutil.rmtree("../config/")


def run_config_integrity_test(files):
    create_temp_dir(files)
    suite = unittest.TestLoader().loadTestsFromTestCase(
        config_integrity_test.ReconcilePresubmitTest)
    unittest.TextTestRunner().run(suite)


def run_data_integrity_test(files):
    create_temp_dir(files)
    suite = unittest.TestLoader().loadTestsFromTestCase(
        data_integrity_test.DataIntegrityTest)
    unittest.TextTestRunner().run(suite)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--dry-run', dest='dry_run',
                        help='Runs tests to make sure input files are valid, \
                        and runs a dry run of the reconciler',
                        action='store_true', default=False)
    parser.add_argument('files',
                        help='The files to run the reconciler on',
                        nargs='+')
    parser.add_argument('--data-integrity', dest='data_integrity',
                        help='Runs a test to make sure the data in the input \
                        files is the same as in prod',
                        action='store_true', default=False)
    args = parser.parse_args()
    logging.basicConfig(level=logging.DEBUG)

    if args.data_integrity:
        try:
            run_data_integrity_test(args.files)
        finally:
            delete_temp_dir()
        return

    if args.dry_run:
        try:
            run_config_integrity_test(args.files)
        finally:
            delete_temp_dir()

    r = tag_reconciler.TagReconciler()
    for f in args.files:
        logging.debug('---Processing {0}---'.format(f))
        with open(f) as tag_map:
            data = json.load(tag_map)
            r.reconcile_tags(data, args.dry_run)


if __name__ == '__main__':
    with patched.Httplib2():
        main()

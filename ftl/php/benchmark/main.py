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

import subprocess
import argparse
import datetime
import time
import os
import logging
from google.cloud import bigquery

DATASET_NAME = 'ftl_benchmark'
TABLE_NAME = 'ftl_php_benchmark'
PROJECT_NAME = 'ftl-node-test'

parser = argparse.ArgumentParser(
    description='Run FTL php benchmarks.')

parser.add_argument(
    '--base', action='store', help=('The name of the docker base image.'))

parser.add_argument(
    '--name', action='store', help=('The name of the docker image to push.'))

parser.add_argument(
    '--directory',
    action='store',
    help='The path where the application data sits.')

parser.add_argument(
    '--iterations',
    action='store',
    type=int,
    help='Number of times to build the image')

parser.add_argument(
    '--description', action='store',
    help=('Description of the app being benchmarked.'))


def _record_build_times_to_bigquery(build_times, description):
    current_date = datetime.datetime.now()
    logging.info('Retrieving bigquery client')
    client = bigquery.Client(project=PROJECT_NAME)

    dataset_ref = client.dataset(DATASET_NAME)
    table_ref = dataset_ref.table(TABLE_NAME)
    table = client.get_table(table_ref)

    logging.info('Adding build time data to bigquery table')
    rows = [(current_date, description, bt[0], bt[1]) for bt in build_times]
    client.create_rows(table, rows)
    logging.info('Finished adding build time data to bigquery table')


def main():
    args = parser.parse_args()
    logging.getLogger().setLevel("NOTSET")
    logging.basicConfig(
        format='%(asctime)s.%(msecs)03d %(levelname)-8s %(message)s',
        datefmt='%Y-%m-%d,%H:%M:%S')
    build_times = []
    logging.info('Beginning building php images')
    for _ in range(args.iterations):
        start_time = time.time()

        # For the binary
        php_builder_path = 'ftl/php_builder.par'
        # For the image
        if not os.path.isfile(php_builder_path):
            php_builder_path = ("./ftl/php/benchmark/php_image."
                                "binary.runfiles/__main__/ftl/"
                                "php_builder.par")
        # For container builder
        if not os.path.isfile(php_builder_path):
            php_builder_path = 'bazel-bin/ftl/php_builder.par'

        cmd = subprocess.Popen([php_builder_path,
                                '--base', args.base,
                                '--name', args.name,
                                '--directory', args.directory,
                                '--no-cache'], stderr=subprocess.PIPE)
        _, output = cmd.communicate()

        build_time = round(time.time() - start_time, 2)
        build_times.append((build_time, output))

    logging.info('Beginning recording build times to bigquery')
    _record_build_times_to_bigquery(build_times, args.description)


if __name__ == '__main__':
    main()

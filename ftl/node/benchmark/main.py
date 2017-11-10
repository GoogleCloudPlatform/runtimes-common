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
TABLE_NAME = 'ftl_benchmark'
PROJECT_NAME = 'ftl-node-test'
NUM_ITERATIONS = 1

parser = argparse.ArgumentParser(
    description='Run FTL node benchmarks.')

parser.add_argument(
    '--base', action='store', help=('The name of the docker base image.'))

parser.add_argument(
    '--name', action='store', help=('The name of the docker image to push.'))

parser.add_argument(
    '--directory',
    action='store',
    help='The path where the application data sits.')

parser.add_argument(
    '--repo', action='store', help=('The repo being tested on.'))


def _record_build_times_to_bigquery(build_times, repo):
    current_date = datetime.datetime.now()
    logging.info('Retrieving bigquery client')
    client = bigquery.Client(project=PROJECT_NAME)

    dataset_ref = client.dataset(DATASET_NAME)
    table_ref = dataset_ref.table(TABLE_NAME)
    table = client.get_table(table_ref)

    logging.info('Adding build time data to bigquery table')
    rows = [(current_date, repo, build_time) for build_time in build_times]
    client.create_rows(table, rows)


def _print_data_in_table():
    client = bigquery.Client(project=PROJECT_NAME)
    dataset_ref = client.dataset(DATASET_NAME)
    table_ref = dataset_ref.table(TABLE_NAME)
    table = client.get_table(table_ref)
    for row in client.list_rows(table):  # API request
        print(row)


def main():
    args = parser.parse_args()
    logging.getLogger().setLevel("NOTSET")
    logging.basicConfig(
        format='%(asctime)s.%(msecs)03d %(levelname)-8s %(message)s',
        datefmt='%Y-%m-%d,%H:%M:%S')
    build_times = []
    logging.info('Beginning building node images')
    for _ in range(NUM_ITERATIONS):
        start_time = time.time()

        # Path for the binary
        node_builder_path = 'ftl/node_builder.par'

        # Path for the image
        if not os.path.isfile(node_builder_path):
            node_builder_path = ("./ftl/node/benchmark/node_benchmark_image."
                                "binary.runfiles/__main__/ftl/"
                                "node_builder.par")

        subprocess.check_call([node_builder_path,
                              '--base', args.base,
                               '--name', args.name,
                               '--directory', args.directory,
                               '--no-cache'])
        build_time = round(time.time() - start_time, 2)
        build_times.append(build_time)
    logging.info('Beginning recording build times to bigquery')
    _record_build_times_to_bigquery(build_times, args.repo)
    _print_data_in_table()


if __name__ == '__main__':
    main()

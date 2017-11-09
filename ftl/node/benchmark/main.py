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
import uuid
import logging
from google.cloud import bigquery

DATASET_NAME = 'ftl_benchmark'
TABLE_NAME = 'ftl_benchmark_timestamp'
PROJECT_NAME = 'priya-wadhwa'
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
    '--benchmark', action='store', help=('The size of the app.'))


def _record_build_times_to_bigquery(build_times, benchmark):
    current_date = datetime.datetime.now()
    client = bigquery.Client(project=PROJECT_NAME)

    dataset_ref = client.dataset(DATASET_NAME)
    table_ref = dataset_ref.table(TABLE_NAME)
    table = client.get_table(table_ref)

    print('Adding build time data to bigquery table')
    rows = [(current_date, benchmark, build_time) for build_time in build_times]
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
    build_times = []
    for _ in range(NUM_ITERATIONS):
        start_time = time.time()
        subprocess.check_call(['./ftl/node_builder',
                              '--base', args.base,
                              '--name', args.name,
                              '--directory', args.directory, 
                              '--no-cache'])
        build_time = round(time.time() - start_time, 2)
        build_times.append(build_time)
    _record_build_times_to_bigquery(build_times, args.benchmark)
    _print_data_in_table()
        
if __name__ == '__main__':
    main()
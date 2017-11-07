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
from google.cloud import bigquery

DATASET_NAME = 'ftl_benchmark'
TABLE_NAME = 'ftl_benchmark'

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

def _record_build_times_to_bigquery(build_times, bechmark):
    current_date = datetime.datetime.now()

    for build_time in build_times:
        row = [(current_date, benchmark, build_time)]
        project = "priya-wadhwa"
        logging.debug('Fetching bigquery client for project %s', project)
        client = bigquery.Client(project=project)
        dataset = client.dataset(DATASET_NAME)
        logging.debug('Writing bigquery data to table %s in dataset %s',
                    TABLE_NAME, dataset)
        table = bigquery.Table(name=TABLE_NAME, dataset=dataset)
        table.reload()
        return table.insert_data(row)

def main():
    args = parser.parse_args()
    build_times = []
    for _ in range(2):
        start_time = time.time()
        subprocess.check_call(['./ftl/node_builder',
                              '--base', args.base,
                              '--name', args.name,
                              '--directory', args.directory, 
                              '--no-cache'])
        build_time = round(time.time() - start_time, 2)
        build_times.append(build_time)
        _record_build_times_to_bigquery(build_times, args.benchmark)
        
if __name__ == '__main__':
    main()
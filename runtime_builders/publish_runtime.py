#!/usr/bin/python

# Copyright 2016 Google Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import argparse
from datetime import datetime
import json
import logging
import regex
from ruamel import yaml
import subprocess
import sys

from google.cloud import storage

TAG_REGEX = '(?<=:|@)(.*)'


def main():
    logging.getLogger().setLevel(logging.INFO)
    parser = argparse.ArgumentParser()
    parser.add_argument('--infile', '-i',
                        help='templated cloudbuild config file.',
                        required=True)
    parser.add_argument('--bucket', '-b',
                        help='GCS bucket to publish runtime to',
                        default='runtime-builders')
    parser.add_argument('--builder-name',
                        help='the name of the runtime or project '
                        'associated with this builder',
                        required=True)
    args = parser.parse_args()

    templated_file_contents = _resolve_tags(args.infile)
    logging.info(templated_file_contents)
    logging.info(_publish_to_gcs(templated_file_contents,
                                 args.builder_name,
                                 args.bucket))


def _resolve_tags(config_file):
    """
    Given a templated YAML cloudbuild config file, parse it, resolve image tags
    on each build step's image to the corresponding digest, and write new
    config with fully qualified images to temporary file for upload to GCS.

    Keyword arguments:
    config_file -- string representing path to
    templated cloudbuild YAML config file

    Return value:
    path to temporary file containing fully qualified config file, to be
    published to GCS.
    """
    with open(config_file, 'r') as infile:
        try:
            config = yaml.round_trip_load(infile)

            for step in config.get('steps'):
                image = step.get('name')
                templated_step = _resolve_tag(image)
                step['name'] = templated_step

            return yaml.round_trip_dump(config)
        except yaml.YAMLError as e:
            logging.error(e)
            sys.exit(1)


def _resolve_tag(image):
    """
    Given a path to a tagged Docker image in GCR, replace the tag with its
    corresponding sha256 digest.
    """
    if ':' not in image:
        logging.error('Image \'{0}\' must contain explicit tag or '
                      'digest!'.format(image))
        sys.exit(1)
    elif 'sha256' in image:
        return image
    else:
        m = regex.search(TAG_REGEX, image)
        if m:
            target_tag = m.string[m.start():m.end()]
            base_image = m.string[0:m.start()-1]
        else:
            logging.error('Error when parsing image tag!')
            sys.exit(1)

    command = ['gcloud', 'beta', 'container', 'images',
               'list-tags', base_image, '--format=json']

    try:
        output = subprocess.check_output(command)
        entries = json.loads(output)
        for image in entries:
            for tag in image.get('tags'):
                if tag == target_tag:
                    digest = image.get('digest')
                    return base_image + '@' + digest
        logging.error('Tag {0} not found on image {1}!'
                      .format(target_tag, base_image))
        sys.exit(1)
    except subprocess.CalledProcessError as e:
        logging.error(e)

    logging.error('No digest found for tag {0} on '
                  'image {1}'.format(target_tag, base_image))


def _publish_to_gcs(builder_file_contents, builder_name, bucket):
    """
    Given a cloudbuild YAML config file, publish the file to a bucket in GCS.
    """
    client = storage.Client()
    runtime_bucket = client.get_bucket(bucket)

    if runtime_bucket is None:
        logging.error('Bucket {0} not found!'.format(bucket))

    file_name = '{0}-builder-{1}.yaml'.format(
        builder_name,
        datetime.now().strftime('%Y%m%d%H%M%S'))

    blob = storage.Blob(file_name, runtime_bucket)
    blob.upload_from_string(builder_file_contents)

    full_path = 'gs://{0}/{1}'.format(bucket, file_name)

    return full_path


if __name__ == '__main__':
    sys.exit(main())

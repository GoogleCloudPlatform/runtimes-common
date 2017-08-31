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

"""Reads json files mapping docker digests to tags and reconciles them.

Reads all json files in current directory and parses it into repositories
and tags. Calls gcloud container images add-tag on each entry.
If there are no changes that api call is no-op.
"""

import argparse
import json
import logging
import os
import subprocess
import httplib2


from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_image_list
from containerregistry.client.v2_2 import docker_session
from containerregistry.client import docker_creds
from containerregistry.transport import transport_pool
from containerregistry.client import docker_name

class TagReconciler:
    def call(self, command, dry_run, fmt="json"):
        command += " --format=" + fmt
        if not dry_run:
            logging.debug('Running {0}'.format(command))
            output = subprocess.check_output([command], shell=True)
            logging.debug(output)
            return output
        else:
            logging.debug('Would have run {0}'.format(command))

    def add_tags(self, digest, tag, dry_run): 
        if not dry_run:
            src_name = docker_name.Digest(digest)
            dest_name = docker_name.Tag(tag)
            creds = docker_creds.DefaultKeychain.Resolve(src_name)
            transport = transport_pool.Http(httplib2.Http)

            with docker_image.FromRegistry(src_name, creds, transport) as src_img:
                if src_img.exists():
                    logging.debug('Tagging {0} with {1}'.format(digest, tag))
                    with docker_session.Push(dest_name, creds, transport) as push:
                        push.upload(src_img)
        else:
            logging.debug("Would have tagged {0} with {1}".format(digest, tag))

    # This turns a list of lists into one flat list of tags
    def flatten_tags_list(self, list_of_lists):
        flat_tags_list = []
        for sublist in list_of_lists:
            for tag in sublist:
                if tag:
                    flat_tags_list.append(tag)
        return flat_tags_list
    
    def get_existing_tags(self, full_repo, digest):
        full_digest = full_repo + '@sha256:' + digest
        existing_tags = []

        name = docker_name.Digest(full_digest)
        creds = docker_creds.DefaultKeychain.Resolve(name)
        transport = transport_pool.Http(httplib2.Http)

        with docker_image.FromRegistry(name, creds, transport) as img:
            if img.exists():
                existing_tags = img.tags()
        
        return existing_tags
    
    def get_latest_digest(self, full_repo, manifests):
        for digest in manifests:
            if "latest" in manifests[digest]['tag']:
                return digest
    
    def reconcile_tags(self, data, dry_run):
        self.call('gcloud config list', False)
        for project in data['projects']:
            default_registry = project['base_registry']

            registries = project.get('additional_registries', [])
            registries.append(default_registry)
            for registry in registries:
                full_repo = os.path.join(registry, project['repository'])
                default_repo = os.path.join(default_registry, project['repository'])

                for image in project['images']:
                    full_digest = full_repo + '@sha256:' + image['digest']
                    full_tag = full_repo + ':' + image['tag']
                    existing_tags = []

                    name = docker_name.Digest(full_digest)
                    creds = docker_creds.DefaultKeychain.Resolve(name)
                    transport = transport_pool.Http(httplib2.Http)

                    with docker_image.FromRegistry(name, creds, transport) as img:
                        if img.exists():
                            existing_tags = img.tags()
                            logging.debug("Existing Tags: {0}".format(existing_tags))
                            logging.debug(full_repo)

                            manifests = img.manifests()
                            latest = self.get_latest_digest(full_repo, manifests)
                            
                            # Don't retag latest if it's already latest
                            if latest:
                                if latest.startswith('sha256:'):
                                    latest = latest[len('sha256:'):]
                                if (image['tag'] == 'latest'
                                and latest.startswith(image['digest'])):
                                    logging.debug('Skipping tagging %s as latest as '
                                                'it is already latest.',
                                                image['digest'])
                                    continue
                            
                        self.add_tags(full_digest, full_tag, dry_run)

                logging.debug(self.get_existing_tags(full_repo, project['images'][0]['digest']))

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--dry-run', dest='dry_run',
                        action='store_true', default=False)
    parser.add_argument('files',
                        help='The files to run the reconciler on',
                        nargs='+')
    args = parser.parse_args()
    logging.basicConfig(level=logging.DEBUG)
    r = TagReconciler()
    for f in args.files:
        logging.debug('---Processing {0}---'.format(f))
        with open(f) as tag_map:
            data = json.load(tag_map)
            r.reconcile_tags(data, args.dry_run)


if __name__ == '__main__':
    main()

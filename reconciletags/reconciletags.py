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
from containerregistry.client import docker_creds
from containerregistry.client import docker_name
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_session
from containerregistry.transport import transport_pool
from containerregistry.tools import patched
import httplib2


class TagReconciler:

    def add_tags(self, digest, tag, dry_run):
        if dry_run:
            logging.debug("Would have tagged {0} with {1}".format(digest, tag))
            return

        src_name = docker_name.Digest(digest)
        dest_name = docker_name.Tag(tag)
        creds = docker_creds.DefaultKeychain.Resolve(src_name)
        transport = transport_pool.Http(httplib2.Http)

        with docker_image.FromRegistry(src_name, creds, transport) as src_img:
            if src_img.exists():
                creds = docker_creds.DefaultKeychain.Resolve(dest_name)
                logging.debug('Tagging {0} with {1}'.format(digest, tag))
                with docker_session.Push(dest_name, creds, transport) as push:
                        push.upload(src_img)
            else:
                logging.debug("""Unable to tag {0}
                    as the image can't be found""".format(digest))

    def get_existing_tags(self, full_repo, digest):
        full_digest = full_repo + '@sha256:' + digest
        existing_tags = []

        name = docker_name.Digest(full_digest)
        creds = docker_creds.DefaultKeychain.Resolve(name)
        transport = transport_pool.Http(httplib2.Http)

        with docker_image.FromRegistry(name, creds, transport) as img:
            if img.exists():
                existing_tags = img.tags()
            else:
                logging.debug(
                    """Unable to get existing tags for {0}
                        as the image can't be found""".format(full_digest))
        return existing_tags

    def get_tagged_digest(self, manifests, tag):
        for digest in manifests:
            if tag in manifests[digest]['tag']:
                return digest
        return ""

    def get_digest(self, repo, prefix):
        name = docker_name.Repository(repo)
        creds = docker_creds.DefaultKeychain.Resolve(name)
        transport = transport_pool.Http(httplib2.Http)

        with docker_image.FromRegistry(name, creds, transport) as img:
            digests = [d[len('sha256:'):] for d in img.manifests()]
            matches = [d for d in digests if d.startswith(prefix)]
            if len(matches) == 1:
                return matches.pop()
            elif len(matches) > 1:
                logging.debug('{0} is not a unique digest prefix'.format(
                                                                  prefix))
        return ""

    def reconcile_tags(self, data, dry_run):
        for project in data['projects']:

            default_registry = project['base_registry']
            registries = project.get('additional_registries', [])
            registries.append(default_registry)

            for registry in registries:
                full_repo = os.path.join(registry, project['repository'])

                for image in project['images']:
                    digest = self.get_digest(full_repo, image['digest'])
                    if not digest:
                        continue
                    full_digest = full_repo + '@sha256:' + digest
                    full_tag = full_repo + ':' + image['tag']

                    name = docker_name.Digest(full_digest)
                    creds = docker_creds.DefaultKeychain.Resolve(name)
                    transport = transport_pool.Http(httplib2.Http)

                    with docker_image.FromRegistry(name, creds,
                                                   transport) as img:
                            if not img.exists():
                                logging.debug('Could not retrieve  ' +
                                              '{0}'.format(full_digest))
                                return

                            existing_tags = img.tags()
                            logging.debug('Existing Tags: ' +
                                          '{0}'.format(existing_tags))

                            manifests = img.manifests()
                            tagged_digest = self.get_tagged_digest(
                                manifests, image['tag'])

                            # Don't retag an image if the tag already exists
                            if tagged_digest.startswith('sha256:'):
                                tagged_digest = tagged_digest[len('sha256:'):]
                            if tagged_digest.startswith(digest):
                                logging.debug('Skipping tagging %s with %s as '
                                              'that tag already exists.',
                                              digest, image['tag'])
                                continue

                            self.add_tags(full_digest, full_tag, dry_run)

                    logging.debug(self.get_existing_tags(full_repo, digest))


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
    with patched.Httplib2():
        main()

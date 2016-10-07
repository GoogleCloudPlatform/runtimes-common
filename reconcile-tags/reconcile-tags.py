"""Reads json files mapping docker digests to tags and reconciles them.

Reads all json files in current directory and parses it into repositories
and tags. Calls gcloud beta container images add-tag on each entry.
If there are no changes that api call is no-op.
"""

import argparse
import glob
import json
import logging
import subprocess


def call(command, dry_run):
    if not dry_run:
        logging.debug('Running {0}'.format(command))
        output = subprocess.check_output([command], shell=True)
        logging.debug(output)
        return output
    else:
        logging.debug('Would have run {0}'.format(command))


def add_tags(digest, tag, dry_run):
    logging.debug('Tagging {0} with {1}'.format(digest, tag))
    command = ('gcloud beta container images add-tag {0} {1} '
               '-q'.format(digest, tag))
    call(command, dry_run)


def delete_tag(repo, tag, dry_run):
    full_tag = repo + ":" + tag
    logging.debug('Removing {0}'.format(full_tag))
    command = ('gcloud beta container images delete {0} -q'.format(full_tag))
    call(command, dry_run)


def reconcile_tags(dry_run):
    call('gcloud config list', False)
    files = glob.glob('./*.json')
    for f in files:
        logging.debug('---Processing {0}---'.format(f))
        with open(f) as tag_map:
            data = json.load(tag_map)
            for repo, images in data.items():
                # list-tags doesn't allow hyphens in repository names
                # so convert those to underscores since they're
                # guaranteed to be the same in GCR
                output = call('gcloud beta container images list-tags '
                              '--format=\'value(tags)\' {0}'
                              .format(repo.replace('-', '_')), False)

                list_of_tags = [tag.split(',') for tag in output.split('\n')]
                existing_tags = list(filter(None, (tag.rstrip()
                                                   for sublist in list_of_tags
                                                   for tag in sublist)))
                logging.debug(existing_tags)

                reconciled_tags = list()
                for image in images:
                    digest = image['digest']
                    tag = image['tag']
                    full_digest = repo + '@sha256:' + digest
                    full_tag = repo + ':' + tag
                    add_tags(full_digest, full_tag, dry_run)
                    reconciled_tags.append(tag)

                for t in list(set(existing_tags) - set(reconciled_tags)):
                    delete_tag(repo, t, dry_run)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--dry-run', dest="dry_run",
                        action='store_true', default=False)
    args = parser.parse_args()
    logging.basicConfig(level=logging.DEBUG)
    reconcile_tags(args.dry_run)

if __name__ == '__main__':
    main()

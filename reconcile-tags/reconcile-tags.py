"""Reads a json file mapping docker digests to tags and reconciles them.

Reads all json files in current directory and parses it into repositories
and tags. Calls gcloud beta container images add-tag on each entry.
If there are no changes that api call is no-op.
"""

import argparse
import glob
import json
from subprocess import call


def reconcile_tags(dry_run):
    call(['gcloud', 'config', 'list'])
    files = glob.glob('./*.json')
    for f in files:
        print '\nProcessing {0}'.format(f)
        with open(f) as tag_map:
            data = json.load(tag_map)
            for repo, images in data.items():
                for image in images:
                    digest = image['digest']
                    tag = image['tag']
                    full_digest = repo + '@sha256:' + digest
                    full_tag = repo + ':' + tag
                    print ('\nTagging {0} with {1}'
                           .format(full_digest, full_tag))
                    if not dry_run:
                        print ('Running gcloud beta container images '
                               'add-tag {0} {1} -q'
                               .format(full_digest, full_tag))
                        call(['gcloud', 'beta', 'container', 'images',
                              'add-tag', full_digest, full_tag, '-q'])


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--dryrun', action='store_true', default=False)
    args = parser.parse_args()
    reconcile_tags(args.dryrun)

if __name__ == '__main__':
    main()

"""Reads json files mapping docker digests to tags and reconciles them.

Reads all json files in current directory and parses it into repositories
and tags. Calls gcloud beta container images add-tag on each entry.
If there are no changes that api call is no-op.
"""

import argparse
import glob
import json
import subprocess


def call(digest, tag, dry_run):
    print '\nTagging {0} with {1}'.format(digest, tag)
    if not dry_run:
        command = 'gcloud beta container images add-tag'
        ' {0} {1} -q'.format(digest, tag)
        print 'Running {0}'.format(command)
        output = subprocess.check_output([command], shell=True)
        print output


def reconcile_tags(dry_run):
    subprocess.call(['gcloud', 'config', 'list'])
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
                    call(full_digest, full_tag, dry_run)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--dry-run', dest="dry_run",
                        action='store_true', default=False)
    args = parser.parse_args()
    reconcile_tags(args.dry_run)

if __name__ == '__main__':
    main()

"""Reads a json file mapping docker digests to tags and reconciles them.

Reads tap-map.json, currently a file in the same folder as this script
and parses the json. Calls gcloud beta container images add-tag on each entry.
If there are no changes that api call is no-op.
"""

import glob
import json
from subprocess import call


def main():
    files = glob.glob('./*.json')
    for f in files:
        print f
        with open(f) as tag_map:
            data = json.load(tag_map)
            for repo, images in data.items():
                print repo
                for image in images:
                    print image
                    digest = image['digest']
                    tag = image['tag']
                    print digest
                    print tag
                    call(['gcloud', 'beta', 'container', 'images',
                          'add-tag', repo+'@sha256:'+digest,
                          repo+':'+tag, '-q'])

if __name__ == '__main__':
    main()

"""Checks the specified image for security vulnerabilities."""

import argparse
import json
import logging
import sys
import subprocess

_GCLOUD_CMD = ['gcloud', 'beta', 'container', 'images', '--format=json']


def _run_gcloud(cmd):
    full_cmd = _GCLOUD_CMD + cmd
    output = subprocess.check_output(full_cmd)
    return json.loads(output)


def _check_image(image):
    digest = _resolve_latest(image)
    full_name = '%s@%s' % (image, digest)
    parsed = _run_gcloud(['describe', full_name])
    vulnz = (parsed['total_vulnerability_found'] -
             parsed['not_fixed_vulnerability_count'])
    if vulnz:
        logging.info('Found %s unpatched vulnerabilities in %s. Run '
                     '[gcloud beta container images describe %s] '
                     'to see the full list.',
                     vulnz, image, full_name)
    return vulnz


def _resolve_latest(image):
    parsed = _run_gcloud(['list-tags', image, '--no-show-occurrences'])
    for digest in parsed:
        if 'latest' in digest['tags']:
            return digest['digest']
    raise Exception("Unable to find digest of 'latest' tag for %s" % image)


def _main():
    parser = argparse.ArgumentParser()
    parser.add_argument('image', help='The image to test')
    args = parser.parse_args()

    logging.basicConfig(level=logging.DEBUG)
    return _check_image(args.image)


if __name__ == '__main__':
    sys.exit(_main())

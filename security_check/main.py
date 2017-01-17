"""Checks the specified image for security vulnerabilities."""

import argparse
import json
import logging
import sys
import subprocess


_GCLOUD_CMD = ['gcloud', 'beta', 'container', 'images', '--format=json']


# Severities
_LOW = 'LOW'
_MEDIUM = 'MEDIUM'
_HIGH = 'HIGH'
_CRITICAL = 'CRITICAL'

_SEV_MAP = {
    _LOW: 0,
    _MEDIUM: 1,
    _HIGH: 2,
    _CRITICAL: 3,
}


def _run_gcloud(cmd):
    full_cmd = _GCLOUD_CMD + cmd
    output = subprocess.check_output(full_cmd)
    return json.loads(output)


def _check_image(image, severity, whitelist):
    digest = _resolve_latest(image)
    full_name = '%s@%s' % (image, digest)
    parsed = _run_gcloud(['describe', full_name])

    unpatched = 0
    for vuln in parsed.get('vulz_analysis', []):
        if vuln.get('patch_not_available'):
            continue
        if vuln.get('vulnerability') in whitelist:
            continue
        if _filter_severity(vuln['severity'], severity):
            unpatched += 1

    if unpatched:
        base_unpatched = 0
        img = parsed.get('image_analysis')
        if img:
            base_img_url = img[0]['base_image_url']
            base_image = base_img_url[len('https://'):base_img_url.find('@')]
            base_unpatched = _check_image(base_image, severity, whitelist)
        unpatched -= base_unpatched
        logging.info('Found %s unpatched vulnerabilities in %s. Run '
                     '[gcloud beta container images describe %s] '
                     'to see the full list.',
                     unpatched, image, full_name)
    return unpatched


def _resolve_latest(image):
    parsed = _run_gcloud(['list-tags', image, '--no-show-occurrences'])
    for digest in parsed:
        if 'latest' in digest['tags']:
            return digest['digest']
    raise Exception("Unable to find digest of 'latest' tag for %s" % image)


def _filter_severity(sev1, sev2):
    """Returns whether sev1 is higher than sev2"""
    DEFAULT = _SEV_MAP.get(_LOW)
    return _SEV_MAP.get(sev1, DEFAULT) >= _SEV_MAP.get(sev2, DEFAULT)


def _main():
    parser = argparse.ArgumentParser()
    parser.add_argument('image', help='The image to test')
    parser.add_argument('--severity',
                        choices=[_LOW, _MEDIUM, _HIGH, _CRITICAL],
                        default=_MEDIUM,
                        help='The minimum severity to filter on.')
    parser.add_argument('--whitelist-file', dest='whitelist',
                        help='The path to the whitelist json file',
                        default='whitelist.json')
    args = parser.parse_args()

    logging.basicConfig(level=logging.DEBUG)

    try:
        whitelist = json.load(open(args.whitelist, 'r'))
    except IOError:
        whitelist = []
    logging.info(whitelist)

    return _check_image(args.image, args.severity, whitelist)


if __name__ == '__main__':
    sys.exit(_main())

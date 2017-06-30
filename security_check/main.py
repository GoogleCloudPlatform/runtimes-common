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

"""Checks the specified image for security vulnerabilities."""

import apt_pkg
import argparse
import json
import logging
import re
import sys
import subprocess


_GCLOUD_CMD = ['gcloud', 'container', 'images', '--format=json']


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

_SUB_MAP = {
    'launcher.gcr.io/google': 'gcr.io/google-appengine',
    'l.gcr.io/google': 'gcr.io/google-appengine'
}


def _sub_image(full_image):
    repo, image = full_image.rsplit('/', 1)
    repo = _SUB_MAP.get(repo, repo)
    new_image = '/'.join((repo, image))
    if new_image != full_image:
        logging.info('Checking %s instead of %s', new_image, full_image)
    return new_image


def _run_gcloud(cmd):
    full_cmd = _GCLOUD_CMD + cmd
    output = subprocess.check_output(full_cmd)
    return json.loads(output)


def _find_base_image(image):
    parsed = _run_gcloud(['describe', image])
    img = parsed.get('image_analysis')
    if img:
        base_img_url = img[0]['base_image_url']
        base_img = base_img_url[len('https://'):base_img_url.find('@')]
        return _sub_image(base_img)


def _check_for_vulnz(image, severity, whitelist):
    unpatched = _check_image(image, severity, whitelist)
    if not unpatched:
        return unpatched

    base_image = _find_base_image(image)
    base_unpatched = {}
    if base_image:
        base_unpatched = _check_image(base_image, severity, whitelist)
    else:
        logging.info("Could not find base image for %s", image)

    count = 0
    for k, vuln in unpatched.items():
        if k not in base_unpatched.keys():
            count += 1
            logging.info(_format_vuln(vuln))
        else:
            logging.info('Vulnerability %s exists in the base '
                         'image. Skipping.', vuln)

    if count > 0:
        logging.info('Found %s unpatched vulnerabilities in %s. Run '
                     '[gcloud container images describe %s] '
                     'to see the full list.', count, image, image)

    return unpatched


def _format_vuln(vuln):
    return '''
Vulnerability found.
CVE: {0}
SEVERITY: {1}
PACKAGES: {2}
FIXED PACKAGES: {3}
    '''.format(
        vuln['vulnerability'],
        vuln['severity'],
        ' '.join([v['affected_package'] for v in vuln['pkg_vulnerabilities']]),
        ' '.join(v['fixed_package'] for v in vuln['pkg_vulnerabilities']))


def _check_image(image, severity, whitelist):
    parsed = _run_gcloud(['describe', image])

    unpatched = {}
    for vuln in parsed.get('vulz_analysis', []):
        if vuln.get('patch_not_available'):
            continue
        if not _check_vuln_is_valid(vuln):
            continue
        if vuln.get('vulnerability') in whitelist:
            logging.info('Vulnerability %s is whitelisted. Skipping.',
                         vuln.get('vulnerability'))
            continue
        if _filter_severity(vuln['severity'], severity):
            unpatched[vuln['vulnerability']] = vuln

    return unpatched


def _filter_severity(sev1, sev2):
    """Returns whether sev1 is higher than sev2"""
    DEFAULT = _SEV_MAP.get(_LOW)
    return _SEV_MAP.get(sev1, DEFAULT) >= _SEV_MAP.get(sev2, DEFAULT)


def _check_vuln_is_valid(vuln):
    for pkg in vuln.get('pkg_vulnerabilities', []):
        if 'affected_package' in pkg and 'fixed_package' in pkg:
            # Parse the version out of the "package_name (version)" string.
            version_re = r'.*\((.*)\)'
            affected_version = re.match(
                version_re,
                pkg.get('affected_package')).groups()[0]
            fixed_version = re.match(
                version_re,
                pkg.get('fixed_package')).groups()[0]
            if apt_pkg.version_compare(fixed_version, affected_version) > 0:
                return True
    logging.info('Vulnerability %s is already fixed. '
                 'The affected package: %s is greater '
                 'than the fixed package: %s',
                 vuln.get('vulnerability'),
                 affected_version,
                 fixed_version)
    return False


def security_check(image, severity=_MEDIUM, whitelist_file='whitelist.json'):
    apt_pkg.init()

    try:
        whitelist = json.load(open(whitelist_file, 'r'))
    except IOError:
        whitelist = []
    logging.info("whitelist=%s", whitelist)

    return _check_for_vulnz(_sub_image(image), severity, whitelist)


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
    return len(security_check(args.image, args.severity, args.whitelist))


if __name__ == '__main__':
    sys.exit(_main())

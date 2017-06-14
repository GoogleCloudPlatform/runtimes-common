#!/usr/bin/python

import argparse
import json
import logging
import semver
import sys
import requests
from requests.auth import HTTPBasicAuth

requests.packages.urllib3.disable_warnings()

# OWNER = 'GoogleCloudPlatform'
OWNER = 'nkubala'
REPO = 'runtimes-common'
USER = 'XXX'
PW = 'XXX'


API_BASE = 'https://api.github.com'
REPO_BASE = '{api_base}/repos/{owner}/{repo}'
REPO_URL = REPO_BASE.format(api_base=API_BASE, owner=OWNER, repo=REPO)

COMMIT_URL = REPO_URL + '/commits'
COMMITISH_URL = COMMIT_URL + '/master'

SHA_URL = COMMIT_URL + '/master'
RELEASE_URL = REPO_URL + '/releases'
LATEST_RELEASE_URL = RELEASE_URL + '/latest'

COMMITISH_HEADER = {'Accept': 'application/vnd.github.VERSION.sha'}

sess = requests.Session()
sess.auth = HTTPBasicAuth(USER, PW)


def main():
    logging.getLogger().setLevel(logging.INFO)

    parser = argparse.ArgumentParser()
    parser.add_argument('--old', '-o', help='old image')
    parser.add_argument('--new', '-n', help='new image')
    parser.add_argument('--verbose', '-v', action='store_true')
    args = parser.parse_args()
    # image_diff = _diff_images(args.old, args.new)
    commits = _retrieve_commit_messages()
    # package_diff = _run_package_analysis(args.old, args.new)

    # TODO: combine retrieved info into real release notes 
    full_notes = "here are some notes"

    # _create_release(full_notes)


def _diff_images(old_image, new_image):
    '''
    most likely use the image differ tool created by the interns here
    we could potentially remove the drydock querying in favor of using this
    depending on what kind of language level diffing the differ tools has
    '''
    return 0


def _retrieve_commit_messages():
    # TODO: only make this call once and cache somewhere
    latest_release = sess.get(LATEST_RELEASE_URL).content
    latest_release = json.loads(latest_release)
    prev_sha = latest_release.get('target_commitish')
    logging.info(prev_sha)

    prev_commit = json.loads(sess.get(COMMIT_URL + '/{0}'.format(prev_sha)).content)
    prev_timestamp = prev_commit.get('commit').get('author').get('date')
    commits = sess.get(COMMIT_URL, params={'since':prev_timestamp}).content
    logging.info(json.dumps(json.loads(commits), indent=4))

    '''
    1. use releases API to retrieve commit hash from last release
    2. retrieve all commit hashes between then and HEAD
    3. loop through commits, retrieve messages
    4. (maybe) parse through each commit and retrieve relevant information?
    '''
    return []


def _run_package_analysis(old_image, new_image):
    '''
    here we want to use drydock to retrieve package analysis information for
    each image, then diff this info to see differences in installed packages

    unclear if we want to do language specific diffing here: this may just come
    from information provided manually by the maintainers
    '''
    return 0


def _create_release(release_notes):
    release_payload = _generate_release_payload(release_notes)
    logging.info('Posting to url {0}'.format(RELEASE_URL))
    response = sess.post(RELEASE_URL, data=json.dumps(release_payload))
    # response = requests.get(RELEASE_URL)
    if response.status_code < 200 or response.status_code > 299:
        logging.error('Error when creating release (code {0})'
                      .format(response.status_code))
        logging.error(response.text)
        sys.exit(1)
    logging.info(json.dumps(json.loads(response.text), indent=4))
    return 0


def _generate_release_payload(release_notes):
    try:
        logging.debug('getting latest release from url: %s', LATEST_RELEASE_URL)
        latest_release = sess.get(LATEST_RELEASE_URL).content
        latest_release = json.loads(latest_release)
        logging.debug(latest_release)

        prev_tag = latest_release['tag_name'].replace('v', '')
        tag = 'v' + semver.bump_minor(prev_tag)

        commitish = sess.get(COMMITISH_URL, headers=COMMITISH_HEADER).content
    except (TypeError, KeyError, ValueError) as e:
        logging.error('Error encountered when retrieving latest version! %s', e)
        sys.exit(1)


    return {
        "tag_name": tag,
        "target_commitish": commitish,
        "name": tag,
        "body": release_notes,
        "draft": True
    }


if __name__ == '__main__':
    main()

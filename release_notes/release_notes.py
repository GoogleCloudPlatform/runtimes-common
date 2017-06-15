#!/usr/bin/python

import argparse
import json
import logging
import semver
import sys
import requests
from requests.auth import HTTPBasicAuth

requests.packages.urllib3.disable_warnings()


class ReleaseNotes():

    def __init__(self, args):
        self.owner = args.owner
        self.repo = args.repo
        self.old_image = args.old
        self.new_image = args.new

        API_BASE = 'https://api.github.com'
        REPO_BASE = '{api_base}/repos/{owner}/{repo}'
        self.REPO_URL = REPO_BASE.format(api_base=API_BASE,
                                         owner=args.owner,
                                         repo=args.repo)

        self.COMMIT_URL = self.REPO_URL + '/commits'
        self.COMMITISH_URL = self.COMMIT_URL + '/master'

        self.SHA_URL = self.COMMIT_URL + '/master'
        self.RELEASE_URL = self.REPO_URL + '/releases'
        self.LATEST_RELEASE_URL = self.RELEASE_URL + '/latest'

        self.COMMITISH_HEADER = {'Accept':
                                 'application/vnd.github.VERSION.sha'}

        self.sess = requests.Session()
        self.sess.auth = HTTPBasicAuth(args.user, args.password)

        # retrieve latest release for use later on
        latest_release = self.sess.get(self.LATEST_RELEASE_URL).content
        self.latest_release = json.loads(latest_release)

    def diff_images(self):
        '''
        most likely use the image differ tool created by the interns here
        we could potentially remove the drydock querying in favor of using this
        depending on what kind of language level diffing the differ tools has
        '''
        # TODO: implement
        return ''

    def retrieve_commit_messages(self):
        prev_sha = self.latest_release.get('target_commitish')

        # get timestamp of last release,then retrieve all commits since
        commit = json.loads(self.sess.get(self.COMMIT_URL +
                                          '/{0}'.format(prev_sha)).content)
        prev_timestamp = commit.get('commit').get('author').get('date')
        params = {'since': prev_timestamp}
        raw_commits = json.loads(self.sess.get(self.COMMIT_URL,
                                               params=params).content)

        # TODO: possibly parse each commit and retrieve relevant information?
        return '\n'.join(['* ' + str(c.get('commit').get('message'))
                          for c in raw_commits])

    def run_package_analysis(self):
        '''
        use drydock to retrieve package analysis information for each image,
        then diff this info to see differences in installed packages

        unclear if we want to do language specific diffing here: this may just
        come from information provided manually by the maintainers
        '''
        # TODO: implement
        return ''

    def create_release(self, release_notes):
        release_payload = self._generate_release_payload(release_notes)
        logging.debug('Posting to url {0}'.format(self.RELEASE_URL))
        response = self.sess.post(self.RELEASE_URL,
                                  data=json.dumps(release_payload))
        if response.status_code < 200 or response.status_code > 299:
            logging.error('Error when creating release (code {0})'
                          .format(response.status_code))
            logging.error(response.text)
            sys.exit(1)
        logging.info('Drafted release created. View at: {0}'
                     .format(json.loads(response.text).get('html_url')))

    def _generate_release_payload(self, release_notes):
        try:
            # TODO: support non-semver versions and tags
            prev_tag = self.latest_release['tag_name'].replace('v', '')
            tag = 'v' + semver.bump_minor(prev_tag)

            commitish = self.sess.get(self.COMMITISH_URL,
                                      headers=self.COMMITISH_HEADER).content
        except (TypeError, KeyError, ValueError) as e:
            logging.error('Error encountered when retrieving '
                          'latest version! %s', e)
            sys.exit(1)

        return {
            "tag_name": tag,
            "target_commitish": commitish,
            "name": tag,
            "body": release_notes,
            "draft": True
        }


def main():
    logging.getLogger().setLevel(logging.INFO)

    parser = argparse.ArgumentParser()
    parser.add_argument('--old', '-o', help='old image', required=True)
    parser.add_argument('--new', '-n', help='new image', required=True)
    parser.add_argument('--user',
                        help='Github robot username',
                        required=True)
    parser.add_argument('--password',
                        help='Github robot password',
                        required=True)
    parser.add_argument('--owner',
                        help='Github project owner',
                        required=True)
    parser.add_argument('--repo', '-r',
                        help='target Github repository to publish release',
                        required=True)
    parser.add_argument('--verbose', '-v', action='store_true')
    args = parser.parse_args()

    release_notes = ReleaseNotes(args)

    if args.verbose:
        logging.getLogger().setLevel(logging.DEBUG)

    image_diff = release_notes.diff_images()
    commits = release_notes.retrieve_commit_messages()
    package_diff = release_notes.run_package_analysis()

    full_notes = '''**Notes:**

**Image Diff Information**:
{image_diff}

**Package Diff Information**:
{package_diff}

**Commits since last release**:
{commits}
'''.format(image_diff=image_diff,
           package_diff=package_diff,
           commits=commits)

    release_notes.create_release(full_notes)


if __name__ == '__main__':
    main()

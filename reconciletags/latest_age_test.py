"""Latest age tests.

Checks the build date of the image marked as latest for a repository and fails
if it's over two weeks old."""

import datetime
import glob
import json
import logging
import os
import subprocess
import unittest


class LatestAgeTest(unittest.TestCase):

    def _get_latest_timestamp(self, repo):
        images = json.loads(
                  subprocess.check_output(['gcloud', 'beta', 'container',
                                           'images', 'list-tags',
                                           '--no-show-occurrences',
                                           '--format=json', repo]))
        for image in images:
            if 'latest' in image['tags']:
                return datetime.datetime.strptime(
                    image['timestamp']['datetime'], '%Y-%m-%d %H:%M:%S')
        return None

    def test_latest_age(self):
        old_repos = []
        invalid_repos = []
        for f in glob.glob('../config/tag/*.json'):
            # We don't care about how old these images are
            if f.endswith('runtimes_common.json'):
                continue
            logging.debug('Testing {0}'.format(f))
            with open(f) as tag_map:
                data = json.load(tag_map)
                for project in data['projects']:
                    full_repo = os.path.join(project['base_registry'],
                                             project['repository'])
                    last_deploy = self._get_latest_timestamp(full_repo)
                    if last_deploy is None:
                        invalid_repos.append(full_repo)
                    threshold = last_deploy + datetime.timedelta(weeks=2)
                    if threshold < datetime.datetime.now():
                        old_repos.append(full_repo)

        if len(old_repos) > 0 or len(invalid_repos) > 0:
            msg = ''
            if len(old_repos) > 0:
                msg += ('The following repos have not been deployed in '
                        'over two weeks: {0}. '.format(str(old_repos)))

            if len(invalid_repos) > 0:
                msg += ('The following repos have no latest tag: {0}.'
                        .format(str(invalid_repos)))

            self.fail(msg)


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    unittest.main()

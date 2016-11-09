"""Tests for data_integrity."""


import glob
import json
import logging
import os
import subprocess
import unittest

class DataIntegrityTest(unittest.TestCase):

    def _get_real_data(self, repo):
        return json.loads(
                subprocess.check_output(['gcloud', 'beta', 'container',
                                         'images', 'list-tags',
                                         '--no-show-occurrences',
                                         '--format=json', repo]))

    def test_data_consistency(self):
        failed_digests = []
        for f in glob.glob('../config/*.json'):
            logging.debug('Testing {0}'.format(f))
            with open(f) as tag_map:
                data = json.load(tag_map)
                for project in data['projects']:
                    full_repo = os.path.join(project['base_registry'],
                                             project['repository'])
                    real_digests = self._get_real_data(full_repo)
                    for image in project['images']:
                        for i in real_digests:
                            digest = i['digest']
                            if i['digest'].split(':')[1].startswith(
                                image['digest']):
                                if image['tag'] not in i['tags']:
                                    failed_digests.append({full_repo: image})

        if len(failed_digests) > 0:
            self.fail('These entries do not correspond with what is'
                      ' currently live:' + str(failed_digests))


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    unittest.main()

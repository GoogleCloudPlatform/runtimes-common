"""Tests to check the integrity of json config files.

These tests assume that the json configs live in a top-level
folder named config."""

import glob
import json
import logging
import os
import subprocess
import unittest


class ReconcilePresubmitTest(unittest.TestCase):

    # This function requires gcloud be installed and authenticated
    # with your project
    def _get_digests(self, repo):
        try:
            output = json.loads(
                subprocess.check_output(['gcloud', 'beta', 'container',
                                         'images', 'list-tags',
                                         '--no-show-occurrences',
                                         '--format=json', repo]))
            # grab the digest for each image and strip off the 'sha256:'
            # for matching purposes
            digests = [image['digest'].split(':')[1] for image in output]
            return digests
        except OSError as e:
            logging.error(e)
            self.fail('Make sure gcloud is installed and properly '
                      'authenticated')

    def test_json_structure(self):
        for f in glob.glob('../config/*.json'):
            logging.debug('Testing {0}'.format(f))
            with open(f) as tag_map:
                data = json.load(tag_map)
                for project in data['projects']:
                    self.assertEquals(project['base_registry'], 'gcr.io')
                    for registry in project['additional_registries']:
                        self.assertRegexpMatches(registry, '^.*gcr.io$')
                    self.assertIsNotNone(project['repository'])
                    for image in project['images']:
                        self.assertIsInstance(image, dict)
                        self.assertIsNotNone(image['digest'])
                        self.assertIsNotNone(image['tag'])

    def test_digests_are_real(self):
        for f in glob.glob('../config/*.json'):
            logging.debug('Testing {0}'.format(f))
            with open(f) as tag_map:
                data = json.load(tag_map)
                for project in data['projects']:
                    default_registry = project['base_registry']
                    full_repo = os.path.join(default_registry,
                                             project['repository'])
                    logging.debug('Checking {0}'.format(full_repo))
                    digests = self._get_digests(full_repo)
                    for image in project['images']:
                        logging.debug('Checking {0}'
                                      .format(image['digest']))
                        self.assertTrue(any(
                                        digest.startswith(image['digest'])
                                        for digest in digests))


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    unittest.main()

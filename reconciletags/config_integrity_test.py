"""Tests to check the integrity of json config files.

These tests assume that the json configs live in a top-level
folder named config."""

import glob
import json
import subprocess
import unittest


class ReconcilePresubmitTest(unittest.TestCase):

    def _get_digests(self, repo):
        output = json.loads(
            subprocess.check_output(['gcloud', 'beta', 'container',
                                     'images', 'list-tags',
                                     '--no-show-occurrences',
                                     '--format=json', repo]))
        # grab the digest for each image and strip off the 'sha256:'
        # for matching purposes
        digests = [image['digest'][7:] for image in output]
        return digests

    def test_json_structure(self):
        for f in glob.glob('../config/*.json'):
            print "Testing {0}".format(f)
            with open(f) as tag_map:
                data = json.load(tag_map)
                for repo, images in data.items():
                    self.assertRegexpMatches(repo, 'gcr.io/.*')
                    for image in images:
                        self.assertIsInstance(image, dict)
                        self.assertIsNotNone(image['digest'])
                        self.assertIsNotNone(image['tag'])

    def test_digests_are_real(self):
        for f in glob.glob('../config/*.json'):
            print "Testing {0}".format(f)
            with open(f) as tag_map:
                data = json.load(tag_map)
                for repo, images in data.items():
                    digests = self._get_digests(repo)
                    for image in images:
                        print "Checking {0}".format(image['digest'])
                        self.assertTrue(any(digest.startswith(image['digest'])
                                            for digest in digests))


if __name__ == '__main__':
    unittest.main()

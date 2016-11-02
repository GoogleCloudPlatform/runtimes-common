"""End to end test for the tag reconciler."""

import json
import os
import subprocess
import sys
import unittest


class ReconciletagsE2eTest(unittest.TestCase):

    _BUCKET = ''
    _FILE_NAME = 'e2e_test.json'
    _DIR = 'tiny_docker_image/'
    _REPO = 'gcr.io/google-appengine-qa/e2etest'
    _TAG = 'initial'
    _TEST_JSON = """
    {{
      "projects":[
        {{
          "registries": ["gcr.io"],
          "repository": "google-appengine-qa/e2etest",
          "images": [
                  {{
                      "digest": "{0}",
                      "tag": "testing"
                  }}
          ]
        }}
      ]
    }}
    """

    def _ListTags(self, repo):
        output = json.loads(
            subprocess.check_output(['gcloud', 'beta', 'container',
                                     'images', 'list-tags',
                                     '--no-show-occurrences',
                                     '--format=json', repo]))
        return output

    def _BuildImage(self, full_image_name, bucket):
        # create a non-functional but tiny docker image
        subprocess.call(['gcloud', 'beta', 'container', 'builds',
                         'submit', self._DIR, '-q', '--tag', full_image_name,
                         '--gcs-log-dir', bucket + '/logs',
                         '--gcs-source-staging-dir', bucket + '/staging'])

        # create the json config file and write to it
        test_json = open(self._FILE_NAME, 'w')

        # grab the just created digest
        output = self._ListTags(self._REPO)
        digests = [image['digest'].split(':')[1] for image in output]
        self.assertEqual(len(digests), 1)
        self.digest = digests[0]

        # write the proper json to the config file
        test_json.write(self._TEST_JSON.format(self.digest))

    def setUp(self):
        self._BuildImage(self._REPO + ':' + self._TAG, self._BUCKET)

    def tearDown(self):
        subprocess.call(['gcloud', 'beta', 'container', 'images',
                         'delete', self._REPO + '@sha256:' + self.digest,
                         '-q'])
        os.remove(self._FILE_NAME)

    def testTagReconciler(self):
        # run the reconciler
        subprocess.check_output(['python', 'reconciletags.py',
                                 self._FILE_NAME])

        # check list-tags to see if it added the correct tag
        output = self._ListTags(self._REPO)
        for image in output:
            if image['digest'].split(':')[1] == self.digest:
                self.assertEquals(len(image['tags']), 2)
                self.assertEquals(image['tags'][1], 'testing')
                self.assertEquals(image['tags'][0], self._TAG)


def usage():
    print "Usage: python reconciletags_e2e_test.py <bucket>"
    sys.exit(1)

if __name__ == '__main__':
    if (len(sys.argv) < 2):
        usage()
    ReconciletagsE2eTest._BUCKET = sys.argv.pop(1)
    unittest.main()

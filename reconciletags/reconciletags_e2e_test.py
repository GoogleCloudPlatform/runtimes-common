"""End to end test for the tag reconciler."""

import json
import os
import reconciletags
import shutil
import subprocess
import tempfile
import unittest


class ReconciletagsE2eTest(unittest.TestCase):

    _FILE_NAME = 'e2e_test.json'
    _DIR = 'tiny_docker_image/'
    _REPO = 'gcr.io/gcp-runtimes/reconciler-e2etest'
    _TAG = 'initial'
    _TEST_JSON = {
      "projects": [
        {
          "registries": ["gcr.io"],
          "repository": "gcp-runtimes/reconciler-e2etest",
          "images": [
                  {
                      "digest": "",
                      "tag": "testing"
                  }
          ]
        }
      ]
    }

    def _ListTags(self, repo):
        output = json.loads(
            subprocess.check_output(['gcloud', 'beta', 'container',
                                     'images', 'list-tags',
                                     '--no-show-occurrences',
                                     '--format=json', repo]))
        return output

    def _BuildImage(self, full_image_name):
        # create a non-functional but tiny docker image
        subprocess.call(['gcloud', 'beta', 'container', 'builds',
                         'submit', self._DIR, '-q', '--tag', full_image_name])

        # create the json config file and write to it
        self.tmpdir = tempfile.mkdtemp()
        self.full_filename = os.path.join(self.tmpdir, self._FILE_NAME)
        test_json = open(self.full_filename, 'w')

        # grab the just created digest
        output = self._ListTags(self._REPO)
        self.assertEqual(len(output), 1)
        self.digest = output.pop()['digest'].split(':')[1]

        # write the proper json to the config file
        self._TEST_JSON['projects'][0]['images'][0]['digest'] = self.digest
        json.dump(self._TEST_JSON, test_json)

    def setUp(self):
        self._BuildImage(self._REPO + ':' + self._TAG)

    def tearDown(self):
        subprocess.call(['gcloud', 'beta', 'container', 'images',
                         'delete', self._REPO + '@sha256:' + self.digest,
                         '-q'])
        shutil.rmtree(self.tmpdir)

    def testTagReconciler(self):
        # run the reconciler
        subprocess.check_output(['python', 'reconciletags.py',
                                 self.full_filename])

        # check list-tags to see if it added the correct tag
        output = self._ListTags(self._REPO)
        for image in output:
            if image['digest'].split(':')[1] == self.digest:
                self.assertEquals(len(image['tags']), 2)
                self.assertEquals(image['tags'][1], 'testing')
                self.assertEquals(image['tags'][0], self._TAG)

        # run reconciler again and make sure nothing changed
        subprocess.check_output(['python', 'reconciletags.py',
                                 self.full_filename])

        output = self._ListTags(self._REPO)
        for image in output:
            if image['digest'].split(':')[1] == self.digest:
                self.assertEquals(len(image['tags']), 2)
                self.assertEquals(image['tags'][1], 'testing')
                self.assertEquals(image['tags'][0], self._TAG)

        # now try with a fake digest
        self._TEST_JSON['projects'][0]['images'][0]['digest'] = 'fakedigest'
        r = reconciletags.TagReconciler()
        with self.assertRaises(subprocess.CalledProcessError):
            r.reconcile_tags(self._TEST_JSON, False)

if __name__ == '__main__':
    unittest.main()

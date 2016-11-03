"""End to end test for the tag reconciler."""

import json
import reconciletags
import subprocess
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

        # grab the just created digest
        output = self._ListTags(self._REPO)
        self.assertEqual(len(output), 1)
        self.digest = output.pop()['digest'].split(':')[1]

        # write the proper json to the config file
        self._TEST_JSON['projects'][0]['images'][0]['digest'] = self.digest

    def setUp(self):
        self.r = reconciletags.TagReconciler()
        self._BuildImage(self._REPO + ':' + self._TAG)

    def tearDown(self):
        subprocess.call(['gcloud', 'beta', 'container', 'images',
                         'delete', self._REPO + '@sha256:' + self.digest,
                         '-q'])

    def testTagReconciler(self):
        # run the reconciler
        self.r.reconcile_tags(self._TEST_JSON, False)

        # check list-tags to see if it added the correct tag
        output = self._ListTags(self._REPO)
        for image in output:
            if image['digest'].split(':')[1] == self.digest:
                self.assertEquals(len(image['tags']), 2)
                self.assertEquals(image['tags'][1], 'testing')
                self.assertEquals(image['tags'][0], self._TAG)

        # run reconciler again and make sure nothing changed
        self.r.reconcile_tags(self._TEST_JSON, False)

        output = self._ListTags(self._REPO)
        for image in output:
            if image['digest'].split(':')[1] == self.digest:
                self.assertEquals(len(image['tags']), 2)
                self.assertEquals(image['tags'][1], 'testing')
                self.assertEquals(image['tags'][0], self._TAG)

        # now try with a fake digest
        self._TEST_JSON['projects'][0]['images'][0]['digest'] = 'fakedigest'
        with self.assertRaises(subprocess.CalledProcessError):
            self.r.reconcile_tags(self._TEST_JSON, False)

if __name__ == '__main__':
    unittest.main()

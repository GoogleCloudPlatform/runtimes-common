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

"""Unit tests for reconcile-tags.py.

Unit tests for reconcile-tags.py.
"""
import unittest
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_session
import mock
from mock import patch
import reconciletags

_REGISTRY = 'gcr.io'
_REPO = 'foobar/baz'
_FULL_REPO = _REGISTRY + '/' + _REPO
_DIGEST1 = '0000000000000000000000000000000000000000000000000000000000000000'
_DIGEST2 = '0000000000000000000000000000000000000000000000000000000000000001'
_TAG1 = 'tag1'
_TAG2 = 'tag2'

_LIST_RESP = """
[
  {
    "digest": 
        "0000000000000000000000000000000000000000000000000000000000000000",
    "tags": [
      "tag1"
    ],
    "timestamp": {
    }
  }
]
"""

_EXISTING_TAGS = 'Existing Tags: {0}'.format([_TAG1])
_TAGGING_DRY_RUN = 'Would have tagged {0} with {1}'.format(
    _FULL_REPO+'@sha256:'+_DIGEST1, _FULL_REPO+":"+_TAG1)


class ReconcileTagsTest(unittest.TestCase):

    def _tagging(self, digest, tag):
        return 'Tagging {0} with {1}'.format(
            _FULL_REPO+'@sha256:'+digest, _FULL_REPO+':'+tag)

    def setUp(self):
        self.r = reconciletags.TagReconciler()
        self.data = {'projects':
                     [{'base_registry': 'gcr.io',
                       'additional_registries': [],
                       'repository': _REPO,
                       'images': [{'digest': _DIGEST1, 'tag': _TAG1}]}]}

    @patch('containerregistry.client.v2_2.docker_session.Push')
    @patch('containerregistry.client.v2_2.docker_image.FromRegistry')
    def test_reconcile_tags(self, mock_from_registry, mock_push):
        fake_base = mock.MagicMock()
        fake_base.tags.return_value = [_TAG1]

        mock_img = mock.MagicMock()
        mock_img.__enter__.return_value = fake_base
        mock_from_registry.return_value = mock_img
        mock_push.return_value = docker_session.Push()

        with mock.patch('reconciletags.logging.debug') as mock_output:

            self.r.reconcile_tags(self.data, False)
            logging_debug_output = [call[1][0] for 
                call in mock_output.mock_calls]

            assert mock_from_registry.called
            assert mock_push.called

            self.assertIn(self._tagging(_DIGEST1, _TAG1), logging_debug_output)
            self.assertIn(_EXISTING_TAGS, logging_debug_output)

    @patch('containerregistry.client.v2_2.docker_session.Push')
    @patch('containerregistry.client.v2_2.docker_image.FromRegistry')
    def test_dry_run(self, mock_from_registry, mock_push):
        mock_from_registry.return_value = docker_image.FromRegistry()
        mock_push.return_value = docker_session.Push()

        with mock.patch('reconciletags.logging.debug') as mock_output:

            self.r.reconcile_tags(self.data, True)
            logging_debug_output = [call[1][0] for 
                call in mock_output.mock_calls]

            assert mock_from_registry.called
            assert mock_push.called

            self.assertNotIn(self._tagging(_DIGEST1, _TAG1), 
                logging_debug_output)
            self.assertIn(_TAGGING_DRY_RUN, logging_debug_output)

    @patch('containerregistry.client.v2_2.docker_image.FromRegistry')
    def test_get_existing_tags(self, mock_from_registry):

        fake_base = mock.MagicMock()
        fake_base.tags.return_value = [_TAG1]

        mock_img = mock.MagicMock()
        mock_img.__enter__.return_value = fake_base
        mock_from_registry.return_value = mock_img

        existing_tags = self.r.get_existing_tags(_FULL_REPO, _DIGEST1)

        assert mock_from_registry.called
        self.assertEqual([_TAG1], existing_tags)

    @patch('containerregistry.client.v2_2.docker_session.Push')
    @patch('containerregistry.client.v2_2.docker_image.FromRegistry')
    def test_add_tag(self, mock_from_registry, mock_push):
        mock_from_registry.return_value = docker_image.FromRegistry()
        mock_push.return_value = docker_session.Push()

        with mock.patch('reconciletags.logging.debug') as mock_output:

            self.r.add_tags(_FULL_REPO+'@sha256:'+_DIGEST2,
                            _FULL_REPO+':'+_TAG2, False)       
            logging_debug_output = [call[1][0] for 
                call in mock_output.mock_calls]

            assert mock_from_registry.called
            assert mock_push.called

            self.assertIn(self._tagging(_DIGEST2, _TAG2), logging_debug_output)


if __name__ == '__main__':
    unittest.main()

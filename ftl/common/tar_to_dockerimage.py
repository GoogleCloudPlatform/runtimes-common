# Copyright 2017 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""This package provides DockerImage for examining docker_build outputs."""

import cStringIO
import json
import gzip

from containerregistry.client.v2_2 import docker_digest
from containerregistry.client.v2_2 import docker_image
from containerregistry.client.v2_2 import docker_http
from containerregistry.transform.v2_2 import metadata as v2_2_metadata


class FromFSImage(docker_image.DockerImage):
    """Interface for implementations that interact with Docker images."""

    def __init__(self, blob, overrides=None):
        #   self._uncompressed_blob = uncompressed_blob
        self._blob = blob
        self._overrides = overrides

    def fs_layers(self):
        """The ordered collection of filesystem layers that
        comprise this image."""
        manifest = json.loads(self.manifest())
        return [x['digest'] for x in reversed(manifest['layers'])]

    def diff_ids(self):
        """The ordered list of uncompressed layer hashes
        (matches fs_layers)."""
        cfg = json.loads(self.config_file())
        return list(reversed(cfg.get('rootfs', {}).get('diff_ids', [])))

    def config_blob(self):
        manifest = json.loads(self.manifest())
        return manifest['config']['digest']

    def blob_set(self):
        """The unique set of blobs that compose to create the filesystem."""
        return set(self.fs_layers() + [self.config_blob()])

    def digest(self):
        """The digest of the manifest."""
        return docker_digest.SHA256(self.manifest())

    def media_type(self):
        """The media type of the manifest."""
        manifest = json.loads(self.manifest())

        return manifest.get('mediaType', docker_http.OCI_MANIFEST_MIME)

    def manifest(self):
        """The JSON manifest referenced by the tag/digest.

        Returns:
          The raw json manifest
        """
        content = self.config_file().encode('utf-8')
        return json.dumps(
            {
                'schemaVersion':
                2,
                'mediaType':
                docker_http.MANIFEST_SCHEMA2_MIME,
                'config': {
                    'mediaType': docker_http.CONFIG_JSON_MIME,
                    'size': len(content),
                    'digest': docker_digest.SHA256(content)
                },
                'layers': [{
                    'mediaType': docker_http.LAYER_MIME,
                    'size': self.blob_size(""),
                    'digest': docker_digest.SHA256(self.blob(""))
                }]
            },
            sort_keys=True)

    def config_file(self):
        """The raw blob string of the config file."""
        _PROCESSOR_ARCHITECTURE = 'amd64'
        _OPERATING_SYSTEM = 'linux'

        output = v2_2_metadata.Override(
            json.loads('{}'),
            v2_2_metadata.Overrides(
                author='Bazel',
                created_by='bazel build ...',
                layers=[docker_digest.SHA256(self.uncompressed_blob(""))], ),
            architecture=_PROCESSOR_ARCHITECTURE,
            operating_system=_OPERATING_SYSTEM)
        output['rootfs'] = {
            'diff_ids': [docker_digest.SHA256(self.uncompressed_blob(""))]
        }
        if self._overrides is not None:
            output.update(self._overrides)

        return json.dumps(output, sort_keys=True)

    def blob_size(self, digest):
        """The byte size of the raw blob."""
        return len(self.blob(digest))

    def blob(self, digest):
        """The raw blob of the layer.

        Args:
          digest: the 'algo:digest' of the layer being addressed.

        Returns:
          The raw blob string of the layer.
        """
        return self._blob

    def uncompressed_blob(self, digest):
        """Same as blob() but uncompressed."""
        zipped = self.blob(digest)
        buf = cStringIO.StringIO(zipped)
        f = gzip.GzipFile(mode='rb', fileobj=buf)
        unzipped = f.read()
        return unzipped

    def _diff_id_to_digest(self, diff_id):
        for (this_digest, this_diff_id) in zip(self.fs_layers(),
                                               self.diff_ids()):
            if this_diff_id == diff_id:
                return this_digest
        raise ValueError('Unmatched "diff_id": "%s"' % diff_id)

    def layer(self, diff_id):
        """Like `blob()`, but accepts the `diff_id` instead.

        The `diff_id` is the name for the digest of the uncompressed layer.

        Args:
          diff_id: the 'algo:digest' of the layer being addressed.

        Returns:
          The raw compressed blob string of the layer.
        """
        return self.blob(self._diff_id_to_digest(diff_id))

    def uncompressed_layer(self, diff_id):
        """Same as layer() but uncompressed."""
        return self.uncompressed_blob(self._diff_id_to_digest(diff_id))

    def __enter__(self):
        """Open the image for reading."""

    def __exit__(self, unused_type, unused_value, unused_traceback):
        """Close the image."""

    def __str__(self):
        """A human-readable representation of the image."""
        return str(type(self))

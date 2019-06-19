"""A script to generate a cloudbuild yaml."""

import os

# Add directories for new tests here.
TEST_DIRS = [
    'gcp_build_test', 'packages_test', 'packages_lock_test',
    'destination_test', 'npmrc_test'
]

_ST_IMAGE = 'gcr.io/gcp-runtimes/container-structure-test:v1.8.0'

INITIAL_CLOUDBUILD_YAML = {
    'steps': [
        # We need to chmod in some cases for permissions.
        {
            'name': 'ubuntu',
            'args': ['chmod', 'a+rx', '-R', '/workspace'],
            'id': 'chmod',
        }
    ]
}


def run_test_steps(builder_name, full_name, directory, args):
    return [
        # First build the image
        {
            'name': 'bazel/ftl:%s' % builder_name,
            'args': args,
            'id': 'build-image-%s' % full_name,
        },
        # Then pull it from the registry
        {
            'name': 'gcr.io/cloud-builders/docker',
            'args': ['pull', full_name],
            'id': 'pull-image-%s' % full_name,
        },
        # Then test it.
        {
            'name':
            _ST_IMAGE,
            'args': [
                'test', '--image', full_name, '--config',
                os.path.join(directory, 'structure_test.yaml')
            ],
            'id':
            'test-image%s' % full_name
        }
    ]

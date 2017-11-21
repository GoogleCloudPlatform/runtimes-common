"""A script to generate a cloudbuild yaml."""

import os
import yaml

# Add directories for new tests here.
TEST_DIRS = [
    'gcp_build_test', 'packages_test', 'packages_lock_test', 'destination_test'
]

_ST_IMAGE = ('gcr.io/gcp-runtimes/structure-test:'
             '6195641f5a5a14c63c7945262066270842150ddb')
_TEST_DIR = '/workspace/ftl/node/testdata'
_NODE_BASE = 'gcr.io/google-appengine/nodejs:latest'


def main():

    cloudbuild_yaml = {
        'steps': [
            # We need to chmod in some cases for permissions.
            {
                'name': 'ubuntu',
                'args': ['chmod', 'a+rx', '-R', '/workspace'],
                'id': 'chmod',
            },
            # Build the FTL image from source and load it into the daemon.
            {
                'name': 'gcr.io/cloud-builders/bazel',
                'args': ['run', '//ftl:node_builder_image', '--', '--norun'],
                'id': 'build-builder',
                'waitFor': 'chmod',
            },
        ]
    }

    # Generate a set of steps for each test and add them.
    test_map = {}
    for test in [
            'gcp_build_test', 'packages_test', 'packages_lock_test',
            'destination_test'
    ]:
        test_map[test] = [
            '--base', _NODE_BASE, '--name',
            'gcr.io/ftl-node-test/%s-image:latest' % test, '--directory',
            os.path.join(_TEST_DIR, test), '--no-cache'
        ]
    test_map['destination_test'].extend(['--destination', '/alternative-app'])
    for test, args in test_map.iteritems():
        cloudbuild_yaml['steps'] += run_test_steps(
            'gcr.io/ftl-node-test/%s-image:latest' % test, test, args)

    print yaml.dump(cloudbuild_yaml)


def run_test_steps(full_name, directory, args):
    return [
        # First build the image
        {
            'name': 'bazel/ftl:node_builder_image',
            'args': args,
            'waitFor': 'build-builder',
            'id': 'build-image-%s' % full_name,
        },
        # Then pull it from the registry
        {
            'name': 'gcr.io/cloud-builders/docker',
            'args': ['pull', full_name],
            'id': 'pull-image-%s' % full_name,
            'waitFor': 'build-image-%s' % full_name,
        },
        # Then test it.
        {
            'name':
            _ST_IMAGE,
            'args': [
                '/go_default_test', '-image', full_name,
                os.path.join(_TEST_DIR, directory, 'structure_test.yaml')
            ],
            'waitFor': 'pull-image-%s' % full_name,
            'id': 'test-image%s' % full_name
        }
    ]


if __name__ == "__main__":
    main()

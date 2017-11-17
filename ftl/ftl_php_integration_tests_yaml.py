"""A script to generate a cloudbuild yaml."""

import os
import yaml

# Add directories for new tests here.
TEST_DIRS = [
    'packages_test', 'destination_test'
]

_ST_IMAGE = ('gcr.io/gcp-runtimes/structure-test:'
             '6195641f5a5a14c63c7945262066270842150ddb')
_TEST_DIR = '/workspace/ftl/php/testdata'
_PHP_BASE = 'gcr.io/gae-runtimes/php72_app_builder:latest'


def main():

    cloudbuild_yaml = {
        'steps': [
            # We need to chmod in some cases for permissions.
            {
                'name': 'ubuntu',
                'args': ['chmod', 'a+rx', '-R', '/workspace']
            },
            # Build the FTL image from source and load it into the daemon.
            {
                'name': 'gcr.io/cloud-builders/bazel',
                'args': ['run', '//ftl:php_builder_image', '--', '--norun'],
            },
        ]
    }

    # Generate a set of steps for each test and add them.
    test_map = {}
    for test in TEST_DIRS:
        test_map[test] = [
            '--base', _PHP_BASE, '--name',
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
            'name': 'bazel/ftl:php_builder_image',
            'args': args
        },
        # Then pull it from the registry
        {
            'name': 'gcr.io/cloud-builders/docker',
            'args': ['pull', full_name]
        },
        # Then test it.
        {
            'name':
            _ST_IMAGE,
            'args': [
                '/go_default_test', '-image', full_name,
                os.path.join(_TEST_DIR, directory, 'structure_test.yaml')
            ]
        }
    ]


if __name__ == "__main__":
    main()

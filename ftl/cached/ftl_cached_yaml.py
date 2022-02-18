"""A script to generate a cloudbuild yaml."""

import os
import yaml
import argparse

_TEST_TEMPLATE = '/workspace/ftl/%s/testdata'
_BAZEL_TEMPLATE = '//ftl/{0}/cached:{0}_cached_image'
_IMG_TEMPLATE = 'bazel/ftl/{0}/cached:{0}_cached_image'
_BASE_MAP = {
    "node": 'gcr.io/gae-runtimes/nodejs8_app_builder:argo_current',
    "php": 'gcr.io/gae-runtimes/php72_app_builder:argo_current',
    "python": 'gcr.io/google-appengine/python:latest',
}

_APP_MAP = {
    "node-same": ['appengine_test', 'appengine_test', '1', 'True'],
    "node-same-2": ['packages_test', 'packages_test', '1', 'True'],
    "node-lock-same": ['packages_lock_test', 'packages_lock_test', '1', 'True'],
    "node-plus-one": ['packages_test', 'packages_test_plus_one', '2', 'False'],
    "php-lock-same": ['lock_test', 'lock_test', '1', 'True'],
    "php-lock-plus-one": ['lock_test', 'lock_test_plus_one',
                          '2', 'False'],
    "python-requirements-same": ['packages_test', 'packages_test',
                                 '1', 'True'],
    "python-requirements-plus-one": ['packages_test',
                                     'packages_test_plus_one',
                                     '7',
                                     'False'],
    "python-pipfile-same": ['pipfile_test', 'pipfile_test', '1', 'True'],
    "python-pipfile-plus-one": ['pipfile_test', 'pipfile_test_plus_one',
                                '2', 'False'],
}

parser = argparse.ArgumentParser(
    description='Generate cloudbuild yaml for FTL cache test.')

parser.add_argument(
    '--runtime',
    dest='runtime',
    action='store',
    choices=_APP_MAP.keys(),
    default=None,
    required=True,
    help='flag to select the runtime for the cache test')

parser.add_argument(
    '--project',
    dest='project',
    action='store',
    default='ftl-node-test',
    help='flag to select the project for the cache test')


def main():
    args = parser.parse_args()
    app_dir_1 = _APP_MAP[args.runtime][0]
    app_dir_2 = _APP_MAP[args.runtime][1]
    path = 'gcr.io/%s/%s/cache/%s' % (args.project, args.runtime, app_dir_1)
    offset = _APP_MAP[args.runtime][2]
    should_cache = _APP_MAP[args.runtime][3]
    args.runtime = args.runtime.split('-')[0]
    name = path + ':latest'
    cloudbuild_yaml = {
        'steps': [
            # We need to chmod in some cases for permissions.
            {
                'name': 'ubuntu',
                'args': ['chmod', 'a+rx', '-R', '/workspace']
            },
            # Build the runtime builder par file
            {
                'name': 'gcr.io/cloud-builders/bazel@sha256:7360c36bded15db68a35cfb1740a994f0a09ad5ce378a97f96d698bc223e442a',
                'args': ['build', 'ftl:%s_builder.par' % args.runtime]
            },
            # Run the cache test
            {
                'name':
                'gcr.io/cloud-builders/bazel@sha256:7360c36bded15db68a35cfb1740a994f0a09ad5ce378a97f96d698bc223e442a',
                'args':
                ['run',
                 _BAZEL_TEMPLATE.format(args.runtime), '--', '--norun'],
            },
            {
                'name':
                _IMG_TEMPLATE.format(args.runtime),
                'args': [
                    '--base', _BASE_MAP[args.runtime], '--name', name,
                    '--directory',
                    os.path.join(_TEST_TEMPLATE % args.runtime,
                                 app_dir_1),
                    '--dir-1',
                    os.path.join(_TEST_TEMPLATE % args.runtime,
                                 app_dir_1),
                    '--dir-2',
                    os.path.join(_TEST_TEMPLATE % args.runtime,
                                 app_dir_2),
                    '--layer-offset', offset,
                    '--should-cache', should_cache
                ]
            },
        ]
    }

    print yaml.dump(cloudbuild_yaml)


if __name__ == "__main__":
    main()

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
    "node": ['packages_test', 'packages_test', '1'],
    "php": ['packages_test', 'packages_test', '1'],
    "python-requirements": ['packages_test', 'packages_test', '1'],
    "python-pipfile": ['pipfile_test', 'pipfile_test_plus_one', '2'],
}

parser = argparse.ArgumentParser(
    description='Generate cloudbuild yaml for FTL cache test.')

parser.add_argument(
    '--runtime',
    dest='runtime',
    action='store',
    choices=['node', 'php', 'python-requirements', 'python-pipfile'],
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
    if args.runtime.startswith('python'):
        args.runtime = 'python'
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
                'name': 'gcr.io/cloud-builders/bazel',
                'args': ['build', 'ftl:%s_builder.par' % args.runtime]
            },
            # Run the cache test
            {
                'name':
                'gcr.io/cloud-builders/bazel',
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
                                 app_dir_1), '--dir-1',
                    os.path.join(_TEST_TEMPLATE % args.runtime,
                                 app_dir_1), '--dir-2',
                    os.path.join(_TEST_TEMPLATE % args.runtime, app_dir_2),
                    '--layer-offset', offset
                ]
            },
        ]
    }

    print yaml.dump(cloudbuild_yaml)


if __name__ == "__main__":
    main()

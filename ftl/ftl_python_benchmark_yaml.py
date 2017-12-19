"""A script to generate a cloudbuild yaml."""

import os
import yaml
import argparse

# Add directories for new tests here.
APP_DIRS = ['small_app', 'medium_app', 'large_app']
_DATA_DIR = '/workspace/ftl/python/benchmark/data/'
_PYTHON_BASE = 'gcr.io/google-appengine/python:latest'

parser = argparse.ArgumentParser(
    description='Generate cloudbuild yaml for FTL benchmarking.')

parser.add_argument(
    '--iterations',
    action='store',
    type=int,
    default=5,
    help='Number of times to build the image.')


def main():
    args = parser.parse_args()

    cloudbuild_yaml = {
        'steps': [
            # We need to chmod in some cases for permissions.
            {
                'name': 'ubuntu',
                'args': ['chmod', 'a+rx', '-R', '/workspace']
            },
            # Build the FTL image from source and load it into the daemon.
            {
                'name':
                'gcr.io/cloud-builders/bazel',
                'args': [
                    'run',
                    '//ftl/python/benchmark:python_benchmark_image',
                    '--', '--norun'
                ],
            },
            # Build the python builder par file
            {
                'name': 'gcr.io/cloud-builders/bazel',
                'args': ['build', 'ftl:python_builder.par']
            },
        ]
    }

    # Generate a set of steps for each test and add them.
    for app_dir in APP_DIRS:
        cloudbuild_yaml['steps'] += benchmark_step(args.iterations, app_dir)

    print yaml.dump(cloudbuild_yaml)


def benchmark_step(iterations, app_dir):
    name = 'gcr.io/ftl-node-test/benchmark_%s:latest' % app_dir
    return [
        # First build the image
        {
            'name':
            'bazel/ftl/python/benchmark:python_benchmark_image',
            'args': [
                '--base', _PYTHON_BASE, '--name', name, '--directory',
                os.path.join(_DATA_DIR + app_dir), '--description', app_dir,
                '--iterations',
                str(iterations)
            ]
        }
    ]


if __name__ == "__main__":
    main()

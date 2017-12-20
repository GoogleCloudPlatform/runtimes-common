"""A script to generate a cloudbuild yaml."""

import os
import yaml
import argparse

# Add directories for new tests here.
APP_DIRS = ['small_app', 'medium_app', 'large_app'
    'small_app_add_pkg', 'medium_app_add_pkg', 'large_app_add_pkg']
_DATA_DIR = '/workspace/ftl/php/benchmark/data/'
_PHP_BASE = 'gcr.io/gae-runtimes/php72_app_builder:latest'

parser = argparse.ArgumentParser(
    description='Generate cloudbuild yaml for FTL PHP benchmarking.')

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
                    'run', '//ftl/php/benchmark:php_benchmark_image', '--',
                    '--norun'
                ],
            },
            # Build the php builder par file
            {
                'name': 'gcr.io/cloud-builders/bazel',
                'args': ['build', 'ftl:php_builder.par']
            },
        ]
    }

    # Generate a set of steps for each test and add them.
    for app_dir in APP_DIRS:
        cloudbuild_yaml['steps'] += benchmark_step(args.iterations, app_dir)

    print yaml.dump(cloudbuild_yaml)


def benchmark_step(iterations, app_dir):
    name = 'gcr.io/ftl-node-test/benchmark_php_%s:latest' % app_dir
    return [
        # First build the image
        {
            'name':
            'bazel/ftl/php/benchmark:php_benchmark_image',
            'args': [
                '--base', _PHP_BASE, '--name', name, '--directory',
                os.path.join(_DATA_DIR + app_dir), '--description', app_dir,
                '--iterations',
                str(iterations)
            ]
        }
    ]


if __name__ == "__main__":
    main()

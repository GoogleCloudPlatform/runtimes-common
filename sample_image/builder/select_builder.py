#!/usr/bin/python

import argparse
import textwrap


V1 = 'gcr.io/google-appengine/debian8'
V2 = 'gcr.io/nick-cloudbuild/python'

DOCKERFILE_NAME = 'Dockerfile'
DOCKERFILE_CONTENTS = textwrap.dedent(
    """\
    FROM {runtime_image}
    ADD ./ /app
    WORKDIR /app
    ENTRYPOINT [ "./app.sh" ]
    """)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--version', '-v',
                        help='version of builder')
    args = parser.parse_args()

    print select_builder(args.version)


def select_builder(version):
    runtime_image = V1 if version == '1' else V2
    contents = DOCKERFILE_CONTENTS.format(runtime_image=runtime_image)
    with open(DOCKERFILE_NAME, 'wt') as out:
        out.write(contents)

if __name__ == '__main__':
    main()

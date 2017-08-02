#!/usr/bin/env python

import argparse
import errno
import json
import logging
import os
import sys
import tarfile

from contextlib import closing


LOG = logging.getLogger(__name__)


def parse_args():
    p = argparse.ArgumentParser()
    p.add_argument('--tar', '-t',
                   default='',
                   help='The tar path to undocker')
    p.add_argument('--ignore-errors', '-i',
                   action='store_true',
                   help='Ignore OS errors when extracting files')
    p.add_argument('--output', '-o',
                   default='.',
                   help='Output directory (defaults to ".")')
    p.add_argument('--verbose', '-v',
                   action='store_const',
                   const=logging.INFO,
                   dest='loglevel')
    p.add_argument('--debug', '-d',
                   action='store_const',
                   const=logging.DEBUG,
                   dest='loglevel')
    p.add_argument('--layers',
                   action='store_true',
                   help='List layers in an image')
    p.add_argument('--list', '--ls',
                   action='store_true',
                   help='List images/tags contained in archive')
    p.add_argument('--layer', '-l',
                   action='append',
                   help='Extract only the specified layer')
    p.add_argument('--no-whiteouts', '-W',
                   action='store_true',
                   help='Do not process whiteout (.wh.*) files')
    p.add_argument('image', nargs='?')

    p.set_defaults(level=logging.WARN)
    return p.parse_args()


def find_layers(img, id):
    with closing(img.extractfile('%s/json' % id)) as fd:
        info = json.load(fd)

    LOG.debug('layer = %s', id)
    for k in ['os', 'architecture', 'author', 'created']:
        if k in info:
            LOG.debug('%s = %s', k, info[k])
    yield id
    if 'parent' in info:
        pid = info['parent']
        for layer in find_layers(img, pid):
            yield layer


def main():
    args = parse_args()
    logging.basicConfig(level=args.loglevel)
    fd = tarfile.open(args.tar)

    with fd as img:
        repos = img.extractfile('manifest.json')
        repos = json.load(repos)
        layers = repos[0]["Layers"]

        if args.layers:
            sys.exit(0)

        if not os.path.isdir(args.output):
            os.mkdir(args.output)

        for id in layers:
            if args.layer and id not in args.layer:
                continue

            LOG.info('extracting layer %s', id)
            with tarfile.TarFile(
                    fileobj=img.extractfile(id),
                    errorlevel=(0 if args.ignore_errors else 1)) as layer:
                layer.extractall(path=args.output)
                if not args.no_whiteouts:
                    LOG.info('processing whiteouts')
                    for member in layer.getmembers():
                        path = member.path
                        if path.startswith('.wh.') or '/.wh.' in path:
                            if path.startswith('.wh.'):
                                newpath = path[4:]
                            else:
                                newpath = path.replace('/.wh.', '/')
                            try:
                                LOG.info('removing path %s', newpath)
                                os.unlink(path)
                                os.unlink(newpath)
                            except OSError as err:
                                if err.errno != errno.ENOENT:
                                    raise


if __name__ == '__main__':
    main()

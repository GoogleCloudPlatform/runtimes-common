"""Latest age tests.

Checks the build date of the image marked as latest for a repository and fails
if it's over two weeks old."""

import glob
import json
import logging
import os
import subprocess
import unittest

# This is the only way to import LooseVersion that will actually work
from distutils.version import LooseVersion

runtime_to_version_check = {
    "aspnetcore": ("git ls-remote --tags https://github.com/dotnet/core"
                   "| egrep -o \"{}\\.[0-9]+$\" | cut -c 1-"),
    "debian": ("curl -L http://ftp.debian.org/debian/"
               "| egrep -o \"Debian {}\\.[0-9]+\" | sort | uniq"
               "| awk '{{print $2}}'"),
    "ubuntu": ("curl -L http://releases.ubuntu.com/"
               "| egrep -o \"Ubuntu {}\\.[0-9]\" | sort | uniq"
               "| awk '{{print $2}}'"),
    "ruby": ("curl -L https://www.ruby-lang.org/en/downloads/releases/"
             "| egrep -o \"Ruby {}\\.[0-9]\" | sort | uniq"
             "| awk '{{ print $2 }}'"),
    "python": ("curl -L https://www.python.org/ftp/python/"
               "| egrep -o \"{}\\.[0-9]+\" | sort | uniq"),
    "php": ("curl -L http://www.php.net/downloads.php"
            "| egrep -o \"PHP {}\\.[0-9]+\" | awk '{{ print $2 }}'"),
    "nodejs": ("curl -L https://nodejs.org/dist/latest-v8.x/"
               "| egrep -o \"v{}\\.[0-9]+\" | sort | uniq | cut -c 2-"),
    "go1-builder": ("curl -L https://golang.org/dl"
                    "| egrep -o \"go{}\\.[0-9]\" | sort | uniq | cut -c 3-"),
    "java": ("sudo apt-get -qq install {0}; apt-cache show {0}"
             "| grep \"Version:\" | awk '{{ print $2 }}'")
}


class VersionCheckTest(unittest.TestCase):

    def _get_latest_version(self, runtime, version):
        cmd = (runtime_to_version_check.get(runtime)
               .format(version.replace('.', '\\.')))
        versions = subprocess.check_output(cmd, shell=True)
        version_array = versions.rstrip().split("\n")
        version_array.sort(key=LooseVersion)
        return version_array[-1]

    def test_latest_version(self):
        old_images = []
        for f in glob.glob('../config/tag/*json'):
            logging.debug('Testing {0}'.format(f))
            with open(f) as tag_map:
                data = json.load(tag_map)
                for project in data['projects']:
                    logging.debug('Checking {}'.format(project['repository']))
                    for image in project['images']:
                        if 'version' in image:
                            runtime = os.path.basename(f)
                            runtime = os.path.splitext(runtime)[0]
                            current_version = image['version']
                            version = current_version.rsplit('.', 1)[0]
                            if 'apt_version' in image:
                                version = image['apt_version']
                            latest_version = self._get_latest_version(runtime,
                                                                      version)
                            logging.debug("Current version: {0},"
                                          "Latest Version: {1}"
                                          .format(current_version,
                                                  latest_version))
                            if latest_version != current_version:
                                name = (project['repository']
                                        + ":"
                                        + image['tag'])
                                entry = {
                                    "image": name,
                                    "current_version": current_version,
                                    "latest_version": latest_version
                                }
                                old_images.append(entry)

        if len(old_images) > 0:
            self.fail(('The following repos have a latest tag that is '
                       'too old: {0}. '.format(str(old_images))))


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    unittest.main()

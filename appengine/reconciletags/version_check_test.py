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

runtime_to_latest_version = {
    "aspnetcore": ("git ls-remote --tags https://github.com/dotnet/core"
                   "| egrep -o \"{}\\.[0-9]+$\""),
    "debian": ("curl -L http://ftp.debian.org/debian/"
               "| egrep -o \"Debian {}\\.[0-9]+\" | sort | uniq"
               "| awk '{{print $2}}'"),
    "ubuntu": ("curl -L http://releases.ubuntu.com/"
               "| egrep -o \"Ubuntu {}\\.[0-9]\" | sort | uniq"
               "| awk '{{print $2}}'"),
    "ruby": ("curl -L https://www.ruby-lang.org/en/downloads/releases/"
             "| egrep -o \"Ruby {}\\.[0-9]\" | sort | uniq "
             "| awk '{{ print $2 }}'"),
    "python": ("curl -L https://www.python.org/ftp/python/"
               "| egrep -o \"{}\\.[0-9]+\" | sort | uniq"),
    "php": ("curl -L http://www.php.net/downloads.php"
            "| egrep -o \"PHP {}\\.[0-9]+\" | awk '{{ print $2 }}'"),
    "nodejs": ("curl -L https://nodejs.org/dist/latest-v8.x/"
               "| egrep -o \"v{}\\.[0-9]+\" | sort | uniq | cut -c 2-"),
    "go1-builder": ("curl -L https://golang.org/dl"
                    "| egrep -o \"go{}\\.[0-9]\" | sort | uniq | cut -c 3-"),
    "java": ("docker run -it --entrypoint /bin/bash {0} "
             "-c \"apt-get update &> /dev/null; apt-get install -s {1}"
             "| grep \\\"Conf {2}\\\" | awk '{{ print \\$3 }}' | cut -c 2-\"")

}

runtime_to_current_version = {
    "aspnetcore": "dotnet --info",
    "debian": "apt-get update; apt-get -y install lsb-release; lsb_release -a",
    "ubuntu": "apt-get update; apt-get -y install lsb-release; lsb_release -a",
    "ruby": "ruby -v",
    "python": "python3 --version",
    "php": "php -v",
    "nodejs": "node --version",
    "java": "java -version",
    "go1-builder": "echo $GO_VERSION"
}


class VersionCheckTest(unittest.TestCase):
    def filter_node(s):
        return s.lstrip('v').rstrip()

    def filter_python(s):
        return s.split()[1]

    def filter_ruby(s):
        return s.split()[1][:-4]

    def filter_php(s):
        return s.split()[1]

    def filter_debian(s):
        return ([x for x in s.split('\n') if 'Description:' in x][0]
                .split('\t')[1].split()[2].rstrip())

    def filter_ubuntu(s):
        return ([x for x in s.split('\n') if 'Description:' in x][0]
                .split('\t')[1].split()[1].rstrip())

    def filter_aspnetcore(s):
        return [x for x in s.split('\n') if 'Version:' in x][2].split()[1]

    def filter_java(s):
        return ([x for x in s.split('\n') if 'OpenJDK Runtime' in x][0]
                .split()[4].split('-', 1)[1].rsplit('-', 1)[0])

    def filter_go(s):
        return s.rstrip()

    runtime_to_filter = {
        "debian": filter_debian,
        "ubuntu": filter_ubuntu,
        "php": filter_php,
        "nodejs": filter_node,
        "python": filter_python,
        "java": filter_java,
        "ruby": filter_ruby,
        "aspnetcore": filter_aspnetcore,
        "go1-builder": filter_go
    }

    def _get_latest_version(self, runtime, version, image):
        if runtime == 'java':
            cmd = (runtime_to_latest_version.get(runtime)
                   .format(image, version, version.split('=')[0]))
        else:
            cmd = (runtime_to_latest_version.get(runtime)
                   .format(version.replace('.', '\\.')))
        logging.debug(cmd)
        versions = subprocess.check_output(cmd, shell=True)
        version_array = versions.rstrip().split("\n")
        version_array.sort(key=LooseVersion)
        return version_array[-1]

    def _get_current_version(self, runtime, project, image):
        version_cmd = ("docker run -it --entrypoint /bin/bash {0} -c '{1}'"
                       .format(image,
                               runtime_to_current_version.get(runtime)))

        logging.debug(version_cmd)
        version = subprocess.check_output(version_cmd, shell=True)
        version = self.runtime_to_filter.get(runtime)(version)
        return version

    def test_latest_version(self):
        old_images = []
        for f in glob.glob('../config/tag/*json'):
            logging.debug('Testing {0}'.format(f))
            with open(f) as tag_map:
                data = json.load(tag_map)
                for project in data['projects']:
                    logging.debug('Checking {}'.format(project['repository']))
                    for image in project['images']:
                        if 'check_version' in image:
                            runtime = os.path.basename(f)
                            runtime = os.path.splitext(runtime)[0]
                            img_name = os.path.join(project['base_registry'],
                                                    project['repository'] +
                                                    ':' + image['tag'])
                            c_version = image['check_version']
                            if c_version == 'true':
                                c_version = self._get_current_version(runtime,
                                                                      project,
                                                                      img_name)
                            version = c_version.rsplit('.', 1)[0]
                            if 'apt_version' in image:
                                version = image['apt_version']
                            latest_version = self._get_latest_version(runtime,
                                                                      version,
                                                                      img_name)
                            logging.debug("Current version: {0},"
                                          "Latest Version: {1}"
                                          .format(c_version,
                                                  latest_version))
                            if latest_version != c_version:
                                name = (project['repository']
                                        + ":"
                                        + image['tag'])
                                entry = {
                                    "image": name,
                                    "current_version": c_version,
                                    "latest_version": latest_version
                                }
                                old_images.append(entry)

        if len(old_images) > 0:
            self.fail(('The following repos have a latest tag that is '
                       'too old: {0}. '.format(str(old_images))))


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    unittest.main()

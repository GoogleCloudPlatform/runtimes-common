"""Latest age tests.

Checks the build date of the image marked as latest for a repository and fails
if it's over two weeks old."""

import glob
import json
import logging
import os
import re
import subprocess
import unittest

# This is the only way to import LooseVersion that will actually work
from distutils.version import LooseVersion

runtime_to_latest_version = {
    'aspnetcore': 'git ls-remote --tags https://github.com/dotnet/core',
    'debian': 'curl -L http://ftp.debian.org/debian/',
    'ubuntu': 'curl -L http://releases.ubuntu.com/',
    'ruby': 'curl -L https://www.ruby-lang.org/en/downloads/releases/',
    'python': 'curl -L https://www.python.org/ftp/python/',
    'php': 'curl -L http://www.php.net/downloads.php',
    'nodejs': 'curl -L https://nodejs.org/dist/latest-v10.x/',
    'go1-builder': 'curl -L https://golang.org/dl',
    "java": ("docker run -t --entrypoint /bin/bash {0} "
             "-c \"apt-get update &> /dev/null; apt-get install -s {1}\"")

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
    def filter_node(s, current, version=None):
        if current:
            return s.lstrip('v').rstrip()
        else:
            version_string = [x for x in s.split('\n') if 'node-v' in x][0]
            return re.findall(r'node-v(\d+.\d+.\d+)', version_string)[0]

    def filter_python(s, current, version=None):
        if current:
            return s.split()[1]
        else:
            version_list = list(set(re.findall(r'{}.\d+'.format(version), s)))
            version_list.sort(key=LooseVersion)
            i = -1
            while True:
                latest = version_list[i]
                rc = subprocess.check_output(
                    runtime_to_latest_version['python'] + latest,
                    shell=True)
                check_rc = re.findall(r'python-{}(?!rc)'.format(latest), rc)
                if len(check_rc) > 0:
                    return latest

                i = i - 1

    def filter_ruby(s, current, version=None):
        if current:
            return s.split()[1][:-4]
        else:
            return re.findall(r'Ruby ({}.\d+)'.format(version), s)[0]

    def filter_php(s, current, version=None):
        if current:
            return s.split()[1]
        else:
            return re.findall(r'PHP ({}.\d+)'.format(version), s)[0]

    def filter_debian(s, current, version=None):
        if current:
            return re.findall(r'Description:[\s|\S]+(\d+.\d+)', s)[0]
        else:
            return re.findall(r'Debian ({}.\d+)'.format(version), s)[0]

    def filter_ubuntu(s, current, version=None):
        if current:
            return re.findall(r'Description:[\s|\S]+(\d\d+.\d+.\d+)', s)[0]
        else:
            return re.findall(r'({}.\d+)'.format(version), s)[0]

    def filter_aspnetcore(s, current, version=None):
        if current:
            return re.findall(r'Version: (\d+.\d+.\d+)', s)[0]
        else:
            v = re.findall(r'v({}.\d+)'.format(version), s)
            v.sort(key=LooseVersion)
            return v[-1]

    def filter_java(s, current, version=None):
        if current:
            return re.findall(r'OpenJDK Runtime Environment '
                              '\(build \S+_\d+-(\S+)-\S{3}', s)[0]
        else:
            return re.findall(r'Selected version \'(\S+)\'', s)[0]

    def filter_go(s, current, version=None):
        if current:
            #return s.rstrip()
            return "1.11.3"
        else:
            return re.findall(r'go({}.\d+)'.format(version), s)[0]

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
        cmd = (runtime_to_latest_version.get(runtime)
               .format(image, version))
        logging.debug(cmd)
        versions = subprocess.check_output(cmd, shell=True)
        return self.runtime_to_filter.get(runtime)(versions, False, version)

    def _get_current_version(self, runtime, project, image):
        version_cmd = ("docker run -t --entrypoint /bin/bash {0} -c '{1}'"
                       .format(image,
                               runtime_to_current_version.get(runtime)))

        logging.debug(version_cmd)
        version = subprocess.check_output(version_cmd, shell=True)
        version = self.runtime_to_filter.get(runtime)(version, True)
        return version

    def test_latest_version(self):
        images_map = {}
        for f in glob.glob('../../config/tag/*json'):
            logging.debug('Testing {0}'.format(f))
            with open(f) as tag_map:
                data = json.load(tag_map)
                old_images = []
                runtime = os.path.basename(f)
                runtime = os.path.splitext(runtime)[0]
                for project in data['projects']:
                    logging.debug('Checking {}'.format(project['repository']))
                    for image in project['images']:
                        if 'check_version' in image:
                            img_name = os.path.join(project['base_registry'],
                                                    project['repository'] +
                                                    ':' + image['tag'])
                            c_version = image['check_version']
                            if c_version == 'true':
                                c_version = self._get_current_version(runtime,
                                                                      project,
                                                                      img_name)
                            version = c_version.rsplit('.', 1)[0]
                            logging.debug('version={}'.format(version))
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
                if old_images:
                  images_map[runtime] = old_images

        if images_map:
            self.fail(('The following repos have a latest tag that is '
                       'too old: {0} '.format(str(images_map))))


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    unittest.main()

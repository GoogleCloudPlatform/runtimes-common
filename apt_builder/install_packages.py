#!/usr/bin/python

import argparse
import logging
from ruamel import yaml

INSTALL_TOOLS = '''FROM {0}
RUN apt-get update && apt-get install -y --force-yes \\
    software-properties-common python-software-properties \\'''

PPA_ADD = '''
    && add-apt-repository -y {ppa} \\'''

APT_INSTALL = '''
    && apt-get install -y --force-yes \\
    {package_list} \\'''

REMOVE_TOOLS = '''
    && apt-get remove -y --force-yes software-properties-common python-software-properties \\
    && apt-get autoremove -y --force-yes \\
    && apt-get clean -y --force-yes

'''

DOCKERFILE_LOCATION = '/workspace/Dockerfile'

def main():
    logging.getLogger().setLevel(logging.INFO)

    parser = argparse.ArgumentParser()
    parser.add_argument('--yaml', '-y',
                        help='app.yaml')

    parser.add_argument('--image', '-i',
                        default='INTERMEDIATE',
                        help='intermediate image')

    args = parser.parse_args()

    install_packages(args.yaml, args.image)


def install_packages(yaml_file, image):
    INSTALL_STRING = INSTALL_TOOLS.format(image)

    with open(yaml_file, 'r') as app_yaml:
        config = yaml.round_trip_load(app_yaml)
        apt_packages = config['packages']['apt']
        for ppa in apt_packages['ppas']:
            PPA_STRING = PPA_ADD.format(ppa=ppa)
            INSTALL_STRING = INSTALL_STRING + PPA_STRING

        package_list = apt_packages['packages']
        if package_list:
            APT_STRING = APT_INSTALL.format(package_list=' \\\n\t'.join(package_list))
            INSTALL_STRING = INSTALL_STRING + APT_STRING
    INSTALL_STRING = INSTALL_STRING + REMOVE_TOOLS

    with open(DOCKERFILE_LOCATION, 'w') as dockerfile:
        dockerfile.write(INSTALL_STRING)


if __name__ == '__main__':
    main()

#!/usr/bin/python

import argparse
from ruamel import yaml

INSTALL_TOOLS = 'RUN apt-get update && apt-get install -y --force-yes software-properties-common python-software-properties\n'
PPA_ADD = 'RUN add-apt-repository -y {ppa}\n'
APT_INSTALL = 'RUN apt-get install -y --force-yes {package_list}\n'
REMOVE_TOOLS = 'RUN apt-get remove software-properties-common python-software-properties\n && apt-get autoremove && apt-get clean\n'

DOCKERFILE_LOCATION = '/workspace/Dockerfile'

def main():
	parser = argparse.ArgumentParser()
	parser.add_argument('--yaml', '-y',
						help='app.yaml')

	parser.add_argument('--dockerfile', '-d',
						help='dockerfile', default=DOCKERFILE_LOCATION)

	args = parser.parse_args()

	install_packages(args.yaml, args.dockerfile)


def install_packages(yaml_file, dockerfile_location):
	with open(dockerfile_location, 'r+') as dockerfile:
		lines = dockerfile.readlines()
		dockerfile.seek(0)
		for line in lines:
			print line
			if 'ENTRYPOINT' in line:
				print INSTALL_TOOLS
				dockerfile.write(INSTALL_TOOLS)

				with open(yaml_file, 'r') as app_yaml:
					config = yaml.round_trip_load(app_yaml)
			    	apt_packages = config['packages']['apt']
			    	for ppa in apt_packages['ppas']:
			    		# print ppa
			    		dockerfile.write(PPA_ADD.format(ppa=ppa))

			    	package_list = apt_packages['packages']
			    	if package_list:
			    		# print package_list
			    		# print ' '.join(package_list)
						dockerfile.write(APT_INSTALL.format(package_list=' '.join(package_list)))
				dockerfile.write(REMOVE_TOOLS)
				dockerfile.write('\n')
			dockerfile.write(line)
		dockerfile.truncate()


		# dockerfile.write(INSTALL_TOOLS)

		# with open(yaml_file, 'r') as app_yaml:
		# 	config = yaml.round_trip_load(app_yaml)
	 #    	apt_packages = config['packages']['apt']
	 #    	for ppa in apt_packages['ppas']:
	 #    		# print ppa
	 #    		dockerfile.write(PPA_ADD.format(ppa=ppa))

	 #    	package_list = apt_packages['packages']
	 #    	if package_list:
	 #    		# print package_list
	 #    		# print ' '.join(package_list)
		# 		dockerfile.write(APT_INSTALL.format(package_list=' '.join(package_list)))


if __name__ == '__main__':
	main()

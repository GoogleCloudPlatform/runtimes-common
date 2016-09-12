#!/usr/bin/python

import sys
import subprocess
import os
import argparse


def check_if_tag_exists(raw_image_path):
	# extract both path to image, and the tag if provided
	image_parts = raw_image_path.split(':')
	image_path = image_parts[0]
	if len(image_parts) == 2:
		image_tag = image_parts[1]
	else:
		image_tag = 'latest'

	p = subprocess.Popen(["/builder/google-cloud-sdk/bin/gcloud alpha container images list-tags --format='value(tags)' {0}".format(image_path)],
		shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

	output, error = p.communicate()
	if p.returncode != 0:
		sys.exit("Error encountered when retrieving existing image tags! Full log: \n\n" + output)

	existing_tags = set(tag.rstrip() for tag in output.split('\n'))
	print "Existing tags for image {0}:".format(image_path)
	for tag in existing_tags:
		print tag

	return image_tag in existing_tags


def main():
	parser = argparse.ArgumentParser()
	parser.add_argument('--image', type=str, help='Fully qualified remote path for the target image')
	args = parser.parse_args()

	if check_if_tag_exists(args.image):
		sys.exit("Tag already exists in remote repository! Exiting build.")
	print "Tag does not exist in remote repository! Continuing with build."

if __name__ == "__main__":
	main()

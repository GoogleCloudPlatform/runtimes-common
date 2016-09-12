#!/usr/bin/python

import sys
import subprocess
import os
from sets import Set

AUTH_FILE_PATH_REMOTE = "gs://nick-cloudbuild-auth/nick-cloudbuild-743a54b95573.json"
AUTH_FILE_PATH_LOCAL = "/auth.json"

def dedupe(input_tag, existing_tags):
	for tag in existing_tags:
		if input_tag == tag:
			sys.exit("Tag already exists in remote repository! Exiting build.")
	print "Tag does not exist in remote repository! Continuing with build."


# Since gcloud auth information doesn't transfer automatically into a new container in prod, 
# we use a placeholder service account so we can make authenticated calls to the gcloud API.
def login():
	#######################################################################################
	# TODO: once we can make unauthenticated gsutil calls in prod, 
	# uncomment this and remove the temporary auth.json file from the docker image.

	# p1 = subprocess.Popen(["/builder/google-cloud-sdk/bin/gsutil cp {0} {1}".format(AUTH_FILE_PATH_REMOTE, AUTH_FILE_PATH_LOCAL)],
	# 	shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

	# for line in p1.stdout.readlines():
	# 	print line
	# 	if "ERROR" in line or "Exception" in line:
	# 		sys.exit("Error encountered when retrieving gcloud credentials. Exiting build. ({0})".format(line))
	#######################################################################################

	p2 = subprocess.Popen(["/builder/google-cloud-sdk/bin/gcloud auth activate-service-account --key-file {0}".format(AUTH_FILE_PATH_LOCAL)],
		shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

	for line in p2.stdout.readlines():
		print line
		if "ERROR" in line or "Exception" in line: sys.exit("Error encountered when logging into gcloud. Exiting build.")

	os.remove(AUTH_FILE_PATH_LOCAL)

def main():
	args = sys.argv[1:] # ignore script name
	if len(args) < 1:
		sys.exit('Please provide the qualified path to your image, in the form: $REPO_ID/$PROJECT_ID/$IMAGE_NAME:$TAG')

	login()

	# extract both path to image, and the tag if provided
	image_parts = args[0].split(':')
	image_path = image_parts[0]
	if len(image_parts) == 2:
		image_tag = image_parts[1]
	else:
		image_tag = 'latest'

	p = subprocess.Popen(["/builder/google-cloud-sdk/bin/gcloud alpha container images list-tags --format='value(tags)' {0}".format(image_path)],
		shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

	existing_tags = Set(map(lambda tag : tag.rstrip('\n'), p.stdout.readlines()))

	if "ERROR" in existing_tags:
		sys.exit("Error encountered when retrieving existing image tags. Full log:" + existing_tags)

	print "Existing tags for image {0}:".format(image_path)
	for tag in existing_tags:
		print tag

	dedupe(image_tag, existing_tags)

if __name__ == "__main__":
	main()

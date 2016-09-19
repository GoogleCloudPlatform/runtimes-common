#!/usr/bin/python

import sys
import subprocess
import argparse

AUTH_FILE_PATH_LOCAL = "/auth.json"


def check_if_tag_exists(raw_image_path, force_build):
    # extract both path to image, and the tag if provided
    image_parts = raw_image_path.split(':')
    image_path = image_parts[0]
    if len(image_parts) == 2:
        image_tag = image_parts[1]
    else:
        image_tag = 'latest'

    p = subprocess.Popen(["/builder/google-cloud-sdk/bin/gcloud "
                          + "alpha container images list-tags "
                          + "--format='value(tags)' {0}".format(image_path)],
                         shell=True, stdout=subprocess.PIPE,
                         stderr=subprocess.STDOUT)

    output, error = p.communicate()
    if p.returncode != 0:
        sys.exit("Error encountered when retrieving existing image tags! "
                 + "Full log: \n\n" + output)

    existing_tags = set(tag.rstrip() for tag in output.split('\n'))
    print "Existing tags for image {0}:".format(image_path)
    for tag in existing_tags:
        print tag

    if image_tag in existing_tags:
        print "Tag \'" + image_tag + "\' already exists in remote repository!"
        if not force_build:
            sys.exit("Exiting build.")
        else:
            print "Forcing build. Tag \'" + image_tag + "\' will be overwritten!"
            return
    print "Tag \'" + image_tag + "\' does not exist in remote repository! Continuing with build."


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--image', type=str,
                        help='Fully qualified remote path for the '
                        + 'target image')
    parser .add_argument('--force', action='store_true', default=False)
    args = parser.parse_args()

    check_if_tag_exists(args.image, args.force)

if __name__ == "__main__":
    main()

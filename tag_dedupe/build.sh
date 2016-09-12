#!/bin/bash

set -e

IMAGE=$1
REMOTE_AUTH_FILE_PATH=$2

if [ -z "$1" ]; then
	echo "Usage: ./build.sh [image_path] [auth_file_remote_path]"
	echo "Please provide fully qualified path to target image."
	exit 1
fi

if [ -z "$2" ]; then
	echo "Path to remote auth file path not provided; using auth file bundled into container!"
fi

sed -i "s|%IMAGE%|$IMAGE|g" cloudbuild.yaml
sed -i "s|%AUTH_FILE%|$AUTH_FILE|g" cloudbuild.yaml
gcloud alpha container builds create . --config=cloudbuild.yaml

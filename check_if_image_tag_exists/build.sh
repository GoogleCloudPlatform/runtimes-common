#!/bin/bash

set -e

IMAGE=$1

if [ -z "$1" ]; then
  echo "Usage: ./build.sh [image_path] [auth_file_remote_path]"
  echo "Please provide fully qualified path to target image."
  exit 1
fi

sed -i "s|%IMAGE%|$IMAGE|g" cloudbuild.yaml
gcloud alpha container builds create . --config=cloudbuild.yaml

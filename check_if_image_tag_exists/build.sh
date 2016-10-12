#!/bin/bash

set -e

export IMAGE=$1

if [ -z "$1" ]; then
  echo "Usage: ./build.sh [image_path] [auth_file_remote_path]"
  echo "Please provide fully qualified path to target image."
  exit 1
fi

envsubst '${IMAGE}' < cloudbuild.yaml.in > cloudbuild.yaml
gcloud alpha container builds create . --config=cloudbuild.yaml

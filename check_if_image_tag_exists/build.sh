#!/bin/bash

usage() { echo "Usage: ./build.sh [target_image_path] [schema_version]"; exit 1; }

set -e

export IMAGE=$1
export VERSION=$2

if [ -z $IMAGE ]; then
  usage
else if [ -z $VERSION ]; then
	echo "Defaulting to latest JSON schema version..."
	export VERSION="latest"
fi

envsubst < cloudbuild.yaml.in > cloudbuild.yaml
gcloud alpha container builds create . --config=cloudbuild.yaml

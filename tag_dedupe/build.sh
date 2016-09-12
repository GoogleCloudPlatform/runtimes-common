#!/bin/bash

set -e
set -x

IMAGE=$1

sed -i "s|%IMAGE%|$IMAGE|g" cloudbuild.yaml
gcloud alpha container builds create . --config=cloudbuild.yaml

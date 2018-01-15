#!/bin/sh

# Copyright 2017 Google Inc. All rights reserved.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

usage() { echo "Usage: ./build.sh [target_image_path]"; exit 1; }

set -e

export IMAGE=$1

if [ -z "$IMAGE" ]; then
  usage
fi

cd ..
gcloud container builds submit . --config apt_installer/cloudbuild.yaml --substitutions=_IMAGE="$IMAGE"

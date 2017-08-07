#!/bin/bash

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

set -euo pipefail

VERBOSE=0
PULL=1
IMAGE_TAR=""
declare -a CMD_STRING
CMD_STRING=()
ENTRYPOINT="/test/structure_test"
ST_IMAGE="gcr.io/gcp-runtimes/structure_test"
USAGE_STRING="Usage: $0 [-i <image>] [-t <image.tar>] [-c <config>] [-w <workspace>] [-v] [-e <entrypoint>] [--no-pull]"

CONFIG_DIR=$(pwd)/.cfg
mkdir -p "$CONFIG_DIR"

declare -a VOLUME_STR
VOLUME_STR=()
# With pipefail and set -e we just silently exit here instead of printing the error message.
set +e
command -v docker > /dev/null 2>&1 || { echo "Docker is required to run GCP structure tests, but is not installed on this host."; exit 1; }
command docker ps > /dev/null 2>&1 || { echo "Cannot connect to the Docker daemon"; exit 1; }
set -e

cleanup() {
	rm -rf "$CONFIG_DIR"
}

usage() {
	echo "$USAGE_STRING"
	cleanup
	exit 1
}

helper() {
	echo "$USAGE_STRING"
	echo
	echo "    -i, --image          Image to run tests on"
	echo "    -t, --image-tar      Image tarball to run tests on"
	echo "    -c, --config         Path to JSON/YAML config file"
	echo "    -w, --workspace      Path to directory to be mounted as"
	echo "                         /workspace in remote container."
	echo "                         This is most likely the project root for the image, as all"
	echo "                         resources in the directory will be included in the test run."
	echo "                         (e.g. -w ../../python-runtime)"
	echo "    -v                   Display verbose testing output"
	echo "    -e, --entrypoint     Specify custom docker entrypoint for image"
	echo "    --no-pull            Don't pull latest structure test image"
	exit 0
}

while test $# -gt 0; do
	case "$1" in
		--image|-i)
			shift
			if test $# -gt 0; then
				IMAGE_NAME=$1
			else
				usage
			fi
			shift
			;;
		--image-tar|-t)
			shift
			if test $# -gt 0; then
				IMAGE_TAR=$1
			else
				usage
			fi
			shift
			;;
		--verbose|-v)
			VERBOSE=1
			shift
			;;
		--no-pull)
			PULL=0
			shift
			;;
		--help|-h)
			helper
			;;

		--workspace|-w)
			shift
			if test $# -eq 0; then
				usage
			else
				if [ ! -d "$1" ] || [ ! -d "$(readlink -f "$1")" ]; then
				        echo "$1 is not a valid directory."
				        cleanup
				        exit 1
				fi
				FULLPATH=$(readlink -f "$1")
				VOLUME_STR+=(-v "$FULLPATH:/workspace")
			fi
			shift
			;;
		--config|-c)
			shift
			if test $# -eq 0; then
				usage
			else
				if [ ! -f "$1" ]; then
					echo "$1 is not a valid file."
					cleanup
					exit 1
				fi
				# structure tests allow specifying any number of configs,
				# which can live anywhere on the host file system. to simplify
				# the docker volume mount, we copy all of these configs into
				# a /tmp directory and mount this single directory into the
				# test image. this directory is cleaned up after testing.
				filename=$(basename "$1")
				cp "$1" "$CONFIG_DIR"/"$filename"
				# Make the config file world-readable in-case the image has a non-root user specified.
				chmod +r "$CONFIG_DIR"/"$filename"
				CMD_STRING+=(--config "/cfg/$filename")
			fi
			shift
			;;
		*)
			usage
			;;
	esac
done

if [ -z "$IMAGE_NAME" ]; then
	usage
fi

if [ $VERBOSE -eq 1 ]; then
	CMD_STRING+=(-test.v)
fi

if [ $PULL -eq 1 ]; then
	docker pull "$ST_IMAGE"
fi

if [ -n "$IMAGE_TAR" ]; then
	docker load -i "$IMAGE_TAR"
fi

st_container=$(docker run -d --entrypoint="/bin/sh" "$ST_IMAGE" 2>/dev/null)
VOLUME_STR+=(--volumes-from "$st_container" -v "$CONFIG_DIR:/cfg")

# shellcheck disable=SC2086
docker run --rm --entrypoint="$ENTRYPOINT" "${VOLUME_STR[@]}" "$IMAGE_NAME" "${CMD_STRING[@]}"

docker rm "$st_container" > /dev/null 2>&1
cleanup

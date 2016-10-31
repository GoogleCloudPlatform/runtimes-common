#!/bin/sh

command -v docker > /dev/null 2>&1 || { echo "Docker is required to run GCP structure tests, but is not installed on this host."; exit 1; }

command docker ps > /dev/null 2>&1 || { echo "Cannot connect to the Docker daemon!"; exit 1; }

usage() {
	echo "Usage: $0 [-i <image>] [-c <config>] [-v] [-e <entrypoint>]"
	exit 1
}

export VERBOSE=0
export CMD_STRING=""
export ENTRYPOINT="./test/structure_test"
export ST_IMAGE="gcr.io/nick-cloudbuild/structure_test"

declare -i CONFIG_COUNTER=0

mkdir /tmp/st_configs

while test $# -gt 0; do
	case "$1" in
		--image|-i)
			shift
			if test $# -gt 0; then
				export IMAGE_NAME=$1
			else
				usage
			fi
			shift
			;;
		--verbose|-v)
			export VERBOSE=1
			shift
			;;
		--entrypoint|-e)
			shift
			if test $# -gt 0; then
				export ENTRYPOINT=$1
			else
				usage
			fi
			shift
			;;
		--config|-c)
			shift
			if test $# -eq 0; then
				usage
			else
				cp "$1" /tmp/st_configs/cfg_$CONFIG_COUNTER.json
				CMD_STRING=$CMD_STRING" --config /cfg/cfg_$CONFIG_COUNTER.json"
				CONFIG_COUNTER+=1
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
	CMD_STRING=$CMD_STRING" -test.v"
fi

docker pull "$ST_IMAGE"

ST_CONTAINER=$(docker run -d --entrypoint="/bin/sh" --name st_container "$ST_IMAGE")
if [ -z "$ST_CONTAINER" ]; then
	exit 1
fi

TEST_CONTAINER=$(docker create --entrypoint="$ENTRYPOINT" --volumes-from st_container -v /tmp/st_configs:/cfg "$IMAGE_NAME" $CMD_STRING)
if [ -z "$TEST_CONTAINER" ]; then
	exit 1
fi

docker start -a -i "$TEST_CONTAINER"

docker logs -f "$TEST_CONTAINER" | while read LINE
do
	continue
done

docker rm "$TEST_CONTAINER" > /dev/null 2>&1
docker rm "$ST_CONTAINER" > /dev/null 2>&1
rm -rf /tmp/st_configs

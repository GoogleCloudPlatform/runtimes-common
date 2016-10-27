#!/bin/sh

set -x

command -v docker > /dev/null 2>&1 || { echo "Docker is required to run GCP structure tests, but is not installed on this host."; exit 1; }

usage() {
	echo "Usage: $0 [-i <image>] [-c <config>] [-v] [-e <entrypoint>]"
	exit 1
}

export VERBOSE=0
export CMD_STRING=""
export ENTRYPOINT="./structure_test"

declare -i CONFIG_COUNTER=0

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
				cp $1 ./cfg_$CONFIG_COUNTER.json
				CMD_STRING=$CMD_STRING" --config cfg_$CONFIG_COUNTER.json"
				CONFIG_COUNTER+=1
			fi
			shift
			;;
		*)
			usage
			;;
	esac
done

if [ -z $IMAGE_NAME ]; then
	usage
fi

if [ $VERBOSE -eq 1 ]; then
	CMD_STRING=$CMD_STRING" -test.v"
fi

docker pull gcr.io/gcp-runtimes/structure_test

CONTAINER=$(docker run -d --entrypoint="/bin/sh" gcr.io/gcp-runtimes/structure_test)
if [ -z $CONTAINER ]; then
	exit 1
fi

docker cp $CONTAINER:/test/structure_test .
docker rm $CONTAINER

CONTAINER=$(docker create --entrypoint=$ENTRYPOINT $IMAGE_NAME $CMD_STRING)
if [ -z $CONTAINER ]; then
	exit 1
fi

docker cp structure_test $CONTAINER:/structure_test
for f in cfg_*.json; do
	docker cp $f $CONTAINER:/
	rm $f
done
# docker cp cfg_*.json $CONTAINER:/
rm structure_test
# rm cfg_*.json
docker start -a -i $CONTAINER

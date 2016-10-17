#!/bin/sh

export DOCKER_API_VERSION="1.21"

export VERBOSE=0
export CONFIG=()
while test $# -gt 0; do
	case "$1" in
		--image|-i)
			shift
			if test $# -gt 0; then
				export IMAGE_NAME=$1
			else
				echo "Please provide fully qualified path to image under test."
				exit 1
			fi
			shift
			;;
		--verbose|-v)
			export VERBOSE=1
			shift
			;;
		--config|-c)
			shift
			if test $# -eq 0; then
				echo "Please provide fully qualified path to config file."
				exit 1
			else
				while test $# -gt 0; do
					CONFIG=("${CONFIG[@]}" "$1")
					shift
				done
			fi
			shift
			;;
		*)
			shift
			;;
	esac
done

cp /test/* /workspace/

export CMD_STRING="/workspace/structure_test"

if [ $VERBOSE -eq 1 ]; then
	CMD_STRING=$CMD_STRING" -v"
fi

if [ -n "$CONFIG" ]; then
	CMD_STRING=$CMD_STRING" --config ${CONFIG[@]}"
fi
echo "$CMD_STRING"

docker run --privileged=true -v /workspace:/workspace "$IMAGE_NAME" /bin/sh -c "$CMD_STRING"

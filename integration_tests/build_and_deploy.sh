#!/bin/sh

set -ex

USAGE_STRING="Usage: $0 [-d <directory>]"
# VERSION="integration-test-$(date +'%s')"
# TODO: deploy with version
VERSION=""

WORKDIRS=()

# trap cleanup 0 1 2 3 13 15 # EXIT HUP INT QUIT PIPE TERM

# cleanup() {
# 	gcloud app versions stop "$VERSION"
# }

usage() {
	echo "$USAGE_STRING"
	exit 1
}

while test $# -gt 0; do
	case "$1" in
		--directory|-d)
			shift
			if test $# -gt 0; then
				directory=$(readlink -f "$1")
				WORKDIRS+=(${directory})
			else
				usage
			fi
			shift
			;;
		--image|-i)
			shift
			if test $# -gt 0; then
				export STAGING_IMAGE="$1"
			else
				usage
			fi
			shift
			;;
		*)
			usage
			;;
	esac
done

if [ ${#WORKDIRS[@]} -eq 0 ]; then
	usage
fi

deploy_and_test_app() {
	workdir="$1"
	if [ -z "$workdir" ]; then
		usage
	fi

	echo $workdir

	# cd "$workdir"
	# cp /app.yaml .
	# envsubst < Dockerfile.in > Dockerfile

	# # docker build -t sample_app .
	# gcloud auth activate-service-account --key-file=/auth.json
	# gcloud auth list
	# gcloud config list
	# # gcloud config list project
	# gcloud app deploy --stop-previous-version --verbosity=debug
	# # gcloud app deploy --stop-previous-version --version "$VERSION" --verbosity=debug

	# PROJECT_ID="nick-cloudbuild"
	# # URL="https://$PROJECT_ID_$VERSION.appspot.com"
	# URL="https://$PROJECT_ID.appspot.com"

	# echo "giving app 20 seconds to deploy and become responsive..."
	# sleep 20
	# echo "awake! hitting $URL now"

	# APP_RESPONSE=$(curl --silent -L "$URL")

	# echo "response from app was: $APP_RESPONSE"
}

for directory in ${WORKDIRS[@]}; do
	deploy_and_test_app ${directory}
done

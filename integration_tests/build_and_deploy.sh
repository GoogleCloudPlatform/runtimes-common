#!/bin/sh

set -e

USAGE_STRING="Usage: $0 [-d <directory>]"
# VERSION="integration-test-$(date +'%s')"
# TODO: deploy with version
VERSION=""

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
				WORKDIR=$(readlink -f "$1")
			else
				usage
			fi
			shift
			;;
		--runtime|-r)
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

if [ -z "$WORKDIR" ]; then
	usage
fi

cd "$WORKDIR"
cp /app.yaml .
envsubst < Dockerfile.in > Dockerfile

# docker build -t sample_app .
gcloud auth activate-service-account --key-file=/auth.json
gcloud auth list
gcloud config list
# gcloud config list project
gcloud app deploy --stop-previous-version --verbosity=debug
# gcloud app deploy --stop-previous-version --version "$VERSION" --verbosity=debug

PROJECT_ID="nick-cloudbuild"
# URL="https://$PROJECT_ID_$VERSION.appspot.com"
URL="https://$PROJECT_ID.appspot.com"

echo "sleeping for 20 seconds..."
sleep 20
echo "awake! hitting $URL now"

APP_RESPONSE=$(curl -L "$URL")

echo "$APP_RESPONSE"

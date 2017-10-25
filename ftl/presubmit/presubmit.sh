#!/bin/bash
# shellcheck disable=SC1090
set -ex
source "$KOKORO_GFILE_DIR/common.sh"
cd github/runtimes-common
test_tag="ftl-integration-kokoro-presubmit-$KOKORO_BUILD_NUMBER"
TAG=$test_tag gcloud container builds submit --config ftl/ftl_node_integration_tests.yaml .

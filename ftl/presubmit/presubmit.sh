#!/bin/bash
# shellcheck disable=SC1090
set -ex
source "$KOKORO_GFILE_DIR/common.sh"
cd github/runtimes-common
gcloud container builds submit --config ftl/ftl_node_integration_tests.yaml .

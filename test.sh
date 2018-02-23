#!/bin/bash

set -ex
echo "Checking formatting..."
find . -name "*.sh" | grep -v vendor/ | xargs shellcheck
flake8 .
./.gofmt.sh
./.buildifier.sh

echo "Running tests..."
bazel test --test_output=errors appengine/reconciletags:reconciletags_test
bazel test --test_output=errors appengine/reconciletags:reconciletags_e2e_test
bazel test --test_output=errors ftl/... --deleted_packages=ftl/node/benchmark,ftl/php/benchmark,ftl/python/benchmark,ftl/benchmark
bazel test --test_output=errors testing/lib:mock_registry_tests
pushd appengine/runtime_builders && py.test test_manifest.py && popd

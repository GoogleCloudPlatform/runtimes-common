#!/bin/bash

set -ex
echo "Checking formatting..."
find . -name "*.sh" | grep -v vendor/ | xargs shellcheck
flake8 .
./.gofmt.sh
./.buildifier.sh

echo "Running tests..."
bazel test --test_output=errors reconciletags:reconciletags_test
bazel test --test_output=errors ftl/... --deleted_packages=ftl/node/benchmark,ftl/php/benchmark,ftl/benchmark
bazel test --test_output=errors testing/lib:mock_registry_tests
cd runtime_builders && py.test test_manifest.py && cd ..

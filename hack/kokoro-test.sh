#!/bin/bash
set -ex

echo "Installing dependencies..."

# shellcheck source=/dev/null
source "$KOKORO_GFILE_DIR/common.sh"

# Grab the latest version of shellcheck and add it to PATH
if [ -f "$KOKORO_GFILE_DIR"/shellcheck-latest.linux ]; then
    sudo cp "$KOKORO_GFILE_DIR"/shellcheck-latest.linux /usr/local/bin/shellcheck
    sudo cp "$KOKORO_GFILE_DIR"/shellcheck-latest.linux /usr/local/bin/shellcheck
    sudo chmod +x /usr/local/bin/shellcheck
fi

# Add the keys.json from keystore to expected bazel target
if [ -f "$KOKORO_KEYSTORE_DIR"/72508_tuf_integration_test ]; then
  cp "$KOKORO_KEYSTORE_DIR"/72508_tuf_integration_test tuf/keys.json
fi

pushd github/runtimes-common
# Install deps.
sudo pip install --upgrade -r requirements.txt

echo "Running unit tests..."
# Run the tests.
./test.sh

echo "Running integration tests..."
# Run these in parallel.
declare -a pids=()
# Generate the integration test yaml and pass it to cloud build via stdin.
python ftl/integration_tests/ftl_node_integration_tests_yaml.py | gcloud container builds submit --config /dev/fd/0 . > node.log &
pids+=($!)

python ftl/integration_tests/ftl_php_integration_tests_yaml.py | gcloud container builds submit --config /dev/fd/0 . > php.log &
pids+=($!)

python ftl/integration_tests/ftl_python_integration_tests_yaml.py | gcloud container builds submit --config /dev/fd/0 . > python.log &
pids+=($!)

gcloud container builds submit --config ftl/integration_tests/ftl_python_error_test.yaml . > python_error.log &
pids+=($!)

python ftl/cached/ftl_cached_yaml.py --runtime=node-same | gcloud container builds submit --config /dev/fd/0 . > node_same_cached.log &
pids+=($!)

python ftl/cached/ftl_cached_yaml.py --runtime=node-plus-one | gcloud container builds submit --config /dev/fd/0 . > node_plus_one_cached.log &
pids+=($!)

python ftl/cached/ftl_cached_yaml.py --runtime=php-same | gcloud container builds submit --config /dev/fd/0 . > php_same_cached.log &
pids+=($!)

python ftl/cached/ftl_cached_yaml.py --runtime=php-lock-same | gcloud container builds submit --config /dev/fd/0 . > php_lock_same_cached.log &
pids+=($!)

python ftl/cached/ftl_cached_yaml.py --runtime=php-lock-plus-one | gcloud container builds submit --config /dev/fd/0 . > php_lock_plus_one_cached.log &
pids+=($!)

python ftl/cached/ftl_cached_yaml.py --runtime=python-requirements-same | gcloud container builds submit --config /dev/fd/0 . > python_requirements_same_cached.log &
pids+=($!)

python ftl/cached/ftl_cached_yaml.py --runtime=python-requirements-plus-one | gcloud container builds submit --config /dev/fd/0 . > python_requirements_plus_one_cached.log &
pids+=($!)

python ftl/cached/ftl_cached_yaml.py --runtime=python-pipfile-same | gcloud container builds submit --config /dev/fd/0 . > python_pipfile_same_cached.log &
pids+=($!)

python ftl/cached/ftl_cached_yaml.py --runtime=python-pipfile-plus-one | gcloud container builds submit --config /dev/fd/0 . > python_pipfile_plus_one_cached.log &
pids+=($!)

# Wait for them to finish, and check the exit codes.

failures=0
set +e

for pid in "${pids[@]}"; do
    wait "$pid"
    status=$?
    failures+=$status
done
set -e

if [[ $failures -gt 0 ]]; then
    echo "Integration test failure."
    cat node.log
    cat python.log
    cat php.log
    cat node_same_cached.log
    cat node_plus_one_cached.log
    cat php_same_cached.log
    cat php_lock_same_cached.log
    cat php_lock_plus_one_cached.log
    cat python_requirements_same_cached.log
    cat python_requirements_plus_one_cached.log
    cat python_pipfile_same_cached.log
    cat python_pipfile_plus_one_cached.log
    exit 1
fi

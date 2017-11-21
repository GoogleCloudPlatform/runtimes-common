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
python ftl/ftl_node_integration_tests_yaml.py | gcloud container builds submit --config /dev/fd/0 . > node.log &
pids+=($!)

python ftl/ftl_php_integration_tests_yaml.py | gcloud container builds submit --config /dev/fd/0 . > php.log &
pids+=($!)

gcloud container builds submit --config ftl/ftl_python_integration_tests.yaml . > python.log &
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
    exit 1
fi

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
# Generate the integration test yaml and pass it to cloud build via stdin.
python ftl/ftl_node_integration_tests_yaml.py | gcloud container builds submit --config /dev/fd/0 .
python ftl/ftl_php_integration_tests_yaml.py | gcloud container builds submit --config /dev/fd/0 .
gcloud container builds submit --config ftl/ftl_python_integration_tests.yaml .

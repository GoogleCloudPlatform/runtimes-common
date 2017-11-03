#!/bin/bash
set -ex

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

# Run the tests.
./test.sh

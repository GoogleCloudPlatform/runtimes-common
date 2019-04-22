#!/bin/bash
set -ex

if [ -z "${CONFIG_FILE}" ]
then
  echo "CONFIG_FILE is not set"
  exit 1
fi

# shellcheck source=/dev/null
source "$KOKORO_GFILE_DIR/common.sh"

cd github
yes | sudo pip install ruamel.yaml
python runtimes-common/appengine/runtime_builders/template_builder.py -f "${KOKORO_GFILE_DIR}/${CONFIG_FILE}"

#!/bin/bash
set -ex

if [ -z ${CONFIG_FILE} ]
then
  echo "CONFIG_FILE is not set"
  exit 1
fi

source $KOKORO_GFILE_DIR/common.sh

cd github/runtimes-common/appengine/runtime_builders
yes | sudo pip install ruamel.yaml
python template_builder.py -f "${KOKORO_GFILE_DIR}/${CONFIG_FILE}

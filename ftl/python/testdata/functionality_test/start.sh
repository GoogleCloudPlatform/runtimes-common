#!/bin/bash
set -xe
if [ "$(cat /workspace/ftl/python/testdata/additional_directory/child/b.txt)" = "barfoo" ]; then
  python /srv/app.py &
  wget --tries=20 localhost:8080
else
  exit 1
fi

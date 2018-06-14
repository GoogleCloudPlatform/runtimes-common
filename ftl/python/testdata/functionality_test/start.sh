#!/bin/bash
set -xe
python /srv/app.py &
wget --tries=20 localhost:8080 &
[ "$(cat /workspace/ftl/python/testdata/additional_directory/child/b.txt)" == "barfoo" ]

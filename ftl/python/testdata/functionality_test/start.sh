#!/bin/bash
set -xe
[ "$(cat /workspace/ftl/python/testdata/additional_directory/child/b.txt)" == "barfoo" ] &&
python /srv/app.py &
wget --tries=20 localhost:8080

#!/bin/bash
set -xe
python /srv/app.py &
wget --tries=20 localhost:8080

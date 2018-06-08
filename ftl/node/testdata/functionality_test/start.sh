#!/bin/bash
set -xe
cd /srv
npm start &
wget --tries=20 localhost:5000

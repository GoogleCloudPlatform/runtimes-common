#!/bin/bash
set -xe
cd /srv
php artisan serve &
sleep 5s
curl --connect-timeout 5 \
     --max-time 10 \
     --retry 5 \
     --retry-delay 0 \
     --retry-max-time 60 \
     localhost:8000

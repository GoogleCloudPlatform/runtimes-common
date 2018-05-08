#!/bin/bash

if [ -e "tuf/keys.json" ]
then
    echo "PASS"
else
    echo "FAIL"
    exit 1
fi

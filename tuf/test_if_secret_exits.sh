#!/bin/bash

if [ -e keys.json ]
then
    echo "PASS"
else
    echo "FAIL"
    exit 1
fi

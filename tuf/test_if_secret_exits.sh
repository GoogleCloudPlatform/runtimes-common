#!/bin/bash

if [[ ! -v TUF_INTEGRATION_TEST ]]; then
    echo "TUF_INTEGRATION_TEST is not set"
    exit 1
else
    echo "TUF_INTEGRATION_TEST is set to $TUF_INTEGRATION_TEST"
fi

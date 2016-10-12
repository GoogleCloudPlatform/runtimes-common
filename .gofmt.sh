#!/bin/bash

set -e
files=$(gofmt -l -s .)
if [[ $files ]]; then
    echo "Gofmt errors in files: $files"
    exit 1
fi

files=$(go vet ./structure_tests)
if [[ $files ]]; then
   echo "Go vet errors in files: $files"
   exit 1
fi

#!/bin/bash

set -e
fmt_errors=$(gofmt -d .)
if [[ $fmt_errors ]]; then
    echo "$fmt_errors"
    exit 1
fi

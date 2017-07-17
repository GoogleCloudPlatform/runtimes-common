#!/bin/bash
go run iDiff/main.go iDiff gcr.io/google_containers/busybox:1.24 gcr.io/google_containers/busybox:latest dir -j > iDiff/tests/busybox_diff_actual.json
if [[ $? -ne 0 ]]; then
  echo "iDiff simple run failed"
  exit 1
fi

diff=$(diff iDiff/tests/busybox_diff_expected.json iDiff/tests/busybox_diff_actual.json)
if [[ -n "$diff" ]]; then
  echo $diff
  echo "iDiff output is not as expected"
  exit 1
fi

go test `go list ./... | grep iDiff | grep -v iDiff/vendor`


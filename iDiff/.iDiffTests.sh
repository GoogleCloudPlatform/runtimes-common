#!/bin/bash
go run iDiff/main.go iDiff 0cb40641836c e7d168d7db45 dir -j > iDiff/tests/runs/busybox_diff_actual.json
diff=$(diff iDiff/tests/busybox_diff_expected.json iDiff/tests/runs/busybox_diff_actual.json)
if [ $diff ]; then
  echo "iDiff output is not as expected"
  exit 1
fi

go test `go list ./... | grep iDiff | grep -v iDiff/vendor`


#!/bin/bash
go run iDiff/main.go iDiff file gcr.io/google_containers/busybox:1.24 gcr.io/google_containers/busybox:latest -j > iDiff/tests/busybox_diff_actual.json
if [[ $? -ne 0 ]]; then
  echo "iDiff simple run failed"
  exit 1
fi

python iDiff/fileDiff_test_processor.py iDiff/tests/busybox_diff_expected.json
if [[ $? -ne 0 ]]; then
  echo "Could not process expected test file for file diff comparison"
  exit 1
fi
python iDiff/fileDiff_test_processor.py iDiff/tests/busybox_diff_actual.json
if [[ $? -ne 0 ]]; then
  echo "Could not process actual test file for file diff comparison"
  exit 1
fi
diff=$(diff iDiff/tests/busybox_diff_expected.json iDiff/tests/busybox_diff_actual.json)
if [[ -n "$diff" ]]; then
  echo "iDiff file diff output is not as expected"
  echo $diff
  exit 1
fi

go test `go list ./... | grep iDiff | grep -v iDiff/vendor`


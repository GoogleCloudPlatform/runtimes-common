#!/bin/bash
while read -r differ image1 image2 file; do
  go run iDiff/main.go iDiff $differ $image1 $image2 -j > $file
  if [[ $? -ne 0 ]]; then
    echo "iDiff" "$differ" "differ failed"
    exit 1
  fi
done < iDiff/tests/differ_runs.txt

python iDiff/tests/fileDiff_test_processor.py iDiff/tests/file_diff_expected.json
if [[ $? -ne 0 ]]; then
  echo "Could not process expected test file for file diff comparison"
  exit 1
fi
python iDiff/tests/fileDiff_test_processor.py iDiff/tests/file_diff_actual.json
if [[ $? -ne 0 ]]; then
  echo "Could not process actual test file for file diff comparison"
  exit 1
fi

while read -r differ actual expected; do
  success=0
  diff=$(diff "$actual" "$expected")
  if [[ -n "$diff" ]]; then
    echo "iDiff" "$differ" "diff output is not as expected"
    echo $diff
    success=1
  fi
done < iDiff/tests/diff_comparisons.txt
if [[ "$success" -ne 0 ]]; then
  exit 1
fi

go test `go list ./... | grep iDiff | grep -v iDiff/vendor`


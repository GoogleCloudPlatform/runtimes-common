#!/bin/bash
while IFS=$' \n\r' read -r differ image1 image2 file; do
  go run iDiff/main.go $image1 $image2 --$differ -j > $file
  if [[ $? -ne 0 ]]; then
    echo "iDiff" "$differ" "differ failed"
    exit 1
  fi
done < iDiff/tests/differ_runs.txt

while IFS=$' \n\r' read -r preprocess json; do
  python $preprocess $json
  if [[ $? -ne 0 ]]; then
  echo "Could not preprocess" "$json" "for diff comparison"
  exit 1
fi
done < iDiff/tests/preprocess_files.txt

while IFS=$' \n\r' read -r differ actual expected; do
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


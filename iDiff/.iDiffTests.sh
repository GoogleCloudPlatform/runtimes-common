#!/bin/bash
while IFS=$' \n\r' read -r flag differ image1 image2 file; do
  go run iDiff/main.go $image1 $image2 $flag -j > $file
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

success=0
while IFS=$' \n\r' read -r differ actual expected; do
  diff=$(jq --argfile a "$actual" --argfile b "$expected" -n 'def walk(f): . as $in | if type == "object" then reduce keys[] as $key ( {}; . + { ($key):  ($in[$key] | walk(f)) } ) | f elif type == "array" then map( walk(f) ) | f else f end; ($a | walk(if type == "array" then sort else . end)) as $a | ($b | walk(if type == "array" then sort else . end)) as $b | $a == $b')
  if ! "$diff" ; then
    echo "iDiff" "$differ" "diff output is not as expected"
    cat "$actual"
    success=1
  fi
done < iDiff/tests/diff_comparisons.txt
if [[ "$success" -ne 0 ]]; then
  exit 1
fi

go test `go list ./... | grep iDiff | grep -v iDiff/vendor`

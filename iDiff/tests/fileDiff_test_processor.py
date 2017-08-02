import json
import re
import sys


def _process_test_diff(file_path):
    with open(file_path) as f:
        diff = json.load(f)

    diff_result = diff["Diff"]
    diff_result["Adds"] = _trim_file_names(diff_result["Adds"])
    diff_result["Dels"] = _trim_file_names(diff_result["Dels"])

    with open(file_path, 'w') as f:
        json.dump(diff, f, indent=4)


def _trim_file_names(files):
    trimmed_files = []
    for file in files:
        trimmed_file = _trim_layer_hash(file)
        trimmed_files.append(trimmed_file)
    return sorted(trimmed_files)


def _trim_layer_hash(filename):
    hash_match = re.match(r'^([a-z|0-9])/', filename)
    if hash_match:
        hash = hash_match.group(1)
        return re.sub(hash, "", filename)
    return filename


if __name__ == '__main__':
    sys.exit(_process_test_diff(sys.argv[1]))

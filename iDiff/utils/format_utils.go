package utils

import "encoding/json"

func JSONify(diff interface{}) (string, error) {
	diffBytes, err := json.MarshalIndent(diff, "", "  ")
	if err != nil {
		return "", err
	}
	return string(diffBytes), nil
}

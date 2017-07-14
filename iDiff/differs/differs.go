package differs

import (
	"bytes"
	"errors"
	"os"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

var diffs = map[string]func(string, string, bool) (string, error){
	"hist": History,
	"dir":  Package,
	"apt":  AptDiff,
	"node": NodeDiff,
	"pip":  PipDiff,
}

func Diff(arg1, arg2, differ string, json bool) (string, error) {
	if f, exists := diffs[differ]; exists {
		if differ == "hist" {
			return f(arg1, arg2, json)
		}
		return specificDiffer(f, arg1, arg2, json)
	}
	return "", errors.New("Unknown differ")
}

func specificDiffer(f func(string, string, bool) (string, error), img1, img2 string, json bool) (string, error) {
	var buffer bytes.Buffer
	validDiff := true
	jsonPath1, dirPath1, err := utils.ImageToDir(img1)
	if err != nil {
		buffer.WriteString(err.Error())
		validDiff = false
	}
	jsonPath2, dirPath2, err := utils.ImageToDir(img2)
	if err != nil {
		buffer.WriteString(err.Error())
		validDiff = false
	}

	var diff string
	if validDiff {
		output, err := f(jsonPath1, jsonPath2, json)
		if err != nil {
			buffer.WriteString(err.Error())
		}
		diff = output
	}

	errStr := remove(dirPath1, true, "")
	errStr = remove(dirPath2, true, errStr)
	errStr = remove(jsonPath1, false, errStr)
	errStr = remove(jsonPath2, false, errStr)
	if errStr != "" {
		buffer.WriteString(errStr)
	}

	if buffer.String() != "" {
		return diff, errors.New(buffer.String())
	}
	return diff, nil
}

func remove(path string, dir bool, errStr string) string {
	if path == "" {
		return ""
	}

	var err error
	if dir {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		errStr += "\nUnable to remove " + path
	}
	return errStr
}

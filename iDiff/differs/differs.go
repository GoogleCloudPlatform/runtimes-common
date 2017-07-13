package differs

import (
	"errors"
	"os"
	"reflect"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

var diffs = map[string]func(string, string, bool) (string, error){
	"hist":    HistoryDiff,
	"history": HistoryDiff,
	"file":    FileDiff,
	"apt":     AptDiff,
	"linux":   AptDiff,
}

func Diff(arg1, arg2, differ string, json bool) (string, error) {
	if f, exists := diffs[differ]; exists {
		fValue := reflect.ValueOf(f)
		histValue := reflect.ValueOf(HistoryDiff)
		if fValue.Pointer() == histValue.Pointer() {
			return f(arg1, arg2, json)
		}
		return specificDiffer(f, arg1, arg2, json)
	}
	return "", errors.New("Unknown differ")
}

func specificDiffer(f func(string, string, bool) (string, error), img1, img2 string, json bool) (string, error) {
	jsonPath1, dirPath1, err := utils.ImageToDir(img1)
	if err != nil {
		return "", err
	}
	jsonPath2, dirPath2, err := utils.ImageToDir(img2)
	if err != nil {
		return "", err
	}
	diff, err := f(jsonPath1, jsonPath2, json)
	if err != nil {
		return "", err
	}

	errStr := remove(dirPath1, true, "")
	errStr = remove(dirPath2, true, errStr)
	errStr = remove(jsonPath1, false, errStr)
	errStr = remove(jsonPath2, false, errStr)

	if errStr != "" {
		return diff, errors.New(errStr)
	}

	return diff, err
}

func remove(path string, dir bool, errStr string) string {
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

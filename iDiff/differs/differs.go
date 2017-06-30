package differs

import (
	"errors"
	"os"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

var diffs = map[string]func(string, string) (string, error){
	"hist": History,
	"dir":  Package,
	"apt":  AptDiff,
	"node": NodeDiff,
}

func Diff(arg1, arg2, differ string) (string, error) {
	if f, exists := diffs[differ]; exists {
		if differ == "hist" {
			return f(arg1, arg2)
		} else {
			return specificDiffer(f, arg1, arg2)
		}

	} else {
		return "", errors.New("Unknown differ.")
	}
}

func specificDiffer(f func(string, string) (string, error), img1, img2 string) (string, error) {
	jsonPath1, dirPath1, err := utils.ImageToDir(img1)
	if err != nil {
		return "", err
	}
	jsonPath2, dirPath2, err := utils.ImageToDir(img2)
	if err != nil {
		return "", err
	}
	diff, err := f(jsonPath1, jsonPath2)
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
	return diff, nil
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

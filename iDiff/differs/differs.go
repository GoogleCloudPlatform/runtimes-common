package differs

import (
	"bytes"
	"errors"
	"os"
	"reflect"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

var diffs = map[string]func(string, string, bool, bool) (string, error){
	"hist":    HistoryDiff,
	"history": HistoryDiff,
	"file":    FileDiff,
	"apt":     AptDiff,
	"linux":   AptDiff,
	"pip":     PipDiff,
	"node":    NodeDiff,
}

func Diff(arg1, arg2, differ string, json bool, eng bool) (string, error) {
	if f, exists := diffs[differ]; exists {
		fValue := reflect.ValueOf(f)
		histValue := reflect.ValueOf(HistoryDiff)
		fileValue := reflect.ValueOf(FileDiff)
		if fValue.Pointer() == histValue.Pointer() || fValue.Pointer() == fileValue.Pointer() {
			return f(arg1, arg2, json, eng)
		}
		return specificDiffer(f, arg1, arg2, json, eng)
	}
	return "", errors.New("Unknown differ")
}

func specificDiffer(f func(string, string, bool, bool) (string, error), img1, img2 string, json bool, eng bool) (string, error) {
	var buffer bytes.Buffer
	validDiff := true
	imgPath1, err := utils.ImageToFS(img1, eng)
	if err != nil {
		buffer.WriteString(err.Error())
		validDiff = false
	}
	imgPath2, err := utils.ImageToFS(img2, eng)
	if err != nil {
		buffer.WriteString(err.Error())
		validDiff = false
	}

	var diff string
	if validDiff {
		output, err := f(imgPath1, imgPath2, json, eng)
		if err != nil {
			buffer.WriteString(err.Error())
		}
		diff = output
	}

	errStr := remove(imgPath1, true)
	errStr += remove(imgPath2, true)
	if errStr != "" {
		buffer.WriteString(errStr)
	}

	if buffer.String() != "" {
		return diff, errors.New(buffer.String())
	}
	return diff, nil
}

func remove(path string, dir bool) string {
	var errStr string
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
		errStr = "\nUnable to remove " + path
	}
	return errStr
}

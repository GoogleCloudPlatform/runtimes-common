package differs

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

var diffs = map[string]func(string, string) (string, error){
	"hist": History,
	"dir":  Package,
	"apt":  AptDiff,
}

func Diff(arg1, arg2, differ, output string) (string, error) {
	if f, exists := diffs[differ]; exists {
		if differ == "hist" {
			return f(arg1, arg2)
		}
		return specificDiffer(f, arg1, arg2, output)
	}
	return "", errors.New("Unknown differ")
}

func specificDiffer(f func(string, string) (string, error), img1, img2, output string) (string, error) {
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

	var returnVal string

	if output != "" {
		f, _ := os.Create(output)
		defer f.Close()
		w := bufio.NewWriter(f)
		err = writeDiff(w, diff)
		if err != nil {
			return "", err
		}
		returnVal = "Image diff successfully written to " + output
		w.Flush()
	} else {
		err = writeDiff(os.Stdout, diff)
	}

	return returnVal, err
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

func writeDiff(w io.Writer, diff string) error {
	diffBytes := []byte(diff)
	var buf bytes.Buffer
	_, err := buf.Write(diffBytes)
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}
	return nil
}

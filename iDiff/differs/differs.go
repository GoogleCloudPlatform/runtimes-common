package differs

import (
	"bytes"
	"errors"
	"os"
	"reflect"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

type ImageDiff struct {
	Image1 Image
	Image2 Image
	DiffType Differ
	UseDocker bool
}

type Differ interface {
	Diff(diff ImageDiff) (string, error)
}

var diffs = map[string]Differ{
	"hist":    HistoryDiffer,
	"history": HistoryDiffer,
	"file":    FileDiffer,
	"apt":     AptDiffer,
	"linux":   AptDiffer,
	"pip":     PipDiffer,
	"node":    NodeDiffer,
}

func (diff ImageDiff) GetDiff() (string, error) {
	img1 := diff.Image1
	img2 := diff.Image2
	differ := diff.DiffType
	eng := diff.UseDocker
	return differ.Diff(image1, image2, true, eng) //TODO: eliminate JSON param and eventually bool
}

func getDiffer(differ string) (differ Differ, err error) {
	if d, exists := diffs[differ]; exists {
		differ = d
	} else {
		errors.New("Unknown differ")
	}
	return
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

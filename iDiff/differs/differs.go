package differs

import (
	"errors"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

type ImageDiff struct {
	Image1 utils.Image
	Image2 utils.Image
	DiffType Differ
	UseDocker bool
}

type Differ interface {
	Diff(image1, image2 utils.Image, json, eng bool) (string, error)
}

var diffs = map[string]Differ{
	"hist":    HistoryDiffer{},
	"history": HistoryDiffer{},
	"file":    FileDiffer{},
	"apt":     AptDiffer{},
	"linux":   AptDiffer{},
	"pip":     PipDiffer{},
	"node":    NodeDiffer{},
}

func (diff ImageDiff) GetDiff() (string, error) {
	img1 := diff.Image1
	img2 := diff.Image2
	differ := diff.DiffType
	eng := diff.UseDocker
	return differ.Diff(img1, img2, true, eng) //TODO: eliminate JSON param and eventually bool
}

func GetDiffer(diffName string) (differ Differ, err error) {
	if d, exists := diffs[diffName]; exists {
		differ = d
	} else {
		errors.New("Unknown differ")
	}
	return
}

/*func specificDiffer(f func(string, string, bool, bool) (string, error), img1, img2 string, json bool, eng bool) (string, error) {
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
}*/



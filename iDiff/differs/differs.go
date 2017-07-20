package differs

import (
	"errors"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

type DiffRequest struct {
	Image1   utils.Image
	Image2   utils.Image
	DiffType Differ
}

type DiffResult interface {
	OutputJSON() error
	OutputText() error
}

type Differ interface {
	Diff(image1, image2 utils.Image) (DiffResult, error)
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

func (diff DiffRequest) GetDiff() (DiffResult, error) {
	img1 := diff.Image1
	img2 := diff.Image2
	differ := diff.DiffType
	return differ.Diff(img1, img2)
}

func GetDiffer(diffName string) (differ Differ, err error) {
	if d, exists := diffs[diffName]; exists {
		differ = d
	} else {
		errors.New("Unknown differ")
	}
	return
}

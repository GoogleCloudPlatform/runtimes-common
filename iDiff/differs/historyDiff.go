package differs

import (
	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

type HistoryDiffer struct {
}

func (d HistoryDiffer) Diff(image1, image2 utils.Image) (utils.DiffResult, error) {
	diff, err := getHistoryDiff(image1, image2)
	return &utils.HistDiffResult{DiffType: "HistoryDiffer", Diff: diff}, err
}

func getHistoryDiff(image1, image2 utils.Image) (utils.HistDiff, error) {
	history1 := image1.History
	history2 := image2.History

	adds := utils.GetAdditions(history1, history2)
	dels := utils.GetDeletions(history1, history2)
	diff := utils.HistDiff{image1.Source, image2.Source, adds, dels}
	return diff, nil
}

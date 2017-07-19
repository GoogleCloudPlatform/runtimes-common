package differs

import (
	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
)

type HistoryDiffer struct {
}

func (d HistoryDiffer) Diff(image1, image2 utils.Image, eng bool) (DiffResult, error) {
	diff, err := getHistoryDiff(image1, image2, eng)
	return &HistDiffResult{Diff: diff}, err
}

type HistDiff struct {
	Image1 string
	Image2 string
	Adds   []string
	Dels   []string
}

type HistDiffResult struct {
	Diff HistDiff
}

func (m *HistDiffResult) OutputJSON() error {
	return utils.JSONify(m.Diff)
}

func (m *HistDiffResult) OutputText() error {
	return utils.TemplateOutput(m.Diff)
}

func getHistoryDiff(image1, image2 utils.Image, eng bool) (HistDiff, error) {
	history1 := image1.History
	history2 := image2.History

	adds := utils.GetAdditions(history1, history2)
	dels := utils.GetDeletions(history1, history2)
	diff := HistDiff{image1.Source, image2.Source, adds, dels}
	return diff, nil
}

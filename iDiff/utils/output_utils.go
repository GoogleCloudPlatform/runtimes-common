package utils

type DiffResult interface {
	OutputJSON(diffType string) error
	OutputText(diffType string) error
}

type MultiVersionPackageDiffResult struct {
	DiffType string
	Diff     MultiVersionPackageDiff
}

func (m *MultiVersionPackageDiffResult) OutputJSON(diffType string) error {
	return JSONify(*m)
}

func (m *MultiVersionPackageDiffResult) OutputText(diffType string) error {
	return TemplateOutput(*m)
}

type PackageDiffResult struct {
	DiffType string
	Diff     PackageDiff
}

func (m *PackageDiffResult) OutputJSON(diffType string) error {
	return JSONify(*m)
}

func (m *PackageDiffResult) OutputText(diffType string) error {
	return TemplateOutput(*m)
}

type HistDiffResult struct {
	DiffType string
	Diff     HistDiff
}

func (m *HistDiffResult) OutputJSON(diffType string) error {
	return JSONify(*m)
}

func (m *HistDiffResult) OutputText(diffType string) error {
	return TemplateOutput(*m)
}

type DirDiffResult struct {
	DiffType string
	Diff     DirDiff
}

func (m *DirDiffResult) OutputJSON(diffType string) error {
	return JSONify(*m)
}

func (m *DirDiffResult) OutputText(diffType string) error {
	return TemplateOutput(*m)
}

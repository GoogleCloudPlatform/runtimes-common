package main

type FileContentTestv0 struct {
	Name             string   // name of test
	Path             string   // file to check existence of
	ExpectedContents []string // list of expected contents of file
	ExcludedContents []string // list of excluded contents of file
}

func (t FileContentTestv0) GetName() string {
	return t.Name
}

func (t FileContentTestv0) GetPath() string {
	return t.Path
}

func (t FileContentTestv0) GetExpectedContents() []string {
	return t.ExpectedContents
}

func (t FileContentTestv0) GetExcludedContents() []string {
	return t.ExcludedContents
}

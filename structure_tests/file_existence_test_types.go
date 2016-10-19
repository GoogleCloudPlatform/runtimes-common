package main

type FileExistenceTest interface {
	GetName() string      // name of test
	GetPath() string      // file to check existence of
	GetIsDirectory() bool // whether or not the path points to a directory
	GetShouldExist() bool // whether or not the file should exist
}

type FileExistenceTestv0 struct {
	Name        string // name of test
	Path        string // file to check existence of
	IsDirectory bool   // whether or not the path points to a directory
	ShouldExist bool   // whether or not the file should exist
}

func (t FileExistenceTestv0) GetName() string {
	return t.Name
}

func (t FileExistenceTestv0) GetPath() string {
	return t.Path
}

func (t FileExistenceTestv0) GetIsDirectory() bool {
	return t.IsDirectory
}

func (t FileExistenceTestv0) GetShouldExist() bool {
	return t.ShouldExist
}

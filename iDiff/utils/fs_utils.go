package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Directory struct {
	Name string
	Files []string
	Dirs []Directory
}

func GetDirectory(dirpath string) Directory {
	dirfile, e := ioutil.ReadFile(dirpath)
	if e != nil {
		panic(e)
		os.Exit(1)
	}

	var dir Directory
	e = json.Unmarshal(dirfile, &dir)
	if e != nil {
		panic(e)
		os.Exit(1)
	}
	return dir
}

// Checks for content differences between files of the same name from different directories
func getModifiedFiles(d1, d2 Directory) []string {
	d1files := d1.Files
	d2files := d2.Files

	filematches := GetMatches(d1files, d2files)

	var modified []string
	for _, f := range filematches {
		f1path := fmt.Sprintf("%s%s", d1.Name, f)
		f2path := fmt.Sprintf("%s%s", d2.Name, f)
		if !checkSameFile(f1path, f2path) {
			modified = append(modified, f)
		}
	}
	return modified
}

func getAddedFiles(d1, d2 Directory) []string {
	return GetAdditions(d1.Files, d2.Files)
}

func getDeletedFiles(d1, d2 Directory) []string {
	return GetDeletions(d1.Files, d2.Files)
}

func compareFileEntries(d1, d2 Directory) ([]string, []string, []string) {
	adds := getAddedFiles(d1, d2)
	dels := getDeletedFiles(d1, d2)
	mods := getModifiedFiles(d1, d2)

	return adds, dels, mods
}

func checkSameFile(f1name, f2name string) bool {
	// Check first if files differ in size and immediately return
	f1stat, err := os.Stat(f1name)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	f2stat, err := os.Stat(f2name)
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	if f1stat.Size() != f2stat.Size() {
		return false
	}

	// Next, check file contents
	f1, err := ioutil.ReadFile(f1name)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	f2, err := ioutil.ReadFile(f2name)
	if err != nil {
		panic(err)
		os.Exit(1)
	}

	if !bytes.Equal(f1, f2) {
		return false
	}
	return true
}

func DiffDirectory(d1, d2 Directory) ([]string, []string, []string) {
	// Diff file entries in the directories
	adds, dels, mods := compareFileEntries(d1, d2)
	return adds, dels, mods

	// TODO: Diff subdirectories within the directories
}


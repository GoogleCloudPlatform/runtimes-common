package differs

import (
	"bytes"
	"fmt"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/golang/glog"
)

// Package diffs two packages and compares their contents
func Package(dir1, dir2 string, json bool) (string, error) {
	diff, err := getDiffOutput(dir1, dir2, json)
	if err != nil {
		return "", err
	}

	return diff, nil
}

func getDiffOutput(d1file, d2file string, json bool) (string, error) {
	d1, err := utils.GetDirectory(d1file)
	if err != nil {
		glog.Errorf("Error reading directory structure from file %s: %s\n", d1file, err)
		return "", err
	}
	d2, err := utils.GetDirectory(d2file)
	if err != nil {
		glog.Errorf("Error reading directory structure from file %s: %s\n", d2file, err)
		return "", err
	}

	d1name := d1.Root
	d2name := d2.Root

	dirDiff := utils.DiffDirectory(d1, d2)

	if json {
		return utils.JSONify(dirDiff)
	}

	var buffer bytes.Buffer

	s := fmt.Sprintf("These entries have been added to %s\n", d1name)
	buffer.WriteString(s)
	if len(dirDiff.Adds) == 0 {
		buffer.WriteString("\tNo files have been added\n")
	} else {
		for _, f := range dirDiff.Adds {
			s = fmt.Sprintf("\t%s\n", f)
			buffer.WriteString(s)
		}
	}

	s = fmt.Sprintf("These entries have been deleted from %s\n", d1name)
	buffer.WriteString(s)
	if len(dirDiff.Dels) == 0 {
		buffer.WriteString("\tNo files have been deleted\n")
	} else {
		for _, f := range dirDiff.Dels {
			s = fmt.Sprintf("\t%s\n", f)
			buffer.WriteString(s)
		}
	}
	s = fmt.Sprintf("These entries have been changed between %s and %s\n", d1name, d2name)
	buffer.WriteString(s)
	if len(dirDiff.Mods) == 0 {
		buffer.WriteString("\tNo files have been modified\n")
	} else {
		for _, f := range dirDiff.Mods {
			s = fmt.Sprintf("\t%s\n", f)
			buffer.WriteString(s)
		}
	}

	return buffer.String(), nil
}

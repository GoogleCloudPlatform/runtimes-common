package differs

import (
	"bytes"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/golang/glog"
)

//  Diffs two packages and compares their contents
func Package(img1, img2 string) (string, error) {
	dir1, path1, err := utils.ImageToDir(img1)
	if err != nil {
		return "", err
	}
	dir2, path2, err := utils.ImageToDir(img2)
	if err != nil {
		return "", err
	}
	diff, err := getDiffOutput(dir1, dir2)

	defer os.RemoveAll(path1)
	defer os.RemoveAll(path2)
	defer os.Remove(dir1)
	defer os.Remove(dir2)

	if err != nil {
		return "", err
	}

	return diff, nil
}

func getDiffOutput(d1file, d2file string) (string, error) {
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

	adds, dels, mods := utils.DiffDirectory(d1, d2)

	var buffer bytes.Buffer

	s := fmt.Sprintf("These entries have been added to %s\n", d1name)
	buffer.WriteString(s)
	if len(adds) == 0 {
		buffer.WriteString("\tNo files have been added\n")
	} else {
		for _, f := range adds {
			s = fmt.Sprintf("\t%s\n", f)
			buffer.WriteString(s)
		}
	}

	s = fmt.Sprintf("These entries have been deleted from %s\n", d1name)
	buffer.WriteString(s)
	if len(dels) == 0 {
		buffer.WriteString("\tNo files have been deleted\n")
	} else {
		for _, f := range dels {
			s = fmt.Sprintf("\t%s\n", f)
			buffer.WriteString(s)
		}
	}
	s = fmt.Sprintf("These entries have been changed between %s and %s\n", d1name, d2name)
	buffer.WriteString(s)
	if len(mods) == 0 {
		buffer.WriteString("\tNo files have been modified\n")
	} else {
		for _, f := range mods {
			s = fmt.Sprintf("\t%s\n", f)
			buffer.WriteString(s)
		}
	}

	return buffer.String(), nil
}

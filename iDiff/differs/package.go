package differs

import (
	"bytes"
	"fmt"
	"os"
	"testing/runtimes-common/iDiff/utils"
)

//  Diffs two packages to see if they have the same contents
func Package(d1file, d2file string) string {
	d1, err := utils.GetDirectory(d1file)
	if err != nil {
		fmt.Errorf("Error reading directory structure from file %s: %s", d1file, err)
		os.Exit(1)
	}
	d2, err := utils.GetDirectory(d2file)
	if err != nil {
		fmt.Errorf("Error reading directory structure from file %s: %s", d2file, err)
		os.Exit(1)
	}

	d1name := d1.Name
	d2name := d2.Name

	adds, dels, mods := utils.DiffDirectory(d1, d2)

	var buffer bytes.Buffer
	if adds == nil {
		buffer.WriteString("No files to diff\n")
	} else {
		s := fmt.Sprintf("These files have been added to %s\n", d1name)
		buffer.WriteString(s)
		if len(adds) == 0 {
			buffer.WriteString("none\n")
		} else {
			for _, f := range adds {
				s = fmt.Sprintf("%s\n", f)
				buffer.WriteString(s)
			}
		}

		s = fmt.Sprintf("These files have been deleted from %s\n", d1name)
		buffer.WriteString(s)
		if len(dels) == 0 {
			buffer.WriteString("none\n")
		} else {
			for _, f := range dels {
				s = fmt.Sprintf("%s\n", f)
				buffer.WriteString(s)
			}
		}
		s = fmt.Sprintf("These files have been changed between %s and %s\n", d1name, d2name)
		buffer.WriteString(s)
		if len(mods) == 0 {
			buffer.WriteString("none\n")
		} else {
			for _, f := range mods {
				s = fmt.Sprintf("%s\n", f)
				buffer.WriteString(s)
			}
		}
	}
	return buffer.String()
}

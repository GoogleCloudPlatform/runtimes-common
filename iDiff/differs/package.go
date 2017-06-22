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
		fmt.Printf("Error reading directory structure from file %s: %s\n", d1file, err)
		os.Exit(1)
	}
	d2, err := utils.GetDirectory(d2file)
	if err != nil {
		fmt.Printf("Error reading directory structure from file %s: %s\n", d2file, err)
		os.Exit(1)
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

	return buffer.String()
}

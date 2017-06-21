package differs

import (
	"bytes"
	"fmt"
	"testing/runtimes-common/iDiff/utils"
)

// History compares the Docker history for each image.
func Package(d1file, d2file string) string {
	d1 := utils.GetDirectory(d1file)
	d2 := utils.GetDirectory(d2file)

	d1name := d1.Name
	d2name := d2.Name

	adds, dels, mods := utils.DiffDirectory(d1, d2)

	var buffer bytes.Buffer
	s := fmt.Sprintf("These files have been deleted from %s\n", d1name)
	buffer.WriteString(s)
	if adds == nil {
		buffer.WriteString("none\n")
	}else {
		for _, f := range adds {
			s = fmt.Sprintf("%s\n", f)
			buffer.WriteString(s)
		}
	}

	s = fmt.Sprintf("These files have been deleted from %s\n", d1name)
	buffer.WriteString(s)
	if dels == nil {
		buffer.WriteString("none\n")
	}else {
		for _, f := range dels {
			s = fmt.Sprintf("%s\n", f)
			buffer.WriteString(s)
		}
	}
	s = fmt.Sprintf("These files have been changed between %s and %s\n", d1name, d2name)
	buffer.WriteString(s)
	if mods == nil {
		buffer.WriteString("none\n")
	}else {
		for _, f := range mods {
			s = fmt.Sprintf("%s\n", f)
			buffer.WriteString(s)
		}
	}
	return buffer.String()
}

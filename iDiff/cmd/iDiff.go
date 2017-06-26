package cmd

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/runtimes-common/iDiff/differs"
	"github.com/spf13/cobra"
)

// iDiff represents the iDiff command
var iDiffCmd = &cobra.Command{
	Use:   "iDiff [container1] [container2] [differ]",
	Short: "Compare two images.",
	Long:  `Compares two images using the specifed differ. `,
	Run: func(cmd *cobra.Command, args []string) {
		if valid, err := checkArgNum(args); !valid {
			glog.Fatalf(err.Error())
		}
		if diff, err := differs.Diff(args[0], args[1], args[2]); err == nil {
			fmt.Println(diff)
		} else {
			glog.Fatalf(err.Error())
		}
	},
}

func checkArgNum(args []string) (bool, error) {
	var err_message string
	if len(args) < 3 {
		err_message = "Please have two image IDs and one differ as arguments."
		return false, errors.New(err_message)
	} else if len(args) > 3 {
		err_message = "Too many arguments."
		return false, errors.New(err_message)
	} else {
		return true, nil
	}
}

func init() {
	RootCmd.AddCommand(iDiffCmd)
}

package cmd

import (
	"errors"
	"fmt"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/differs"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var json bool

// iDiff represents the iDiff command
var iDiffCmd = &cobra.Command{
	Use:   "iDiff [container1] [container2] [differ]",
	Short: "Compare two images.",
	Long:  `Compares two images using the specifed differ. `,
	Run: func(cmd *cobra.Command, args []string) {
		if validArgs, err := validateArgs(args); !validArgs {
			glog.Fatalf(err.Error())
		}
		if diff, err := differs.Diff(args[0], args[1], args[2], json); err == nil {
			fmt.Println(diff)
		} else {
			glog.Fatalf(err.Error())
		}
	},
}

func validateArgs(args []string) (bool, error) {
	if validArgNum, err := checkArgNum(args); !validArgNum {
		return false, err
	}
	return true, nil
}

func checkArgNum(args []string) (bool, error) {
	var errMessage string
	if len(args) < 3 {
		errMessage = "Too few arguments. Should have three: [IMAGE] [IMAGE] [DIFFER]."
		return false, errors.New(errMessage)
	} else if len(args) > 3 {
		errMessage = "Too many arguments. Should have three: [IMAGE] [IMAGE] [DIFFER]."
		return false, errors.New(errMessage)
	} else {
		return true, nil
	}
}

func init() {
	RootCmd.AddCommand(iDiffCmd)
	iDiffCmd.Flags().BoolVarP(&json, "json", "j", false, "JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).")
}

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

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
	if validArgType, err := checkArgType(args); !validArgType {
		return false, err
	}
	return true, nil
}

func checkArgNum(args []string) (bool, error) {
	var errMessage string
	if len(args) < 3 {
		errMessage = "Too few arguments. Should have three: [IMAGE ID] [IMAGE ID] [DIFFER]."
		return false, errors.New(errMessage)
	} else if len(args) > 3 {
		errMessage = "Too many arguments. Should have three: [IMAGE ID] [IMAGE ID] [DIFFER]."
		return false, errors.New(errMessage)
	} else {
		return true, nil
	}
}

func checkArgType(args []string) (bool, error) {
	var buffer bytes.Buffer
	valid := true
	if !checkImageID(args[0]) {
		valid = false
		errMessage := fmt.Sprintf("Argument %s is not an image ID\n", args[0])
		buffer.WriteString(errMessage)
	}
	if !checkImageID(args[1]) {
		valid = false
		errMessage := fmt.Sprintf("Argument %s is not an image ID\n", args[1])
		buffer.WriteString(errMessage)
	}
	if checkImageID(args[2]) {
		valid = false
		buffer.WriteString("Do not provide more than two image IDs\n")
	} else if !checkDiffer(args[2]) {
		valid = false
		buffer.WriteString("Please provide a differ name as the third argument")
	}
	if !valid {
		return false, errors.New(buffer.String())
	}
	return true, nil
}

func checkImageID(arg string) bool {
	pattern := regexp.MustCompile("[a-z|0-9]{12}")
	if exp := pattern.FindString(arg); exp != arg {
		return false
	}
	return true
}

func checkDiffer(arg string) bool {
	pattern := regexp.MustCompile("[a-z|A-Z]*")
	if exp := pattern.FindString(arg); exp != arg {
		return false
	}
	return true
}

func init() {
	RootCmd.AddCommand(iDiffCmd)
	iDiffCmd.Flags().BoolVarP(&json, "JSON Output", "j", false, "JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).")
}

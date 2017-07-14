package cmd

import (
	"bytes"
	"errors"
	goflag "flag"
	"fmt"
	"os"
	"regexp"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/differs"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var json bool

var RootCmd = &cobra.Command{
	Use:   "iDiff [differ] [container1] [container2]",
	Short: "Compare two images.",
	Long:  `Compares two images using the specifed differ (hist, dir, or apt).`,
	Run: func(cmd *cobra.Command, args []string) {
		if validArgs, err := validateArgs(args); !validArgs {
			glog.Error(err.Error())
			os.Exit(1)
		}
		if diff, err := differs.Diff(args[1], args[2], args[0], json); err == nil {
			fmt.Println(diff)
		} else {
			glog.Error(err.Error())
		}
	},
}

func validateArgs(args []string) (bool, error) {
	validArgNum, err := checkArgNum(args)
	if err != nil {
		return false, err
	} else if !validArgNum {
		return false, nil
	}
	validArgType, err := checkArgType(args)
	if err != nil {
		return false, err
	} else if !validArgType {
		return false, nil
	}
	return true, nil
}

func checkArgNum(args []string) (bool, error) {
	var errMessage string
	if len(args) < 3 {
		errMessage = "Too few arguments. Should have three: [DIFFER] [IMAGE ID] [IMAGE ID]."
		return false, errors.New(errMessage)
	} else if len(args) > 3 {
		errMessage = "Too many arguments. Should have three: [DIFFER] [IMAGE ID] [IMAGE ID]."
		return false, errors.New(errMessage)
	} else {
		return true, nil
	}
}

func checkArgType(args []string) (bool, error) {
	var buffer bytes.Buffer
	valid := true
	if !checkDiffer(args[0]) {
		valid = false
		buffer.WriteString("Please provide a differ name as the third argument (hist, dir, or apt)\n")
	}
	if !checkImageID(args[1]) {
		valid = false
		errMessage := fmt.Sprintf("Argument %s is not an image ID\n", args[1])
		buffer.WriteString(errMessage)
	}
	if !checkImageID(args[2]) {
		valid = false
		errMessage := fmt.Sprintf("Argument %s is not an image ID\n", args[2])
		buffer.WriteString(errMessage)
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
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	RootCmd.Flags().BoolVarP(&json, "json", "j", false, "JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).")
}

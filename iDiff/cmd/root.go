package cmd

import (
	"bytes"
	"errors"
	goflag "flag"
	"fmt"
	"os"
	"regexp"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/differs"
	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/utils"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var json bool
var eng bool

var RootCmd = &cobra.Command{
	Use:   "[differ] [image1] [image2]",
	Short: "Compare two images.",
	Long:  `Compares two images using the specifed differ (see iDiff documentation for available differs).`,
	Run: func(cmd *cobra.Command, args []string) {
		if validArgs, err := validateArgs(args[1:]); !validArgs {
			glog.Error(err.Error())
			os.Exit(1)
		}
		image1, err := utils.ImagePrepper{args[2], eng}.GetImage()
		if err != nil {
			glog.Error(err.Error())
			os.Exit(1)
		}
		image2, err := utils.ImagePrepper{args[3], eng}.GetImage()
		if err != nil {
			glog.Error(err.Error())
			os.Exit(1)
		}
		differ, err := differs.GetDiffer(args[1])
		if err != nil {
			glog.Error(err.Error())
			os.Exit(1)
		}

		diff := differs.DiffRequest{image1, image2, differ, eng}
		if diff, err := diff.GetDiff(); err == nil {
			if json {
				err = diff.OutputJSON()
				if err != nil {
					glog.Error(err)
				}
			} else {
				err = diff.OutputText()
				if err != nil {
					glog.Error(err)
				}
			}

			errMsg := remove(image1.FSPath, true)
			errMsg += remove(image2.FSPath, true)
			if errMsg != "" {
				glog.Error(errMsg)
			}
		} else {
			glog.Error(err.Error())
			os.Exit(1)
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
		errMessage = "Too few arguments. Should have three: [DIFFER] [IMAGE] [IMAGE]."
		return false, errors.New(errMessage)
	} else if len(args) > 3 {
		errMessage = "Too many arguments. Should have three: [DIFFER] [IMAGE] [IMAGE]."
		return false, errors.New(errMessage)
	} else {
		return true, nil
	}
}

func checkImage(arg string) bool {
	if !utils.CheckImageID(arg) && !utils.CheckImageURL(arg) && !utils.CheckTar(arg) {
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

func checkArgType(args []string) (bool, error) {
	var buffer bytes.Buffer
	valid := true
	if !checkDiffer(args[0]) {
		valid = false
		buffer.WriteString("Please provide a differ name as the first argument")
	}
	if !checkImage(args[1]) {
		valid = false
		errMessage := fmt.Sprintf("Argument %s is not an image ID, URL, or tar\n", args[1])
		buffer.WriteString(errMessage)
	}
	if !checkImage(args[2]) {
		valid = false
		errMessage := fmt.Sprintf("Argument %s is not an image ID, URL, or tar\n", args[2])
		buffer.WriteString(errMessage)
	}
	if !valid {
		return false, errors.New(buffer.String())
	}
	return true, nil
}

func remove(path string, dir bool) string {
	var errStr string
	if path == "" {
		return ""
	}

	var err error
	if dir {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		errStr = "\nUnable to remove " + path
	}
	return errStr
}

func init() {
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	RootCmd.Flags().BoolVarP(&json, "json", "j", false, "JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).")
	RootCmd.Flags().BoolVarP(&eng, "eng", "e", false, "By default the docker calls are shelled out locally, set this flag to use the Docker Engine Client (version compatibility required).")
}

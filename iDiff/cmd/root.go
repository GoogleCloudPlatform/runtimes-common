package cmd

import (
	goflag "flag"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var RootCmd = &cobra.Command{
	Use:   "iDiff",
	Short: "iDiff is an image differ tool.",
	Long:  `iDiff is an image differ tool.`,
	Run: func(command *cobra.Command, args []string) {
		glog.Info("Root command started")
	},
}

func init() {
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

package cmd

import (
	goflag "flag"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var RootCmd = &cobra.Command{
	Use:   "iDiff",
	Short: "iDiff is an image differ tool.",
	Long:  `iDiff is an image differ tool.`,
	Run: func(command *cobra.Command, args []string) {
		fmt.Println("Root command started")
	},
}

func init() {
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

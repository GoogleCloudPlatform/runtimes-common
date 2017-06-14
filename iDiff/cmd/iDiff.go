package cmd

import (
	"fmt"
	"testing/runtimes-common/iDiff/differs"

	"github.com/spf13/cobra"
)

var container1, container2, differ string

// iDiff represents the iDiff command
var iDiffCmd = &cobra.Command{
	Use:   "iDiff [container1] [container2] [differ]",
	Short: "Compare two images.",
	Long:  `Compares two images using the specifed differ. `,
	Run: func(cmd *cobra.Command, args []string) {
		if args[2] == "hist" {
			diff := differs.History(args[0], args[1])
			fmt.Println(diff)
		} else {
			fmt.Println("Unknown differ")
		}
	},
}

func init() {
	RootCmd.AddCommand(iDiffCmd)
}

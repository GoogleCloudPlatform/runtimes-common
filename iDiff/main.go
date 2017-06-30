package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/runtimes-common/iDiff/cmd"
)

func main() {
	flag.Parse()
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

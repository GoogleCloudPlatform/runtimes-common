package main

import (
	"fmt"
	"os"

	"runtimes-common/iDiff/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

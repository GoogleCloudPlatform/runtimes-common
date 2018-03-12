/*
Copyright 2018 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ctc_lib

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

type CLIInterface interface {
	PrintO(c *cobra.Command, args []string)
	SetRun(func(c *cobra.Command, args []string))
	GetCommand() *cobra.Command
	ValidateCommand() error
	IsRunODefined() bool
	Init()
}

func Execute(ctb CLIInterface) (err error) {
	defer errRecover(&err)

	if err := ctb.ValidateCommand(); err != nil {
		CommandExit(err)
	}
	ctb.Init()
	if ctb.IsRunODefined() {
		cobraRun := func(c *cobra.Command, args []string) {
			ctb.PrintO(c, args)
		}
		ctb.SetRun(cobraRun)
	}

	err = ctb.GetCommand().Execute()

	//Add empty line as template.Execute does not print an empty line
	ctb.GetCommand().Println()
	CommandExit(err)
	return err
}

// errRecover is the handler that turns panics into returns from the top
// level of Parse.
func errRecover(errp *error) {
	if e := recover(); e != nil {
		// TODO: Change this to Log.Error once Logging is introduced.
		fmt.Println(e)
		*errp = errors.New(fmt.Sprintf("%v", e))
		CommandExit(*errp)
	}
}

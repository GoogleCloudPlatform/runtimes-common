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
	"os"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/flags"
	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/util"
	"github.com/spf13/cobra"
)

type ContainerToolCommand struct {
	*ContainerToolCommandBase
	Output interface{}
	RunO   func(command *cobra.Command, args []string) (interface{}, error)
}

type ContainerToolListCommand struct {
	*ContainerToolCommandBase
	OutputList []interface{}
	RunO       func(command *cobra.Command, args []string) ([]interface{}, error)
}

func (ctc *ContainerToolCommand) Execute() (err error) {
	defer errRecover(&err)
	isRunODefined := ctc.RunO != nil

	if err := ctc.CheckValidCommand(isRunODefined); err != nil {
		panic(err.Error())
	}
	ctc.init()
	if isRunODefined {
		cobraRun := func(c *cobra.Command, args []string) {
			obj, _ := ctc.RunO(c, args)
			ctc.Output = obj
			util.ExecuteTemplate(flags.TemplateString, obj, ctc.OutOrStdout())
		}
		ctc.Command.Run = cobraRun
	}

	err = ctc.Command.Execute()
	//Add empty line.
	fmt.Println()

	if err != nil {
		panic(err)
	}
	return err
}

// errRecover is the handler that turns panics into returns from the top
// level of Parse.
func errRecover(errp *error) {
	if e := recover(); e != nil {
		// TODO: Change this to Log.Error once Logging is introduced.
		fmt.Println(e)
		if !flags.NoExit {
			os.Exit(1)
		}
		*errp = errors.New(fmt.Sprintf("%v", e))
	}
}

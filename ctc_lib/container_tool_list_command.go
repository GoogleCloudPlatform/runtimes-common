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

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/flags"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/constants"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/util"
	"github.com/spf13/cobra"
)

type ListCommandOutputObject struct {
	OutputList    []interface{}
	SummaryObject interface{}
}

type ContainerToolListCommand struct {
	*ContainerToolCommandBase
	OutputList []interface{}
	// RunO Executes cobra.Command.Run and returns an List[Output]
	RunO func(command *cobra.Command, args []string) ([]interface{}, error)
	// When defined, StreamO Executes cobra.Command.Run and streams each item in the List as its added.
	// This will ignore the RunO function.
	StreamO func(command *cobra.Command, args []string)
	// This function will execute over the output list and return a Summary Object which can be printed.
	// The SummaryTemplate Field can be used to print the Object.
	TotalO          func(list []interface{}) (interface{}, error)
	SummaryObject   interface{}
	SummaryTemplate string
	// Defines the stream to write objects to when using StreamO
	Stream chan interface{}
}

func (commandList ContainerToolListCommand) ReadFromStream(streamOutput bool) ([]interface{}, error) {
	results := make([]interface{}, 0)
	for obj := range commandList.Stream {
		if _, ok := obj.(string); ok {
			// Display any Arbitary strings written to the channel as is.
			// These could be headers or any text.
			// TODO: Provide a callback for users to overwrite this default behavior.
			if streamOutput {
				util.ExecuteTemplate(constants.EmptyTemplate,
					obj, commandList.TemplateFuncMap, commandList.OutOrStdout())
				continue
			}
		}
		results = append(results, obj)
		if streamOutput {
			err := util.ExecuteTemplate(commandList.ReadTemplateFromFlagOrCmdDefault(),
				obj, commandList.TemplateFuncMap, commandList.OutOrStdout())
			if err != nil {
				return nil, err
			}
		}
	}
	return results, nil
}

func (commandList ContainerToolListCommand) isRunODefined() bool {
	return commandList.RunO != nil || commandList.StreamO != nil
}

func (ctc *ContainerToolListCommand) ValidateCommand() error {
	if (ctc.Run != nil || ctc.RunE != nil) && ctc.isRunODefined() {
		return errors.New(`Cannot provide both Command.Run and RunO implementation.
Either implement Command.Run implementation or RunO implemetation`)
	}
	return nil
}

func (ctc *ContainerToolListCommand) printO(c *cobra.Command, args []string) error {
	var commandError error = nil
	if ctc.StreamO != nil {
		// Stream Objects only when outputJson = False
		ctc.StreamO(c, args)
		ctc.OutputList, commandError = ctc.ReadFromStream(!flags.OutputJson)
	} else {
		// Run RunO function.
		ctc.OutputList, commandError = ctc.RunO(c, args)
	}
	// If TotalO function defined and Summary Template provided, print the summary.
	if ctc.TotalO != nil && ctc.SummaryTemplate != "" {
		ctc.SummaryObject, commandError = ctc.TotalO(ctc.OutputList)
	}
	if commandError != nil {
		Log.Errorf("%v", commandError)
	}
	displayError := ctc.printResult()
	if displayError != nil {
		Log.Errorf("%v", commandError)
	}
	if commandError != nil && displayError != nil {
		return errors.New("One or more errors")
	}
	return nil
}

func (ctc *ContainerToolListCommand) printResult() error {
	if flags.OutputJson {
		data := ListCommandOutputObject{
			OutputList:    ctc.OutputList,
			SummaryObject: ctc.SummaryObject,
		}
		return util.ExecuteTemplate("", data, nil, ctc.OutOrStdout())
	}
	var err error
	// Do not display the object list again.
	if ctc.StreamO != nil {
		err = util.ExecuteTemplate(ctc.ReadTemplateFromFlagOrCmdDefault(),
			ctc.OutputList, ctc.TemplateFuncMap, ctc.OutOrStdout())
		if err != nil {
			Log.Errorf("%v", err)
		}
	}
	totalErr := util.ExecuteTemplate(ctc.SummaryTemplate,
		ctc.SummaryObject, ctc.TemplateFuncMap, ctc.OutOrStdout())
	if totalErr != nil && err != nil {
		return nil
	} else if totalErr != nil {
		return totalErr
	} else {
		return err
	}
}

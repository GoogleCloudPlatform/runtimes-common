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

	"github.com/spf13/cobra"
)

type ContainerToolListCommand struct {
	*ContainerToolCommandBase
	OutputList []interface{}
	RunO       func(command *cobra.Command, args []string) ([]interface{}, error)
}

func (commandList ContainerToolListCommand) isRunODefined() bool {
	return commandList.RunO != nil
}

func (ctc *ContainerToolListCommand) ValidateCommand() error {
	if (ctc.Run != nil || ctc.RunE != nil) && ctc.isRunODefined() {
		return errors.New("Cannot provide both Command.Run and RunO implementation." +
			"\nEither implement Command.Run implementation or RunO implemetation")
	}
	return nil
}

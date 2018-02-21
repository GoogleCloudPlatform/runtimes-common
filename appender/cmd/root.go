/*
Copyright 2018 Google LLC
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

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/GoogleCloudPlatform/runtimes-common/appender/layer"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "appender",
	Short: "Creates a new docker image by appending a tarball to an existing one.",
	Long: `Creates a new docker image by appending a tarball to an existing one.
	The base and created images may be in different repositories, and may require different permissions.
	The tarball may be compressed on uncompressed, and is specified as a path to a local file.`,
	Run: func(cmd *cobra.Command, args []string) {
		l, err := ioutil.ReadFile(tarballPath)
		if err != nil {
			logrus.Fatalf("error reading tarball-path: %s", err)
		}

		if err := layer.AppendLayer(baseImage, image, l); err != nil {
			logrus.Fatalf("error appending layer: %s", err)
		}
	},
}

var (
	tarballPath string
	baseImage   string
	image       string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.appender.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&tarballPath, "tarball-path", "t", "", "Path to .tar or .tar.gz.")
	rootCmd.Flags().StringVarP(&baseImage, "base-image", "b", "", "Fully qualified name of base image.")
	rootCmd.Flags().StringVarP(&image, "image", "i", "", "Fully qualified name of image to create.")
}

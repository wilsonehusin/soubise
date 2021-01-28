/*
Copyright © 2021 Wilson Husin <wilsonehusin@gmail.com>

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
	"bytes"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"

	"github.com/wilsonehusin/soubise/internal/client"
	"github.com/wilsonehusin/soubise/internal/printer"
)

const getCmdName = "get"

type getOptions struct {
	RefPath string
}

var getOpts = &getOptions{}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   getCmdName,
	Short: "get",
	Long:  `get a shared file`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := envconfig.Process(progName+"_"+getCmdName, getOpts); err != nil {
			return err
		}

		if getOpts.RefPath == "" {
			return fmt.Errorf("no reference path specified -- what are you trying to get?")
		}

		return nil
	},
	Run: func(*cobra.Command, []string) {
		if err := client.Get(getOpts.RefPath); err != nil {
			printer.Stderr("unable to get content: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	var optionsUsage bytes.Buffer
	if err := envconfig.Usagef(progName+"_"+getCmdName, getOpts, &optionsUsage, optionsUsageTemplate); err != nil {
		panic(err)
	}
	getCmd.SetUsageTemplate(getCmd.UsageTemplate() + optionsUsageHeader + optionsUsage.String() + rootCmdOptionsUsage())

	getCmd.Flags().StringVarP(&getOpts.RefPath, "path", "p", getOpts.RefPath, "reference path to retrieve from, prefixed with soubise://")

	rootCmd.AddCommand(getCmd)
}

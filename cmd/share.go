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
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"

	"github.com/wilsonehusin/soubise/internal/buildinfo"
	"github.com/wilsonehusin/soubise/internal/client"
	"github.com/wilsonehusin/soubise/internal/printer"
)

const shareCmdName = "share"

type shareOptions struct {
	FilePath string
	Lifetime string `default:"24h"`
	Server   string
	//Auth     string
}

var shareOpts = &shareOptions{}

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   shareCmdName,
	Short: "Shares specified file",
	Long: `Sharing end-to-end encrypted file

Soubise generates encryption key from client side and uploads 
encrypted data, therefore Soubise servers never knew about
the encryption keys.

This means that your sharing link contains the encryption key,
be careful with whom and how you share the key!

The same flags available in command line interface are also
configurable through environment variable, providing flexibility
in using Soubise programmatically.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return sharePreCheck()
	},
	Run: func(cmd *cobra.Command, args []string) {
		duration, err := time.ParseDuration(shareOpts.Lifetime)
		if err != nil {
			printer.Stderr("unable to understand provided lifetime: %v\n", err)
			os.Exit(1)
		}
		if err := client.Share(shareOpts.FilePath, duration, shareOpts.Server); err != nil {
			printer.Stderr("unable to share: %v\n", err)
			os.Exit(1)
		}
	},
}

func sharePreCheck() error {
	if err := envconfig.Process(progName+"_"+serverCmdName, serverOpts); err != nil {
		return err
	}

	if shareOpts.FilePath == "" {
		return fmt.Errorf("no filepath specified -- what are you trying to share?")
	}

	if shareOpts.Server == "" {
		shareOpts.Server = buildinfo.Server
	}

	return nil
}

func init() {
	if err := envconfig.Process(progName+"_"+shareCmdName, shareOpts); err != nil {
		panic(err)
	}
	var optionsUsage bytes.Buffer
	if err := envconfig.Usagef(progName+"_"+shareCmdName, shareOpts, &optionsUsage, optionsUsageTemplate); err != nil {
		panic(err)
	}
	shareCmd.SetUsageTemplate(shareCmd.UsageTemplate() + optionsUsageHeader + optionsUsage.String() + rootCmdOptionsUsage())

	shareCmd.Flags().StringVarP(&shareOpts.FilePath, "file", "f", shareOpts.FilePath, "path to file to be shared")
	shareCmd.Flags().StringVarP(&shareOpts.Server, "server", "s", shareOpts.Server, "target server address")
	shareCmd.Flags().StringVarP(&shareOpts.Lifetime, "lifetime", "l", shareOpts.Lifetime, "the lifetime for file to be downloadable")

	rootCmd.AddCommand(shareCmd)
}

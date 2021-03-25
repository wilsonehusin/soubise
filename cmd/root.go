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
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/wilsonehusin/soubise/internal/spinner"
)

type rootOptions struct {
	Trace   bool   `default:"false"` // TODO: implement network call trace logs
	Debug   bool   `default:"false"`
	Json    bool   `default:"false"`
	logPath string //nolint:structcheck,unused // TODO: implement multi-output logger
}

const (
	rootDesc             = "Soubise makes file sharing easy"
	progName             = "soubise"
	optionsUsageHeader   = "\nConfigurable environment variables (with default values):"
	optionsUsageTemplate = `{{range .}}
  {{if usage_required .}}(required) {{else}}           {{end}}{{usage_key .}}={{usage_default .}}{{end}}`
)

var rootOpts = &rootOptions{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               progName,
	Short:             rootDesc,
	Long:              rootDesc,
	PersistentPreRunE: rootCmdInit,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	rootCmd.SetUsageTemplate(rootCmd.UsageTemplate() + optionsUsageHeader + rootCmdOptionsUsage())
}

func rootCmdOptionsUsage() string {
	var optionsUsage bytes.Buffer
	if err := envconfig.Usagef(progName, rootOpts, &optionsUsage, optionsUsageTemplate); err != nil {
		panic(err)
	}
	return optionsUsage.String()
}

func rootCmdInit(cmd *cobra.Command, args []string) error {
	if err := envconfig.Process(progName, rootOpts); err != nil {
		return err
	}

	if rootOpts.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.With().Caller().Logger()
		zerolog.CallerMarshalFunc = func(file string, line int) string {
			return path.Base(file) + ":" + strconv.Itoa(line)
		}

	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if rootOpts.Json {
		spinner.Disable()
	} else {
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Debug().Bool("Debug", rootOpts.Debug).Send()

	// TODO: handle LogPath for dual log output to support users watching from tty
	// and collect their usage / execution metrics at the same time

	log.Debug().Str("Subcommand", cmd.Use).Send()
	return nil
}

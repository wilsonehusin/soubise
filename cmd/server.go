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
	"os"
	"os/signal"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/wilsonehusin/soubise/internal/resolve"
	"github.com/wilsonehusin/soubise/internal/server"
	"github.com/wilsonehusin/soubise/internal/server/router"
	"github.com/wilsonehusin/soubise/pkg/storage"
)

const serverCmdName = "server"

type serverOptions struct {
	Host        string `default:"pub.soubise.org"`
	Port        int    `default:"8080"`
	StoragePath string `default:"inmemory"`
	BrokerPath  string
}

var serverOpts = &serverOptions{}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   serverCmdName,
	Short: "Starts target server",
	Long:  `Starts target server`, // TODO: explain parameters
	PreRunE: func(*cobra.Command, []string) error {
		return serverPreCheck()
	},
	Run: func(cmd *cobra.Command, args []string) {
		serverRun()
	},
}

func serverPreCheck() error {
	if err := envconfig.Process(progName+"_"+serverCmdName, serverOpts); err != nil {
		return err
	}

	log.Info().
		Str("Address", serverOpts.Host).
		Int("Port", serverOpts.Port).
		Msg("booting server")
	return nil
}

func init() {
	var optionsUsage bytes.Buffer
	if err := envconfig.Usagef(progName+"_"+serverCmdName, serverOpts, &optionsUsage, optionsUsageTemplate); err != nil {
		panic(err)
	}
	serverCmd.SetUsageTemplate(serverCmd.UsageTemplate() + optionsUsageHeader + optionsUsage.String() + rootCmdOptionsUsage())

	rootCmd.AddCommand(serverCmd)
}

func serverRun() {
	userStop := make(chan os.Signal)
	signal.Notify(userStop, os.Interrupt)

	serverWaiter := make(chan bool)

	brokerProvider := resolve.NewBrokerFromPath(serverOpts.BrokerPath)
	storageProvider := resolve.NewStorageFromPath(serverOpts.StoragePath, brokerProvider)
	if err := storage.SetStorage(storageProvider); err != nil {
		log.Fatal().Err(err).Msg("storage initialization")
	}

	mux := router.NewMux()
	webserver := server.HttpServer{
		Router: mux,
		Host:   serverOpts.Host,
		Port:   serverOpts.Port,
	}

	if err := webserver.Start(); err != nil {
		log.Fatal().Err(err).Msg("failed to launch http server")
	}
	log.Info().
		Int("Port", serverOpts.Port).
		Str("Storage", storageProvider.Kind()).
		Msg("server running")

	go func() {
		<-userStop
		log.Info().Msg("received shutdown request, allowing maximum 10s for server to shutdown gracefully")
		if err := webserver.Stop(); err != nil {
			log.Error().Err(err).Msg("shutdown attempt")
		}
		log.Info().Msg("server has stopped")
		serverWaiter <- true
	}()

	<-serverWaiter
}

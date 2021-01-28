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

package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type HttpServer struct {
	Host   string
	Port   int
	Router http.Handler
	server *http.Server
}

func (h *HttpServer) PreCheck() error {
	errs := []string{}

	if h.Host == "" {
		errs = append(errs, "Host cannot be empty")
	}
	if h.Port == 0 {
		errs = append(errs, "Port cannot be empty")
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed HttpServer PreCheck: %v", strings.Join(errs, ", "))
	}
	return nil
}

func (h *HttpServer) Start() error {
	if err := h.PreCheck(); err != nil {
		return err
	}

	h.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", h.Port),
		Handler: h.Router,
	}

	go func() {
		if err := h.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("http server aborted")
		}
	}()

	return nil
}

func (h *HttpServer) Stop() error {
	ctx, forceStopServer := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer forceStopServer()
	if err := h.server.Shutdown(ctx); err != nil {
		return err
	}
	h.server = nil
	return nil
}

func (h *HttpServer) IsActive() bool {
	return h.server != nil
}

func (h *HttpServer) SetAddr(host string, port int) error {
	h.Host = host
	h.Port = port
	return nil
}

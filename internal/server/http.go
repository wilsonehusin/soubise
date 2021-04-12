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
	"container/heap"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/wilsonehusin/soubise/internal/storage"
)

type Config struct {
	Host           string
	Port           int
	PreCheckExpiry bool
	ActiveExpiry   bool
	TickExpiry     time.Duration
}

type HttpServer struct {
	Config     Config
	Router     http.Handler
	ctx        context.Context
	cancelFunc context.CancelFunc
	server     *http.Server
}

func (h *HttpServer) PreCheck() error {
	errs := []string{}

	if h.Config.Host == "" {
		errs = append(errs, "Host cannot be empty")
	}
	if h.Config.Port == 0 {
		errs = append(errs, "Port cannot be empty")
	}
	if h.cancelFunc == nil {
		h.ctx, h.cancelFunc = context.WithCancel(context.Background())
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
		Addr:    fmt.Sprintf(":%d", h.Config.Port),
		Handler: h.Router,
	}

	if h.Config.ActiveExpiry {
		interval := 1 * time.Second
		if h.Config.TickExpiry != 0 {
			interval = h.Config.TickExpiry
		}
		log.Info().Int64("Duration", int64(interval)).Msg("actively checking expired archives")
		go func() {
			heap.Init(storage.ExpiryHeap)
			ticker := time.NewTicker(interval)
			for {
				select {
				case <-ticker.C:
					for {
						if storage.ExpiryHeap.Len() > 0 && (*storage.ExpiryHeap)[0].HasExpired() {
							expiredTag := heap.Pop(storage.ExpiryHeap).(storage.ExpiryTag)
							log.Debug().
								Time("Expiry", expiredTag.Expiry).
								Str("Id", expiredTag.Id).
								Msg("found expired archive, deleting")
							err := storage.Delete(expiredTag.Id)
							log.Err(err).Str("Id", expiredTag.Id).Msg("delete expired archive")
						} else {
							break
						}
					}
				case <-h.ctx.Done():
					// TODO: save heap state
					return
				}
			}
		}()
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

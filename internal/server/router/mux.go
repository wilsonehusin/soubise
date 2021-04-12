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

package router

import (
	"bytes"
	"container/heap"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"github.com/wilsonehusin/soubise/internal/archive"
	"github.com/wilsonehusin/soubise/internal/server/middleware"
	"github.com/wilsonehusin/soubise/internal/server/routes"
	"github.com/wilsonehusin/soubise/internal/storage"
)

func NewMux() http.Handler {
	router := mux.NewRouter()

	router.Use(middleware.RequestIdentifier)
	router.Use(middleware.Logger)
	// TODO: router.Use(middleware.Authenticator)

	router.HandleFunc(routes.CreateObject, createObject).Methods("POST")
	router.HandleFunc(routes.GetObjectId, getObject).Methods("GET")

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: redirect to product landing page / GitHub repository
		if _, err := w.Write([]byte("soubise")); err != nil {
			requestLogger(r).Error().Err(err).Send()
		}
	})
	return router
}

func requestLogger(r *http.Request) *zerolog.Logger {
	return r.Context().Value(middleware.RequestLogger{}).(*zerolog.Logger)
}

func createObject(w http.ResponseWriter, r *http.Request) {
	bodyBuffer := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(bodyBuffer, r.Body); err != nil {
		requestLogger(r).Error().Err(err).Send()
		return
	}

	requestLogger(r).Debug().
		Dict("Storage", zerolog.Dict().
			Str("Action", "create")).
		Msg("processing archive")

	bodyBytes := bodyBuffer.Bytes()
	toStore, err := archive.LoadArchive(bodyBytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		requestLogger(r).Error().
			Err(err).Msg("request does not form Archive")
		return
	}

	id, err := storage.Create(bodyBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		requestLogger(r).Error().
			Err(err).Msg("create object")
		return
	}

	heap.Push(storage.ExpiryHeap, storage.ExpiryTag{Id: id, Expiry: toStore.Expiry})

	requestLogger(r).Info().
		Dict("Storage", zerolog.Dict().
			Str("Id", id).
			Str("Action", "create")).
		Msg("created archive")

	if _, err := w.Write([]byte(id)); err != nil {
		requestLogger(r).Error().Err(err).Send()
		return
	}
}

func getObject(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["Id"]
	requestLogger(r).Debug().
		Dict("Storage", zerolog.Dict().
			Str("Id", id).
			Str("Action", "get")).
		Msg("processing archive")

	obj, err := storage.Get(id)
	if len(obj) == 0 || err != nil {
		w.WriteHeader(http.StatusNotFound)
		requestLogger(r).Error().
			Err(err).Send()
		return
	}

	objArchive, err := archive.LoadArchive(obj)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		requestLogger(r).Error().
			Err(err).Send()
		return
	}

	if objArchive.HasExpired() {
		w.WriteHeader(http.StatusNotFound)
		requestLogger(r).Error().Err(fmt.Errorf("expired object was requested")).Send()
		requestLogger(r).Info().Msg("deleting expired object")
		if err := storage.Delete(id); err != nil {
			requestLogger(r).Error().Err(err).Msg("unsuccessful deletion")
		}
		return
	}

	requestLogger(r).Info().
		Dict("Storage", zerolog.Dict().
			Str("Id", id).
			Str("Action", "get")).
		Msg("found archive")

	if _, err := w.Write(obj); err != nil {
		requestLogger(r).Fatal().Err(err).Send()
	}
}

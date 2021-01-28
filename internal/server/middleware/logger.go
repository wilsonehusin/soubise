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

package middleware

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type RequestLogger struct{}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.With().Str("RequestId", r.Header.Get(RequestIdKey)).Logger()
		logger.Debug().
			Dict("request",
				zerolog.Dict().
					Str("Host", r.Host).
					Str("RemoteAddr", r.RemoteAddr).
					Str("RequestURI", r.RequestURI).
					Str("Method", r.Method).
					Str("Proto", r.Proto).
					Int64("ContentLength", r.ContentLength).
					Dict("Header", enumerateHeader(r.Header))).
			Msg("received request")
		ctx := context.WithValue(r.Context(), RequestLogger{}, &logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func enumerateHeader(h http.Header) *zerolog.Event {
	result := zerolog.Dict()
	for k, v := range h {
		result = result.Strs(k, v)
	}
	return result
}

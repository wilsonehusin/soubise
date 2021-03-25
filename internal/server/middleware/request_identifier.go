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
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	RequestIdKey = "X-SOUBISE-REQUEST-ID"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RequestIdentifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(RequestIdKey, uuid.NewString())
		next.ServeHTTP(w, r)
	})
}

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

package resolve

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/wilsonehusin/soubise/pkg/broker"
	"github.com/wilsonehusin/soubise/pkg/storage"
)

func NewStorageFromPath(storagePath string, b broker.Broker) storage.Storage {
	if b == nil {
		log.Fatal().Msg("undefined broker")
	}

	var s storage.Storage
	switch {
	case storagePath == "inmemory":
		s = storage.NewInMemoryStorage(b)
		log.Warn().Dict("Storage", zerolog.Dict().Str("Kind", s.Kind())).Msg("do NOT use in production")
	case strings.HasPrefix(storagePath, "s3://"):
		log.Fatal().Msg("s3 storage not implemented")
	case strings.HasPrefix(storagePath, "gcs://"):
		log.Fatal().Msg("gcs storage not implemented")
	default:
		log.Fatal().Msg("filesystem storage not implemented")
	}
	return s
}

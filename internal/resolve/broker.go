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

	"github.com/wilsonehusin/soubise/internal/broker"
)

func NewBrokerFromPath(brokerPath string) broker.Broker {
	var b broker.Broker
	switch {
	case strings.HasPrefix(brokerPath, "redis://"):
		log.Fatal().Dict("Broker", zerolog.Dict().Str("Kind", "redis")).Msg("not implemented")
	default:
		b = &broker.InMemoryBroker{}
		log.Warn().Dict("Broker", zerolog.Dict().Str("Kind", b.Kind())).
			Msg("unkown brokerPath, fallback to inmemory (NOT replica-safe)")
	}
	return b
}

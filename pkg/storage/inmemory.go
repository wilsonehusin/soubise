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

package storage

import (
	"github.com/wilsonehusin/soubise/pkg/broker"
)

type InMemoryStorage struct {
	broker broker.Broker
	data   map[string][]byte
}

func NewInMemoryStorage(b broker.Broker) Storage {
	return &InMemoryStorage{
		broker: b,
		data:   map[string][]byte{},
	}
}

func (s *InMemoryStorage) Create(id string, value []byte) error {
	s.broker.Lock()
	s.data[id] = value
	s.broker.Unlock()
	return nil
}

func (s *InMemoryStorage) Get(id string) ([]byte, error) {
	s.broker.RLock()
	value := s.data[id]
	s.broker.RUnlock()
	if value == nil {
		return []byte{}, &StorageNotFoundError{}
	}
	return value, nil
}

func (s *InMemoryStorage) Delete(id string) error {
	s.broker.Lock()
	delete(s.data, id)
	s.broker.Unlock()
	return nil
}

func (s *InMemoryStorage) Kind() string {
	return "inmemory"
}

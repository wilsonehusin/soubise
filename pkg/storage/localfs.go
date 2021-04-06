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
	"path"

	"github.com/peterbourgon/diskv/v3"

	"github.com/wilsonehusin/soubise/pkg/broker"
)

type LocalFsStorage struct {
	broker  broker.Broker
	backend *diskv.Diskv
}

func NewLocalFsStorage(b broker.Broker, basePath string) Storage {
	return &LocalFsStorage{
		broker: b,
		backend: diskv.New(diskv.Options{
			BasePath: path.Join(basePath, "soubisedata"),
			Transform: func(s string) []string {
				top := s[0:2]
				sub := s[2:4]
				f := s[4 : len(s)-1]
				return []string{top, sub, f}
			},
		}),
	}
}

func (s *LocalFsStorage) Create(id string, data []byte) error { // TODO: use stream?
	s.broker.Lock()
	s.backend.Write(id, data)
	s.broker.Unlock()
	return nil
}
func (s *LocalFsStorage) Get(id string) ([]byte, error) { // TODO: use stream?
	s.broker.RLock()
	val, err := s.backend.Read(id)
	s.broker.RUnlock()
	if err != nil {
		return []byte{}, err
	}
	return val, nil
}
func (s *LocalFsStorage) Delete(id string) error {
	s.broker.Lock()
	err := s.backend.Erase(id)
	s.broker.Unlock()
	if err != nil {
		return err
	}
	return nil
}
func (s *LocalFsStorage) Kind() string {
	return "localfs"
}

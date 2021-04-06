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
	"bytes"
	"fmt"
	"testing"

	"github.com/wilsonehusin/soubise/pkg/broker"
)

var inMemoryBroker = &broker.InMemoryBroker{}

var backends = map[string]Storage{
	"inmemory": NewInMemoryStorage(inMemoryBroker),
	"localfs":  NewLocalFsStorage(inMemoryBroker, "/tmp"),
}

func TestStorageBackends(t *testing.T) {
	k := "thequickbrownfox"
	v := []byte("jumpsoverthelazydog")

	for _, s := range backends {
		if err := s.Create(k, v); err != nil {
			t.Fatal(err)
		}

		val, err := s.Get(k)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(val, v) {
			t.Fatal(fmt.Errorf("expected %v, received %v", v, val))
		}

		if err := s.Delete(k); err != nil {
			t.Fatal(err)
		}

		val, err = s.Get(k)
		if err == nil || len(val) != 0 {
			t.Fatal(fmt.Errorf("expected key-value pair to have been deleted, but found value (%v) or no error thrown", val))
		}
	}
}

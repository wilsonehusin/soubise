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

package archive

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type Archive interface {
	ToBytes() ([]byte, error)
	GetName() string
	GetContent() []byte
	HasExpired() bool
}

type ArchiveObject struct {
	Name    string
	Content []byte
	Expiry  time.Time
}

func LoadArchiveObject(bin []byte) (*ArchiveObject, error) {
	var data *ArchiveObject
	decoder := gob.NewDecoder(bytes.NewBuffer(bin))

	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding archive to data: %w", err)
	}

	return data, nil
}

func (a *ArchiveObject) ToBytes() ([]byte, error) {
	var bin bytes.Buffer
	encoder := gob.NewEncoder(&bin)

	if err := encoder.Encode(a); err != nil {
		return []byte{}, fmt.Errorf("encoding data to archive: %w", err)
	}

	return bin.Bytes(), nil
}

func (a *ArchiveObject) GetName() string {
	return a.Name
}

func (a *ArchiveObject) GetContent() []byte {
	return a.Content
}

func (a *ArchiveObject) HasExpired() bool {
	return a.Expiry.Before(time.Now())
}

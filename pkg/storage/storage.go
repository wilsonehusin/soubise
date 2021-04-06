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
	"github.com/wilsonehusin/soubise/pkg/crypto"
)

var storageProvider Storage

type Storage interface {
	Create(id string, data []byte) error
	Get(id string) ([]byte, error)
	Delete(id string) error
	Kind() string
}

func SetStorage(s Storage) error {
	if storageProvider != nil {
		return &InitializedStorageError{}
	}
	storageProvider = s
	return nil
}

func Create(data []byte) (string, error) {
	if storageProvider == nil {
		return "", &UninitializedStorageError{}
	}
	id := crypto.RandLen(18).String()
	if err := storageProvider.Create(id, data); err != nil {
		return "", err
	}
	return id, nil
}

func Get(id string) ([]byte, error) {
	if storageProvider == nil {
		return []byte{}, &UninitializedStorageError{}
	}
	return storageProvider.Get(id)
}

func Delete(id string) error {
	if storageProvider == nil {
		return &UninitializedStorageError{}
	}
	return storageProvider.Delete(id)
}

func Kind() string {
	if storageProvider == nil {
		return ""
	}
	return storageProvider.Kind()
}

type InitializedStorageError struct{}

const InitializedStorageErrorString = "modifications to initialized storage is not allowed"

func (i *InitializedStorageError) Error() string {
	return InitializedStorageErrorString
}

type UninitializedStorageError struct{}

const UninitializedStorageErrorString = "storageProvider has not been initalized"

func (u *UninitializedStorageError) Error() string {
	return UninitializedStorageErrorString
}

type StorageNotFoundError struct{}

const StorageNotFoundErrorString = "unable to find archive with such key"

func (u *StorageNotFoundError) Error() string {
	return StorageNotFoundErrorString
}

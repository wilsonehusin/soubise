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

package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

const keyByteLength = 44 // encryption key and nonce

func GenerateKey() *Base64Data {
	return RandLen(keyByteLength)
}

func EncryptBlob(blob []byte, key *Base64Data) (*[]byte, error) {
	compoundKey := key.Bytes()
	if len(compoundKey) != keyByteLength {
		return nil, fmt.Errorf("key does not meet standards")
	}

	enckey := compoundKey[0:32]
	nonce := compoundKey[32:44]

	block, err := aes.NewCipher(enckey)
	if err != nil {
		return nil, fmt.Errorf("initiating cipher block: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("initiating AED cipher block: %w", err)
	}

	encrypted := aesgcm.Seal(nil, nonce, blob, nil)
	return &encrypted, nil
}

func DecryptBlob(blob []byte, key *Base64Data) (*[]byte, error) {
	compoundKey := key.Bytes()
	if len(compoundKey) != keyByteLength {
		return nil, fmt.Errorf("malformed key")
	}

	enckey := compoundKey[0:32]
	nonce := compoundKey[32:44]

	block, err := aes.NewCipher(enckey)
	if err != nil {
		return nil, fmt.Errorf("initiating cipher block: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("initiating AED cipher block: %w", err)
	}

	decrypted, err := aesgcm.Open(nil, nonce, blob, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failure: %w", err)
	}
	return &decrypted, nil
}

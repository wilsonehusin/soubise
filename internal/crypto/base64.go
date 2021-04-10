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
	"encoding/base64"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Base64Data struct {
	data []byte
}

func (b *Base64Data) Bytes() []byte {
	return b.data
}

func (b *Base64Data) String() string {
	return base64.URLEncoding.EncodeToString(b.data)
}

func RandLen(length int) *Base64Data {
	values := make([]byte, length)
	for i := range values {
		values[i] = uint8(rand.Intn(math.MaxUint8))
	}
	return &Base64Data{data: values}
}

func Base64FromString(str string) (*Base64Data, error) {
	values, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return &Base64Data{data: values}, nil
}

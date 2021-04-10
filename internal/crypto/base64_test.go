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
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	data := "dGVzdAo="

	b64, err := Base64FromString(data)
	if err != nil {
		t.Fatal(err)
	}

	if b64.String() != data {
		t.Fatal("inconstent encode decode")
	}
}
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
	"testing"
	"time"
)

var tomorrow time.Time

func init() {
	tomorrow = time.Now().Add(24 * time.Hour)
}

func (a *ArchiveObject) Equals(other ArchiveObject) bool {
	if a.Name != other.Name {
		return false
	}

	if !bytes.Equal(a.Content, other.Content) {
		return false
	}

	if a.Expiry != other.Expiry {
		return false
	}

	return true
}

func TestEncodeDecode(t *testing.T) {
	obj := ArchiveObject{
		Name:    "fakeFile_here.txt",
		Content: []byte("not much how about you"),
		Expiry:  tomorrow,
	}

	receivedBin, err := obj.ToBytes()
	if err != nil {
		t.Fatalf("%v", err)
	}

	receivedData, err := LoadArchiveObject(receivedBin)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if obj.Equals(*receivedData) {
		t.Fatalf("decoded object (%v) does not match encoded object (%v)", receivedData, obj)
	}
}

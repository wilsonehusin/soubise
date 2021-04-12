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
	"time"
)

var ExpiryHeap = &ExpiryTags{}

type ExpiryTag struct {
	Id     string
	Expiry time.Time
}

func (e *ExpiryTag) HasExpired() bool {
	return e.Expiry.Before(time.Now())
}

type ExpiryTags []ExpiryTag

func (e ExpiryTags) Len() int {
	return len(e)
}

func (e ExpiryTags) Less(i, j int) bool {
	return e[i].Expiry.Before(e[j].Expiry)
}

func (e ExpiryTags) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e *ExpiryTags) Push(x interface{}) {
	*e = append(*e, x.(ExpiryTag))
}

func (e *ExpiryTags) Pop() interface{} {
	// while the ExpiryTags slice itself has the earliest expiry on [0] index,
	// as a heap (datastructure), the same [0] is at the end of slice, ref:
	// - https://golang.org/pkg/container/heap/#Pop
	// - https://play.golang.org/p/PF1BteQxdqU
	original := *e
	count := len(original)
	earliest := original[count-1]
	*e = original[0 : count-1]
	return earliest
}

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

package shareablepath

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

type RefPath struct {
	Server        string
	Id            string
	EncryptionKey string
	OwnerKey      string
}

const Prefix = "soubise://"

var matcher = regexp.MustCompile(`soubise://(?P<host>[\w-_=]+)/(?P<id>[\w-_=]+)/(?P<encryptionKey>[\w-_=]+)/?(?P<ownerKey>[\w-_=]+)?`)

func Parse(str string) (*RefPath, error) {
	if !strings.HasPrefix(str, Prefix) {
		return nil, &RefPathParseError{invalidRefPath: str}
	}

	ownerKey := ""
	match := matcher.FindStringSubmatch(str)
	if len(match) < 4 {
		return nil, &RefPathParseError{invalidRefPath: str}
	} else if len(match) < 5 {
		ownerKey = match[4]
	}

	decodedServer, err := base64.URLEncoding.DecodeString(match[1])
	if err != nil {
		return nil, fmt.Errorf("decoding failed: %w", err)
	}
	return &RefPath{
		Server:        string(decodedServer),
		Id:            match[2],
		EncryptionKey: match[3],
		OwnerKey:      ownerKey,
	}, nil
}

func (r *RefPath) String() string {
	suffix := ""
	if r.OwnerKey != "" {
		suffix = fmt.Sprintf("/%s", r.OwnerKey)
	}

	encodedServer := base64.URLEncoding.EncodeToString([]byte(r.Server))

	return fmt.Sprintf("%s%s/%s/%s%s", Prefix, encodedServer, r.Id, r.EncryptionKey, suffix)
}

type RefPathParseError struct {
	invalidRefPath string
}

func (r *RefPathParseError) Error() string {
	return fmt.Sprintf("provided %s is not a valid RefPath", r.invalidRefPath)
}

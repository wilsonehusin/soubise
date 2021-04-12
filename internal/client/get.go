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

package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/dustin/go-humanize"

	"github.com/wilsonehusin/soubise/internal"
	"github.com/wilsonehusin/soubise/internal/archive"
	"github.com/wilsonehusin/soubise/internal/buildinfo"
	"github.com/wilsonehusin/soubise/internal/crypto"
	"github.com/wilsonehusin/soubise/internal/printer"
	"github.com/wilsonehusin/soubise/internal/server/routes"
	"github.com/wilsonehusin/soubise/internal/spinner"
)

func Get(refPath string) error {
	claimTag, err := internal.Parse(refPath)
	if err != nil {
		return err
	}
	uriBuilder, err := url.Parse(claimTag.Server)
	if err != nil {
		return fmt.Errorf("unable to parse server: %w", err)
	}
	printer.Stdout("  Server: %v\n\n", uriBuilder.String())

	key64, err := crypto.Base64FromString(claimTag.EncryptionKey)
	if err != nil {
		return fmt.Errorf("unable to decode encryption key: %w", err)
	}

	archiveBlob, err := downloadShareable(claimTag)
	if err != nil {
		return fmt.Errorf("unable to download file: %w", err)
	}

	spinner.Start(" unpack", "reconstructing")
	archiveToStore, err := archive.LoadArchive(*archiveBlob)
	if err != nil {
		spinner.StopFail("failed")
		return fmt.Errorf("unable to understand archive: %w", err)
	}
	spinner.Stop("done")

	spinner.Start(" decrypt", "doing math")
	decryptedContent, err := crypto.DecryptBlob(archiveToStore.Content, key64)
	if err != nil {
		spinner.StopFail("failed")
		return fmt.Errorf("unable to decrypt file: %w", err)
	}
	spinner.Stop("done")

	if err := writeToFile(archiveToStore.Name, *decryptedContent); err != nil {
		return fmt.Errorf("unable to write downloaded archive: %w", err)
	}

	return nil
}

func writeToFile(name string, content []byte) error {
	// TODO: check if user wants to overwrite
	spinner.Start(" write to file", "creating file")
	dest, err := os.Create(name)
	if err != nil {
		spinner.StopFail("unable to create file")
		return err
	}
	defer dest.Close()
	size, err := dest.Write(content)
	if err != nil {
		spinner.StopFail("unable to write content to file")
		return err
	}
	spinner.Stop(fmt.Sprintf("%s (%s)", name, humanize.Bytes(uint64(size))))
	return nil
}

func downloadShareable(claimTag *internal.ClaimTag) (*[]byte, error) {
	spinner.Start(" download", "resolving path")
	uriBuilder, err := url.Parse(claimTag.Server)
	if err != nil {
		spinner.StopFail("unable to parse server")
		return nil, fmt.Errorf("unable to parse %s as url: %w", claimTag.Server, err)
	}
	uriBuilder.Path = path.Join(uriBuilder.Path, routes.GetObjectWithId(claimTag.Id))
	request, err := http.NewRequest("GET", uriBuilder.String(), nil)
	if err != nil {
		spinner.StopFail("unable to compose request")
		return nil, fmt.Errorf("unable to compose request ot server: %w", err)
	}
	request.Header.Set("User-Agent", fmt.Sprintf("Soubise/%v", buildinfo.Version))

	client := &http.Client{}
	spinner.Update("pulling data")
	response, err := client.Do(request)
	if err != nil {
		spinner.StopFail("failed to download")
		return nil, fmt.Errorf("unable to download from server: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		spinner.StopFail("error")
		return nil, fmt.Errorf("server did not process request successfully: %v", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		spinner.StopFail("failed")
		return nil, fmt.Errorf("unable to parse response: %w", err)
	}
	spinner.Stop("done")

	return &body, nil
}

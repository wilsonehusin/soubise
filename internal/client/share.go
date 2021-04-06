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
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/wilsonehusin/soubise/internal/buildinfo"
	"github.com/wilsonehusin/soubise/internal/printer"
	"github.com/wilsonehusin/soubise/internal/server/routes"
	"github.com/wilsonehusin/soubise/internal/spinner"
	"github.com/wilsonehusin/soubise/pkg/archive"
	"github.com/wilsonehusin/soubise/pkg/crypto"
	"github.com/wilsonehusin/soubise/pkg/shareablepath"
)

func Share(pathToFile string, lifetime time.Duration, server string) error {
	uriBuilder, err := url.Parse(server)
	if err != nil {
		return err
	}
	uriBuilder.Path = path.Join(uriBuilder.Path, routes.CreateObject)
	printer.Stdout("   Server: %v\n", server)

	encryptionKey := crypto.GenerateKey()

	archiveToShare, err := prepareShareable(pathToFile, encryptionKey, lifetime)
	if err != nil {
		return err
	}

	encoded, err := archiveToShare.ToBytes()
	if err != nil {
		return err
	}

	// TODO: compression? ¯\_(ツ)_/¯
	buf := bytes.NewBuffer(encoded)

	request, err := http.NewRequest("POST", uriBuilder.String(), buf)
	if err != nil {
		return err
	}

	request.Header.Set("User-Agent", fmt.Sprintf("Soubise/%v", buildinfo.Version))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	spinner.Start("  upload", "sending to server")
	response, err := client.Do(request)
	if err != nil {
		spinner.StopFail("failed to upload\n")
		return err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		spinner.StopFail("error")
		return fmt.Errorf("server did not process request successfully: %v", response.Status)
	}
	spinner.Stop("done")
	printer.Stdout("\n")

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unable to read response from server: %w", err)
	}
	shareId := string(rawBody)
	shareableRef := &shareablepath.RefPath{
		Server:        server,
		Id:            shareId,
		EncryptionKey: encryptionKey.String(),
	}

	printer.Stdout("Encrypted file has been stored successfully! Use the following to share:\n")
	printer.Stdout("  %v\n", shareableRef.String())

	return nil
}

func prepareShareable(pathToFile string, encryptionKey *crypto.Base64Data, lifetime time.Duration) (*archive.ArchiveObject, error) {
	finfo, err := os.Stat(pathToFile)
	if err != nil {
		return nil, fmt.Errorf("unable to find %v: %w\n", pathToFile, err)
	}
	if finfo.IsDir() {
		return nil, fmt.Errorf("%v looks like a directory, Soubise can only process a specific file for now\n", pathToFile)
	}

	name := filepath.Base(pathToFile)
	printer.Stdout("     File: %v\n", name)
	printer.Stdout("     Size: %v\n", humanize.Bytes(uint64(finfo.Size())))

	fd, err := os.Open(pathToFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer fd.Close()

	content, err := io.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}
	checksum := sha256.Sum256(content)
	printer.Stdout("   SHA256: %x\n", checksum)

	expiry := time.Now().Add(lifetime)
	printer.Stdout("  Expires: %v (%v from now)\n\n", expiry.Format(time.RFC1123), lifetime)

	spinner.Start(" encrypt", "running some math")
	encrypted, err := crypto.EncryptBlob(content, encryptionKey)
	if err != nil {
		spinner.StopFail("failed to encrypt content")
		return nil, fmt.Errorf("unable to encrypt file: %w", err)
	}
	spinner.Update("validating")
	decrypted, err := crypto.DecryptBlob(*encrypted, encryptionKey)
	if err != nil {
		spinner.StopFail("failed to validate")
		return nil, fmt.Errorf("unable to validate encrypted file: %w", err)
	}
	if !bytes.Equal(content, *decrypted) {
		spinner.StopFail("inconsistent")
		return nil, fmt.Errorf("inconsistent encryption behavior")
	}
	spinner.Stop("done")

	return &archive.ArchiveObject{
		Name:    name,
		Content: *encrypted,
		Expiry:  expiry,
	}, nil
}

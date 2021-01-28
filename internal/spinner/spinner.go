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

package spinner

import (
	"time"

	"github.com/theckman/yacspin"

	"github.com/wilsonehusin/soubise/internal/printer"
)

var config = yacspin.Config{
	Frequency:         100 * time.Millisecond,
	CharSet:           yacspin.CharSets[69],
	SuffixAutoColon:   true,
	StopCharacter:     "✓",
	StopColors:        []string{"fgGreen"},
	StopFailCharacter: "✗",
	StopFailColors:    []string{"fgRed"},
}

var spinner *yacspin.Spinner
var enabled = true

func Disable() {
	if spinner != nil {
		printer.Stderr("spinner was initialized, forcefully stopping")
		_ = spinner.Stop()
		spinner = nil
	}
	enabled = false
}

func Start(status, message string) {
	if !enabled {
		return
	}
	if spinner != nil {
		_ = spinner.Stop()
		spinner = nil
	}

	s, err := yacspin.New(config)
	if err != nil {
		printer.Stderr("unable to start spinner, progress will be reported through log instead")
		return
	}

	spinner = s
	spinner.Suffix(status)
	spinner.Message(message)

	if err = spinner.Start(); err != nil {
		printer.Stderr("spinner was not initialized, continuing...")
	}
}

func Update(message string) {
	if !enabled {
		return
	}
	spinner.Message(message)
}

func Stop(message string) {
	if !enabled {
		return
	}
	if message != "" {
		spinner.StopMessage(message)
	}
	if err := spinner.Stop(); err != nil {
		printer.Stderr("spinner was not initialized, continuing...")
	}
}

func StopFail(status string) {
	if !enabled {
		return
	}
	if status != "" {
		spinner.StopFailMessage(status)
	}
	if err := spinner.StopFail(); err != nil {
		printer.Stderr("spinner was not initialized, continuing...")
	}
}

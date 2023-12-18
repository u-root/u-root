// Copyright 2023 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package eventchannel holds the JSON definition of an event.
package eventchannel

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// Action are the actions a guest can send.
type Action string

const (
	// ActionGuestEvent is used for a payload event.
	ActionGuestEvent Action = "guestevent"

	// ActionDone is used to signal no more events will be sent.
	ActionDone Action = "done"
)

// Event is an event channel event.
type Event[T any] struct {
	GuestAction Action `json:"hugelgupf_vmtest_guest_action"`
	Actual      T      `json:",omitempty"`
}

// ProcessJSONByLine reads JSON events from r separated by new lines.
func ProcessJSONByLine[T any](r io.Reader, callback func(T)) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		var e T
		if err := json.Unmarshal(line, &e); err != nil {
			return fmt.Errorf("JSON error (line: %s): %w", line, err)
		}
		callback(e)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}
	return nil
}

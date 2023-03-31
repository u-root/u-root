// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package progress

import (
	"bytes"
	"testing"
	"time"
)

func TestProgressBegin(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		sendQuit bool
		wait     time.Duration
	}{
		{
			name:     "Progress Begin",
			mode:     "none",
			sendQuit: false,
			wait:     0,
		},
		{
			name:     "Progress mode progress",
			mode:     "progress",
			sendQuit: true,
			wait:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			someVariable := int64(1)
			b := &bytes.Buffer{}
			p := New(b, tt.mode, &someVariable)
			p.Begin()
			time.Sleep(tt.wait * time.Second)

			if tt.sendQuit {
				p.quit <- struct{}{}
			}
		})
	}
}

func TestProgressEnd(t *testing.T) {
	tests := []struct {
		name string
		mode string
		wait time.Duration
	}{
		{
			name: "Mode none",
			mode: "none",
			wait: 1,
		},
		{
			name: "Mode progress",
			mode: "progress",
			wait: 1,
		},
		{
			name: "Mode xfer",
			mode: "xfer",
			wait: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			someVariable := int64(1)
			b := &bytes.Buffer{}
			p := New(b, tt.mode, &someVariable)
			p.Begin()
			time.Sleep(tt.wait * time.Second)

			p.End()

			// Looks like this check is sometimes faster than the channel
			time.Sleep(50 * time.Millisecond)

			p.endTimeMutex.Lock()
			if p.end.IsZero() {
				t.Errorf(`Begin(%q, %v).end.IsZero() = %t, want false`, tt.mode, &someVariable, p.end.IsZero())
			}
			p.endTimeMutex.Unlock()
		})
	}
}

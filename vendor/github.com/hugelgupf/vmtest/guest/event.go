// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package guest

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hugelgupf/vmtest/internal/eventchannel"
)

// Emitter is an event channel emitter.
type Emitter[T any] struct {
	file  *os.File
	w     *io.PipeWriter
	errCh chan error
}

// EventChannel opens an event channel to the host over the given device.
//
// Callers must call Close on Emitter to publish a final "done" event to signal
// the host no more events are coming. If the "done" event is not published,
// qemu.EventChannel is configured to return an error on VM exit on the host.
//
// T should be the type of a JSON event being sent, matching the host
// configuration on qemu.EventChannel reading from this channel.
func EventChannel[T any](path string) (*Emitter[T], error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0o777)
	if err != nil {
		return nil, err
	}

	emit := &Emitter[T]{
		file: f,
	}

	r, w := io.Pipe()
	errCh := make(chan error)
	go func() {
		defer r.Close()
		err := eventchannel.ProcessJSONByLine[T](r, func(t T) {
			if err := emit.Emit(t); err != nil {
				log.Printf("Error emitting event: %v", err)
			}
		})
		errCh <- err
	}()
	emit.w = w
	emit.errCh = errCh
	return emit, nil
}

// Write writes JSON bytes on the event channel. Write expects events to be
// separated by new lines. Callers may chunk their writes.
//
// This makes the Emitter compatible with exec.Cmd.Stdout/Stderr for commands
// that emit JSON one line at a time.
func (e *Emitter[T]) Write(p []byte) (int, error) {
	return e.w.Write(p)
}

// Emit emits one T event.
func (e *Emitter[T]) Emit(t T) error {
	return e.sendEvent(eventchannel.Event[T]{
		Actual:      t,
		GuestAction: eventchannel.ActionGuestEvent,
	})
}

func (e *Emitter[T]) sendEvent(event eventchannel.Event[T]) error {
	b, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	b = append(b, '\n')
	if n, err := e.file.Write(b); err != nil {
		return err
	} else if n != len(b) {
		return fmt.Errorf("incomplete write: want %d, sent %d", len(b), n)
	}
	return nil
}

// Close sends the "done" event to assure the host there will be no more events
// and closes the event channel.
func (e *Emitter[T]) Close() error {
	// Ensure that ActionDone is the last event we send by waiting for
	// Goroutine to exit first.
	e.w.Close()
	err := <-e.errCh

	if werr := e.sendEvent(eventchannel.Event[T]{GuestAction: eventchannel.ActionDone}); werr != nil && err != nil {
		err = werr
	}
	_ = e.file.Sync()
	e.file.Close()
	return err
}

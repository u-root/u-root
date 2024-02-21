// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package qevent implements a JSON-based event channel between guest and host.
package qevent

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"

	"github.com/creack/pty"
	"github.com/hugelgupf/vmtest/internal/eventchannel"
	"github.com/hugelgupf/vmtest/qemu"
)

type ptmClosedErrorConverter struct {
	r io.Reader
}

// "read /dev/ptmx: input/output error" error occufs on Linux while reading
// from the ptm after the pts is closed.
var ptmClosed = os.PathError{
	Op:   "read",
	Path: "/dev/ptmx",
	Err:  syscall.EIO,
}

func (c ptmClosedErrorConverter) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	var perr *os.PathError
	if errors.As(err, &perr) && *perr == ptmClosed {
		return n, io.EOF
	}
	return n, err
}

// ErrEventChannelMissingDoneEvent is returned when the final event channel
// event is not received.
var ErrEventChannelMissingDoneEvent = errors.New("never received the final event channel event (did you call Close() on the guest event channel emitter?)")

// EventChannel adds a virtio-serial-backed channel between host and guest to
// send JSON events (T).
//
// Use guest.SerialEventChannel with the same name to get access to the emitter
// in the guest.
//
// Guest events will be sent on the supplied channel. The channel will be
// closed when the guest exits or indicates that no more events are coming. If
// the guest exits without indicating that no more events are coming, the VM
// exit will return an error. (guest.SerialEventChannel.Close emits this "done"
// event.)
//
// If the channel is blocking, guest event processing is blocked as well.
func EventChannel[T any](name string, events chan<- T) qemu.Fn {
	return func(alloc *qemu.IDAllocator, opts *qemu.Options) error {
		pipeID := alloc.ID("pipe")

		ptm, pts, err := pty.Open()
		if err != nil {
			return err
		}
		fd := opts.AddFile(pts)
		opts.AppendQEMU(
			"-device", "virtio-serial",
			"-device", fmt.Sprintf("virtserialport,chardev=%s,name=%s", pipeID, name),
			"-chardev", fmt.Sprintf("pipe,id=%s,path=/proc/self/fd/%d", pipeID, fd),
		)

		var gotDone bool
		opts.Tasks = append(opts.Tasks, qemu.WaitVMStarted(func(ctx context.Context, n *qemu.Notifications) error {
			// Close ptm if it isn't already closed due to the VM
			// exiting.
			defer ptm.Close()

			// Close write-end on parent side.
			pts.Close()

			err := eventchannel.ProcessJSONByLine[eventchannel.Event[T]](ptmClosedErrorConverter{ptm}, func(c eventchannel.Event[T]) {
				switch c.GuestAction {
				case eventchannel.ActionGuestEvent:
					events <- c.Actual

				case eventchannel.ActionDone:
					close(events)
					gotDone = true
				}
			})
			if err != nil {
				if !gotDone {
					close(events)
				}
				return err
			}
			if !gotDone {
				close(events)
				return ErrEventChannelMissingDoneEvent
			}
			return nil
		}))
		return nil
	}
}

// EventChannelCallback adds a virtio-serial-backed channel between host and
// guest to send JSON events (T).
//
// Use guest.SerialEventChannel with the same name to get access to the emitter
// in the guest.
//
// When a guest event occurs, the callback is called.
func EventChannelCallback[T any](name string, callback func(T)) qemu.Fn {
	ch := make(chan T)
	return func(alloc *qemu.IDAllocator, opts *qemu.Options) error {
		opts.Tasks = append(opts.Tasks, func(ctx context.Context, n *qemu.Notifications) error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()

				case e, ok := <-ch:
					if !ok {
						return nil
					}
					callback(e)
				}
			}
		})
		return EventChannel[T](name, ch)(alloc, opts)
	}
}

// ReadFile reads events from a file that was written to using
// guest.EventChannel.
func ReadFile[T any](path string) ([]T, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var t []T
	var gotDone bool
	err = eventchannel.ProcessJSONByLine[eventchannel.Event[T]](f, func(c eventchannel.Event[T]) {
		switch c.GuestAction {
		case eventchannel.ActionGuestEvent:
			t = append(t, c.Actual)

		case eventchannel.ActionDone:
			gotDone = true
		}
	})
	if err != nil {
		return nil, err
	}
	if !gotDone {
		return nil, ErrEventChannelMissingDoneEvent
	}
	return t, nil
}

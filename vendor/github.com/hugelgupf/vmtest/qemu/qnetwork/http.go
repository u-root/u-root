// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qnetwork

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/hugelgupf/vmtest/qemu"
)

// ServeHTTP serves s on l until the VM guest exits.
func ServeHTTP(s *http.Server, l net.Listener) qemu.Fn {
	return qemu.All(
		qemu.WithTask(func(ctx context.Context, n *qemu.Notifications) error {
			if err := s.Serve(l); !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return nil
		}),
		qemu.WithTask(qemu.Cleanup(func() error {
			// Stop HTTP server.
			return s.Close()
		})),
	)
}

// Copyright 2017-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tftp

package curl

import (
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/u-root/u-root/pkg/uio"
	"pack.ag/tftp"
)

var (
	// DefaultTFTPClient is the default TFTP FileScheme.
	DefaultTFTPClient = NewTFTPClient(tftp.ClientMode(tftp.ModeOctet), tftp.ClientBlocksize(1450), tftp.ClientWindowsize(64))
)

func init() {
	DefaultSchemes["tftp"] = DefaultTFTPClient
}

// TFTPClient implements FileScheme for TFTP files.
type TFTPClient struct {
	opts []tftp.ClientOpt
}

// NewTFTPClient returns a new TFTP client based on the given tftp.ClientOpt.
func NewTFTPClient(opts ...tftp.ClientOpt) FileScheme {
	return &TFTPClient{
		opts: opts,
	}
}

func tftpFetch(_ context.Context, t *TFTPClient, u *url.URL) (io.Reader, error) {
	// TODO(hugelgupf): These clients are basically stateless, except for
	// the options. Figure out whether you actually have to re-establish
	// this connection every time. Audit the TFTP library.
	c, err := tftp.NewClient(t.opts...)
	if err != nil {
		return nil, err
	}

	r, err := c.Get(u.String())
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Fetch implements FileScheme.Fetch for TFTP.
func (t *TFTPClient) Fetch(ctx context.Context, u *url.URL) (io.ReaderAt, error) {
	r, err := tftpFetch(ctx, t, u)
	if err != nil {
		return nil, err
	}
	return uio.NewCachingReader(r), nil
}

// FetchWithoutCache implements FileScheme.FetchWithoutCache for TFTP.
func (t *TFTPClient) FetchWithoutCache(ctx context.Context, u *url.URL) (io.Reader, error) {
	return tftpFetch(ctx, t, u)
}

// RetryTFTP retries downloads if the error does not contain FILE_NOT_FOUND.
//
// pack.ag/tftp does not export the necessary structs to get the
// code out of the error message cleanly, but it does embed FILE_NOT_FOUND in
// the error string.
func RetryTFTP(u *url.URL, err error) bool {
	return !strings.Contains(err.Error(), "FILE_NOT_FOUND")
}

// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftp

import (
	"io"

	"pack.ag/tftp"
)

// Response serves as interface which allows mocking of tftp.Response for testing and
// usage of the real implementation, depending on the use-case.
type Response interface {
	Read([]byte) (int, error)
	Size() (int64, error)
}

// ClientIf serves as interface which allows mocking of tftp.Client for testing and
// usage of the real implementation, depending on the use-case.
type ClientIf interface {
	Put(string, io.Reader, int64) error
	Get(string) (Response, error)
}

// ClientMock serves as the Mock structure of Client for testing.
type ClientMock struct{}

// DummyResp serves as the mock structure of Response for testing.
type DummyResp struct{}

// Size mocks the Size function of tftp.Response for testing.
func (d *DummyResp) Size() (int64, error) {
	return 0, nil
}

// Read mocks the Read function of tftp.Response for testing.
func (d *DummyResp) Read(b []byte) (int, error) {
	return 0, nil
}

// Get mocks the Get method of tftp.Client.
func (c *ClientMock) Get(url string) (Response, error) {
	return &DummyResp{}, nil
}

// Put mocks the Put method of tftp.Client.
func (c *ClientMock) Put(url string, r io.Reader, size int64) error {
	return nil
}

// Client implements the ClientIf and uses the tftp.Client as member to interact with the real library.
type Client struct {
	*tftp.Client
}

// RealResponse implements the Response interface and uses the tftp.Response as member to interact with the real library.
type RealResponse struct {
	*tftp.Response
}

// Read provides the Read function of tftp.Response.
func (r *RealResponse) Read(b []byte) (int, error) {
	return r.Response.Read(b)
}

// Size provides the Size function of tftp.Response.
func (r *RealResponse) Size() (int64, error) {
	return r.Response.Size()
}

// NewClient sets up a new tftp.Client according to the given ClientCfg struct.
func NewClient(ccfg *ClientCfg) (*Client, error) {
	c, err := tftp.NewClient(tftp.ClientMode(ccfg.Mode), ccfg.Rexmt, ccfg.Timeout)
	return &Client{
		Client: c,
	}, err
}

// Get provides the Get method of tftp.Client.
func (c *Client) Get(url string) (Response, error) {
	r, err := c.Client.Get(url)
	return &RealResponse{r}, err
}

// Put provides the Put method of tftp.Client.
func (c *Client) Put(url string, r io.Reader, size int64) error {
	return c.Client.Put(url, r, size)
}

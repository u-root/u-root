// Copyright 2012-2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tftp

import (
	"io"

	"pack.ag/tftp"
)

type Response interface {
	Read([]byte) (int, error)
	Size() (int64, error)
}

type ClientIf interface {
	Put(string, io.Reader, int64) error
	Get(string) (Response, error)
}

type ClientMock struct{}

type DummyResp struct{}

func (d *DummyResp) Size() (int64, error) {
	return 0, nil
}

func (d *DummyResp) Read(b []byte) (int, error) {
	return 0, nil
}

func (c *ClientMock) Get(url string) (Response, error) {
	return &DummyResp{}, nil
}

func (c *ClientMock) Put(url string, r io.Reader, size int64) error {
	return nil
}

type Client struct {
	*tftp.Client
}

type RealResponse struct {
	*tftp.Response
}

func (r *RealResponse) Read(b []byte) (int, error) {
	return r.Response.Read(b)
}

func (r *RealResponse) Size() (int64, error) {
	return r.Response.Size()
}

func NewClient(ccfg *clientCfg) (*Client, error) {
	c, err := tftp.NewClient(tftp.ClientMode(ccfg.mode), ccfg.rexmt, ccfg.timeout)
	return &Client{
		Client: c,
	}, err
}

func (c *Client) Get(url string) (Response, error) {
	r, err := c.Client.Get(url)
	return &RealResponse{r}, err
}

func (c *Client) Put(url string, r io.Reader, size int64) error {
	return c.Client.Put(url, r, size)
}

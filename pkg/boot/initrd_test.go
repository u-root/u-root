// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package boot

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/uio"
)

type file struct {
	name    string
	content []byte
}

func (f file) String() string {
	return f.name
}

func (f file) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(f.content)) {
		return 0, io.EOF
	}
	return copy(p, f.content[off:]), nil
}

func TestCatInitrds(t *testing.T) {
	for _, tt := range []struct {
		readers     []io.ReaderAt
		wantName    string
		wantContent []byte
		wantSize    int
	}{
		{
			readers: []io.ReaderAt{
				bytes.NewReader(make([]byte, 512)),
				bytes.NewReader(make([]byte, 512)),
			},
			wantName:    "*bytes.Reader,*bytes.Reader",
			wantContent: make([]byte, 1024),
			wantSize:    1024,
		},
		{
			readers: []io.ReaderAt{
				strings.NewReader("yay"),
				bytes.NewReader(make([]byte, 777)),
			},
			wantName:    "*strings.Reader,*bytes.Reader",
			wantContent: append([]byte("yay"), make([]byte, 509+777)...),
			wantSize:    3 + 509 + 777,
		},
		{
			readers: []io.ReaderAt{
				strings.NewReader("yay"),
			},
			wantName:    "*strings.Reader",
			wantContent: []byte("yay"),
			wantSize:    3,
		},
		{
			readers: []io.ReaderAt{
				strings.NewReader("foo"),
				strings.NewReader("bar"),
			},
			wantName:    "*strings.Reader,*strings.Reader",
			wantContent: append(append([]byte("foo"), make([]byte, 509)...), []byte("bar")...),
			wantSize:    3 + 509 + 3,
		},
		{
			readers: []io.ReaderAt{
				file{
					name:    "/bar/foo",
					content: []byte("foo"),
				},
				file{
					name:    "/bar/bar",
					content: []byte("bar"),
				},
			},
			wantName:    "/bar/foo,/bar/bar",
			wantContent: append(append([]byte("foo"), make([]byte, 509)...), []byte("bar")...),
			wantSize:    3 + 509 + 3,
		},
	} {
		got := CatInitrds(tt.readers...)

		by, err := uio.ReadAll(got)
		if err != nil {
			t.Errorf("CatInitrdReader errored: %v", err)
		}

		if len(by) != tt.wantSize {
			t.Errorf("Cat(%v) = len %d, want len %d", tt.readers, len(by), tt.wantSize)
		}
		if !bytes.Equal(by, tt.wantContent) {
			t.Errorf("Cat(%v) = %v, want %v", tt.readers, by, tt.wantContent)
		}
		s := fmt.Sprintf("%s", got)
		if s != tt.wantName {
			t.Errorf("Cat(%v) = name %s, want %s", tt.readers, s, tt.wantName)
		}
	}
}

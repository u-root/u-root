// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kmodule

import (
	"bytes"
	"compress/gzip"
	"errors"
	"os"
	"path"
	"testing"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
	"golang.org/x/sync/errgroup"
)

var procModsMock = `hid_generic 16384 0 - Live 0x0000000000000000
usbhid 49152 0 - Live 0x0000000000000000
ccm 20480 6 - Live 0x0000000000000000
`

func TestGenLoadedMods(t *testing.T) {
	m := depMap{
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko":   &dependency{},
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko": &dependency{},
		"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko":                &dependency{},
	}
	br := bytes.NewBufferString(procModsMock)
	err := genLoadedMods(br, m)
	if err != nil {
		t.Fatalf("fail to genLoadedMods: %v\n", err)
	}
	for mod, d := range m {
		if d.state != loaded {
			t.Fatalf("mod %q should have been loaded", path.Base(mod))
		}
	}
}

func TestParallelLoad(t *testing.T) {
	loadTime := 100 * time.Millisecond

	m := depMap{
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko":   &dependency{},
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko": &dependency{},
		"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko":                &dependency{},
		"/lib/modules/6.6.6-generic/kernel/tests/depmod.ko": &dependency{
			deps: []string{"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko",
				"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko",
				"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko",
			},
		},
		"/lib/modules/6.6.6-generic/kernel/tests/depmod2.ko": &dependency{
			deps: []string{"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko",
				"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko",
			},
		},
	}

	var eg errgroup.Group

	prober := Prober{
		deps: m,
		// Wait time encourages racing between dependencies.
		opts: ProbeOpts{DryRunCB: func(path string) { time.Sleep(loadTime) }},
	}

	eg.Go(func() error {
		return prober.Probe("depmod", "")
	})
	eg.Go(func() error {
		return prober.Probe("depmod2", "")
	})

	start := time.Now()
	err := eg.Wait()

	if err != nil {
		t.Fatalf("Probing failed: %v", err)
	}

	// Racy... longest parallel chain is 2, and we load 5 total modules.
	if time.Since(start) > loadTime*3 {
		t.Fatalf("module loading slow")
	}

	for mod, d := range m {
		if d.state != loaded {
			t.Fatalf("mod %q should have been loaded", path.Base(mod))
		}
	}
}

func TestInvalidCircularLoad(t *testing.T) {
	m := depMap{
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/hid-generic.ko":   &dependency{},
		"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko": &dependency{},
		"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko":                &dependency{},
		"/lib/modules/6.6.6-generic/kernel/tests/circlemod.ko": &dependency{
			deps: []string{"/lib/modules/6.6.6-generic/kernel/tests/depmod.ko"},
		},
		"/lib/modules/6.6.6-generic/kernel/tests/depmod.ko": &dependency{
			deps: []string{"/lib/modules/6.6.6-generic/kernel/crypto/ccm.ko",
				"/lib/modules/6.6.6-generic/kernel/drivers/hid/usbhid/usbhid.ko",
				"/lib/modules/6.6.6-generic/kernel/tests/circlemod.ko",
			},
		},
	}

	prober := Prober{
		deps: m,
		opts: ProbeOpts{DryRunCB: func(path string) {}},
	}

	err := prober.Probe("depmod", "")

	if err == nil {
		// If we reach this, we're probably hung.
		t.Fatalf("Circular dep should have errored...")
	}
}

// Helper function to generate compression test data for TestCompression.
// Generates a map with the name of the file as key and the compressed data as value.
// The data is compressed using xz, gzip, and zstd.
// There is also one file with a bad extension to test the error handling.
// Testing compression itself is out of scope.
func generateCompressionTestData(data []byte) (map[string][]byte, error) {
	var compressionBuffer bytes.Buffer

	tData := make(map[string][]byte)

	// 0. ko
	tData["test.ko"] = make([]byte, len(data))
	copy(tData["test.xz"], data)

	// 1. xz
	wXZ, err := xz.NewWriter(&compressionBuffer)
	if err != nil {
		return nil, err
	}
	_, err = wXZ.Write(data)
	if err != nil {
		return nil, err
	}
	if err = wXZ.Close(); err != nil {
		return nil, err
	}

	tData["test.xz"] = make([]byte, compressionBuffer.Len())
	copy(tData["test.xz"], compressionBuffer.Bytes())
	compressionBuffer.Reset()

	// 2. gzip
	wGZ := gzip.NewWriter(&compressionBuffer)
	if _, err = wGZ.Write(data); err != nil {
		return nil, err
	}
	if err = wGZ.Close(); err != nil {
		return nil, err
	}

	tData["test.gz"] = make([]byte, compressionBuffer.Len())
	copy(tData["test.gz"], compressionBuffer.Bytes())
	compressionBuffer.Reset()

	// 3. zstd
	wZST, err := zstd.NewWriter(&compressionBuffer)
	if err != nil {
		return nil, err
	}

	if _, err = wZST.Write(data); err != nil {
		return nil, err
	}

	if err = wZST.Close(); err != nil {
		return nil, err
	}

	tData["test.zst"] = make([]byte, compressionBuffer.Len())
	copy(tData["test.zst"], compressionBuffer.Bytes())
	compressionBuffer.Reset()

	// 4. bad
	tData["test.bad"] = []byte{'b', 'a', 'd'}
	return tData, nil
}

// Since we don't need to test the compression function, we just check
// for validity of file extension detection.
func TestCompression(t *testing.T) {
	const compressionTestString = "test\x00"

	tDir := t.TempDir()
	tFd := make(map[string]*os.File, 4)
	tFiles, err := generateCompressionTestData([]byte(compressionTestString))
	if err != nil {
		t.Fatalf("failed to generate test data: '%v'\n", err)
	}

	for name, data := range tFiles {
		tFd[name], err = os.Create(path.Join(tDir, name))
		if err != nil {
			t.Fatalf("failed to create test file %q: '%v'\n", name, err)
		}
		defer tFd[name].Close()

		n, err := tFd[name].Write(data)
		if err != nil {
			t.Fatalf("failed to write to test file %q: '%v'\n", name, err)
		}

		if err = tFd[name].Sync(); err != nil {
			t.Fatalf("failed to sync test file %q: '%v'\n", name, err)
		}

		if _, err := tFd[name].Seek(0, 0); err != nil {
			t.Fatalf("failed to seek to beginning of test file %q: '%v'\n", name, err)
		}

		if n != len(data) {
			t.Fatalf("failed to write all data to test file %q. Expected %d bytes, wrote %d\n", name, len(data), n)
		}
	}

	// defer func() {
	// 	for _, f := range tFd {
	// 		f.Close()
	// 	}
	// }()

	testCases := map[string]struct {
		file    *os.File
		ext     string
		isError bool
		err     error
	}{
		"test.ko": {
			file:    tFd["test.ko"],
			ext:     ".ko",
			isError: false,
			err:     nil,
		},
		"test.xz": {
			file:    tFd["test.xz"],
			ext:     ".xz",
			isError: false,
			err:     nil,
		},
		"test.gz": {
			file:    tFd["test.gz"],
			ext:     ".gz",
			isError: false,
			err:     nil,
		},
		"test.zst": {
			file:    tFd["test.zst"],
			ext:     ".zst",
			isError: false,
			err:     nil,
		},
		"test.bad": {
			file:    tFd["test.bad"],
			ext:     ".bad",
			isError: true,
			err:     os.ErrNotExist,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := compressionReader(tc.file)
			if tc.isError {
				if errors.Is(err, tc.err) {
					return
				}
				t.Fatalf("expected error %v but got '%v'\n", tc.err, err)
			}
			if !tc.isError && err != nil {
				t.Fatalf("expected no error but got '%v'\n", err)
			}
		})
	}
}

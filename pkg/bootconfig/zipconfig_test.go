// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bootconfig

import (
	"archive/zip"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// the sample ZIP file contains the following structure:
// test/
// test/a
//
// where the file "test/a" contains the ASCII string "blah"
var sampleZIP = []byte("PK\x03\x04\n\x00\x00\x00\x00\x00\xa6\x858M\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05\x00\x1c\x00test/UT\t\x00\x03\x88\x06\xa9[\x8b\x06\xa9[ux\x0b\x00\x01\x04\xe8\x03\x00\x00\x04\xe8\x03\x00\x00PK\x03\x04\n\x00\x00\x00\x00\x00\xa6\x858M-2\xc4P\x05\x00\x00\x00\x05\x00\x00\x00\x06\x00\x1c\x00test/aUT\t\x00\x03\x88\x06\xa9[\x88\x06\xa9[ux\x0b\x00\x01\x04\xe8\x03\x00\x00\x04\xe8\x03\x00\x00blah\nPK\x01\x02\x1e\x03\n\x00\x00\x00\x00\x00\xa6\x858M\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05\x00\x18\x00\x00\x00\x00\x00\x00\x00\x10\x00\xfdA\x00\x00\x00\x00test/UT\x05\x00\x03\x88\x06\xa9[ux\x0b\x00\x01\x04\xe8\x03\x00\x00\x04\xe8\x03\x00\x00PK\x01\x02\x1e\x03\n\x00\x00\x00\x00\x00\xa6\x858M-2\xc4P\x05\x00\x00\x00\x05\x00\x00\x00\x06\x00\x18\x00\x00\x00\x00\x00\x01\x00\x00\x00\xb4\x81?\x00\x00\x00test/aUT\x05\x00\x03\x88\x06\xa9[ux\x0b\x00\x01\x04\xe8\x03\x00\x00\x04\xe8\x03\x00\x00PK\x05\x06\x00\x00\x00\x00\x02\x00\x02\x00\x97\x00\x00\x00\x84\x00\x00\x00\x00\x00")

func TestMemoryZipReader(t *testing.T) {
	r, err := zip.NewReader(&memoryZipReader{Content: sampleZIP}, int64(len(sampleZIP)))
	require.NoError(t, err)
	// this is ugly, but we expect exactly this sequence of files (and the
	// parent directory has to be first)
	var numEntries int
	for idx, f := range r.File {
		numEntries++
		switch idx {
		case 0:
			require.Equal(t, "test/", f.Name)
		case 1:
			require.Equal(t, "test/a", f.Name)
			fd, err := f.Open()
			require.NoError(t, err)
			buf, err := ioutil.ReadAll(fd)
			require.NoError(t, err)
			require.Equal(t, []byte("blah\n"), buf)
		}
	}
	// exactly two entries in the zip file
	require.Equal(t, 2, numEntries)
}

func TestFromZip(t *testing.T) {
	manifest, tempdir, err := FromZip("testdata/bootconfig.zip", nil)
	defer func() {
		if tempdir != "" {
			if err := os.RemoveAll(tempdir); err != nil {
				log.Printf("Cannot remove temp dir %s: %v", tempdir, err)
			}
		}
	}()
	require.NoError(t, err)
	require.NotEqual(t, "", tempdir)
	require.NotNil(t, manifest)
	require.Equal(t, 1, manifest.Version)
	require.Equal(t, 1, len(manifest.Configs))
	bc := manifest.Configs[0]
	require.Equal(t, "first boot entry", bc.Name)
	require.Equal(t, "/path/to/kernel", bc.Kernel)
	require.Equal(t, "console=ttyS0", bc.KernelArgs)
}

func TestFromZipWithSignature(t *testing.T) {
	pubkey := "testdata/pubkey"
	manifest, tempdir, err := FromZip("testdata/bootconfig_signed.zip", &pubkey)
	defer func() {
		if tempdir != "" {
			if err := os.RemoveAll(tempdir); err != nil {
				log.Printf("Cannot remove temp dir %s: %v", tempdir, err)
			}
		}
	}()
	require.NoError(t, err)
	require.NotEqual(t, "", tempdir)
	require.NotNil(t, manifest)
	require.Equal(t, 1, manifest.Version)
	require.Equal(t, 1, len(manifest.Configs))
	bc := manifest.Configs[0]
	require.Equal(t, "boot entry 0", bc.Name)
	require.Equal(t, "/path/to/kernel", bc.Kernel)
}

func TestFromZipWithMissingSignature(t *testing.T) {
	pubkey := "testdata/pubkey"
	_, tempdir, err := FromZip("testdata/bootconfig.zip", &pubkey)
	defer func() {
		// called just in case FromZip does not return an error
		if tempdir != "" {
			if err := os.RemoveAll(tempdir); err != nil {
				log.Printf("Cannot remove temp dir %s: %v", tempdir, err)
			}
		}
	}()
	require.Error(t, err)
}

func TestFromZipNoSuchFile(t *testing.T) {
	_, _, err := FromZip("testdata/nonexisting_bootconfig.zip", nil)
	require.Error(t, err)
}

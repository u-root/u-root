// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bzimage

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os/exec"

	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/therootcompany/xz"
	"github.com/ulikunitz/xz/lzma"
)

// stripSize returns a decompressor which strips off the last 4 bytes of the
// data from the reader and copies the bytes to the writer.
func stripSize(d decompressor) decompressor {
	return func(w io.Writer, r io.Reader) error {
		// Read all of the bytes so that we can determine the size.
		allBytes, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("error reading all bytes: %w", err)
		}
		strippedLen := int64(len(allBytes) - 4)
		Debug("Stripped reader is of length %d bytes", strippedLen)

		reader := bytes.NewReader(allBytes)
		return d(w, io.LimitReader(reader, strippedLen))
	}
}

// execer returns a decompressor which executes the command that reads
// compressed bytes from stdin and writes the decompressed bytes to stdout.
func execer(command string, args ...string) decompressor {
	return func(w io.Writer, r io.Reader) error {
		cmd := exec.Command(command, args...)
		cmd.Stdin = r
		cmd.Stdout = w

		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			return fmt.Errorf("error creating Stderr pipe: %w", err)
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("error starting decompressor: %w", err)
		}

		stderr, err := io.ReadAll(stderrPipe)
		if err != nil {
			return fmt.Errorf("error reading stderr: %w", err)
		}

		if err := cmd.Wait(); err != nil || len(stderr) > 0 {
			return fmt.Errorf("decompressor failed: err=%w, stderr=%q", err, stderr)
		}
		return nil
	}
}

// gunzip reads compressed bytes from the io.Reader and writes the uncompressed bytes to the
// writer. gunzip satisfies the decompressor interface.
func gunzip(w io.Writer, r io.Reader) error {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("error creating gzip reader: %w", err)
	}

	if _, err := io.Copy(w, gzipReader); err != nil {
		return fmt.Errorf("failed writing decompressed bytes to writer: %w", err)
	}
	return nil
}

// unlzma reads compressed bytes from the io.Reader and writes the uncompressed bytes to the
// writer. unlzma satisfies the decompressor interface.
func unlzma(w io.Writer, r io.Reader) error {
	lzmaReader, err := lzma.NewReader(r)
	if err != nil {
		return fmt.Errorf("error creating lzma reader: %w", err)
	}

	if _, err := io.Copy(w, lzmaReader); err != nil {
		return fmt.Errorf("failed writing decompressed bytes to writer: %w", err)
	}
	return nil
}

// unlz4 reads compressed bytes from the io.Reader and writes the uncompressed bytes to the
// writer. unlz4 satisfies the decompressor interface.
func unlz4(w io.Writer, r io.Reader) error {
	lz4Reader := lz4.NewReader(r)

	if _, err := io.Copy(w, lz4Reader); err != nil {
		return fmt.Errorf("failed writing decompressed bytes to writer: %w", err)
	}
	return nil
}

// unbzip2 reads compressed bytes from the io.Reader and writes the uncompressed bytes to the
// writer. unbzip2 satisfies the decompressor interface.
func unbzip2(w io.Writer, r io.Reader) error {
	bzip2Reader := bzip2.NewReader(r)

	if _, err := io.Copy(w, bzip2Reader); err != nil {
		return fmt.Errorf("failed writing decompressed bytes to writer: %w", err)
	}
	return nil
}

// unzstd reads compressed bytes from the io.Reader and writes the uncompressed bytes to the
// writer. unzstd satisfies the decompressor interface.
func unzstd(w io.Writer, r io.Reader) error {
	zstdReader, err := zstd.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to create new reader: %w", err)
	}
	defer zstdReader.Close()

	if _, err := io.Copy(w, zstdReader); err != nil {
		return fmt.Errorf("failed writing decompressed bytes to writer: %w", err)
	}
	return nil
}

// unxz reads compressed bytes from the io.Reader and writes the uncompressed bytes to the
// writer. unxz satisfies the decompressor interface.
func unxz(w io.Writer, r io.Reader) error {
	unxzReader, err := xz.NewReader(r, 0)
	if err != nil {
		return fmt.Errorf("failed to create new reader: %w", err)
	}

	if _, err := io.Copy(w, unxzReader); err != nil {
		return fmt.Errorf("failed writing decompressed bytes to writer: %w", err)
	}
	return nil
}

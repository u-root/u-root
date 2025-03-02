// Copyright 2012-2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tail prints the lasts 10 lines of a file. Can additionally follow the
// the end of the file as it grows.
//
// Synopsis:
//     tail [-f] [-n lines_to_show] [FILE]
//
// Description:
//     If no files are specified, read from stdin.
//
// Options:
//     -f: follow the end of the file as it grows
//     -n: specify the number of lines to show (default: 10)

// Missing features:
// - follow-mode (i.e. tail -f)

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"syscall"
	"time"
)

var (
	flagFollow   = flag.Bool("f", false, "follow the end of the file")
	flagNumLines = flag.Int("n", 10, "specify the number of lines to show")
)

type readAtSeeker interface {
	io.ReaderAt
	io.Seeker
}

// tailConfig is a configuration object for the Tail function
type tailConfig struct {
	// enable follow-mode (-f)
	follow bool

	// specifies the number of lines to print (-n)
	numLines int

	// time to wait for new write in follow mode
	followDuration time.Duration
}

// getBlockSize returns the number of bytes to read for each ReadAt call. This
// helps minimize the number of syscalls to get the last N lines of the file.
func getBlockSize(numLines int) int64 {
	// This is currently computed as 81 * N, where N is the requested number of
	// lines, and 81 is a relatively generous estimation of the average line
	// length.
	return 81 * int64(numLines)
}

// lastNLines finds the n-th-to-last line in `buf`, and returns a new slice
// containing only the last `n` lines. If less lines are found, the input slice
// is returned unmodified.
func lastNLines(buf []byte, n int) []byte {
	slice := buf
	// `data` contains up to `n` lines of the file
	var data []byte
	if len(slice) != 0 {
		if slice[len(slice)-1] == '\n' {
			// don't consider the last new line for the line count
			slice = slice[:len(slice)-1]
		}
		var (
			foundLines int
			idx        int
		)
		for foundLines < n {
			// find newlines backwards from the end of `slice`
			idx = bytes.LastIndexByte(slice, '\n')
			if idx == -1 {
				// there are less than `n` lines
				break
			}
			foundLines++
			if len(slice) > 1 && slice[idx-1] == '\n' {
				slice = slice[:idx]
			} else {
				slice = slice[:idx-1]
			}
		}
		if idx == -1 {
			// if there are less than `numLines` lines, use all what we have read
			data = buf
		} else {
			data = buf[idx+1:] // +1 to skip the newline belonging to the previous line
		}
	}
	return data
}

// readLastLinesBackwards reads the last N lines from the provided file, reading
// backwards from the end of the file. This is more efficient than reading from
// the beginning, but can only be done on seekable files, (e.g. this won't work
// on stdin). For non-seekable files see readLastLinesFromBeginning.
// It returns an error, if any. If no error is encountered, the File object's
// offset is positioned after the last read location.
func readLastLinesBackwards(input readAtSeeker, writer io.Writer, numLines int) error {
	blkSize := getBlockSize(numLines)
	// go to the end of the file
	lastPos, err := input.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	// read block by block backwards until `numLines` lines are found
	readData := make([]byte, 0)
	buf := make([]byte, blkSize)
	pos := lastPos
	var foundLines int
	// for each block, count how many new lines, until they add up to `numLines`
	for pos != 0 {
		thisChunkSize := min(pos, blkSize)
		pos -= thisChunkSize
		n, err := input.ReadAt(buf, pos)
		if err != nil && err != io.EOF {
			return err
		}
		// avoid problems where `n` is smaller or bigger than `thisChunkSize`
		n = min(n, int(thisChunkSize))
		// merge this block to what was read so far
		readData = slices.Concat(buf[:n], readData)
		// count how many lines we have so far, and stop reading if we have
		// enough
		foundLines += bytes.Count(buf[:n], []byte{'\n'})
		if foundLines >= numLines {
			break
		}
	}
	// find the start of the n-th to last line
	data := lastNLines(readData, numLines)
	// write the requested lines to the writer
	if _, err = writer.Write(data); err != nil {
		return err
	}
	// reposition the stream at the end, so the caller can keep reading the file
	// (e.g. when using follow-mode)
	_, err = input.Seek(lastPos, io.SeekStart)
	return err
}

// readLastLinesFromBeginning reads the last N lines from the provided file,
// reading from the beginning of the file and keeping track of the last N lines.
// This is necessary for files that are not seekable (e.g. stdin), but it's less
// efficient. For an efficient alternative that works on seekable files see
// readLastLinesBackwards.
// It returns an error, if any. If no error is encountered, the File object's
// offset is positioned after the last read location.
func readLastLinesFromBeginning(input io.ReadSeeker, writer io.Writer, numLines int) error {
	blkSize := getBlockSize(numLines)
	// read block by block until EOF and store a reference to the last lines
	buf := make([]byte, blkSize)
	var (
		slice      []byte // will hold the final data, after moving line by line
		foundLines int
	)
	for {
		n, err := io.ReadFull(input, buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			if err != io.ErrUnexpectedEOF {
				return err
			}
		}
		// look for newlines and keep a slice starting at the n-th to last line
		// (no further than numLines)
		foundLines += bytes.Count(buf[:n], []byte{'\n'})
		slice = append(slice, buf[:n]...) // this is the slice that points to the wanted lines
		// process the current slice
		slice = lastNLines(slice, numLines)
	}
	if _, err := writer.Write(slice); err != nil {
		return err
	}
	return nil
}

func isTruncated(file *os.File) (bool, error) {
	// current read position in a file
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}
	// file stat to get the size
	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	return currentPos > fileInfo.Size(), nil
}

// tail reads the last N lines from the input File and writes them to the Writer.
// The tailConfig object allows to specify the precise behaviour.
func tail(inFile *os.File, writer io.Writer, config tailConfig) error {
	// try reading from the end of the file
	retryFromBeginning := false
	err := readLastLinesBackwards(inFile, writer, config.numLines)
	if err != nil {
		// if it failed because it couldn't seek, mark it for retry reading from
		// the beginning
		if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == syscall.ESPIPE {
			retryFromBeginning = true
		} else {
			return err
		}
	}
	// if reading backwards failed because the file is not seekable,
	// retry from the beginning
	if retryFromBeginning {
		if err = readLastLinesFromBeginning(inFile, writer, config.numLines); err != nil {
			return err
		}
	}
	if config.follow {
		blkSize := getBlockSize(1)
		// read block by block until EOF and store a reference to the last lines
		buf := make([]byte, blkSize)
		for {
			n, err := inFile.Read(buf)
			if err == io.EOF {
				// without this sleep you would hogg the CPU
				time.Sleep(config.followDuration)
				// truncated ?
				truncated, errTruncated := isTruncated(inFile)
				if errTruncated != nil {
					break
				}
				if truncated {
					// seek from start
					_, errSeekStart := inFile.Seek(0, io.SeekStart)
					if errSeekStart != nil {
						break
					}
				}
				continue
			}
			if err == nil {
				_, err := writer.Write(buf[:n])
				if err != nil {
					return err
				}
				continue
			}
			break
		}
	}
	return nil
}

func run(reader *os.File, writer io.Writer, follow bool, numLines int, followDuration time.Duration, args []string) error {
	var (
		inFile *os.File
		err    error
	)

	// TODO: add support for parsing + (from beggining of the file)
	// negative sign is the same as none
	if numLines < 0 {
		numLines = -1 * numLines
	}
	config := tailConfig{follow: follow, numLines: numLines, followDuration: followDuration}

	switch len(args) {
	case 0:
		return tail(reader, writer, config)
	case 1:
	default:
		if follow {
			return fmt.Errorf("tail: can only read one file at a time if follow true")
		}
	}

	for i, file := range args {
		if len(args) > 1 {
			fmt.Fprintf(writer, "==> %s <==\n", file)
		}

		inFile, err = os.Open(file)
		if err != nil {
			return err
		}

		err := tail(inFile, writer, config)
		if err != nil {
			return err
		}

		if len(args) > 1 && i < len(args)-1 {
			fmt.Fprintln(writer)
		}
	}

	return nil
}

func main() {
	flag.Parse()
	if err := run(os.Stdin, os.Stdout, *flagFollow, *flagNumLines, 500*time.Millisecond, flag.Args()); err != nil {
		log.Fatalf("tail: %v", err)
	}
}

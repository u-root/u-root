package pxe

import (
	"bytes"
	"io"
	"math"
)

// LazyOpener is a lazy io.Reader.
//
// LazyOpener will use a given open function to derive an io.Reader when Read
// is first called on the LazyOpener.
type LazyOpener struct {
	r    io.Reader
	err  error
	open func() (io.Reader, error)
}

// NewLazyOpener returns a lazy io.Reader based on `open`.
func NewLazyOpener(open func() (io.Reader, error)) io.Reader {
	return &LazyOpener{open: open}
}

// Read implements io.Reader.Read lazily.
//
// If called for the first time, the underlying reader will be obtained and
// then used for the first and subsequent calls to Read.
func (lr *LazyOpener) Read(p []byte) (int, error) {
	if lr.r == nil && lr.err == nil {
		lr.r, lr.err = lr.open()
	}
	if lr.err != nil {
		return 0, lr.err
	}
	return lr.r.Read(p)
}

// CachingReader is a lazily caching wrapper of an io.Reader.
//
// The wrapped io.Reader is only read from on demand, not upfront.
type CachingReader struct {
	buf bytes.Buffer
	r   io.Reader
	pos int
	eof bool
}

// NewCachingReader buffers reads from r.
//
// r is only read from when Read() is called.
func NewCachingReader(r io.Reader) *CachingReader {
	return &CachingReader{
		r: r,
	}
}

func (cr *CachingReader) read(p []byte) (int, error) {
	n, err := cr.r.Read(p)
	cr.buf.Write(p[:n])
	if err == io.EOF || (n == 0 && err == nil) {
		cr.eof = true
		return n, io.EOF
	}
	return n, err
}

// NewReader returns a new io.Reader that reads cr from offset 0.
func (cr *CachingReader) NewReader() io.Reader {
	return io.NewSectionReader(cr, 0, math.MaxInt64)
}

// Read reads from cr; implementing io.Reader.
//
// TODO(chrisko): Decide whether to keep this or only keep NewReader().
func (cr *CachingReader) Read(p []byte) (int, error) {
	n, err := cr.ReadAt(p, int64(cr.pos))
	cr.pos += n
	return n, err
}

// ReadAt reads from cr; implementing io.ReaderAt.
func (cr *CachingReader) ReadAt(p []byte, off int64) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	end := int(off) + len(p)

	// Is the caller asking for some uncached bytes?
	unread := end - cr.buf.Len()
	if unread > 0 {
		// Avoiding allocations: use `p` to read more bytes.
		for unread > 0 {
			toRead := unread % len(p)
			if toRead == 0 {
				toRead = len(p)
			}

			m, err := cr.read(p[:toRead])
			unread -= m
			if err == io.EOF {
				break
			}
			if err != nil {
				return 0, err
			}
		}
	}

	// If this is true, the entire file was read just to find out, but the
	// offset is beyond the end of the file.
	if off > int64(cr.buf.Len()) {
		return 0, io.EOF
	}

	var err error
	// Did the caller ask for more than was available?
	//
	// Note that any io.ReaderAt implementation *must* return an error for
	// short reads.
	if cr.eof && unread > 0 {
		err = io.EOF
	}
	return copy(p, cr.buf.Bytes()[off:]), err
}

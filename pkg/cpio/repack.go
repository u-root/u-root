// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/u-root/u-root/pkg/gzip"
	"github.com/ulikunitz/xz"
)

// Archiver represents a CPIO archive (in-memory or on-disk, optionally
// compressed) to which records and files can be written and closed.
type Archiver interface {
	RecordWriter
	AddFile(filePath string, rename func(string) string) error
	ConcatCPIOWithSkip(src string, skip map[string]bool) error
	Symlink(name, target string) error
	Directory(name string, mode uint64) error
	StaticFile(name, content string, mode uint64) error
	Close() error
}

// AddFileToRoot adds filePath to the root (/) of the archive using its base
// name.
func AddFileToRoot(a Archiver, filePath string) error {
	return a.AddFile(filePath, func(oldname string) string {
		return filepath.Base(oldname)
	})
}

// AddFiles appends a list of files from the host filesystem to the archive.
// For each file, leading slashes are stripped so it is stored as a relative
// path in the archive.
func AddFiles(a Archiver, files []string) error {
	for _, f := range files {
		if err := a.AddFile(f, func(name string) string {
			return strings.TrimLeft(name, "/")
		}); err != nil {
			return fmt.Errorf("add file %s: %w", f, err)
		}
	}
	return nil
}

// AddFilesToRoot appends a list of files to the root of the archive using
// their base names.
func AddFilesToRoot(a Archiver, files []string) error {
	for _, f := range files {
		if err := AddFileToRoot(a, f); err != nil {
			return fmt.Errorf("add file %s to root: %w", f, err)
		}
	}
	return nil
}

// ConcatCPIO appends the contents of src to the archive.
func ConcatCPIO(a Archiver, src string) error {
	return a.ConcatCPIOWithSkip(src, nil)
}

// MemArchive is an in-memory CPIO archive implementing Archiver on top of
// Archive.
type MemArchive struct {
	baseArchive
	*Archive
}

var _ Archiver = (*MemArchive)(nil)

// NewMemArchive creates an in-memory CPIO archive.
func NewMemArchive() *MemArchive {
	mem := InMemArchive()
	return &MemArchive{
		baseArchive: newBaseArchive(memRecordWriter{Archive: mem}, ""),
		Archive:     mem,
	}
}

// WriteRecord writes a record to the in-memory archive, automatically creating
// any missing intermediate directory records in archive order.
func (m *MemArchive) WriteRecord(rec Record) error {
	return m.baseArchive.WriteRecord(rec)
}

// Close is a no-op for an in-memory archive.
func (m *MemArchive) Close() error {
	return nil
}

// DiskArchive is an uncompressed on-disk CPIO archive file.
type DiskArchive struct {
	baseArchive
	cpioPath string
	file     *os.File
	closed   bool
}

var _ Archiver = (*DiskArchive)(nil)

// NewDiskArchive creates an uncompressed CPIO archive file at archivePath.
// Call Close() when done writing records.
func NewDiskArchive(archivePath string) (*DiskArchive, error) {
	cf, err := os.Create(archivePath)
	if err != nil {
		return nil, fmt.Errorf("newcpio: %w", err)
	}
	cw := Newc.Writer(cf)
	return &DiskArchive{
		baseArchive: newBaseArchive(cw, filepath.Dir(archivePath)),
		cpioPath:    archivePath,
		file:        cf,
	}, nil
}

func newTempDiskArchive(archivePath string) (*DiskArchive, error) {
	dir := filepath.Dir(archivePath)
	base := filepath.Base(trimCompressionSuffix(archivePath))
	cf, err := os.CreateTemp(dir, ".tmp_cpio_*-"+base)
	if err != nil {
		return nil, fmt.Errorf("newcpio: %w", err)
	}
	cw := Newc.Writer(cf)
	return &DiskArchive{
		baseArchive: newBaseArchive(cw, dir),
		cpioPath:    cf.Name(),
		file:        cf,
	}, nil
}

// Close writes the CPIO trailer record and closes the underlying file.
// Subsequent calls to Close are no-ops.
func (d *DiskArchive) Close() error {
	if d == nil || d.file == nil || d.closed {
		return nil
	}
	d.closed = true
	trailerErr := WriteTrailer(d)
	closeErr := d.file.Close()
	if err := errors.Join(trailerErr, closeErr); err != nil {
		return fmt.Errorf("writeTrailer: %w", err)
	}
	return nil
}

func (d *DiskArchive) abort() {
	if d == nil || d.file == nil || d.closed {
		return
	}
	d.closed = true
	_ = d.file.Close()
	_ = os.Remove(d.cpioPath)
}

// GzipArchive is a gzip-compressed on-disk CPIO archive (.cpio.gz) that embeds
// DiskArchive and compresses the CPIO file into a .gz archive on Close.
type GzipArchive struct {
	*DiskArchive
	gzPath string
}

var _ Archiver = (*GzipArchive)(nil)

// NewGzipArchive creates a gzip-compressed CPIO archive at archivePath (which
// must end with ".gz"). An intermediate uncompressed CPIO file is written in
// the destination directory until Close() is called.
func NewGzipArchive(archivePath string) (*GzipArchive, error) {
	if !strings.HasSuffix(archivePath, ".gz") {
		return nil, fmt.Errorf("gzip archive path %q must end with .gz", archivePath)
	}
	da, err := newTempDiskArchive(archivePath)
	if err != nil {
		return nil, err
	}
	return &GzipArchive{DiskArchive: da, gzPath: archivePath}, nil
}

// Close writes the CPIO trailer record, closes the intermediate CPIO file, and
// compresses it into a .gz file. Subsequent calls to Close are no-ops.
func (g *GzipArchive) Close() error {
	if g == nil || g.DiskArchive == nil || g.closed {
		return nil
	}
	if err := g.DiskArchive.Close(); err != nil {
		_ = os.Remove(g.cpioPath)
		return err
	}
	return compressToFile(g.cpioPath, g.gzPath, func(dst io.Writer, src io.Reader) error {
		return gzip.Compress(src, dst, 6, 32, runtime.NumCPU())
	})
}

// XZArchive is an xz-compressed on-disk CPIO archive (.cpio.xz) that embeds
// DiskArchive and compresses the CPIO file into an .xz archive using CRC32
// checksums (compatible with the Linux kernel's xz decompressor) on Close.
type XZArchive struct {
	*DiskArchive
	xzPath string
}

var _ Archiver = (*XZArchive)(nil)

// NewXZArchive creates an xz-compressed CPIO archive at archivePath (which
// must end with ".xz"). An intermediate uncompressed CPIO file is written in
// the destination directory until Close() is called.
func NewXZArchive(archivePath string) (*XZArchive, error) {
	if !strings.HasSuffix(archivePath, ".xz") {
		return nil, fmt.Errorf("xz archive path %q must end with .xz", archivePath)
	}
	da, err := newTempDiskArchive(archivePath)
	if err != nil {
		return nil, err
	}
	return &XZArchive{DiskArchive: da, xzPath: archivePath}, nil
}

// Close writes the CPIO trailer record, closes the intermediate CPIO file, and
// compresses it into an .xz file using CRC32 checksums. Subsequent calls to
// Close are no-ops.
func (x *XZArchive) Close() error {
	if x == nil || x.DiskArchive == nil || x.closed {
		return nil
	}
	if err := x.DiskArchive.Close(); err != nil {
		_ = os.Remove(x.cpioPath)
		return err
	}
	return compressToFile(x.cpioPath, x.xzPath, func(dst io.Writer, src io.Reader) error {
		w, err := (xz.WriterConfig{CheckSum: xz.CRC32}).NewWriter(dst)
		if err != nil {
			return fmt.Errorf("create xz writer: %w", err)
		}
		_, copyErr := io.Copy(w, src)
		closeErr := w.Close()
		return errors.Join(copyErr, closeErr)
	})
}

// NewArchive creates an on-disk CPIO archive at archivePath, returning a
// GzipArchive if archivePath ends with ".gz", an XZArchive if archivePath ends
// with ".xz", or a DiskArchive otherwise.
//
// Close() the archive when you're done.
func NewArchive(archivePath string) (Archiver, error) {
	switch {
	case strings.HasSuffix(archivePath, ".gz"):
		return NewGzipArchive(archivePath)
	case strings.HasSuffix(archivePath, ".xz"):
		return NewXZArchive(archivePath)
	default:
		return NewDiskArchive(archivePath)
	}
}

// ConcatArchives merges one or more source CPIO archives (.cpio, .cpio.gz, or
// .cpio.xz) into a new archive at dst, stripping intermediate trailer records.
// Records are written in the order of srcs, so when multiple source archives
// contain a record at the same normalized path, the last archive's record
// overwrites earlier ones on unpack.
func ConcatArchives(dst string, srcs ...string) error {
	if len(srcs) == 0 {
		return errors.New("at least one source archive must be specified")
	}
	if dst == "" {
		return errors.New("destination archive (dst) must be specified")
	}
	return Repack(RepackOptions{
		Src:       srcs[0],
		ExtraSrcs: srcs[1:],
		Dst:       dst,
	})
}

// RepackOptions configures a Repack operation that copies a source CPIO (or
// CPIO.gz / CPIO.xz) archive and adds, overwrites, or omits files in it.
type RepackOptions struct {
	// Src is the path to the source CPIO archive (required).
	Src string

	// ExtraSrcs is an optional list of additional CPIO archives to
	// concatenate in order after Src. When multiple source archives
	// contain a record at the same normalized path, later archives
	// overwrite earlier ones on unpack.
	ExtraSrcs []string

	// Dst is the path to the destination CPIO archive. If empty (or equal
	// to Src), Src is updated in place.
	Dst string

	// Files is a list of file specifications to write to the archive. Each
	// entry may be either "<host_src>" or "<host_src>:<archive_dst>". Any
	// existing record in Src or ExtraSrcs at the target archive path is
	// overwritten on unpack.
	Files []string

	// ToRoot indicates whether to place appended files (that do not specify
	// an explicit ":<archive_dst>") in the root of the archive using their
	// base names.
	ToRoot bool

	// Skip lists archive paths to omit when copying records from Src and
	// ExtraSrcs.
	Skip []string

	// ModifyArchive is an optional hook invoked on the destination archive
	// after copying Src and adding Files, before Close.
	ModifyArchive func(Archiver) error
}

// Repack copies a source CPIO (or CPIO.gz / CPIO.xz) archive and any ExtraSrcs
// to Dst (or replaces Src in place when Dst is empty or refers to Src),
// omitting any skipped records and appending the specified Files and
// ModifyArchive changes.
//
// Repack always writes to a temporary file in the destination directory and
// atomically renames it to Dst on success, so reading from Src (or ExtraSrcs)
// is safe even when Dst aliases an input archive.
func Repack(opts RepackOptions) (retErr error) {
	if opts.Src == "" {
		return errors.New("source archive (Src) must be specified")
	}
	srcInfo, err := os.Stat(opts.Src)
	if err != nil {
		return fmt.Errorf("source archive: %w", err)
	}

	dst := opts.Dst
	if dst == "" {
		dst = opts.Src
	}

	targetMode := srcInfo.Mode()
	if dstInfo, err := os.Stat(dst); err == nil {
		targetMode = dstInfo.Mode()
	}

	dir := filepath.Dir(dst)
	base := filepath.Base(dst)
	tmpFile, err := os.CreateTemp(dir, ".tmp_cpio_repack_*-"+base)
	if err != nil {
		return fmt.Errorf("create temp archive: %w", err)
	}
	tmpTarget := tmpFile.Name()
	_ = tmpFile.Close()
	defer func() {
		if retErr != nil {
			_ = os.Remove(tmpTarget)
		}
	}()

	if err := repackToFile(opts, tmpTarget); err != nil {
		return err
	}

	if err := os.Chmod(tmpTarget, targetMode); err != nil {
		return fmt.Errorf("preserve permissions: %w", err)
	}

	if err := os.Rename(tmpTarget, dst); err != nil {
		return fmt.Errorf("replace archive: %w", err)
	}
	return nil
}

// baseArchive implements the shared record-writing, parent-directory creation,
// file addition, and archive concatenation logic embedded by MemArchive and
// DiskArchive.
type baseArchive struct {
	rw RecordWriter
	*Recorder
	writtenDirs map[string]bool
	tempDir     string
}

func newBaseArchive(rw RecordWriter, tempDir string) baseArchive {
	return baseArchive{
		rw:          rw,
		Recorder:    NewRecorder(),
		writtenDirs: make(map[string]bool),
		tempDir:     tempDir,
	}
}

func (b *baseArchive) nextIno() uint64 {
	if b.Recorder == nil {
		b.Recorder = NewRecorder()
	}
	ino := b.Recorder.inumber
	b.Recorder.inumber++
	return ino
}

// WriteRecord writes a record to the archive, automatically creating any
// missing intermediate directory records in archive order.
func (b *baseArchive) WriteRecord(rec Record) error {
	if rec.Name == Trailer {
		return b.rw.WriteRecord(rec)
	}

	if b.writtenDirs == nil {
		b.writtenDirs = make(map[string]bool)
	}
	if b.Recorder == nil {
		b.Recorder = NewRecorder()
	}

	rec.Name = Normalize(rec.Name)
	if rec.Name == Trailer {
		return fmt.Errorf("invalid archive path %q", rec.Name)
	}
	if rec.Name == ".." || strings.HasPrefix(rec.Name, "../") {
		return fmt.Errorf("archive path %q escapes root", rec.Name)
	}
	isDir := rec.Info.Mode&S_IFMT == S_IFDIR
	if rec.Name == "." && !isDir {
		return fmt.Errorf("archive path %q cannot be a non-directory", rec.Name)
	}
	if rec.Info.Ino >= b.Recorder.inumber {
		b.Recorder.inumber = rec.Info.Ino + 1
	}

	dir := path.Dir(rec.Name)
	// Skip root-level files (dir is "." or "/").
	if dir != "." && dir != "/" {
		els := strings.Split(dir, "/")
		for i := range els {
			// Use path.Join (not filepath) so separators are
			// always "/".
			d := path.Join(els[:i+1]...)
			if !b.writtenDirs[d] {
				b.writtenDirs[d] = true
				if err := b.rw.WriteRecord(Directory(d, 0o755)); err != nil {
					return fmt.Errorf("write dir %s: %w", d, err)
				}
			}
		}
	}

	if isDir {
		b.writtenDirs[rec.Name] = true
	} else {
		delete(b.writtenDirs, rec.Name)
	}

	return b.rw.WriteRecord(rec)
}

// AddFile adds the file to the archive, possibly renaming the record.
func (b *baseArchive) AddFile(filePath string, rename func(string) string) error {
	if b.Recorder == nil {
		b.Recorder = NewRecorder()
	}
	rec, err := b.GetRecord(filePath)
	if err != nil {
		return fmt.Errorf("getRecord: %w", err)
	}
	if rename != nil {
		rec.Info.Name = rename(rec.Info.Name)
	}
	return b.WriteRecord(rec)
}

// openDecompressed opens a CPIO file at filePath, decompressing it into an
// unlinked temporary file if filePath ends in .gz or .xz (since CPIO reading
// requires io.ReaderAt).
func openDecompressed(filePath, tempDir string) (_ *os.File, retErr error) {
	if !strings.HasSuffix(filePath, ".gz") && !strings.HasSuffix(filePath, ".xz") {
		return os.Open(filePath)
	}
	src, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.CreateTemp(tempDir, "tmpdecompress-*")
	if err != nil {
		return nil, err
	}
	// dst remains open, but unlink its directory entry immediately.
	os.Remove(dst.Name())
	defer func() {
		if retErr != nil {
			_ = dst.Close()
		}
	}()

	if strings.HasSuffix(filePath, ".gz") {
		if err := gzip.Decompress(src, dst, 32, runtime.NumCPU()); err != nil {
			return nil, err
		}
		return dst, nil
	}

	xzr, err := xz.NewReader(src)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(dst, xzr); err != nil {
		return nil, err
	}
	return dst, nil
}

// skipTrailerPadding advances rd.pos past any NUL padding bytes following a
// TRAILER!!! record and reports whether another newc CPIO stream starts after
// the padding.
func skipTrailerPadding(rd *reader) (bool, error) {
	var buf [512]byte
	for {
		n, err := rd.r.ReadAt(buf[:], rd.pos)
		for i := range n {
			if buf[i] != 0 {
				rd.pos += int64(i)
				var magic [magicLen]byte
				if _, mErr := rd.r.ReadAt(magic[:], rd.pos); mErr != nil {
					Debug("ignoring trailing bytes after TRAILER: %v\n", mErr)
					return false, nil
				}
				if string(magic[:]) != rd.n.magic {
					Debug("ignoring trailing non-cpio data after TRAILER\n")
					return false, nil
				}
				return true, nil
			}
		}
		rd.pos += int64(n)
		if err == io.EOF {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
}

// ConcatCPIOWithSkip appends the contents of src to the archive, omitting any
// records whose normalized archive path is present in skip.
func (b *baseArchive) ConcatCPIOWithSkip(src string, skip map[string]bool) error {
	srcCPIO, err := openDecompressed(src, b.tempDir)
	if err != nil {
		return fmt.Errorf("srccpio: %w", err)
	}
	defer srcCPIO.Close()

	rd := &reader{n: newc{magic: newcMagic}, r: srcCPIO}
	inoMap := make(map[uint64]uint64)
	skippedHardlinks := make(map[uint64]Record)
	afterTrailer := false
	for {
		rec, err := rd.ReadRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read cpio record from %s: %w", src, err)
		}
		if rec.Name == Trailer {
			afterTrailer = true
			clear(inoMap)
			clear(skippedHardlinks)
			more, padErr := skipTrailerPadding(rd)
			if padErr != nil {
				return fmt.Errorf("read padding from %s: %w", src, padErr)
			}
			if !more {
				break
			}
			afterTrailer = false
			continue
		}
		afterTrailer = false
		rec.Name = Normalize(rec.Name)
		origIno := rec.Info.Ino
		if skip[rec.Name] {
			if origIno != 0 && rec.Info.FileSize > 0 && rec.Info.Mode&S_IFMT == S_IFREG {
				skippedHardlinks[origIno] = rec
			}
			continue
		}
		if origIno != 0 {
			if prev, ok := skippedHardlinks[origIno]; ok && rec.Info.FileSize == 0 {
				rec.Info.FileSize = prev.Info.FileSize
				rec.ReaderAt = prev.ReaderAt
				delete(skippedHardlinks, origIno)
			}
			mapped, ok := inoMap[origIno]
			if !ok {
				mapped = b.nextIno()
				inoMap[origIno] = mapped
			}
			rec.Info.Ino = mapped
		}
		if err := b.WriteRecord(rec); err != nil {
			return fmt.Errorf("writing cpio record: %w", err)
		}
	}
	if !afterTrailer {
		return fmt.Errorf("archive %s missing CPIO trailer", src)
	}
	return nil
}

// Symlink adds a symlink to the archive.
func (b *baseArchive) Symlink(name, target string) error {
	return b.WriteRecord(Symlink(name, target))
}

// Directory adds a directory to the archive.
func (b *baseArchive) Directory(name string, mode uint64) error {
	return b.WriteRecord(Directory(name, mode))
}

// StaticFile adds an in-memory file record with the given archive path, string
// content, and permission mode.
func (b *baseArchive) StaticFile(name, content string, mode uint64) error {
	return b.WriteRecord(StaticFile(name, content, mode))
}

// memRecordWriter wraps an in-memory *Archive and eagerly copies record
// payloads into memory so records do not retain references to closed or
// removed backing files.
type memRecordWriter struct {
	*Archive
}

func (w memRecordWriter) WriteRecord(rec Record) error {
	if rec.Name == Trailer {
		return nil
	}
	if rec.ReaderAt != nil {
		data, err := io.ReadAll(io.NewSectionReader(rec.ReaderAt, 0, int64(rec.Info.FileSize)))
		if c, ok := rec.ReaderAt.(io.Closer); ok {
			_ = c.Close()
		}
		if err != nil {
			return fmt.Errorf("read record %s: %w", rec.Name, err)
		}
		if uint64(len(data)) != rec.Info.FileSize {
			return fmt.Errorf("read record %s: %w", rec.Name, io.ErrUnexpectedEOF)
		}
		rec.ReaderAt = bytes.NewReader(data)
		rec.Info.FileSize = uint64(len(data))
	}
	return w.Archive.WriteRecord(rec)
}

func compressToFile(srcPath, dstPath string, compress func(dst io.Writer, src io.Reader) error) (retErr error) {
	defer os.Remove(srcPath)

	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open intermediate cpio: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("create compressed archive: %w", err)
	}
	defer func() {
		if retErr != nil {
			_ = dst.Close()
			_ = os.Remove(dstPath)
		}
	}()

	compErr := compress(dst, src)
	closeErr := dst.Close()
	if err := errors.Join(compErr, closeErr); err != nil {
		return fmt.Errorf("compress archive %s: %w", dstPath, err)
	}
	return nil
}

func trimCompressionSuffix(archivePath string) string {
	switch {
	case strings.HasSuffix(archivePath, ".gz"):
		return strings.TrimSuffix(archivePath, ".gz")
	case strings.HasSuffix(archivePath, ".xz"):
		return strings.TrimSuffix(archivePath, ".xz")
	default:
		return archivePath
	}
}

type fileMapping struct {
	src string
	dst string
}

func isValidArchiveDst(normDst string) bool {
	return normDst != "" && normDst != "." && normDst != ".." && !strings.HasPrefix(normDst, "../") && normDst != Trailer
}

// parseFileSpec parses a file specification ("<host_src>" or
// "<host_src>:<archive_dst>") into its host source path and normalized archive
// destination path.
func parseFileSpec(spec string, toRoot bool) (fileMapping, error) {
	if spec == "" {
		return fileMapping{}, errors.New("empty file specification")
	}
	if _, err := os.Lstat(spec); err != nil {
		src, dst, ok := strings.Cut(spec, ":")
		if !ok {
			return fileMapping{}, fmt.Errorf("stat source file %q: %w", spec, err)
		}
		normDst := Normalize(dst)
		if src == "" || !isValidArchiveDst(normDst) {
			return fileMapping{}, fmt.Errorf("invalid file specification %q: expected <host_src>:<archive_dst>", spec)
		}
		if _, statErr := os.Lstat(src); statErr != nil {
			return fileMapping{}, fmt.Errorf("stat source file %q: %w", src, statErr)
		}
		return fileMapping{
			src: src,
			dst: normDst,
		}, nil
	}
	dst := strings.TrimLeft(spec, "/")
	if toRoot {
		dst = filepath.Base(spec)
	}
	normDst := Normalize(dst)
	if !isValidArchiveDst(normDst) {
		return fileMapping{}, fmt.Errorf("invalid file specification %q", spec)
	}
	return fileMapping{
		src: spec,
		dst: normDst,
	}, nil
}

func repackToFile(opts RepackOptions, dstPath string) (retErr error) {
	skip := make(map[string]bool, len(opts.Skip))
	for _, s := range opts.Skip {
		if norm := Normalize(s); norm != "" && norm != "." {
			skip[norm] = true
		}
	}

	var mappings []fileMapping
	for _, spec := range opts.Files {
		m, err := parseFileSpec(spec, opts.ToRoot)
		if err != nil {
			return err
		}
		mappings = append(mappings, m)
	}

	a, err := NewArchive(dstPath)
	if err != nil {
		return fmt.Errorf("create destination archive: %w", err)
	}
	defer func() {
		if retErr != nil {
			if aborter, ok := a.(interface{ abort() }); ok {
				aborter.abort()
			} else {
				_ = a.Close()
			}
		}
	}()

	if err := a.ConcatCPIOWithSkip(opts.Src, skip); err != nil {
		return fmt.Errorf("concat source cpio: %w", err)
	}

	for _, extra := range opts.ExtraSrcs {
		if err := a.ConcatCPIOWithSkip(extra, skip); err != nil {
			return fmt.Errorf("concat extra cpio %s: %w", extra, err)
		}
	}

	for _, m := range mappings {
		if err := a.AddFile(m.src, func(string) string {
			return m.dst
		}); err != nil {
			return fmt.Errorf("add file %s -> %s: %w", m.src, m.dst, err)
		}
	}

	if opts.ModifyArchive != nil {
		if err := opts.ModifyArchive(a); err != nil {
			return err
		}
	}

	if err := a.Close(); err != nil {
		return fmt.Errorf("close destination archive: %w", err)
	}
	return nil
}

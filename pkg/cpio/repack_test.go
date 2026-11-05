// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpio

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ulikunitz/xz"
)

// TestAddFileIntermediateDirectories checks whether adding a file with a
// nested path automatically creates intermediate directory records in the
// CPIO archive.
func TestAddFileIntermediateDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "test.cpio")

	// Create a nested file: <tmpDir>/sub1/sub2/file.txt.
	subDir := filepath.Join(tmpDir, "sub1", "sub2")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	testFile := filepath.Join(subDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	a, err := NewArchive(archivePath)
	if err != nil {
		t.Fatalf("NewArchive: %v", err)
	}

	// Add file as relative path "sub1/sub2/file.txt".
	if err := a.AddFile(testFile, func(name string) string {
		return "sub1/sub2/file.txt"
	}); err != nil {
		t.Fatalf("AddFile: %v", err)
	}

	if err := a.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	// Read all records from the generated archive.
	f, err := os.Open(archivePath)
	if err != nil {
		t.Fatalf("open archive: %v", err)
	}
	defer f.Close()

	rr := Newc.Reader(f)
	records, err := ReadAllRecords(rr)
	if err != nil {
		t.Fatalf("ReadAllRecords: %v", err)
	}

	var names []string
	recordMap := make(map[string]Record)
	for _, r := range records {
		names = append(names, r.Name)
		recordMap[r.Name] = r
	}

	// Check that intermediate directories were created.
	sub1, ok := recordMap["sub1"]
	if !ok {
		t.Errorf("archive missing intermediate directory 'sub1'; all records: %v", names)
	} else if sub1.Mode&S_IFMT != S_IFDIR {
		t.Errorf("'sub1' mode = %o, want S_IFDIR (%o)", sub1.Mode, S_IFDIR)
	}

	sub2, ok := recordMap["sub1/sub2"]
	if !ok {
		t.Errorf("archive missing intermediate directory 'sub1/sub2'; all records: %v", names)
	} else if sub2.Mode&S_IFMT != S_IFDIR {
		t.Errorf("'sub1/sub2' mode = %o, want S_IFDIR (%o)", sub2.Mode, S_IFDIR)
	}

	if _, ok := recordMap["sub1/sub2/file.txt"]; !ok {
		t.Errorf("archive missing file 'sub1/sub2/file.txt'; all records: %v", names)
	}

	// Verify order: sub1 before sub1/sub2, and sub1/sub2 before
	// sub1/sub2/file.txt.
	sub1Idx := slices.Index(names, "sub1")
	sub2Idx := slices.Index(names, "sub1/sub2")
	fileIdx := slices.Index(names, "sub1/sub2/file.txt")
	if sub1Idx < 0 || sub2Idx < 0 || fileIdx < 0 || sub1Idx > sub2Idx || sub2Idx > fileIdx {
		t.Errorf("incorrect order: sub1 at %d, sub1/sub2 at %d, file at %d", sub1Idx, sub2Idx, fileIdx)
	}
}

// TestAddMultipleFilesDeduplicateDirectories verifies that multiple files
// sharing directory prefixes do not create duplicate directory records.
func TestAddMultipleFilesDeduplicateDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	archivePath := filepath.Join(tmpDir, "test_dedup.cpio")

	f1 := filepath.Join(tmpDir, "f1.txt")
	f2 := filepath.Join(tmpDir, "f2.txt")
	f3 := filepath.Join(tmpDir, "f3.txt")
	for _, f := range []string{f1, f2, f3} {
		if err := os.WriteFile(f, []byte("content"), 0o644); err != nil {
			t.Fatalf("WriteFile %s: %v", f, err)
		}
	}

	a, err := NewArchive(archivePath)
	if err != nil {
		t.Fatalf("NewArchive: %v", err)
	}

	if err := a.AddFile(f1, func(string) string { return "a/b/f1.txt" }); err != nil {
		t.Fatalf("AddFile f1: %v", err)
	}
	if err := a.AddFile(f2, func(string) string { return "a/b/f2.txt" }); err != nil {
		t.Fatalf("AddFile f2: %v", err)
	}
	if err := a.AddFile(f3, func(string) string { return "a/c/f3.txt" }); err != nil {
		t.Fatalf("AddFile f3: %v", err)
	}

	if err := a.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	f, err := os.Open(archivePath)
	if err != nil {
		t.Fatalf("open archive: %v", err)
	}
	defer f.Close()

	rr := Newc.Reader(f)
	records, err := ReadAllRecords(rr)
	if err != nil {
		t.Fatalf("ReadAllRecords: %v", err)
	}

	counts := make(map[string]int)
	for _, r := range records {
		counts[r.Name]++
	}

	for _, dir := range []string{"a", "a/b", "a/c"} {
		if counts[dir] != 1 {
			t.Errorf("directory %q count = %d, want 1", dir, counts[dir])
		}
	}
}

func createBaseArchive(t *testing.T, archivePath string) {
	t.Helper()
	a, err := NewArchive(archivePath)
	if err != nil {
		t.Fatalf("NewArchive(%s): %v", archivePath, err)
	}
	if err := a.StaticFile("base_file.txt", "base content", 0o644); err != nil {
		t.Fatalf("StaticFile: %v", err)
	}
	if err := a.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func readArchiveRecords(t *testing.T, archivePath string) map[string]string {
	t.Helper()
	f, err := os.Open(archivePath)
	if err != nil {
		t.Fatalf("open archive %s: %v", archivePath, err)
	}
	defer f.Close()

	var r io.Reader = f
	switch {
	case strings.HasSuffix(archivePath, ".gz"):
		gzr, err := gzip.NewReader(f)
		if err != nil {
			t.Fatalf("gzip.NewReader(%s): %v", archivePath, err)
		}
		defer gzr.Close()
		r = gzr
	case strings.HasSuffix(archivePath, ".xz"):
		xzr, err := xz.NewReader(f)
		if err != nil {
			t.Fatalf("xz.NewReader(%s): %v", archivePath, err)
		}
		r = xzr
	}

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read cpio data: %v", err)
	}
	rAt := strings.NewReader(string(data))

	records := make(map[string]string)
	ar, err := ArchiveFromReader(Newc.Reader(rAt))
	if err != nil {
		t.Fatalf("ArchiveFromReader(%s): %v", archivePath, err)
	}
	rr := ar.Reader()
	for {
		rec, err := rr.ReadRecord()
		if err == io.EOF || rec.Name == Trailer {
			break
		}
		if err != nil {
			t.Fatalf("ReadRecord(%s): %v", archivePath, err)
		}
		content := make([]byte, rec.Info.FileSize)
		if rec.ReaderAt != nil {
			if _, err := rec.ReaderAt.ReadAt(content, 0); err != nil && err != io.EOF {
				t.Fatalf("ReadAt for %s: %v", rec.Name, err)
			}
		}
		if _, exists := records[rec.Name]; exists {
			t.Errorf("archive %s contains duplicate record in Archive.Reader() for %q", archivePath, rec.Name)
		}
		records[rec.Name] = string(content)
	}
	return records
}

func TestRepackGzippedCPIOAndInPlace(t *testing.T) {
	tmpDir := t.TempDir()
	srcArchive := filepath.Join(tmpDir, "src.cpio.gz")
	dstArchive := filepath.Join(tmpDir, "dst.cpio.gz")
	createBaseArchive(t, srcArchive)
	if err := os.Chmod(srcArchive, 0o640); err != nil {
		t.Fatalf("Chmod: %v", err)
	}

	testFile := filepath.Join(tmpDir, "extra_gz.txt")
	if err := os.WriteFile(testFile, []byte("extra gz content"), 0o644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	if err := Repack(RepackOptions{
		Src:    srcArchive,
		Dst:    dstArchive,
		Files:  []string{testFile},
		ToRoot: true,
	}); err != nil {
		t.Fatalf("Repack gzipped with Dst failed: %v", err)
	}

	records := readArchiveRecords(t, dstArchive)
	if records["base_file.txt"] != "base content" {
		t.Errorf("base_file.txt = %q, want %q", records["base_file.txt"], "base content")
	}
	if records["extra_gz.txt"] != "extra gz content" {
		t.Errorf("extra_gz.txt = %q, want %q", records["extra_gz.txt"], "extra gz content")
	}

	// Test in-place (Dst empty) on .cpio.gz and verify permissions are
	// preserved.
	if err := Repack(RepackOptions{
		Src:    srcArchive,
		Files:  []string{testFile},
		ToRoot: true,
	}); err != nil {
		t.Fatalf("Repack gzipped in-place failed: %v", err)
	}
	fi, err := os.Stat(srcArchive)
	if err != nil {
		t.Fatalf("Stat(%s): %v", srcArchive, err)
	}
	if fi.Mode().Perm() != 0o640 {
		t.Errorf("permissions after in-place Repack = %o, want 0640", fi.Mode().Perm())
	}
	inPlaceRecords := readArchiveRecords(t, srcArchive)
	if inPlaceRecords["extra_gz.txt"] != "extra gz content" {
		t.Errorf("extra_gz.txt in-place = %q, want %q", inPlaceRecords["extra_gz.txt"], "extra gz content")
	}
}

func TestRepackOverwriteAndPathMapping(t *testing.T) {
	tmpDir := t.TempDir()
	archive := filepath.Join(tmpDir, "overwrite.cpio")

	a, err := NewArchive(archive)
	if err != nil {
		t.Fatalf("NewArchive(%s): %v", archive, err)
	}
	if err := a.StaticFile("base_file.txt", "original base", 0o644); err != nil {
		t.Fatalf("StaticFile base_file.txt: %v", err)
	}
	if err := a.Symlink("bbin/uinit", "bb"); err != nil {
		t.Fatalf("Symlink bbin/uinit: %v", err)
	}
	if err := a.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	replacementBase := filepath.Join(tmpDir, "base_file.txt")
	if err := os.WriteFile(replacementBase, []byte("updated base"), 0o644); err != nil {
		t.Fatalf("WriteFile replacementBase: %v", err)
	}
	newUinit := filepath.Join(tmpDir, "custom_uinit")
	if err := os.WriteFile(newUinit, []byte("uinit binary"), 0o755); err != nil {
		t.Fatalf("WriteFile newUinit: %v", err)
	}
	oldVtube := filepath.Join(tmpDir, "vtube_old")
	if err := os.WriteFile(oldVtube, []byte("old vtube binary"), 0o755); err != nil {
		t.Fatalf("WriteFile oldVtube: %v", err)
	}
	newVtube := filepath.Join(tmpDir, "vtube_bin")
	if err := os.WriteFile(newVtube, []byte("vtube binary"), 0o755); err != nil {
		t.Fatalf("WriteFile newVtube: %v", err)
	}

	if err := Repack(RepackOptions{
		Src:    archive,
		ToRoot: true,
		Files: []string{
			replacementBase,
			newUinit + ":/bbin/uinit",
			oldVtube + ":bin/vtube",
			newVtube + ":/bin/vtube",
		},
	}); err != nil {
		t.Fatalf("Repack overwrite failed: %v", err)
	}

	if err := Repack(RepackOptions{
		Src:   archive,
		Files: []string{filepath.Join(tmpDir, "nonexistent_file.txt")},
	}); err == nil {
		t.Errorf("Repack with nonexistent source file succeeded, want error")
	}

	records := readArchiveRecords(t, archive)
	if records["base_file.txt"] != "updated base" {
		t.Errorf("base_file.txt = %q, want %q", records["base_file.txt"], "updated base")
	}
	if records["bbin/uinit"] != "uinit binary" {
		t.Errorf("bbin/uinit = %q, want %q", records["bbin/uinit"], "uinit binary")
	}
	if records["bin/vtube"] != "vtube binary" {
		t.Errorf("bin/vtube = %q, want %q", records["bin/vtube"], "vtube binary")
	}

	// Verify that unpacking the archive to disk with CreateFileInRoot
	// overwrites colliding paths with the latest record.
	unpackRoot := t.TempDir()
	data, err := os.ReadFile(archive)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", archive, err)
	}
	rr := Newc.Reader(strings.NewReader(string(data)))
	for {
		rec, err := rr.ReadRecord()
		if err == io.EOF || rec.Name == Trailer {
			break
		}
		if err != nil {
			t.Fatalf("ReadRecord: %v", err)
		}
		if err := CreateFileInRoot(rec, unpackRoot, false); err != nil {
			t.Fatalf("CreateFileInRoot(%q): %v", rec.Name, err)
		}
	}
	for path, want := range map[string]string{
		"base_file.txt": "updated base",
		"bbin/uinit":    "uinit binary",
		"bin/vtube":     "vtube binary",
	} {
		got, err := os.ReadFile(filepath.Join(unpackRoot, path))
		if err != nil {
			t.Fatalf("ReadFile(%s): %v", path, err)
		}
		if string(got) != want {
			t.Errorf("unpacked %s = %q, want %q", path, got, want)
		}
	}
}

func TestConcatArchives_XZAndGzip(t *testing.T) {
	tmpDir := t.TempDir()
	xzArchive := filepath.Join(tmpDir, "test.cpio.xz")
	createBaseArchive(t, xzArchive)

	extraArchive := filepath.Join(tmpDir, "extra.cpio.gz")
	ea, err := NewArchive(extraArchive)
	if err != nil {
		t.Fatalf("NewArchive(%s): %v", extraArchive, err)
	}
	if err := ea.StaticFile("base_file.txt", "overwritten by extra", 0o644); err != nil {
		t.Fatalf("StaticFile base_file.txt: %v", err)
	}
	if err := ea.StaticFile("dir1/dir2/extra.txt", "from extra", 0o644); err != nil {
		t.Fatalf("StaticFile: %v", err)
	}
	if err := ea.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := ea.Close(); err != nil {
		t.Errorf("second Close: %v", err)
	}

	mergedXZ := filepath.Join(tmpDir, "merged.cpio.xz")
	if err := ConcatArchives(mergedXZ, xzArchive, extraArchive); err != nil {
		t.Fatalf("ConcatArchives failed: %v", err)
	}
	got := readArchiveRecords(t, mergedXZ)
	want := map[string]string{
		"base_file.txt":       "overwritten by extra",
		"dir1":                "",
		"dir1/dir2":           "",
		"dir1/dir2/extra.txt": "from extra",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ConcatArchives records diff (-want +got):\n%s", diff)
	}
}

func TestRepack_ExtraSrcsAndIntermediatePathCollision(t *testing.T) {
	tmpDir := t.TempDir()
	bareArchive := filepath.Join(tmpDir, "bare.cpio")
	createBaseArchive(t, bareArchive)

	extraArchive := filepath.Join(tmpDir, "extra.cpio.gz")
	ea, err := NewArchive(extraArchive)
	if err != nil {
		t.Fatalf("NewArchive(%s): %v", extraArchive, err)
	}
	if err := ea.StaticFile("dir1/dir2/extra.txt", "from extra", 0o644); err != nil {
		t.Fatalf("StaticFile: %v", err)
	}
	if err := ea.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	bareGZ := bareArchive + ".gz"
	overrideFile := filepath.Join(tmpDir, "override.txt")
	if err := os.WriteFile(overrideFile, []byte("overridden extra"), 0o644); err != nil {
		t.Fatalf("WriteFile overrideFile: %v", err)
	}
	if err := Repack(RepackOptions{
		Src:       bareArchive,
		ExtraSrcs: []string{extraArchive},
		Dst:       bareGZ,
		Files:     []string{overrideFile + ":dir1/dir2/extra.txt"},
	}); err != nil {
		t.Fatalf("Repack bare.cpio -> bare.cpio.gz with ExtraSrcs failed: %v", err)
	}
	got := readArchiveRecords(t, bareGZ)
	want := map[string]string{
		"base_file.txt":       "base content",
		"dir1":                "",
		"dir1/dir2":           "",
		"dir1/dir2/extra.txt": "overridden extra",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Repack records diff (-want +got):\n%s", diff)
	}

	// Verify ConcatArchives when Dst collides with an entry in ExtraSrcs.
	if err := ConcatArchives(extraArchive, bareArchive, extraArchive); err != nil {
		t.Fatalf("ConcatArchives overwriting ExtraSrcs[0] failed: %v", err)
	}
	gotExtra := readArchiveRecords(t, extraArchive)
	wantExtra := map[string]string{
		"base_file.txt":       "base content",
		"dir1":                "",
		"dir1/dir2":           "",
		"dir1/dir2/extra.txt": "from extra",
	}
	if diff := cmp.Diff(wantExtra, gotExtra); diff != "" {
		t.Errorf("ConcatArchives(extraArchive) records diff (-want +got):\n%s", diff)
	}
}

func TestMemArchive(t *testing.T) {
	tmpDir := t.TempDir()
	srcArchive := filepath.Join(tmpDir, "src.cpio.xz")
	createBaseArchive(t, srcArchive)

	ephemeralFile := filepath.Join(tmpDir, "ephemeral.txt")
	if err := os.WriteFile(ephemeralFile, []byte("from disk before unlink"), 0o644); err != nil {
		t.Fatalf("WriteFile ephemeralFile: %v", err)
	}

	mem := NewMemArchive()
	if err := ConcatCPIO(mem, srcArchive); err != nil {
		t.Fatalf("ConcatCPIO(mem, %s): %v", srcArchive, err)
	}
	if err := AddFileToRoot(mem, ephemeralFile); err != nil {
		t.Fatalf("AddFileToRoot(mem, %s): %v", ephemeralFile, err)
	}
	if err := os.Remove(ephemeralFile); err != nil {
		t.Fatalf("Remove(%s): %v", ephemeralFile, err)
	}
	if err := mem.StaticFile("nested/in/mem.txt", "initial mem content", 0o644); err != nil {
		t.Fatalf("mem.StaticFile: %v", err)
	}
	// Later write at the same path should overwrite the earlier record
	// (last write wins) without duplicating Order.
	if err := mem.StaticFile("nested/in/mem.txt", "mem content", 0o644); err != nil {
		t.Fatalf("mem.StaticFile overwrite: %v", err)
	}
	if err := WriteTrailer(mem); err != nil {
		t.Fatalf("WriteTrailer(mem): %v", err)
	}
	if err := mem.Close(); err != nil {
		t.Fatalf("mem.Close: %v", err)
	}

	recs, err := ReadAllRecords(mem.Reader())
	if err != nil {
		t.Fatalf("ReadAllRecords(mem.Reader()): %v", err)
	}
	got := make(map[string]string, len(recs))
	counts := make(map[string]int, len(recs))
	for _, rec := range recs {
		counts[rec.Name]++
		buf := make([]byte, rec.Info.FileSize)
		if rec.ReaderAt != nil {
			if _, err := rec.ReaderAt.ReadAt(buf, 0); err != nil && err != io.EOF {
				t.Fatalf("ReadAt(%q) failed: %v", rec.Name, err)
			}
		}
		got[rec.Name] = string(buf)
	}
	want := map[string]string{
		"base_file.txt":     "base content",
		"ephemeral.txt":     "from disk before unlink",
		"nested":            "",
		"nested/in":         "",
		"nested/in/mem.txt": "mem content",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("MemArchive records diff (-want +got):\n%s", diff)
	}
	if counts["nested/in/mem.txt"] != 1 {
		t.Errorf("counts[\"nested/in/mem.txt\"] = %d, want 1", counts["nested/in/mem.txt"])
	}
}

func TestArchiveErrorsAndCloseCleanup(t *testing.T) {
	tmpDir := t.TempDir()
	if _, err := NewGzipArchive(filepath.Join(tmpDir, "bad.cpio")); err == nil {
		t.Errorf("NewGzipArchive without .gz suffix succeeded, want error")
	}
	if _, err := NewXZArchive(filepath.Join(tmpDir, "bad.cpio")); err == nil {
		t.Errorf("NewXZArchive without .xz suffix succeeded, want error")
	}
	if err := ConcatArchives(filepath.Join(tmpDir, "out.cpio")); err == nil {
		t.Errorf("ConcatArchives with no srcs succeeded, want error")
	}

	// Closing a GzipArchive should compress it, remove its intermediate
	// temp file, and leave any sibling .cpio file untouched.
	closedGZ := filepath.Join(tmpDir, "closed.cpio.gz")
	siblingCPIO := strings.TrimSuffix(closedGZ, ".gz")
	if err := os.WriteFile(siblingCPIO, []byte("existing sibling"), 0o644); err != nil {
		t.Fatalf("WriteFile siblingCPIO: %v", err)
	}
	gza, err := NewGzipArchive(closedGZ)
	if err != nil {
		t.Fatalf("NewGzipArchive: %v", err)
	}
	intermediatePath := gza.cpioPath
	if err := gza.Close(); err != nil {
		t.Fatalf("gza.Close(): %v", err)
	}
	if _, err := os.Stat(intermediatePath); !os.IsNotExist(err) {
		t.Errorf("intermediate .cpio file %s still exists after Close(): %v", intermediatePath, err)
	}
	if data, err := os.ReadFile(siblingCPIO); err != nil || string(data) != "existing sibling" {
		t.Errorf("sibling .cpio file modified or removed: data=%q, err=%v", data, err)
	}

	// Corrupt source CPIO should return an error from ConcatCPIO.
	corruptFile := filepath.Join(tmpDir, "corrupt.cpio")
	if err := os.WriteFile(corruptFile, []byte("not a valid cpio header"), 0o644); err != nil {
		t.Fatalf("WriteFile corruptFile: %v", err)
	}
	mem := NewMemArchive()
	if err := ConcatCPIO(mem, corruptFile); err == nil {
		t.Errorf("ConcatCPIO on corrupt file succeeded, want error")
	}
	if err := mem.StaticFile("../../etc/passwd", "bad", 0o644); err == nil {
		t.Errorf("StaticFile with path traversal succeeded, want error")
	}
	if err := mem.StaticFile("/TRAILER!!!", "bad", 0o644); err == nil {
		t.Errorf("StaticFile with /TRAILER!!! succeeded, want error")
	}
	if err := mem.StaticFile(".", "bad", 0o644); err == nil {
		t.Errorf("StaticFile with '.' succeeded, want error")
	}
}

func TestRepackPaddedConcatenatedAndHardlinks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two archives and concatenate them on disk with 512-byte NUL
	// padding between the first archive's TRAILER!!! and the second
	// archive's header.
	part1 := filepath.Join(tmpDir, "part1.cpio")
	part2 := filepath.Join(tmpDir, "part2.cpio")
	createBaseArchive(t, part1)

	hlFile1 := filepath.Join(tmpDir, "hl1.txt")
	hlFile2 := filepath.Join(tmpDir, "hl2.txt")
	if err := os.WriteFile(hlFile1, []byte("shared hardlink payload"), 0o644); err != nil {
		t.Fatalf("WriteFile hlFile1: %v", err)
	}
	if err := os.Link(hlFile1, hlFile2); err != nil {
		t.Fatalf("Link: %v", err)
	}

	a2, err := NewArchive(part2)
	if err != nil {
		t.Fatalf("NewArchive(part2): %v", err)
	}
	if err := a2.AddFile(hlFile1, func(string) string { return "bin/link1" }); err != nil {
		t.Fatalf("AddFile link1: %v", err)
	}
	if err := a2.AddFile(hlFile2, func(string) string { return "bin/link2" }); err != nil {
		t.Fatalf("AddFile link2: %v", err)
	}
	if err := a2.Close(); err != nil {
		t.Fatalf("a2.Close: %v", err)
	}

	b1, err := os.ReadFile(part1)
	if err != nil {
		t.Fatalf("ReadFile(part1): %v", err)
	}
	b2, err := os.ReadFile(part2)
	if err != nil {
		t.Fatalf("ReadFile(part2): %v", err)
	}
	var combined []byte
	combined = append(combined, b1...)
	combined = append(combined, make([]byte, 512)...)
	combined = append(combined, b2...)

	concatSrc := filepath.Join(tmpDir, "concat_padded.cpio")
	if err := os.WriteFile(concatSrc, combined, 0o644); err != nil {
		t.Fatalf("WriteFile(concatSrc): %v", err)
	}

	// Repack concatSrc while skipping the first hardlink ("bin/link1") and
	// verify that both streams are read across the NUL padding and that
	// the surviving hardlink ("bin/link2") retains the shared payload.
	dst := filepath.Join(tmpDir, "out.cpio")
	if err := Repack(RepackOptions{
		Src:  concatSrc,
		Dst:  dst,
		Skip: []string{"bin/link1"},
	}); err != nil {
		t.Fatalf("Repack failed: %v", err)
	}

	got := readArchiveRecords(t, dst)
	want := map[string]string{
		"base_file.txt": "base content",
		"bin":           "",
		"bin/link2":     "shared hardlink payload",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Repack records diff (-want +got):\n%s", diff)
	}
}

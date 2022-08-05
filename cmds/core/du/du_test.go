package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestDu(t *testing.T) {
	var want, got bytes.Buffer

	t.Run("directory containing files with different sizes", func(t *testing.T) {
		defer cleanup(&want, &got)
		fp, dirPath, err := setupDifferentFileSizes(t)

		if err != nil {
			t.Fatalf("setup failed")
		}

		fillBuffer(&want, *fp)

		du(&got, []string{dirPath})

		assertEqual(t, &got, &want)
	})

	t.Run("directory with multiple layers", func(t *testing.T) {
		defer cleanup(&want, &got)

		fp, dirPath, err := setupMultipleLayers(t)
		if err != nil {
			t.Fatalf("failed to create directory structur")
		}

		fillBuffer(&want, *fp)
		du(&got, []string{dirPath})

		assertEqual(t, &got, &want)
	})
}

// setup writes a set of files in a directory, putting x bytes (1 <= x <= len(data)) in each file.
//
func setupDifferentFileSizes(t *testing.T) (*[]FileProperties, string, error) {
	var fp []FileProperties

	rand.Seed(time.Now().Unix())
	dataSize := rand.Intn(0xFFFF)
	data := make([]byte, dataSize)
	_, err := rand.Read(data)

	if err != nil {
		return nil, "", err
	}

	dir := t.TempDir()

	i, j := 0, 1
	for i < len(data) {
		sizeOfFile := 1 + rand.Intn(len(data)-i)
		filePath := fmt.Sprintf("%v%d", filepath.Join(dir, "file"), j)
		if err := os.WriteFile(filePath, data[i:i+sizeOfFile], 0o666); err != nil {
			return nil, "", err
		}
		fp = append(fp, FileProperties{filePath, int64(sizeOfFile)})
		i += sizeOfFile
		j++
	}

	sort.Slice(fp, func(i, j int) bool {
		return sort.StringsAreSorted([]string{fp[i].path, fp[j].path})
	})

	dirInfo, err := os.Stat(dir)

	if err != nil {
		return nil, "", nil
	}
	fp = append(fp, FileProperties{dir, int64(dirInfo.Size())})
	return &fp, dir, nil
}

func setupMultipleLayers(t *testing.T) (*[]FileProperties, string, error) {
	var fp []FileProperties

	dirPath, err := os.MkdirTemp(os.TempDir(), "top")
	if err != nil {
		return nil, "", err
	}
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return nil, "", err
	}

	subDirPath1, err := os.MkdirTemp(dirPath, "sub1")
	if err != nil {
		return nil, "", err
	}
	subDirInfo1, err := os.Stat(subDirPath1)
	if err != nil {
		return nil, "", err
	}

	subDirPath2, err := os.MkdirTemp(subDirPath1, "sub2")
	if err != nil {
		return nil, "", err
	}
	subDirInfo2, err := os.Stat(subDirPath2)
	if err != nil {
		return nil, "", err
	}

	file1, err := os.CreateTemp(dirPath, "fileTop")
	if err != nil {
		return nil, "", err
	}
	fileInfo1, err := file1.Stat()
	if err != nil {
		return nil, "", err
	}

	file2, err := os.CreateTemp(subDirPath1, "fileSub1")
	if err != nil {
		return nil, "", err
	}

	fileInfo2, err := file2.Stat()
	if err != nil {
		return nil, "", err
	}

	file3, err := os.CreateTemp(subDirPath2, "fileSub2")
	if err != nil {
		return nil, "", err
	}

	fileInfo3, err := file3.Stat()
	if err != nil {
		return nil, "", err
	}
	fp = append(fp, FileProperties{fileInfo1.Name(), fileInfo1.Size()})
	fp = append(fp, FileProperties{fileInfo2.Name(), fileInfo2.Size()})
	fp = append(fp, FileProperties{fileInfo3.Name(), fileInfo3.Size()})
	fp = append(fp, FileProperties{subDirPath1, subDirInfo1.Size()})
	fp = append(fp, FileProperties{subDirPath2, subDirInfo2.Size()})
	fp = append(fp, FileProperties{dirPath, dirInfo.Size()})
	return &fp, dirPath, nil
}

func fillBuffer(w io.Writer, fps []FileProperties) {
	for _, fp := range fps {
		fmt.Fprintf(w, "%d %v\n", fp.byteSize, fp.path)
	}
}

func cleanup(want *bytes.Buffer, got *bytes.Buffer) {
	want.Reset()
	got.Reset()
}

func assertEqual(t *testing.T, got, want *bytes.Buffer) {
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want\n%vgot\n%v", want, got)
	}
}

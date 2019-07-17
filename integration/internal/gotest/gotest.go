package gotest

import (
	"os"
	"path/filepath"
	"strings"
)

func WalkTests(testRoot string, fn func(i int, path string, pkgName string)) error {
	var i int
	return filepath.Walk(testRoot, func(path string, info os.FileInfo, err error) error {
		if !info.Mode().IsRegular() || !strings.HasSuffix(path, ".test") {
			return nil
		}
		t2, err := filepath.Rel(testRoot, path)
		if err != nil {
			return err
		}
		pkgName := filepath.Dir(t2)

		fn(i, path, pkgName)
		i++
		return nil
	})
}

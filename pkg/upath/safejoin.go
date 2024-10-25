// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package upath

import (
	"fmt"
	"path/filepath"
	"strings"
)

// SafeFilepathJoin safely joins two paths path1+path2. The resulting path will
// always be contained within path1 even if path2 tries to escape with "../".
// If that path is not possible, an error is returned. The resulting path is
// cleaned.
func SafeFilepathJoin(path1, path2 string) (string, error) {
	relPath, err := filepath.Rel(".", path2)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("(zipslip) filepath is unsafe %q: %w", path2, err)
	}
	if path1 == "" {
		path1 = "."
	}
	return filepath.Join(path1, filepath.Join(string(filepath.Separator), relPath)), nil
}

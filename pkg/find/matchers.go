// Copyright 2015-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package find

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
)

// WithBasenameMatch sets up a file base name filter.
func WithBasenameMatch(pattern string) Set {
	return func(f *Finder) {
		f.matchers = append(f.matchers, func(_ string, fe fs.DirEntry) (bool, error) {
			return fe.Name() == pattern, nil
		})
	}
}

// WithRegexPathMatch sets up a path filter using regex.
//
// This design compiles the regex for every file.
// It also has no good way to handle a bad pattern.
// We have to leave it here for backwards compatibility.
func WithRegexPathMatch(pattern string) Set {
	return func(f *Finder) {
		f.matchers = append(f.matchers, func(_ string, fe fs.DirEntry) (bool, error) {
			return regexp.Match(pattern, []byte(fe.Name()))
		})
	}
}

// WithCompiledRegexPathMatch sets up a path filter using a compiled regex.
// It will return an error if the pattern can not be compiled.
func WithCompiledRegexPathMatch(pattern string) (Set, error) {
	re, err := regexp.Compile(pattern)

	if err != nil {
		return nil, fmt.Errorf("%s:%w:%w", pattern, ErrInvalidRegexp, err)
	}

	return func(f *Finder) {
		f.debug("compiled %q to %v", pattern, re)
		f.matchers = append(f.matchers, func(_ string, fe fs.DirEntry) (bool, error) {
			m := re.MatchString(fe.Name())
			f.debug("WithCompiledRegexPathMatch: pattern %q name %q match %v", pattern, fe.Name(), m)
			return m, nil
		})
	}, nil
}

// WithFilenameMatch uses filepath.Match's shell file name matching to filter
// file base names.
func WithFilenameMatch(pattern string) Set {
	return func(f *Finder) {
		pattern = filepath.Base(pattern)
		f.matchers = append(f.matchers, func(_ string, fe fs.DirEntry) (bool, error) {
			m, err := filepath.Match(pattern, fe.Name())
			f.debug("WithFileNameMatch: pattern %q name %q match %v", pattern, fe.Name(), m)
			return m, err

		})
	}
}

// WithModeMatch ensures only files with fileMode & modeMask == mode are returned.
func WithModeMatch(mode, modeMask fs.FileMode) Set {
	return func(f *Finder) {
		f.matchers = append(f.matchers, func(_ string, fe fs.DirEntry) (bool, error) {
			fi, err := fe.Info()
			if err != nil {
				return false, err
			}
			m := fi.Mode()
			f.debug("%s: file mode %v / want mode %s with mask %s", fe.Name(), m, mode, modeMask)
			if masked := m & modeMask; masked != mode {
				f.debug("%s: mode %s (masked %s) does not match expected mode %s", fe.Name(), m, masked, mode)
				return false, nil
			}
			return true, nil
		})
	}
}

// WithDebugLog logs messages to l.
func WithDebugLog(l func(string, ...any)) Set {
	return func(f *Finder) {
		f.debug = l
	}
}


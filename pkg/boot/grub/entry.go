// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grub

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// GRUB block size.
const blockSize = 1024

// EnvFile is a GRUB environment file consisting of key-value pairs akin to the
// GRUB commands load_env and save_env.
type EnvFile struct {
	Vars map[string]string
}

// NewEnvFile allocates a new env file.
func NewEnvFile() *EnvFile {
	return &EnvFile{
		Vars: make(map[string]string),
	}
}

// WriteTo writes key-value pairs to a file, padded to 1024 bytes, as save_env does.
func (env *EnvFile) WriteTo(w io.Writer) (int64, error) {
	var b bytes.Buffer
	b.WriteString("# GRUB Environment Block\n")

	// Sort keys so order is deterministic.
	keys := make([]string, 0, len(env.Vars))
	for k := range env.Vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if len(env.Vars[k]) > 0 {
			b.WriteString(k)
			b.WriteString("=")
			b.WriteString(env.Vars[k])
			b.WriteString("\n")
		}
	}
	length := b.Len()

	// Fill up the file with # until 1024 bytes to make the file size a
	// multiple of the block size.
	remainder := blockSize - length%blockSize
	for i := 0; i < remainder; i++ {
		b.WriteByte('#')
	}

	return b.WriteTo(w)
}

// ParseEnvFile reads a key-value pair GRUB environment file.
//
// ParseEnvFile accepts incorrectly padded GRUB env files, as opposed to GRUB.
func ParseEnvFile(r io.Reader) (*EnvFile, error) {
	s := bufio.NewScanner(r)
	conf := NewEnvFile()

	replacer := strings.NewReplacer(
		"\n", "",
		"\r", "",
	)
	// Best lexer & parser in the world.
	for s.Scan() {
		if len(s.Text()) == 0 {
			continue
		}
		cleanedText := replacer.Replace(s.Text())
		if len(cleanedText) == 0 {
			continue
		}

		// Comments.
		if cleanedText[0] == '#' {
			continue
		}

		tokens := strings.SplitN(cleanedText, "=", 2)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("error parsing %q: must find = or # and key + values in each line", s.Text())
		}

		if tokens[0] == "" || tokens[1] == "" {
			return nil, fmt.Errorf("error parsing %q: either the key or value is empty: %q = %q", s.Text(), tokens[0], tokens[1])
		}

		key, value := tokens[0], tokens[1]
		conf.Vars[key] = value
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return conf, nil
}

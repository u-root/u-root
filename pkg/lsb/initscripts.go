// Copyright 2025 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lsb implements parsing, marshaling, and manipulation of LSB-compliant
// init script metadata blocks, which are used to define dependencies, run levels,
// and other operational properties for initialization scripts in Unix-like systems.
package lsb

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

const (
	blockStartMarker = "### BEGIN INIT INFO"
	blockEndMarker   = "### END INIT INFO"
)

var (
	ErrMarkerMissing = errors.New("lsb block marker missing")
	ErrNilData       = errors.New("data cannot be nil")
)

type (
	RunLevel     = uint8
	BootFacility = string
)

type InitScript struct {
	// Meta
	Provides         string `lsb:"Provides"`
	ShortDescription string `lsb:"Short-Description"`
	Description      string `lsb:"Description"`

	// Operational
	DefaultStart  []RunLevel     `lsb:"Default-Start"`
	DefaultStop   []RunLevel     `lsb:"Default-Stop"`
	RequiredStart []BootFacility `lsb:"Required-Start"`
	RequiredStop  []BootFacility `lsb:"Required-Stop"`
	ShouldStart   []BootFacility `lsb:"Should-Start"`
	ShouldStop    []BootFacility `lsb:"Should-Stop"`

	// Extension
	XStartBefore []BootFacility `lsb:"X-Start-Before"`
	XStopAfter   []BootFacility `lsb:"X-Stop-After"`
	XInteractive bool           `lsb:"X-Interactive"`
}

// Marshal encodes an InitScript into its string representation.
func (s *InitScript) Marshal() (string, error) {
	var sb strings.Builder
	sb.WriteString(blockStartMarker + "\n")

	typeOfScript := reflect.TypeOf(*s)
	valueOfScript := reflect.ValueOf(*s)

	for i := 0; i < typeOfScript.NumField(); i++ {
		field := typeOfScript.Field(i)
		tag := field.Tag.Get("lsb")
		if tag == "" {
			continue
		}
		value := valueOfScript.Field(i).Interface()
		switch v := value.(type) {
		case string:
			if v != "" {
				sb.WriteString(fmt.Sprintf("# %s: %s\n", tag, v))
			}
		case []uint8:
			if len(v) > 0 {
				items := make([]string, len(v))
				for i, item := range v {
					items[i] = strconv.Itoa(int(item))
				}
				sb.WriteString(fmt.Sprintf("# %s: %s\n", tag, strings.Join(items, " ")))
			}
		case []string:
			if len(v) > 0 {
				sb.WriteString(fmt.Sprintf("# %s: %s\n", tag, strings.Join(v, " ")))
			}
		case bool:
			if v {
				sb.WriteString(fmt.Sprintf("# %s: %t\n", tag, v))
			}
		}
	}

	sb.WriteString(blockEndMarker + "\n")
	return sb.String(), nil
}

// Unmarshal decodes a string into an InitScript.
func (s *InitScript) Unmarshal(data io.Reader) error {
	if data == nil {
		return ErrNilData
	}
	scanner := bufio.NewScanner(data)
	inBlock := false

	typeOfScript := reflect.TypeOf(*s)
	valueOfScript := reflect.ValueOf(s).Elem()

	var (
		foundStart = false
		foundEnd   = false
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == blockStartMarker {
			inBlock = true
			foundStart = true
			continue
		} else if line == blockEndMarker {
			foundEnd = true
			break
		}

		if !inBlock || !strings.HasPrefix(line, "#") {
			continue
		}

		line = strings.TrimPrefix(line, "#")
		line = strings.TrimSpace(line)
		colonIndex := strings.Index(line, ":")
		if colonIndex == -1 {
			continue
		}

		tag := strings.TrimSpace(line[:colonIndex])
		value := strings.TrimSpace(line[colonIndex+1:])

		for i := 0; i < typeOfScript.NumField(); i++ {
			field := typeOfScript.Field(i)
			if field.Tag.Get("lsb") == tag {
				fieldValue := valueOfScript.Field(i)
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(value)
				case reflect.Slice:
					if field.Type.Elem().Kind() == reflect.Uint8 {
						items := strings.Fields(value)
						uintValues := make([]uint8, len(items))
						for j, item := range items {
							parsed, err := strconv.ParseUint(item, 10, 8)
							if err != nil {
								return err
							}
							uintValues[j] = uint8(parsed)
						}
						fieldValue.Set(reflect.ValueOf(uintValues))
					} else {
						fieldValue.Set(reflect.ValueOf(strings.Fields(value)))
					}
				case reflect.Bool:
					parsed, err := strconv.ParseBool(value)
					if err != nil {
						return err
					}
					fieldValue.SetBool(parsed)
				}
				break
			}
		}
	}

	var reterr error
	if !foundStart {
		reterr = errors.Join(reterr, fmt.Errorf("%w: %q", ErrMarkerMissing, blockStartMarker))
	}
	if !foundEnd {
		reterr = errors.Join(reterr, fmt.Errorf("%w: %q", ErrMarkerMissing, blockEndMarker))
	}

	if err := scanner.Err(); err != nil {
		return errors.Join(reterr, err)
	}
	return reterr
}

// SequenceNumber calculates sequence number. It prioritizes scripts with fewer dependencies first.
// Required-Start and Should-Start dependencies are weighted heavily. Higher sequence numbers
// indicate later execution.
func (s *InitScript) SequenceNumber() int {
	// Base sequence number (default priority)
	sequence := 50

	// Adjust based on dependencies
	if len(s.RequiredStart) > 0 {
		sequence -= 20 // Higher priority for scripts with explicit dependencies
	}
	if len(s.ShouldStart) > 0 {
		sequence -= 10 // Medium priority for suggested dependencies
	}

	// Ensure sequence number stays within valid range
	if sequence < 1 {
		sequence = 1
	} else if sequence > 99 {
		sequence = 99
	}

	return sequence
}

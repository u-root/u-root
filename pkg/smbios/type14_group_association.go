// Copyright 2016-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// GroupAssociation is defined in DSP0134 7.24.
type GroupAssociation struct {
	Header
	GroupName  string      // 04h
	ItemType   []TableType `smbios:"-"` // 05h
	ItemHandle []uint16    `smbios:"-"` // 06h
}

// MarshalBinary encodes the BaseboardInfo content into a binary
func (ga *GroupAssociation) MarshalBinary() ([]byte, error) {
	t, err := ga.toTable()
	if err != nil {
		return nil, err
	}
	return t.MarshalBinary()
}

func (ga *GroupAssociation) toTable() (*Table, error) {
	h, err := ga.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var d []byte
	var tableStr []string
	id := byte(1)

	d = append(d, h...)
	if ga.GroupName != "" {
		d = append(d, id)
		tableStr = append(tableStr, ga.GroupName)
	} else {
		d = append(d, 0)
	}

	if len(ga.ItemType) != len(ga.ItemHandle) {
		return nil, fmt.Errorf("item type and item handle have different lengths, len of ItemType: %d, len of ItemHandle: %d", len(ga.ItemType), len(ga.ItemHandle))
	}

	if ga.Length != uint8((len(ga.ItemType)*3)+5) {
		return nil, fmt.Errorf("invalid length, length: %d, len of ItemType: %d, len of ItemHandle: %d", ga.Length, len(ga.ItemType), len(ga.ItemHandle))
	}

	for i, it := range ga.ItemType {
		d = append(d, byte(it))
		ih := ga.ItemHandle[i]
		ihb := make([]byte, 2)
		binary.LittleEndian.PutUint16(ihb, ih)
		d = append(d, ihb...)
	}

	t := &Table{
		Header:  ga.Header,
		data:    d,
		strings: tableStr,
	}
	return t, nil
}

// ParseGroupAssociation parses a generic Table into GroupAssociation.
func ParseGroupAssociation(t *Table) (*GroupAssociation, error) {
	return parseGroupAssociation(parseStruct, t)
}

func parseGroupAssociation(parseFn parseStructure, t *Table) (*GroupAssociation, error) {
	if t.Type != TableTypeGroupAssociation {
		return nil, fmt.Errorf("invalid table type %d", t.Type)
	}
	ga := &GroupAssociation{Header: t.Header}
	off, err := parseFn(t, 0 /* off */, false /* complete */, ga)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int((ga.Length-5)/3); i++ {
		itemType, err := t.GetByteAt(off)
		if err != nil {
			return nil, err
		}
		ga.ItemType = append(ga.ItemType, TableType(itemType))
		off++

		ih, err := t.GetWordAt(off)
		if err != nil {
			return nil, err
		}
		ga.ItemHandle = append(ga.ItemHandle, ih)
		off += 2
	}

	if len(ga.ItemType) != len(ga.ItemHandle) {
		return nil, fmt.Errorf("item type and item handle have different lengths, len of ItemType: %d, len of ItemHandle: %d", len(ga.ItemType), len(ga.ItemHandle))
	}

	// Check if the length of ItemType is correct.
	numOfItemType := float64((ga.Length - 5) / 3)
	if float64(len(ga.ItemType)) != numOfItemType {
		return nil, fmt.Errorf("mismatch between length and item type and item handle lengths, len of ItemType: %d, len of ItemHandle: %d, length: %f", len(ga.ItemType), len(ga.ItemHandle), numOfItemType)
	}

	return ga, nil
}

func (ga *GroupAssociation) String() string {
	lines := []string{
		ga.Header.String(),
		fmt.Sprintf("Name: %s", ga.GroupName),
		fmt.Sprintf("Items: %d", int((ga.Length-5)/3)),
	}
	for i, it := range ga.ItemType {
		lines = append(lines, fmt.Sprintf("0x%04X (%s)", ga.ItemHandle[i], it))
	}
	return strings.Join(lines, "\n\t")
}

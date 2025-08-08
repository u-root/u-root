// Copyright 2016-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package smbios

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGroupAssociationString(t *testing.T) {
	tests := []struct {
		name string
		val  GroupAssociation
		want string
	}{
		{
			name: "Fully populated",
			val: GroupAssociation{
				Header: Header{
					Type:   TableTypeGroupAssociation,
					Length: 11,
					Handle: 0,
				},
				GroupName:  "Group",
				ItemType:   []TableType{TableTypeBaseboardInfo, TableTypeProcessorInfo},
				ItemHandle: []uint16{10, 11},
			},
			want: `Handle 0x0000, DMI type 14, 11 bytes
Group Associations
	Name: Group
	Items: 2
	0x000A (Base Board Information)
	0x000B (Processor Information)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.val.String()
			if result != tt.want {
				t.Errorf("BaseboardInfo().String(): '%s', want '%s'", result, tt.want)
			}
		})
	}
}

func TestParseGroupAssociationInfo(t *testing.T) {
	tests := []struct {
		name  string
		val   GroupAssociation
		table Table
		want  error
	}{
		{
			name: "Invalid Type",
			val:  GroupAssociation{},
			table: Table{
				Header: Header{
					Type: TableTypeBIOSInfo,
				},
				data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a},
			},
			want: fmt.Errorf("invalid table type 0"),
		},
		{
			name: "Required fields are missing",
			val:  GroupAssociation{},
			table: Table{
				Header: Header{
					Type: TableTypeGroupAssociation,
				},
				data: []byte{},
			},
			want: fmt.Errorf("required fields missing"),
		},
		{
			name: "Error parsing structure",
			val:  GroupAssociation{},
			table: Table{
				Header: Header{
					Type: TableTypeGroupAssociation,
				},
				data: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19,
					0x1a},
			},
			want: fmt.Errorf("error parsing structure"),
		},
		{
			name: "Parse valid GroupAssociation",
			val:  GroupAssociation{},
			table: Table{
				Header: Header{
					Type:   TableTypeGroupAssociation,
					Length: 8,
					Handle: 0,
				},
				data: []byte{0x00, 0x01, 0x02, 0x03},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseStruct := func(t *Table, off int, complete bool, sp interface{}) (int, error) {
				return 1, tt.want
			}
			_, err := parseGroupAssociation(parseStruct, &tt.table)

			if !checkError(err, tt.want) {
				t.Errorf("parseGroupAssociation(): '%v', want '%v'", err, tt.want)
			}
		})
	}
}

func TestGroupAssociationToTablePass(t *testing.T) {
	tests := []struct {
		name string
		ga   *GroupAssociation
		want *Table
	}{
		{
			name: "Simple Table",
			ga: &GroupAssociation{
				Header: Header{
					Type:   TableTypeGroupAssociation,
					Length: 11,
					Handle: 0,
				},
				GroupName:  "Group",
				ItemType:   []TableType{TableTypeBaseboardInfo, TableTypeProcessorInfo},
				ItemHandle: []uint16{10, 11},
			},
			want: &Table{
				Header: Header{
					Type:   TableTypeGroupAssociation,
					Length: 11,
					Handle: 0,
				},
				data: []byte{
					14, 11, 0, 0, // Header
					1,        // string number
					2, 10, 0, // ItemType, ItemHandle
					4, 11, 0,
				},
				strings: []string{"Group"},
			},
		},
	}

	for _, tt := range tests {
		got, err := tt.ga.toTable()
		if err != nil {
			t.Errorf("toTable() should pass but return error: %v", err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("toTable(): '%v', want '%v'", got, tt.want)
		}
	}
}

func TestGroupAssociationToTableFail(t *testing.T) {

	tests := []struct {
		name string
		ga   *GroupAssociation
	}{
		{
			name: "Invalid Length",
			ga: &GroupAssociation{
				Header: Header{
					Type:   TableTypeGroupAssociation,
					Length: 10,
					Handle: 0,
				},
				ItemType:   []TableType{TableTypeBaseboardInfo, TableTypeProcessorInfo},
				ItemHandle: []uint16{10, 11},
			},
		},
		{
			name: "Mismatch Item Type and Item Handle Lengths",
			ga: &GroupAssociation{
				Header: Header{
					Type:   TableTypeGroupAssociation,
					Length: 11,
					Handle: 0,
				},
				ItemType:   []TableType{TableTypeBaseboardInfo},
				ItemHandle: []uint16{10, 11},
			},
		},
	}

	for _, tt := range tests {
		_, err := tt.ga.toTable()
		if err == nil {
			t.Fatalf("toTable() should fail but pass")
		}
	}
}

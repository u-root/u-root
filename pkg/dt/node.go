// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"unicode"
)

// Empty represents an empty Device Tree value.
type Empty struct{}

// PHandle represents a pointer to another Node.
type PHandle uint32

// PropertyType is an enum of possible property types.
type PropertyType int

// These are the possible values for PropertyType.
const (
	EmptyType PropertyType = iota
	U32Type
	U64Type
	StringType
	PropEncodedArrayType
	PHandleType
	StringListType
)

// StandardPropertyTypes maps properties to values as defined by the spec.
var StandardPropertyTypes = map[string]PropertyType{
	"compatible":     StringListType,
	"model":          StringType,
	"phandle":        PHandleType,
	"status":         StringType,
	"#address-cells": U32Type,
	"#size-cells":    U32Type,
	"reg":            PropEncodedArrayType, // TODO: support cells
	"virtual-reg":    U32Type,
	"ranges":         PropEncodedArrayType, // TODO: or EmptyType
	"dma-ranges":     PropEncodedArrayType, // TODO: or EmptyType
	"name":           StringType,           // deprecated
	"device_tree":    StringType,           // deprecated
}

var (
	errPropertyRegionInvalid = errors.New("property value is not <u64x2>")
)

// Node is one Node in the Device Tree.
type Node struct {
	Name       string
	Properties []Property `json:",omitempty"`
	Children   []*Node    `json:",omitempty"`
}

// Walk calls f on a Node and alls its descendents.
func (n *Node) Walk(f func(*Node) error) error {
	if err := f(n); err != nil {
		return err
	}
	for idx := range n.Children {
		if err := n.Children[idx].Walk(f); err != nil {
			return err
		}
	}
	return nil
}

// Find finds a Node starting at a node, given a matching function.
func (n *Node) Find(f func(*Node) bool) (*Node, bool) {
	if ok := f(n); ok {
		return n, ok
	}
	for idx := range n.Children {
		if nn, ok := n.Children[idx].Find(f); ok {
			return nn, ok
		}
	}
	return nil, false
}

// FindAll returns all Node starting at a node, given a matching function.
func (n *Node) FindAll(f func(*Node) bool) ([]*Node, bool) {
	var nodes []*Node
	if ok := f(n); ok {
		nodes = append(nodes, n)
	}

	for idx := range n.Children {
		if matching, ok := n.Children[idx].FindAll(f); ok {
			nodes = append(nodes, matching...)
		}
	}
	if len(nodes) == 0 {
		return nil, false
	}
	return nodes, true
}

// NodeByName uses Find to find a node by name.
func (n *Node) NodeByName(name string) (*Node, bool) {
	return n.Find(func(n *Node) bool {
		return n.Name == name
	})
}

// LookProperty finds a property by name.
func (n *Node) LookProperty(name string) (*Property, bool) {
	for idx := range n.Properties {
		if n.Properties[idx].Name == name {
			return &n.Properties[idx], true
		}
	}
	return nil, false
}

// RemoveProperty deletes a property by name.
func (n *Node) RemoveProperty(name string) bool {
	for idx := range n.Properties {
		if n.Properties[idx].Name == name {
			lastIdx := len(n.Properties) - 1
			if idx != lastIdx {
				n.Properties[idx] = n.Properties[lastIdx]
			}
			n.Properties = n.Properties[:lastIdx]
			return true
		}
	}
	return false
}

// UpdateProperty updates a property in the node, adding it if it does not exist.
//
// Returning boolean to indicate if the property was found.
func (n *Node) UpdateProperty(name string, value []byte) bool {
	p, found := n.LookProperty(name)
	if found {
		p.Value = value
		return true
	}

	prop := Property{Name: name, Value: value}
	n.Properties = append(n.Properties, prop)
	return false
}

// Property is a name-value pair. Note the PropertyType of Value is not
// encoded.
type Property struct {
	Name  string
	Value []byte
}

// PredictType makes a prediction on what value the property contains based on
// its name and data. The data types are not encoded in the data structure, so
// some heuristics are used.
func (p *Property) PredictType() PropertyType {
	// Standard properties
	if value, ok := StandardPropertyTypes[p.Name]; ok {
		if _, err := p.AsType(value); err == nil {
			return value
		}
	}

	// Heuristic match
	if _, err := p.AsEmpty(); err == nil {
		return EmptyType
	}
	if _, err := p.AsString(); err == nil {
		return StringType
	}
	if _, err := p.AsStringList(); err == nil {
		return StringListType
	}
	if _, err := p.AsU32(); err == nil {
		return U32Type
	}
	if _, err := p.AsU64(); err == nil {
		return U64Type
	}
	return PropEncodedArrayType
}

// AsType converts a Property to a Go type using one of the AsXYX() functions.
// The resulting Go type is as follows:
//
//	AsType(fdt.EmptyType)            -> fdt.Empty
//	AsType(fdt.U32Type)              -> uint32
//	AsType(fdt.U64Type)              -> uint64
//	AsType(fdt.StringType)           -> string
//	AsType(fdt.PropEncodedArrayType) -> []byte
//	AsType(fdt.PHandleType)          -> fdt.PHandle
//	AsType(fdt.StringListType)       -> []string
func (p *Property) AsType(val PropertyType) (interface{}, error) {
	switch val {
	case EmptyType:
		return p.AsEmpty()
	case U32Type:
		return p.AsU32()
	case U64Type:
		return p.AsU64()
	case StringType:
		return p.AsString()
	case PropEncodedArrayType:
		return p.AsPropEncodedArray()
	case PHandleType:
		return p.AsPHandle()
	case StringListType:
		return p.AsStringList()
	}
	return nil, fmt.Errorf("%d not in the PropertyType enum", val)
}

// AsEmpty converts the property to the Go fdt.Empty type.
func (p *Property) AsEmpty() (Empty, error) {
	if len(p.Value) != 0 {
		return Empty{}, fmt.Errorf("property %q is not <empty>", p.Name)
	}
	return Empty{}, nil
}

// AsU32 converts the property to the Go uint32 type.
func (p *Property) AsU32() (uint32, error) {
	if len(p.Value) != 4 {
		return 0, fmt.Errorf("property %q is not <u32>", p.Name)
	}
	var val uint32
	err := binary.Read(bytes.NewBuffer(p.Value), binary.BigEndian, &val)
	return val, err
}

// AsU64 converts the property to the Go uint64 type.
func (p *Property) AsU64() (uint64, error) {
	if len(p.Value) != 8 {
		return 0, fmt.Errorf("property %q is not <u64>", p.Name)
	}
	var val uint64
	err := binary.Read(bytes.NewBuffer(p.Value), binary.BigEndian, &val)
	return val, err
}

// Region represents a memory range.
type Region struct {
	Start uint64
	Size  uint64
}

// AsRegion converts the property to a Region.
func (p *Property) AsRegion() (*Region, error) {
	if len(p.Value) != 16 {
		return nil, errPropertyRegionInvalid
	}
	var start, size uint64
	b := bytes.NewBuffer(p.Value)

	err := binary.Read(b, binary.BigEndian, &start)
	if err != nil {
		return nil, err
	}
	err = binary.Read(b, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	return &Region{Start: start, Size: size}, nil
}

// AsString converts the property to the Go string type. The trailing null
// character is stripped.
func (p *Property) AsString() (string, error) {
	if len(p.Value) == 0 || p.Value[len(p.Value)-1] != 0 {
		return "", fmt.Errorf("property %q is not <string> (0 length or no null)", p.Name)
	}
	str := p.Value[:len(p.Value)-1]
	if !isPrintableASCII(str) {
		return "", fmt.Errorf("property %q is not <string>", p.Name)
	}
	return string(str), nil
}

// AsPropEncodedArray converts the property to the Go []byte type.
func (p *Property) AsPropEncodedArray() ([]byte, error) {
	return p.Value, nil
}

// AsPHandle converts the property to the Go fdt.PHandle type.
func (p *Property) AsPHandle() (PHandle, error) {
	val, err := p.AsU32()
	return PHandle(val), err
}

// AsStringList converts the property to the Go []string type. The trailing
// null character of each string is stripped.
func (p *Property) AsStringList() ([]string, error) {
	if len(p.Value) == 0 || p.Value[len(p.Value)-1] != 0 {
		return nil, fmt.Errorf("property %q is not <stringlist>", p.Name)
	}
	value := p.Value
	strs := []string{}
	for len(p.Value) > 0 {
		nextNull := bytes.IndexByte(value, 0) // cannot be -1
		var str []byte
		str, value = value[:nextNull], value[nextNull+1:]
		if !isPrintableASCII(str) {
			return nil, fmt.Errorf("property %q is not <stringlist>", p.Name)
		}
		strs = append(strs, string(str))
	}
	return strs, nil
}

func isPrintableASCII(s []byte) bool {
	for _, v := range s {
		if v > unicode.MaxASCII || !unicode.IsPrint(rune(v)) {
			return false
		}
	}
	return true
}

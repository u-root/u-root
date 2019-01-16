// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dt

import (
	"bytes"
	"encoding/binary"
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

// Node is one Node in the Device Tree.
type Node struct {
	Name       string
	Properties []*Property `json:",omitempty"`
	Children   []*Node     `json:",omitempty"`
}

// Walk calls f on a Node and alls its descendents.
func (n *Node) Walk(f func(*Node) error) error {
	if err := f(n); err != nil {
		return err
	}
	for _, child := range n.Children {
		if err := child.Walk(f); err != nil {
			return err
		}
	}
	return nil
}

// GetNode returns a subnode with the given name. If the subnode does not
// already exist, nil is returned.
func (n *Node) GetNode(name string) *Node {
	for _, c := range n.Children {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// CreateNode returns a subnode with the given name. If the subnode does not
// already exist, a new one is created.
func (n *Node) CreateNode(name string) *Node {
	c := n.GetNode(name)
	if c == nil {
		c = &Node{Name: name}
		n.Children = append(n.Children, c)
	}
	return c
}

// GetProperty returns a property of the Node with the given name. If the
// property does not exist, nil is returned.
func (n *Node) GetProperty(name string) *Property {
	for _, p := range n.Properties {
		if p.Name == name {
			return p
		}
	}
	return nil
}

// CreateProperty returns a property of the Node with the given name. If the
// property does not already exist, a new one is created.
func (n *Node) CreateProperty(name string) *Property {
	p := n.GetProperty(name)
	if p == nil {
		p = &Property{Name: name}
		n.Properties = append(n.Properties, p)
	}
	return p
}

// DeleteProperty deletes a property from the Node. It returns false iff the
// property did not exist in the first place.
func (n *Node) DeleteProperty(name string) bool {
	for i, p := range n.Properties {
		if p.Name == name {
			n.Properties = append(n.Properties[:i], n.Properties[i+1:]...)
			return true
		}
	}
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
//     AsType(fdt.EmptyType)            -> fdt.Empty
//     AsType(fdt.U32Type)              -> uint32
//     AsType(fdt.U64Type)              -> uint64
//     AsType(fdt.StringType)           -> string
//     AsType(fdt.PropEncodedArrayType) -> []byte
//     AsType(fdt.PHandleType)          -> fdt.PHandle
//     AsType(fdt.StringListType)       -> []string
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

// AsString converts the property to the Go string type. The trailing null
// character is stripped.
func (p *Property) AsString() (string, error) {
	if len(p.Value) == 0 || p.Value[len(p.Value)-1] != 0 {
		return "", fmt.Errorf("property %q is not <string>", p.Name)
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
	return []string{}, nil
}

// SetEmpty sets the property to the empty value.
func (p *Property) SetEmpty() {
	p.Value = []byte{}
}

// SetU32 sets the property to a uint32.
func (p *Property) SetU32(val uint32) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, &val)
	p.Value = buf.Bytes()
}

// SetU64 sets the property to a uint64.
func (p *Property) SetU64(val uint64) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, &val)
	p.Value = buf.Bytes()
}

// SetString sets the property to a string. The string should not contain null
// characters and should be printable for it to be decoded properly.
func (p *Property) SetString(val string) {
	p.Value = append([]byte(val), 0)
}

// SetPropEncodedArray sets the property to a []byte slice.
func (p *Property) SetPropEncodedArray(val []byte) {
	p.Value = val
}

// SetPHandle sets the property to a PHandle.
func (p *Property) SetPHandle(val PHandle) {
	p.SetU32(uint32(val))
}

// SetStringList sets the property to an []string. The strings should not
// contain null characters and should be printable for it to be decoded
// properly.
func (p *Property) SetStringList(val []string) {
	p.Value = []byte{}
	for _, s := range val {
		p.Value = append(p.Value, []byte(s)...)
		p.Value = append(p.Value, 0)
	}
}

func isPrintableASCII(s []byte) bool {
	for _, v := range s {
		if v > unicode.MaxASCII || !unicode.IsPrint(rune(v)) {
			return false
		}
	}
	return true
}

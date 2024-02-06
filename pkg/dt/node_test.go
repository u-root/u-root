// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dt

import (
	"errors"
	"reflect"
	"testing"
)

func TestLookupImmediateChild(t *testing.T) {
	subtree0 := &Node{
		Name: "child0",
		Children: []*Node{
			{
				Name: "child0_1",
			},
		},
	}
	subtree1 := &Node{
		Name: "child1",
		Children: []*Node{
			{
				Name: "child1_1",
			},
		},
	}
	subtree2 := &Node{
		Name: "child2",
		Children: []*Node{
			{
				Name: "child2_1",
			},
		},
	}
	tree := &Node{
		Name: "parent",
		Children: []*Node{
			subtree0,
			subtree1,
			subtree2,
		},
	}
	for _, tc := range []struct {
		name      string
		needle    string
		wantNode  *Node
		wantFound bool
	}{
		{name: "2nd child", needle: "child1", wantNode: subtree1, wantFound: true},
		{name: "3rd child", needle: "child2", wantNode: subtree2, wantFound: true},
		{name: "1st child", needle: "child0", wantNode: subtree0, wantFound: true},
		{name: "exists but not immediate", needle: "child1_1", wantFound: false},
		{name: "missing", needle: "not found", wantFound: false},
		{name: "prefix", needle: "child", wantFound: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			n, found := tree.LookupChildByName(tc.needle)
			if found != tc.wantFound {
				t.Errorf("tree.LookupChildByName(%s) returns found %v, want %v",
					tc.needle, found, tc.wantFound)
			}
			if found && tc.wantFound && !reflect.DeepEqual(n, tc.wantNode) {
				t.Errorf("when looking up %s, got %+v, want %+v", tc.needle, n, tc.wantNode)
			}
		})
	}
}

func TestNewNode(t *testing.T) {
	for _, tt := range []struct {
		name string
		opts []NodeOptioner
		node *Node
	}{
		{
			name: "new",
			opts: nil,
			node: &Node{Name: "new"},
		},
		{
			name: "new-with-property",
			opts: []NodeOptioner{WithProperty(PropertyU64("foo", 1))},
			node: &Node{
				Name: "new-with-property",
				Properties: []Property{
					{
						Name:  "foo",
						Value: []byte{0, 0, 0, 0, 0, 0, 0, 1},
					},
				},
			},
		},
		{
			name: "new-with-children",
			opts: []NodeOptioner{WithChildren(NewNode("child", WithProperty(PropertyString("foo", "abc"))))},
			node: &Node{
				Name: "new-with-children",
				Children: []*Node{
					{
						Name: "child",
						Properties: []Property{
							{
								Name:  "foo",
								Value: []byte{'a', 'b', 'c', 0},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			n := NewNode(tt.name, tt.opts...)
			if !reflect.DeepEqual(n, tt.node) {
				t.Errorf("NewNode = %v, want %v", n, tt.node)
			}
		})
	}
}

func TestPruneSubtree(t *testing.T) {
	subtree1 := &Node{
		Name: "child1",
		Children: []*Node{
			{
				Name: "child1_1",
			},
			{
				Name: "child1_2",
			},
		},
	}
	subtree2 := &Node{
		Name: "child2",
		Children: []*Node{
			{
				Name: "child2_1",
			},
			{
				Name: "child2_2",
			},
		},
	}
	subtree3 := &Node{
		Name: "child3",
		Children: []*Node{
			{
				Name: "child3_1",
			},
			{
				Name: "child3_2",
			},
		},
	}

	for _, tc := range []struct {
		name    string
		idx     int
		want    *Node
		wantErr error
	}{
		{
			name: "remove at head",
			idx:  0,
			want: &Node{
				Name:     "parent",
				Children: []*Node{subtree2, subtree3},
			},
		},
		{
			name: "remove in the middle",
			idx:  1,
			want: &Node{
				Name:     "parent",
				Children: []*Node{subtree1, subtree3},
			},
		},
		{
			name: "remove at tail",
			idx:  2,
			want: &Node{
				Name:     "parent",
				Children: []*Node{subtree1, subtree2},
			},
		},
		{
			name:    "remove invalid index",
			idx:     -1,
			wantErr: errInvalidChildIndex,
		},
		{
			name:    "remove another invalid index",
			idx:     3,
			wantErr: errInvalidChildIndex,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tree := &Node{
				Name: "parent",
				Children: []*Node{
					subtree1,
					subtree2,
					subtree3,
				},
			}
			err := tree.RemoveSubTreeAtIndex(tc.idx)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("t.RemoveSubTreeAtIndex(%d) returned error %v, want error %v",
					tc.idx, err, tc.wantErr)
			}
			if err == nil && tc.wantErr == nil && !reflect.DeepEqual(tree, tc.want) {
				t.Errorf("after removing %d got %+v, want %+v", tc.idx, tree, tc.want)
			}
		})
	}
}

func TestRemoveProperty(t *testing.T) {
	for _, tc := range []struct {
		name     string
		node     *Node
		remove   string
		want     *Node
		wantBool bool
	}{
		{
			name: "empty property list",
			node: &Node{
				Name:       "test node",
				Properties: []Property{},
			},
			remove: "linux,initrd-end",
			want: &Node{
				Name:       "test node",
				Properties: []Property{},
			},
			wantBool: false,
		},
		{
			name: "remove non-exist property",
			node: &Node{
				Name: "test node",
				Properties: []Property{
					{Name: "linux,elfcorehdr", Value: []byte{1, 2, 3}},
					{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
				},
			},
			remove: "linux,initrd-end",
			want: &Node{
				Name: "test node",
				Properties: []Property{
					{Name: "linux,elfcorehdr", Value: []byte{1, 2, 3}},
					{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
				},
			},
			wantBool: false,
		},
		{
			name: "remove middle property, success",
			node: &Node{
				Name: "test node",
				Properties: []Property{
					{Name: "linux,elfcorehdr", Value: []byte{1, 2, 3}},
					{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
					{Name: "kaslr-seed", Value: []byte{1, 2, 3}},
					{Name: "rng-seed", Value: []byte{1, 2, 3}},
					{Name: "linux,initrd-start", Value: []byte{1, 2, 3}},
					{Name: "linux,initrd-end", Value: []byte{1, 2, 3}},
				},
			},
			remove: "linux,initrd-start",
			want: &Node{
				Name: "test node",
				Properties: []Property{
					{Name: "linux,elfcorehdr", Value: []byte{1, 2, 3}},
					{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
					{Name: "kaslr-seed", Value: []byte{1, 2, 3}},
					{Name: "rng-seed", Value: []byte{1, 2, 3}},
					{Name: "linux,initrd-end", Value: []byte{1, 2, 3}},
				},
			},
			wantBool: true,
		},
		{
			name: "remove last property, success",
			node: &Node{
				Name: "test node",
				Properties: []Property{
					{Name: "linux,elfcorehdr", Value: []byte{1, 2, 3}},
					{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
					{Name: "kaslr-seed", Value: []byte{1, 2, 3}},
					{Name: "rng-seed", Value: []byte{1, 2, 3}},
					{Name: "linux,initrd-start", Value: []byte{1, 2, 3}},
					{Name: "linux,initrd-end", Value: []byte{1, 2, 3}},
				},
			},
			remove: "linux,initrd-end",
			want: &Node{
				Name: "test node",
				Properties: []Property{
					{Name: "linux,elfcorehdr", Value: []byte{1, 2, 3}},
					{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
					{Name: "kaslr-seed", Value: []byte{1, 2, 3}},
					{Name: "rng-seed", Value: []byte{1, 2, 3}},
					{Name: "linux,initrd-start", Value: []byte{1, 2, 3}},
				},
			},
			wantBool: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.node.RemoveProperty(tc.remove); got != tc.wantBool {
				t.Errorf("tc.node.RemoveProperty(%s) = %t, want %t", tc.remove, got, tc.wantBool)
			}
			if !reflect.DeepEqual(tc.node, tc.want) {
				t.Errorf("after removing %s got %+v, want %+v", tc.remove, tc.node, tc.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	for _, tt := range []struct {
		node  *Node
		prop  Property
		want  *Node
		found bool
	}{
		{
			node: &Node{
				Name: "test node",
				Properties: []Property{
					{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
					{Name: "kaslr-seed", Value: []byte{1, 2, 3}},
				},
			},
			prop: PropertyU64("kaslr-seed", 1),
			want: &Node{
				Name: "test node",
				Properties: []Property{
					{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
					{Name: "kaslr-seed", Value: []byte{0, 0, 0, 0, 0, 0, 0, 1}},
				},
			},
			found: true,
		},
		{
			node: &Node{
				Name: "test node",
				Properties: []Property{
					{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
				},
			},
			prop: PropertyU64("kaslr-seed", 1),
			want: &Node{
				Name: "test node",
				Properties: []Property{
					{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
					{Name: "kaslr-seed", Value: []byte{0, 0, 0, 0, 0, 0, 0, 1}},
				},
			},
			found: false,
		},
	} {
		got := tt.node.Update(tt.prop)
		if got != tt.found {
			t.Errorf("Node update, want %v, got %v", got, tt.found)
		}
		if !reflect.DeepEqual(tt.node, tt.want) {
			t.Errorf("Node update, want %v, got %v", tt.node, tt.want)
		}
	}
}

func TestUpdateProperty(t *testing.T) {
	node := &Node{
		Name: "test node",
		Properties: []Property{
			{Name: "linux,elfcorehdr", Value: []byte{1, 2, 3}},
			{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
			{Name: "kaslr-seed", Value: []byte{1, 2, 3}},
		},
	}

	// Try update an existing property.
	if got := node.UpdateProperty("kaslr-seed", []byte{3, 4, 5}); !got {
		t.Errorf("node.UpdateProperty(\"kaslr-seed\", []byte{3, 4,5}) = %t, want true", got)
	}
	want1 := &Node{
		Name: "test node",
		Properties: []Property{
			{Name: "linux,elfcorehdr", Value: []byte{1, 2, 3}},
			{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
			{Name: "kaslr-seed", Value: []byte{3, 4, 5}},
		},
	}
	if !reflect.DeepEqual(node, want1) {
		t.Errorf("after updating %s got %+v, want %+v", "kaslr-seed", node, want1)
	}

	// Update an non-existent property.
	if got := node.UpdateProperty("rng-seed", []byte{3, 4, 5}); got {
		t.Errorf("node.UpdateProperty(\"rng-seed\", []byte{3, 4,5}) = %t, want false", got)
	}
	want2 := &Node{
		Name: "test node",
		Properties: []Property{
			{Name: "linux,elfcorehdr", Value: []byte{1, 2, 3}},
			{Name: "linux,usable-memory-range", Value: []byte{1, 2, 3}},
			{Name: "kaslr-seed", Value: []byte{3, 4, 5}},
			{Name: "rng-seed", Value: []byte{3, 4, 5}},
		},
	}
	if !reflect.DeepEqual(node, want2) {
		t.Errorf("after updating %s got %+v, want %+v", "kaslr-seed", node, want2)
	}
}

func TestAsRegion(t *testing.T) {
	for _, tc := range []struct {
		name    string
		p       *Property
		want    *Region
		wantErr error
	}{
		{
			name: "invalid value",
			p: &Property{
				Name:  "linux,initrd-start",
				Value: []byte{},
			},
			want:    &Region{},
			wantErr: ErrPropertyRegionInvalid,
		},
		{
			name: "read start and size, success",
			p: &Property{
				Name:  "linux,initrd-start",
				Value: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0xf},
			},
			want: &Region{
				Start: 0x0001020304050607,
				Size:  0x08090a0b0c0d0e0f,
			},
			wantErr: nil,
		},
		// Given value is of type []byte, and we check length equal to 16
		// at the beginning, it is nearly impossible for binary.Read for
		// 2 uint64 from a fixed size bytes slice of size 16 to go wrong.
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.p.AsRegion()
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("tc.p.AsRegion() returned error %v, want error %v", err, tc.wantErr)
			}
			if err == nil && tc.wantErr == nil {
				if got.Start != tc.want.Start || got.Size != tc.want.Size {
					t.Errorf("got region %v, want region %v", got, tc.want)
				}
			}
		})
	}
}

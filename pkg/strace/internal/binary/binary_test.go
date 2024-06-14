// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binary

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type ts struct {
	I1 int8
	I2 int16
	I3 int32
	I4 int64
	S  tts
	A  [1]int32
}

type tts struct {
	U1 uint8
	U2 uint16
	U3 uint32
	U4 uint64
	A  [1]uint32
}

func TestMarshalUnmarshal(t *testing.T) {
	t1 := ts{
		I1: 1,
		I2: math.MinInt16,
		I3: 3,
		I4: math.MaxInt64,
		S: tts{
			U1: 10,
			U2: math.MaxUint16,
			U3: 12,
			U4: math.MaxUint64,
		},
		A: [1]int32{13},
	}

	for _, order := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		mb := Marshal(nil, order, t1)
		var t2 ts
		Unmarshal(mb, order, &t2)

		diff := cmp.Diff(t1, t2)
		if diff != "" {
			t.Errorf("t1 is not equal to t2 with %v:\n%s", order, diff)
		}
	}
}

func TestUnexport(t *testing.T) {
	type us struct {
		I1 int32
		u  ts
		I2 int32
	}

	t1 := us{I1: 11, I2: 12}

	mb := Marshal(nil, binary.BigEndian, t1)
	var t2 us
	Unmarshal(mb, binary.BigEndian, &t2)

	// cmp.Diff will panic on non exported fields
	if t2.I1 != t1.I1 || t2.I2 != t1.I2 {
		t.Errorf("t2.I1 != t1.I1 or t2.I2 != t1.I2")
	}
}

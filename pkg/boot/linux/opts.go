// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package linux

import (
	"encoding/json"
	"io"
)

// KexecOptions abstract a collection of options to be passed in KexecLoad.
//
// Arch agnostic. Each arch knows to just look for options they care about.
// Alternatively, we could introduce arch specific options, so irrelevant options
// won't be compiled. But for simplification, have one shared struct to begin
// with, we can split when time comes.
type KexecOptions struct {
	// DTB is used as the device tree blob, if specified.
	DTB io.ReaderAt
}

// kexecOptionsJSON is same as KexecOptions, but with transformed fields to help with serialization of KexecOptions.
type kexecOptionsJSON struct {
	dtb string
}

func (ko *KexecOptions) MarshalJSON() ([]byte, error) {
	koJSON := kexecOptionsJSON{}
	return json.Marshal(&koJSON)
}

func (ko *KexecOptions) UnmarshalJSON(b []byte) error {
	koJSON := kexecOptionsJSON{}
	if err := json.Unmarshal(b, &koJSON); err != nil {
		return err
	}
	return nil
}

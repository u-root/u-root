// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package universalpayload

import (
	"fmt"

	"github.com/u-root/u-root/pkg/dt"
)

// Porperties to be fetched from device tree.
const (
	FirstLevelNodeName    = "images"
	SecondLevelNodeName   = "tianocore"
	LoadAddrPropertyName  = "load"
	EntryAddrPropertyName = "entry-start"
)

type FdtLoad struct {
	Load       uint64
	EntryStart uint64
}

// Device Tree Blob resides at the start of FIT binary. In order to
// get the expected load and entry point address, need to walk though
// DTB to get value of properties 'load' and 'entry-start'.
//
// The simplified device tree layout is:
//
//	/{
//	    images {
//	        tianocore {
//	            entry-start = <0x00000000 0x00805ac3>;
//	            load = <0x00000000 0x00800000>;
//	        }
//	    }
//	 }
func getFdtInfo(name string) (*FdtLoad, error) {
	fdt, err := dt.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read fdt file:%s", name)
	}

	firstLevelNode, succeed := fdt.NodeByName(FirstLevelNodeName)
	if succeed != true {
		return nil, fmt.Errorf("failed to find '%s' node", FirstLevelNodeName)
	}

	secondLevelNode, succeed := firstLevelNode.NodeByName(SecondLevelNodeName)
	if succeed != true {
		return nil, fmt.Errorf("failed to find '%s'' node", SecondLevelNodeName)
	}

	loadAddrProp, succeed := secondLevelNode.LookProperty(LoadAddrPropertyName)
	if succeed != true {
		return nil, fmt.Errorf("failed to find get '%s' property", LoadAddrPropertyName)
	}

	loadAddr, err := loadAddrProp.AsU64()
	if err != nil {
		return nil, fmt.Errorf("failed to convert property '%s' to u64", LoadAddrPropertyName)
	}

	entryAddrProp, succeed := secondLevelNode.LookProperty(EntryAddrPropertyName)
	if succeed != true {
		return nil, fmt.Errorf("failed to find get '%s' property", EntryAddrPropertyName)
	}

	entryAddr, err := entryAddrProp.AsU64()
	if err != nil {
		return nil, fmt.Errorf("failed to convert property '%s' to u64", EntryAddrPropertyName)
	}

	return &FdtLoad{
		Load:       loadAddr,
		EntryStart: entryAddr,
	}, nil
}

// Copyright 2010-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

/*
func testTarget(t *testing.T) {
	var tests = []struct {
		n   string
		v   *IBFTTarget
		out []uint8
	}{
		{"Empty", &IBFTTarget{Valid: "0", Boot: "1", CHAP: "0", RCHAP: "1", Index: "0", BootLUN: "8"}, []uint8{}},
	}
	Debug = t.Logf
	for _, tst := range tests {

		b, err := Marshal(tst.v)
		if err != nil {
			t.Errorf("%s: got %v, want nil", tst.n, err)
		}
		t.Logf("%s: (%T, %v) -> (%v, %v)", tst.n, tst.v, tst.v, b, err)
	}
}
*/
func TestIBFTSizes(t *testing.T) {
	if len(rawIBTFHeader) != 48 {
		t.Errorf("length of rawIBTFHeader: got %v, want 48", len(rawIBTFHeader))
	}
}

func TestIBFTMarshal(t *testing.T) {
	i := &IBFT{
		Multi: "1",
		Initiator: IBFTInitiator{
			Valid:                 "1",
			Boot:                  "1",
			SNSServer:             "1.2.3.4",
			SLPServer:             "localhost",
			PrimaryRadiusServer:   "121.1.1.1",
			SecondaryRadiusServer: "222.3.4.5",
			Name:                  "myinitor",
		},
		NIC0: IBFTNIC{
			Valid:        "1",
			Boot:         "1",
			Global:       "1",
			Index:        "0",
			IPAddress:    "5.5.5.5",
			Gateway:      "7.7.7.7",
			PrimaryDNS:   "8.8.8.8",
			SecondaryDNS: "9.9.9.9",
			DHCP:         "11.11.11.11",
			VLAN:         "10",
			MACAddress:   "00:0c:29:12:a4:2e",
			PCIBDF:       "0x18",
			HostName:     "somehost",
		},
		NIC1: IBFTNIC{
			Valid:        "1",
			Boot:         "1",
			Global:       "0",
			Index:        "1",
			IPAddress:    "15.5.5.5",
			Gateway:      "17.7.7.7",
			PrimaryDNS:   "18.8.8.8",
			SecondaryDNS: "19.9.9.9",
			DHCP:         "121.11.11.11",
			VLAN:         "12",
			MACAddress:   "11:22:33:44:55:66",
			PCIBDF:       "0x8",
			HostName:     "otherhost",
		},
		Target0: IBFTTarget{
			Valid:             "1",
			Boot:              "1",
			CHAP:              "1",
			RCHAP:             "0",
			Index:             "1",
			TargetIP:          "1.2.3.4:88",
			BootLUN:           "1234",
			ChapType:          "0",
			TargetName:        "target",
			CHAPName:          "clown",
			CHAPSecret:        "noun",
			ReverseCHAPName:   "verb",
			ReverseCHAPSecret: "adverb",
		},
		Target1: IBFTTarget{
			Valid:             "1",
			Boot:              "1",
			CHAP:              "1",
			RCHAP:             "1",
			Index:             "0",
			TargetIP:          "4.4.4.4:99",
			BootLUN:           "4444",
			ChapType:          "2",
			TargetName:        "bullseye",
			CHAPName:          "bozo",
			CHAPSecret:        "bee",
			ReverseCHAPName:   "barg",
			ReverseCHAPSecret: "arg",
		},
	}

	Debug = t.Logf
	if false {
		b, err := json.MarshalIndent(i, "", "\t")
		if err != nil {
			t.Fatalf("Marshal: got %v, want nil", err)
		}
		t.Logf("%s", string(b))
		j := &IBFT{}
		if err := json.Unmarshal(b, j); err != nil {
			t.Fatalf("Unmarshal: got %v, want nil", err)
		}
		if !reflect.DeepEqual(i, j) {
			t.Fatalf("Reading it in: got %q, want %q", j, i)
		}
	}
	t.Logf("Send it out")
	b, err := Marshal(i)
	if err != nil {
		t.Fatalf("Marshal: got %v, want nil", err)
	}
	t.Logf("Marshal to %v", err)
	if len(b) != 524 {
		t.Fatalf("Marshall: len is %d bytes and should be 2048", len(b))
	}
	f, err := ioutil.TempFile("", "acpi")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	n, err := f.Write(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Wrote %d bytes to %q", n, f.Name())
}

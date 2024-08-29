// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

// func FuzzIPCmd(f *testing.F) {
// 	// no log output
// 	log.SetOutput(io.Discard)
// 	log.SetFlags(0)

// 	// get seed corpora from testdata_new files
// 	seeds, err := filepath.Glob("testdata/fuzz/corpora/*.seed")
// 	if err != nil {
// 		f.Fatalf("failed to find seed corpora files: %v", err)
// 	}

// 	for _, seed := range seeds {
// 		seedBytes, err := os.ReadFile(seed)
// 		if err != nil {
// 			f.Fatalf("failed read seed corpora from files %v: %v", seed, err)
// 		}

// 		f.Add(string(seedBytes))
// 	}

// 	stdout := &bytes.Buffer{}
// 	f.Fuzz(func(t *testing.T, data string) {
// 		stdout.Reset()
// 		arg := strings.Split(data, " ")

// 		handle, err := netlink.NewHandle()
// 		if err != nil {
// 			t.Fatalf("failed to create netlink handle: %v", err)
// 		}

// 		cmd := cmd{
// 			Args:   arg,
// 			out:    stdout,
// 			handle: handle,
// 			Family: netlink.FAMILY_ALL,
// 		}

// 		cmd.run()
// 	})
// }

// Copyright 2018 the LinuxBoot Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//uefifs
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/linuxboot/fiano/pkg/uefi"
)

type uefifs struct {
	f   uefi.Firmware
	res []*regexp.Regexp
}

func (u *uefifs) Run(f uefi.Firmware) error {
	return u.Visit(f)
}

//   File  74696E69-6472-632E-7069-6F2F62696F73  EFI_FV_FILETYPE_FREEFORM           3354174
//    Sec                                        EFI_SECTION_RAW                    3354116
//    Sec  Initrd                                EFI_SECTION_USER_INTERFACE         18
//    Sec                                        EFI_SECTION_VERSION                14
func (u *uefifs) Visit(f uefi.Firmware) error {
	switch f := f.(type) {
	case *uefi.File:
		// I'm leaving this here because it's an example of using
		// json unmarshaling to avoid dealing with all kinds of
		// type coercion foo. You can define the struct you are
		// interested in, with a few elements. The JSON unmarshaler
		// will fill it up and skip the bits you don't want, which
		// you can the examine. It's easier sometimes than
		// navigating a thicket of types.
		// I already know you'll make me remove it but at least
		// it will be there in the PR history :-)
		b, err := uefi.MarshalFirmware(f)
		if err != nil {
			return err
		}
		var xxx struct {
			FirmwareElement struct {
				Header struct {
					UUID struct {
						UUID string
					}
					Type int
				}
			}
		}
		if err := json.Unmarshal(b, &xxx); err != nil {
			// not the droid we're looking for.
			return f.ApplyChildren(u)
		}
		// See if any section matches. If there are no sections,
		// or nothing matches, We Must Go Deeper.
		for _, re := range u.res {
			if !re.MatchString(xxx.FirmwareElement.Header.UUID.UUID) {
				continue
			}
			for _, s := range f.Sections {
				if s.Type != "EFI_SECTION_RAW" {
					continue
				}
				if _, err := os.Stdout.Write(s.Buf()); err != nil {
					log.Fatal(err)
				}
			}
			return nil
		}
		return f.ApplyChildren(u)

	case *uefi.Section:
		return nil

	default:
		return f.ApplyChildren(u)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("ueficat <uefi file> <files ...> # e.g. 74696E69-6472-632E-7069-6F2F62696F73")
	}

	b, err := ioutil.ReadFile(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	f, err := uefi.Parse(b)
	if err != nil {
		log.Fatal(err)
	}
	u := &uefifs{f: f}

	for _, a := range flag.Args()[1:] {
		re, err := regexp.Compile(a)
		if err != nil {
			log.Fatal(err)
		}
		u.res = append(u.res, re)
	}

	if err := u.Visit(f); err != nil {
		log.Fatal(err)
	}

}

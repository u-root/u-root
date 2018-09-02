// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ueficat will cat a set of files from a firmware volume.
// The files are specified with a regular expression that is
// matched against a file GUID.
//
// Synopsis:
//     ueficat <romimage> pat [pat ...]
//
// Description:
//     ueficat reads a firmware volume and, for every file in it, matches it
//     against the set of patterns passed in the command line.
//     The pat can be a simple re matching a GUID, or an re of the form re:re.
//     The optional second re matches a name in in one of the sections, typically
//     an EFI_SECTION_USER_INTERFACE, though we don't currently check for section type.
//     If the pat matches the guid and name, and the file has a section of type EFI_SECTION_RAW,
//     that section is written to os.Stdout.
//     For example, in one UEFI image we have, we can say
//     ueficat 7:Initrd
//     and the initrd
//     is output to stdout. If we have a lot of confidence, we can even say .:Initrd or 7:I
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/linuxboot/fiano/pkg/uefi"
)

type pat struct {
	guid *regexp.Regexp
	name *regexp.Regexp
}

type ueficat struct {
	f    uefi.Firmware
	pats []pat
}

func (u *ueficat) Run(f uefi.Firmware) error {
	return u.Visit(f)
}

//   File  74696E69-6472-632E-7069-6F2F62696F73  EFI_FV_FILETYPE_FREEFORM           3354174
//    Sec                                        EFI_SECTION_RAW                    3354116
//    Sec  Initrd                                EFI_SECTION_USER_INTERFACE         18
//    Sec                                        EFI_SECTION_VERSION                14
func (u *ueficat) Visit(f uefi.Firmware) error {
	switch f := f.(type) {
	case *uefi.File:
		var (
			guid  = f.Header.UUID.String()
			match bool
			dat   []byte
		)
		for _, pat := range u.pats {
			if !pat.guid.MatchString(guid) {
				continue
			}
			// We can't assume the name and the section have an order.
			// So remember if the name matched, and if there was a
			// raw section.
			for _, s := range f.Sections {
				if s.Type == "EFI_SECTION_RAW" {
					dat = s.Buf()
					continue
				}
				if pat.name.MatchString(s.Name) {
					match = true
				}
			}
			if match && dat != nil {
				if _, err := os.Stdout.Write(dat); err != nil {
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
	u := &ueficat{f: f}

	for _, a := range flag.Args()[1:] {
		var p pat
		switch args := strings.Split(a, ":"); len(args) {
		case 2:
			p.name = regexp.MustCompile(args[1])
			p.guid = regexp.MustCompile(args[0])
		case 1:
			p.name = regexp.MustCompile(".")
			p.guid = regexp.MustCompile(args[0])
		default:
			log.Fatalf("%s: more than one :", a)
		}
		u.pats = append(u.pats, p)
	}

	if err := u.Visit(f); err != nil {
		log.Fatal(err)
	}

}

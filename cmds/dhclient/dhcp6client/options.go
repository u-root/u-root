package dhcp6client

import (
	"encoding/binary"
	"sort"

	"github.com/d2g/dhcp4"
	"github.com/mdlayher/dhcp6"
)

type option struct {
	Code dhcp6.OptionCode
	Data []byte
}

type optslice []option
type byOptionCode optslice

func (b byOptionCode) Len() int               { return len(b) }
func (b byOptionCode) Less(i int, j int) bool { return b[i].Code < b[j].Code }
func (b byOptionCode) Swap(i int, j int)      { b[i], b[j] = b[j], b[i] }

func enumerate(o dhcp6.Options) optslice {
	var options optslice
	for k, v := range o {
		for _, vv := range v {
			options = append(options, option{
				Code: k,
				Data: vv,
			})
		}
	}
	sort.Sort(byOptionCode(options))
	return options
}

func (o optslice) count() int {
	var c int
	for _, oo := range o {
		// 2 bytes: option code
		// 2 bytes: option length
		// N bytes: option data
		c += 2 + 2 + len(oo.Data)
	}
	return c
}

func (o optslice) write(p []byte) {
	var i int
	for _, oo := range o {
		// 2 bytes: option code
		binary.BigEndian.PutUint16(p[i:i+2], uint16(oo.Code))
		i += 2

		// 2 bytes: option length
		binary.BigEndian.PutUint16(p[i:i+2], uint16(len(oo.Data)))
		i += 2
		// N bytes: option data
		copy(p[i:i+len(oo.Data)], oo.Data)
		i += len(oo.Data)
	}
}

func addRaw(o dhcp6.Options, key dhcp6.OptionCode, value []byte) {
	o[key] = append(o[key], value)
}

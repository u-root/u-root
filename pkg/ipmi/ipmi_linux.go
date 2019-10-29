package ipmi

import (
	"fmt"
	"os"
)

func Open(devnum int) (*Ipmi, error) {
	d := fmt.Sprintf("/dev/ipmi%d", devnum)

	f, err := os.OpenFile(d, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	return &Ipmi{File: f}, nil
}

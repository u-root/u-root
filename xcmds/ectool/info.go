package main

import (
	"time"
)

func info(ec ec) ([]byte, error) {
	return ec.Command(ecCmdGetChipInfo, 0, []byte{}, 96 /*len(_ ecResponseGetChipInfo)*/, time.Duration(10*time.Second))
}

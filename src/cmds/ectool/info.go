package main

import (
	"time"
)

func info(ec ec) ([]byte, error) {
	return ec.Command(EC_CMD_GET_CHIP_INFO, 0, []byte{}, 96 /*len(_ ec_response_get_chip_info)*/, time.Duration(10*time.Second))
}

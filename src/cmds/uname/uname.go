package main

import (
	"fmt"
	"log"
	"uroot"
)

func main() {
	if u, err := uroot.Uname(); err != nil {
		log.Fatalf("%v", err)
	} else {
        	fmt.Printf("%v", u)
	}
}

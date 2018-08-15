package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"hash/crc32"
	"flag"
	"strconv"
)

func helpPrinter() {

	fmt.Printf("Usage:\ncksum <File Name>\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func versionPrinter() {
	fmt.Println("cksum utility, URoot Version.")
	os.Exit(0)
}



func GetInput(fileName string) (input []byte, err error) {
	if fileName != "" {
		return ioutil.ReadFile(fileName)
	}
	return ioutil.ReadAll(os.Stdin)
}

func printCksum( input []byte ) uint32 {
	// Linux cksum polynomial 04C11DB7
	data := string(input)
	data += strconv.Itoa(len(input))
	return crc32.Checksum([]byte(data), crc32.MakeTable(uint32(0x7BD11C40)))
}



func main() {
	var (
		help      bool
		version   bool
	)
	cliArgs := ""
	flag.BoolVar(&help, "help",false, "Show this help and exit")
	flag.BoolVar(&version, "version", false, "Print Version")
	flag.Parse()

	if help {
		helpPrinter()
	}

	if version {
		versionPrinter()
	}
	if len(os.Args) >= 2 {
		cliArgs = os.Args[1];
	}
	input, err := GetInput(cliArgs)
	if err != nil {
		return
	}
	fmt.Println(printCksum(input),len(input),cliArgs)
}

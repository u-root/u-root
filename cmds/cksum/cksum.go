package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"hash/crc32"
)

func GetInput(fileName string) (input []byte, err error) {

	if fileName != "" {
		data, ioErr := ioutil.ReadFile(fileName)
		if ioErr != nil {
			return nil, ioErr
		}
		dataWithLen := string(data)
		dataWithLen += string(len(data))
		return []byte(dataWithLen[:len(dataWithLen)-2]), ioErr
	}
	scanner := bufio.NewScanner(os.Stdin)
	userInput := ""
	for scanner.Scan() {
		userInput += scanner.Text()
	}
	inputWithLen := fmt.Sprintf( "%s%d", userInput, len(userInput))
	fmt.Println( "in:", userInput ,string(len(userInput)), len(userInput), inputWithLen )
	return []byte(inputWithLen[:len(inputWithLen)-1]), scanner.Err()
}

func main() {
	cliArgs := ""
	if len(os.Args) >= 2 {
		cliArgs = os.Args[1];
	}
	input, err := GetInput(cliArgs)
	if err != nil {
		return
	}
	// 04C11DB7
	fmt.Println("Data:",string(input))
	crc := crc32.Checksum(input, crc32.MakeTable(uint32(0x7BD11C40)))
	fmt.Println(crc,len(input)+1,cliArgs)
}

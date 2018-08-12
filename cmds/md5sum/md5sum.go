package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"crypto/md5"
)

func GetInput(fileName string) (input []byte, err error) {

	if fileName != "" {
		return ioutil.ReadFile(fileName)
	}
	return ioutil.ReadAll(os.Stdin)
}

func main() {
	cliArgs := ""
	if len(os.Args) >= 2 {
		cliArgs = os.Args[1];
	}
	input, err := GetInput(cliArgs)
	if err != nil {
		fmt.Println("Error getting input." );
		os.Exit(-1)
	}
	fmt.Printf("%x ",md5.Sum(input))
	if cliArgs == "" {
		fmt.Printf(" -\n");
	}else{
		fmt.Printf(" %s\n",cliArgs);
	}
	os.Exit(0)
}

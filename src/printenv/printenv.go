package main

import (
	"fmt"
	"os"
);

func main(){
	e := os.Environ()
	
	for _, v := range(e) {
		fmt.Printf("%v\n", v)
	}
}

// Copyright 2013-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Delete files.
//
// Synopsis:
//     rm [-Rrvi] FILE...
//
// Options:
//     -i: interactive mode
//     -v: verbose mode
//     -R: remove file hierarchies
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"strings"
	"log"
)

// You can add more flags to this struct.
type rmFlags struct {
	recursive   bool
	verbose     bool
	interactive bool
}

// Helper method for interactive flag
func iFlag(input *bufio.Scanner, printString string) bool {
	for{	
		//replace Println with Printf
		fmt.Printf("%s (y/n)\n", printString)
		input.Scan()
		if input.Err()!= nil{
			log.Fatalf("Failed at error: %v", input.Err())		
		}
		if(strings.Compare(input.Text(), "y") == 0){	
			return true
		}
		if(strings.Compare(input.Text(), "n") == 0){
			return false
		}
		fmt.Printf("Please only enter y or n. You entered %s. \n", input.Text())	
	}
}


func fileOnly (){
}


func recursiveDelete(file string, flags rmFlags) error {
	input := bufio.NewScanner(os.Stdin)
	statval, err := os.Stat(file)
	if err != nil {
		return err
	}
	if statval.IsDir(){
		// Throws error if flag is not recursive.
		if !flags.recursive {
			return &os.PathError{Op: "\nrm:", Path: file, Err: syscall.EISDIR}
		}
		// At this point, the recursive flag is on.
		if !flags.interactive && !flags.verbose {
			os.RemoveAll(file)
			return nil
		}
		// At this point, either -i or -v are on as well.
		fileList, err := ioutil.ReadDir(file)
		if err != nil {
			return err
		}
		if len(fileList) == 0 {
			if flags.interactive {
				printString :=  strings.Join([]string{"rm: remove directory '", file, "'? "}, "")
				if !iFlag(input,printString){
					return nil				
				}
				if err := os.Remove(file); err != nil {
					return err
				}
			}
			if flags.verbose {
					fmt.Printf("removed directory '%v'\n", file)
			}
		} else {
			if flags.interactive {
				printString :=  strings.Join([]string{"rm: descend into directory '", file, "'? "}, "")
				if !iFlag(input,printString){
					return nil				
				}
			}
			for _, each := range fileList {
				recursiveDelete(filepath.Join(file, each.Name()), flags)
				recursiveDelete(file, flags)
			}
		} 
	} else {
	if flags.interactive {
		if statval.Size() == 0 {
			printString :=  strings.Join([]string{"rm: remove regular empty file '", file, "'? "}, "")
			if !iFlag(input,printString){
					return nil				
				}
		} else {
			printString :=  strings.Join([]string{"rm: remove '", file, "'? "}, "")
			if !iFlag(input,printString){
					return nil				
				}
		}
		
		}
		if err := os.Remove(file); err != nil {
			return err
		}
		if flags.verbose {
			fmt.Printf("removed '%v'\n", file)
		}
	}
	return nil 
}

func rm(files []string, flags rmFlags) error {
	for _, file := range files {
		if err := recursiveDelete(file, flags); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var flags rmFlags
	flag.BoolVar(&flags.verbose, "v", false, "Verbose mode.")
	flag.BoolVar(&flags.recursive, "r", false, "Recursive mode.")
	flag.BoolVar(&flags.interactive, "i", false, "Interactive mode.")
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
	}

	if err := rm(flag.Args(), flags); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

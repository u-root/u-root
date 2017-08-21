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
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type rmFlags struct {
	recursive   bool
	verbose     bool
	interactive bool
}


func isEmpty(file string) (bool, error, []os.FileInfo ) {
	fileList, err := ioutil.ReadDir(file)
	if err != nil {
		//ignore this boolean return value
		return false, err, fileList
	}
	if len(fileList) == 0 {
		return true, nil, fileList
	}
	return false, nil, fileList
}

func interactivePrompt(input *bufio.Scanner, s string, args ...interface{})( bool, error) {
	for {
		fmt.Printf(s+" (y/n)", args...)
		input.Scan()
		if input.Err() != nil {
			return false,input.Err()
		}
		switch input.Text() {
		case "y":
			return true, nil
		case "n":
			return false, nil
		}
		fmt.Printf("Please only enter y or n. You entered %s. \n", input.Text())
	}
}



func printRemove(flags rmFlags, input *bufio.Scanner, s1 string, s2 string,  file string)(bool, error){
	if flags.interactive {
		iVal, err:=interactivePrompt(input, s1, file)
		if  err != nil {
			return false, err
		}
		if !iVal{
			return false, nil		
		}		
	}
	if err := os.Remove(file); err != nil {
		return false, err
	}
	if flags.verbose {
		fmt.Printf(s2, file)
	}
	return true, nil
}


func rmImplement(file string, flags rmFlags) error{
	input := bufio.NewScanner(os.Stdin)
	statval, err := os.Stat(file)
	if err != nil {
		return err
	}
	if statval.IsDir() {
		if !flags.recursive {
			return errors.New(fmt.Sprintf("rm : %s : is a directory", file))
		}
		if !flags.interactive && !flags.verbose {
			if err := os.RemoveAll(file); err != nil {
				return err
			}			
			return nil
		}
		empty, err, fileList := isEmpty(file)
		if err != nil{
			return err		
		}
		if empty{
			cont, err := printRemove(flags, input, "rm: remove directory '%s'? ", "removed directory '%v'\n", file)
			if  err != nil{
				return err
			}
			if !cont{
				return nil
			}
		} else {
			if flags.interactive {
				interact, err := interactivePrompt(input, "rm: descend into directory '%s'? ", file)
				if err != nil {
					return err
				}
				if !interact{
					return nil
				}
			}
			for _, each := range fileList {
				rmImplement(filepath.Join(file, each.Name()), flags)
			}
		}
		empty, err, fileList = isEmpty(file)
		if err != nil{
			return err		
		}
		if empty{
			cont, err := printRemove(flags, input, "rm: remove directory '%s'? ",  "removed directory '%v'\n", file)
			if err != nil{
				return err
			}
			if !cont{
				return nil
			}
		} else{
			fmt.Printf("There are still files in this directory, cannot remove.") 
			return nil		
		}
	} else {
		rmString := "rm: remove file '%s'?"
		if statval.Size() == 0 {
			rmString = "rm: remove regular empty file '%s'?"
		}
		cont, err := printRemove(flags, input, rmString, "removed '%v'\n", file)
		if err != nil{
			return err
		}
		if !cont{
			return nil
		}
	}
	return nil
}

func rm(files []string, flags rmFlags) error {
	for _, file := range files {
		if err := rmImplement(file, flags); err != nil {
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

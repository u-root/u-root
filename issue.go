// Copyright 2012-2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// issue - read an issue file (/etc/issue) or one provided via an argument

// Synopsis:
//     issue [FILE]

// Description:
//    Issue is a command that reads and interprets the contents of /etc/issue
//    or a provided file in the same format.                                -

// Author:
//
//	xplshn <https://github.com/xplshn>
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"syscall"
	"time"
)

// Define functions to generate dynamic replacements
func getOSName() string {
	return os.Getenv("GOOS")
}

func getHostName() string {
	hostName, err := os.Hostname()
	if err != nil {
		return "error: " + err.Error()
	}
	return hostName
}

func getCurrentTTYorTERM() string {
	determineIfTerm := func() string {
		// Attempt to use tput to get the terminal type
		if _, err := exec.Command("tput", "-T", "try", "setaf", "0").Output(); err == nil {
			cmd := exec.Command("tput", "term")
			output, err := cmd.Output()
			if err != nil {
				return ""
			}
			return fmt.Sprintf("%s", output)
		}
		// Fallback to $TERM environment variable
		return os.Getenv("TERM")
	}
	determineIfTTY := func(stdout io.Writer) string {
		fi, err := os.Stdin.Stat()
		if err != nil {
			return ""
		}

		s, ok := fi.Sys().(*syscall.Stat_t)
		if !ok {
			return ""
		}

		var ttyPath string
		filepath.WalkDir("/dev", func(path string, dir os.DirEntry, _ error) error {
			if dir.IsDir() {
				return nil
			}

			if fi, err := os.Stat(path); err == nil {
				stat, ok := fi.Sys().(*syscall.Stat_t)
				if ok {
					if stat.Ino == s.Ino && stat.Dev == s.Dev {
						ttyPath = path
					}
				}
			}
			return nil
		})
		return ttyPath
	}

	currentTTY := determineIfTTY(os.Stdout)
	if currentTTY != "" {
		return currentTTY
	}
	return determineIfTerm()
}

func getCurrentTime() string {
	currentTime := time.Now().Format("15:04:05")
	return currentTime
}
func getCurrentDate() string {
	// Get the current time
	currentDate := time.Now().Format("Monday, 02 Jan 2006")
	return currentDate
}
func getWorkingDirectory() string {
	workingDir, err := os.Getwd()
	if err != nil {
		return "error: " + err.Error()
	}
	return workingDir
}
func getNumberOfCPUs() string {
	return fmt.Sprintf("%d", runtime.NumCPU())
}

// Define a map to hold the sequence patterns and their replacements
var issueSequences = map[string]string{
	"\\l": getCurrentTTYorTERM(),
	"\\t": getCurrentTime(),
	"\\d": getCurrentDate(),
	"\\H": getHostName(),
	"\\w": getWorkingDirectory(),
	"\\c": getNumberOfCPUs(),
	"\\n": "\n",
}

func main() {
	// Check if there is a command-line argument
	args := os.Args[1:]
	var filePath string
	if len(args) > 0 {
		filePath = args[0]
	} else {
		filePath = "/etc/issue"
	}

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Replace special sequences with actual values
		for seq, replacement := range issueSequences {
			// Use raw string literals for regex patterns
			regexPattern := regexp.QuoteMeta(seq)
			regex, err := regexp.Compile(regexPattern)
			if err != nil {
				fmt.Printf("Warning: Ignoring unknown sequence '%s': %v\n", seq, err)
				continue
			}
			line = regex.ReplaceAllString(line, replacement)
		}
		fmt.Println(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

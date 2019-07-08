package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/systemboot/systemboot/pkg/checker"
)

var configFile = flag.String("c", "", "Configuration file that defines checks and remediations")
var verbose = flag.Bool("v", false, "Print verbose messages")

func run() int {
	if *configFile == "" {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: configuration file argument (-c) is required\n")
		flag.Usage()
		return 1
	}

	var checklist []checker.Check

	checkerConfigStr, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Printf("Unable to open config file %s: %v", *configFile, err)
		return 1
	}

	if err = json.Unmarshal(checkerConfigStr, &checklist); err != nil {
		log.Printf("Unable to parse config file %s: %v\n", *configFile, err)
		return 1
	}

	if *verbose {
		log.Printf("Registered Checks: %v\n", checker.ListRegistered())
		log.Printf("Checklist: %s\n", prettyJSON(checklist))
	}

	results, numErrors := checker.Run(checklist)

	fmt.Printf("Checker Results: %s\n", prettyJSON(results))

	if numErrors > 0 {
		return 1
	}

	return 0
}

func prettyJSON(thing interface{}) string {
	formatted, err := json.MarshalIndent(thing, "", "    ")
	if err != nil {
		log.Printf("Error while attempting to pretty print %v: %v", thing, err)
		os.Exit(1)
	}
	return string(formatted)
}

func main() {
	flag.Parse()
	ret := run()
	os.Exit(ret)
}

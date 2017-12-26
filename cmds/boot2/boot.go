package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/u-root/u-root/pkg/diskboot"
)

var (
	v             = flag.Bool("v", false, "Print debug messages")
	verbose       = func(string, ...interface{}) {}
	bootPath      = flag.String("p", "/boot", "Boot path")
	sConfigIndex  = flag.String("c", "", "Config index")
	sEntryIndex   = flag.String("n", "", "Entry index")
	appendCmdline = flag.String("append", "", "Additional kernel params")
)

func getConfig() *diskboot.Config {
	configs := diskboot.FindConfigs(*bootPath)
	if len(configs) == 0 {
		log.Fatal("No config found")
	}

	verbose("Got configs: %v", configs)
	var err error
	configIndex := 0
	if len(configs) > 1 {
		if *sConfigIndex == "" {
			for i, config := range configs {
				log.Printf("Config #%v: %#v", i, config)
			}
			log.Fatal("Multiple configs found - must specify a config index")
		}
		if configIndex, err = strconv.Atoi(*sConfigIndex); err != nil ||
			configIndex < 0 || configIndex >= len(configs) {
			log.Fatal("Invalid config index:", *sConfigIndex)
		}
	}
	return configs[configIndex]
}

func getEntry(config *diskboot.Config) *diskboot.Entry {
	verbose("Got entries: %v", config.Entries)
	var err error
	entryIndex := 0
	if *sEntryIndex != "" {
		if entryIndex, err = strconv.Atoi(*sEntryIndex); err != nil ||
			entryIndex < 0 || entryIndex >= len(config.Entries) {
			log.Fatal("Invalid entry index:", *sEntryIndex)
		}
	} else if config.DefaultEntry >= 0 {
		entryIndex = config.DefaultEntry
	} else {
		for i, entry := range config.Entries {
			log.Printf("Entry #%v: %#v", i, entry)
		}
		log.Fatal("No entry specified")
	}
	return &config.Entries[entryIndex]
}

func bootEntry(config *diskboot.Config, entry *diskboot.Entry) {
	verbose("Booting entry: %v", entry)
	err := entry.KexecLoad(config.MountPath, *appendCmdline)
	if err != nil {
		log.Fatal("Error doing kexec load:", err)
	}

	/*
		err = kexec.Reboot()
		if err != nil {
			log.Fatal("Error doings kexec reboot:", err)
		}
	*/
}

func main() {
	flag.Parse()

	if *v {
		verbose = log.Printf
	}

	config := getConfig()
	entry := getEntry(config)
	bootEntry(config, entry)
}

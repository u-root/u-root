package main

import (
	"flag"
	"log"
	"strconv"

	"github.com/u-root/u-root/pkg/diskboot"
	"github.com/u-root/u-root/pkg/kexec"
)

var (
	bootPath      = flag.String("p", "/boot", "Boot path")
	bootIndex     = flag.String("n", "", "Boot entry index")
	appendCmdline = flag.String("append", "", "additional kernel params")
)

func main() {
	flag.Parse()

	configs := diskboot.FindConfigs(*bootPath)
	if len(configs) == 0 {
		log.Fatal("No config found")
	}

	config := configs[0]
	var entryIndex int
	if *bootIndex != "" {
		index, err := strconv.Atoi(*bootIndex)
		if err != nil || index < 0 || index >= len(config.Entries) {
			log.Fatal("Invalid index:", *bootIndex)
		}
		entryIndex = index
	} else if config.DefaultEntry >= 0 {
		entryIndex = config.DefaultEntry
	} else {
		for i, entry := range config.Entries {
			log.Printf("#%v: %#v", i, entry)
		}
		log.Fatal("No entry specified")
	}

	entry := config.Entries[entryIndex]
	err := entry.KexecLoad(config.MountPath, *appendCmdline)
	if err != nil {
		log.Fatal("Error doing kexec load:", err)
	}

	err = kexec.Reboot()
	if err != nil {
		log.Fatal("Error doing kexec reboot:", err)
	}
}

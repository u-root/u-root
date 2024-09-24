//go:build linux
// +build linux

package core

import (
	"fmt"
	"os"
)

func init() {
	var err error
	clockFactor, tickInUSec, err = readPsched()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}
}

func readPsched() (float64, float64, error) {
	var pschedFile string

	if procRoot := os.Getenv("PROC_ROOT"); procRoot != "" {
		pschedFile = fmt.Sprintf("%s/net/psched", procRoot)
	} else {
		pschedFile = "/proc/net/psched"
	}

	fd, err := os.Open(pschedFile)
	if err != nil {
		return 1.0, 1.0, fmt.Errorf("using default values for clock. could not open /proc/net/psched: %v", err)
	}
	defer fd.Close()

	var t2us, us2t, clockRes, hiClockRes uint32
	_, err = fmt.Fscanf(fd, "%08x %08x %08x %08x", &t2us, &us2t, &clockRes, &hiClockRes)
	if err != nil {
		return 1.0, 1.0, fmt.Errorf("could not read /proc/net/psched: %v", err)
	}

	clockFactor := float64(clockRes) / timeUnitsPerSec
	tickInUSec := float64(t2us) / float64(us2t) * clockFactor

	return clockFactor, tickInUSec, nil
}

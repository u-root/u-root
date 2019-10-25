// Copyright 2017-2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rng

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/u-root/u-root/pkg/recovery"
)

// The concept
//
// Most systems in real world application do not provide enough entropy at boot
// time. Therefore we will seed /dev/random with /dev/hwrng if a HW random
// number generator is available. Entropy is important for cryptographic
// protocols running in network stacks. Also disk encryption can be a problem
// if bad or no entropy is available. It can either block provisioning or makes
// a symmetric key easy to re-calculate.

var (
	// HwRandomCurrentFile shows/sets the current
	// HW random number generator
	HwRandomCurrentFile = "/sys/class/misc/hw_random/rng_current"
	// HwRandomAvailableFile shows the current available
	// HW random number generator
	HwRandomAvailableFile = "/sys/class/misc/hw_random/rng_available"
	// RandomEntropyAvailableFile shows how much of the entropy poolsize is used
	RandomEntropyAvailableFile = "/proc/sys/kernel/random/entropy_avail"
	// EntropyFeedTime sets the loop time for seeding /dev/random by /dev/hwrng
	// in seconds
	EntropyFeedTime = 2 * time.Second
	// EntropyBlockSize sets the bytes to read per Read function call
	EntropyBlockSize = 128
	// EntropyThreshold is used to stop seeding at specific entropy level
	EntropyThreshold uint64 = 3000
	// RandomDevice is the linux random device
	RandomDevice = "/dev/random"
	// HwRandomDevice is the linux hw random device
	HwRandomDevice = "/dev/hwrng"
)

// trngList is a list of hw random number generator
// names used by the Linux kernel.
// Can be extended but keep in mind to priorize
// more secure random sources like hw random over
// timer, jitter based mechanisms. At the top of the array
// is the highest priority.
// <rng-name>
var trngList = []string{
	"tpm-rng",
	"intel-rng",
	"amd-rng",
	"timeriomem-rng",
}

// setAvailableTRNG searches for available True Random Number Generator
// inside the kernel api and sets the most secure on if
// available which seeds /dev/hwrng
func setAvailableTRNG() error {
	var (
		currentRNG    string
		availableRNGs []string
		selectedRNG   string
	)

	availableFileData, err := ioutil.ReadFile(HwRandomAvailableFile)
	if err != nil {
		return err
	}
	availableRNGs = strings.Split(string(availableFileData), " ")

	for _, trng := range trngList {
		for _, rng := range availableRNGs {
			if trng == rng {
				selectedRNG = trng
				break
			}
		}
	}

	if selectedRNG == "" {
		return errors.New("no TRNG found on platform")
	}

	if err = ioutil.WriteFile(HwRandomCurrentFile, []byte(selectedRNG), 0644); err != nil {
		return err
	}

	// Check if the correct TRNG was successful written
	currentFileData, err := ioutil.ReadFile(HwRandomCurrentFile)
	if err != nil {
		return err
	}
	currentRNG = string(currentFileData)

	if currentRNG != selectedRNG {
		return errors.New("Couldn't select TRNG: " + currentRNG)
	}

	return nil
}

// UpdateLinuxRandomness seeds random data from
// /dev/hwrng into /dev/random based on a timer and
// the entropy pool size
func UpdateLinuxRandomness(recoverer recovery.Recoverer) error {
	if err := setAvailableTRNG(); err != nil {
		return err
	}

	hwRng, err := os.OpenFile(HwRandomDevice, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return err
	}

	rng, err := os.OpenFile(RandomDevice, os.O_APPEND|os.O_WRONLY, os.ModeDevice)
	if err != nil {
		return err
	}

	go func() {
		defer hwRng.Close()
		defer rng.Close()

		for {
			time.Sleep(EntropyFeedTime)

			randomEntropyAvailableData, err := ioutil.ReadFile(RandomEntropyAvailableFile)
			if err != nil {
				recoverer.Recover("Can't read entropy pool size")
			}

			formatted := strings.TrimSuffix(string(randomEntropyAvailableData), "\n")
			randomEntropyAvailable, err := strconv.ParseUint(formatted, 10, 32)
			if err != nil {
				recoverer.Recover("Can't parse entropy pool size")
			}

			if randomEntropyAvailable >= EntropyThreshold {
				continue
			}

			var random = make([]byte, EntropyBlockSize)
			length, err := hwRng.Read(random)
			if err != nil {
				recoverer.Recover("Can't open the hardware random device")
			}
			_, err = rng.Write(random[:length])
			if err != nil {
				recoverer.Recover("Can't open the random device")
			}
		}
	}()

	return nil
}

// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package memory

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/unitutil"
	"github.com/jaypipes/ghw/pkg/util"
)

const (
	_WARN_CANNOT_DETERMINE_PHYSICAL_MEMORY = `
Could not determine total physical bytes of memory. This may
be due to the host being a virtual machine or container with no
/var/log/syslog file or /sys/devices/system/memory directory, or
the current user may not have necessary privileges to read the syslog.
We are falling back to setting the total physical amount of memory to
the total usable amount of memory
`
)

var (
	// System log lines will look similar to the following:
	// ... kernel: [0.000000] Memory: 24633272K/25155024K ...
	_REGEX_SYSLOG_MEMLINE = regexp.MustCompile(`Memory:\s+\d+K\/(\d+)K`)
	// regexMemoryBlockDirname matches a subdirectory in either
	// /sys/devices/system/memory or /sys/devices/system/node/nodeX that
	// represents information on a specific memory cell/block
	regexMemoryBlockDirname = regexp.MustCompile(`memory\d+$`)
)

func (i *Info) load() error {
	paths := linuxpath.New(i.ctx)
	tub := memTotalUsableBytes(paths)
	if tub < 1 {
		return fmt.Errorf("Could not determine total usable bytes of memory")
	}
	i.TotalUsableBytes = tub
	tpb := memTotalPhysicalBytes(paths)
	i.TotalPhysicalBytes = tpb
	if tpb < 1 {
		i.ctx.Warn(_WARN_CANNOT_DETERMINE_PHYSICAL_MEMORY)
		i.TotalPhysicalBytes = tub
	}
	i.SupportedPageSizes, _ = memorySupportedPageSizes(paths.SysKernelMMHugepages)
	return nil
}

func AreaForNode(ctx *context.Context, nodeID int) (*Area, error) {
	paths := linuxpath.New(ctx)
	path := filepath.Join(
		paths.SysDevicesSystemNode,
		fmt.Sprintf("node%d", nodeID),
	)

	var err error
	var blockSizeBytes uint64
	var totPhys int64
	var totUsable int64

	totUsable, err = memoryTotalUsableBytesFromPath(filepath.Join(path, "meminfo"))
	if err != nil {
		return nil, err
	}

	blockSizeBytes, err = memoryBlockSizeBytes(paths.SysDevicesSystemMemory)
	if err == nil {
		totPhys, err = memoryTotalPhysicalBytesFromPath(path, blockSizeBytes)
		if err != nil {
			return nil, err
		}
	} else {
		// NOTE(jaypipes): Some platforms (e.g. ARM) will not have a
		// /sys/device/system/memory/block_size_bytes file. If this is the
		// case, we set physical bytes equal to either the physical memory
		// determined from syslog or the usable bytes
		//
		// see: https://bugzilla.redhat.com/show_bug.cgi?id=1794160
		// see: https://github.com/jaypipes/ghw/issues/336
		totPhys = memTotalPhysicalBytesFromSyslog(paths)
	}

	supportedHP, err := memorySupportedPageSizes(filepath.Join(path, "hugepages"))
	if err != nil {
		return nil, err
	}

	return &Area{
		TotalPhysicalBytes: totPhys,
		TotalUsableBytes:   totUsable,
		SupportedPageSizes: supportedHP,
	}, nil
}

func memoryBlockSizeBytes(dir string) (uint64, error) {
	// get the memory block size in byte in hexadecimal notation
	blockSize := filepath.Join(dir, "block_size_bytes")

	d, err := ioutil.ReadFile(blockSize)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(d)), 16, 64)
}

func memTotalPhysicalBytes(paths *linuxpath.Paths) (total int64) {
	defer func() {
		// fallback to the syslog file approach in case of error
		if total < 0 {
			total = memTotalPhysicalBytesFromSyslog(paths)
		}
	}()

	// detect physical memory from /sys/devices/system/memory
	dir := paths.SysDevicesSystemMemory
	blockSizeBytes, err := memoryBlockSizeBytes(dir)
	if err != nil {
		total = -1
		return total
	}

	total, err = memoryTotalPhysicalBytesFromPath(dir, blockSizeBytes)
	if err != nil {
		total = -1
	}
	return total
}

// memoryTotalPhysicalBytesFromPath accepts a directory -- either
// /sys/devices/system/memory (for the entire system) or
// /sys/devices/system/node/nodeX (for a specific NUMA node) -- and a block
// size in bytes and iterates over the sysfs memory block subdirectories,
// accumulating blocks that are "online" to determine a total physical memory
// size in bytes
func memoryTotalPhysicalBytesFromPath(dir string, blockSizeBytes uint64) (int64, error) {
	var total int64
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return -1, err
	}
	// There are many subdirectories of /sys/devices/system/memory or
	// /sys/devices/system/node/nodeX that are named memory{cell} where {cell}
	// is a 0-based index of the memory block. These subdirectories contain a
	// state file (e.g. /sys/devices/system/memory/memory64/state that will
	// contain the string "online" if that block is active.
	for _, file := range files {
		fname := file.Name()
		// NOTE(jaypipes): we cannot rely on file.IsDir() here because the
		// memory{cell} sysfs directories are not actual directories.
		if !regexMemoryBlockDirname.MatchString(fname) {
			continue
		}
		s, err := ioutil.ReadFile(filepath.Join(dir, fname, "state"))
		if err != nil {
			return -1, err
		}
		// if the memory block state is 'online' we increment the total with
		// the memory block size to determine the amount of physical
		// memory available on this system.
		if strings.TrimSpace(string(s)) != "online" {
			continue
		}
		total += int64(blockSizeBytes)
	}
	return total, nil
}

func memTotalPhysicalBytesFromSyslog(paths *linuxpath.Paths) int64 {
	// In Linux, the total physical memory can be determined by looking at the
	// output of dmidecode, however dmidecode requires root privileges to run,
	// so instead we examine the system logs for startup information containing
	// total physical memory and cache the results of this.
	findPhysicalKb := func(line string) int64 {
		matches := _REGEX_SYSLOG_MEMLINE.FindStringSubmatch(line)
		if len(matches) == 2 {
			i, err := strconv.Atoi(matches[1])
			if err != nil {
				return -1
			}
			return int64(i * 1024)
		}
		return -1
	}

	// /var/log will contain a file called syslog and 0 or more files called
	// syslog.$NUMBER or syslog.$NUMBER.gz containing system log records. We
	// search each, stopping when we match a system log record line that
	// contains physical memory information.
	logDir := paths.VarLog
	logFiles, err := ioutil.ReadDir(logDir)
	if err != nil {
		return -1
	}
	for _, file := range logFiles {
		if strings.HasPrefix(file.Name(), "syslog") {
			fullPath := filepath.Join(logDir, file.Name())
			unzip := strings.HasSuffix(file.Name(), ".gz")
			var r io.ReadCloser
			r, err = os.Open(fullPath)
			if err != nil {
				return -1
			}
			defer util.SafeClose(r)
			if unzip {
				r, err = gzip.NewReader(r)
				if err != nil {
					return -1
				}
			}

			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				line := scanner.Text()
				size := findPhysicalKb(line)
				if size > 0 {
					return size
				}
			}
		}
	}
	return -1
}

func memTotalUsableBytes(paths *linuxpath.Paths) int64 {
	amount, err := memoryTotalUsableBytesFromPath(paths.ProcMeminfo)
	if err != nil {
		return -1
	}
	return amount
}

func memoryTotalUsableBytesFromPath(meminfoPath string) (int64, error) {
	// In Linux, /proc/meminfo or its close relative
	// /sys/devices/system/node/node*/meminfo
	// contains a set of memory-related amounts, with
	// lines looking like the following:
	//
	// $ cat /proc/meminfo
	// MemTotal:       24677596 kB
	// MemFree:        21244356 kB
	// MemAvailable:   22085432 kB
	// ...
	// HugePages_Total:       0
	// HugePages_Free:        0
	// HugePages_Rsvd:        0
	// HugePages_Surp:        0
	// ...
	//
	// It's worth noting that /proc/meminfo returns exact information, not
	// "theoretical" information. For instance, on the above system, I have
	// 24GB of RAM but MemTotal is indicating only around 23GB. This is because
	// MemTotal contains the exact amount of *usable* memory after accounting
	// for the kernel's resident memory size and a few reserved bits.
	// Please note GHW cares about the subset of lines shared between system-wide
	// and per-NUMA-node meminfos. For more information, see:
	//
	//  https://www.kernel.org/doc/Documentation/filesystems/proc.txt
	r, err := os.Open(meminfoPath)
	if err != nil {
		return -1, err
	}
	defer util.SafeClose(r)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		key := parts[0]
		if !strings.Contains(key, "MemTotal") {
			continue
		}
		rawValue := parts[1]
		inKb := strings.HasSuffix(rawValue, "kB")
		value, err := strconv.Atoi(strings.TrimSpace(strings.TrimSuffix(rawValue, "kB")))
		if err != nil {
			return -1, err
		}
		if inKb {
			value = value * int(unitutil.KB)
		}
		return int64(value), nil
	}
	return -1, fmt.Errorf("failed to find MemTotal entry in path %q", meminfoPath)
}

func memorySupportedPageSizes(hpDir string) ([]uint64, error) {
	// In Linux, /sys/kernel/mm/hugepages contains a directory per page size
	// supported by the kernel. The directory name corresponds to the pattern
	// 'hugepages-{pagesize}kb'
	out := make([]uint64, 0)

	files, err := ioutil.ReadDir(hpDir)
	if err != nil {
		return out, err
	}
	for _, file := range files {
		parts := strings.Split(file.Name(), "-")
		sizeStr := parts[1]
		// Cut off the 'kb'
		sizeStr = sizeStr[0 : len(sizeStr)-2]
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return out, err
		}
		out = append(out, uint64(size*int(unitutil.KB)))
	}
	return out, nil
}

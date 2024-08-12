// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package cpu

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw/pkg/context"
	"github.com/jaypipes/ghw/pkg/linuxpath"
	"github.com/jaypipes/ghw/pkg/util"
)

var (
	regexForCpulCore = regexp.MustCompile("^cpu([0-9]+)$")
)

func (i *Info) load() error {
	i.Processors = processorsGet(i.ctx)
	var totCores uint32
	var totThreads uint32
	for _, p := range i.Processors {
		totCores += p.NumCores
		totThreads += p.NumThreads
	}
	i.TotalCores = totCores
	i.TotalThreads = totThreads
	return nil
}

func processorsGet(ctx *context.Context) []*Processor {
	paths := linuxpath.New(ctx)

	lps := logicalProcessorsFromProcCPUInfo(ctx)
	// keyed by processor ID (physical_package_id)
	procs := map[int]*Processor{}

	// /sys/devices/system/cpu pseudodir contains N number of pseudodirs with
	// information about the logical processors on the host. These logical
	// processor pseudodirs are of the pattern /sys/devices/system/cpu/cpu{N}
	fnames, err := ioutil.ReadDir(paths.SysDevicesSystemCPU)
	if err != nil {
		ctx.Warn("failed to read /sys/devices/system/cpu: %s", err)
		return []*Processor{}
	}
	for _, fname := range fnames {
		matches := regexForCpulCore.FindStringSubmatch(fname.Name())
		if len(matches) < 2 {
			continue
		}

		lpID, err := strconv.Atoi(matches[1])
		if err != nil {
			ctx.Warn("failed to find numeric logical processor ID: %s", err)
			continue
		}

		procID := processorIDFromLogicalProcessorID(ctx, lpID)
		proc, found := procs[procID]
		if !found {
			proc = &Processor{ID: procID}
			lp, ok := lps[lpID]
			if !ok {
				ctx.Warn(
					"failed to find attributes for logical processor %d",
					lpID,
				)
				continue
			}

			// Assumes /proc/cpuinfo is in order of logical processor id, then
			// lps[lpID] describes logical processor `lpID`.
			// Once got a more robust way of fetching the following info,
			// can we drop /proc/cpuinfo.
			if len(lp.Attrs["flags"]) != 0 { // x86
				proc.Capabilities = strings.Split(lp.Attrs["flags"], " ")
			} else if len(lp.Attrs["Features"]) != 0 { // ARM64
				proc.Capabilities = strings.Split(lp.Attrs["Features"], " ")
			}
			if len(lp.Attrs["model name"]) != 0 {
				proc.Model = lp.Attrs["model name"]
			} else if len(lp.Attrs["uarch"]) != 0 { // SiFive
				proc.Model = lp.Attrs["uarch"]
			}
			if len(lp.Attrs["vendor_id"]) != 0 {
				proc.Vendor = lp.Attrs["vendor_id"]
			} else if len(lp.Attrs["isa"]) != 0 { // RISCV64
				proc.Vendor = lp.Attrs["isa"]
			}
			procs[procID] = proc
		}

		coreID := coreIDFromLogicalProcessorID(ctx, lpID)
		core := proc.CoreByID(coreID)
		if core == nil {
			core = &ProcessorCore{ID: coreID, NumThreads: 1}
			proc.Cores = append(proc.Cores, core)
			proc.NumCores += 1
		} else {
			core.NumThreads += 1
		}
		proc.NumThreads += 1
		core.LogicalProcessors = append(core.LogicalProcessors, lpID)
	}
	res := []*Processor{}
	for _, p := range procs {
		res = append(res, p)
	}
	return res
}

// processorIDFromLogicalProcessorID returns the processor physical package ID
// for the supplied logical processor ID
func processorIDFromLogicalProcessorID(ctx *context.Context, lpID int) int {
	paths := linuxpath.New(ctx)
	// Fetch CPU ID
	path := filepath.Join(
		paths.SysDevicesSystemCPU,
		fmt.Sprintf("cpu%d", lpID),
		"topology", "physical_package_id",
	)
	return util.SafeIntFromFile(ctx, path)
}

// coreIDFromLogicalProcessorID returns the core ID for the supplied logical
// processor ID
func coreIDFromLogicalProcessorID(ctx *context.Context, lpID int) int {
	paths := linuxpath.New(ctx)
	// Fetch CPU ID
	path := filepath.Join(
		paths.SysDevicesSystemCPU,
		fmt.Sprintf("cpu%d", lpID),
		"topology", "core_id",
	)
	return util.SafeIntFromFile(ctx, path)
}

func CoresForNode(ctx *context.Context, nodeID int) ([]*ProcessorCore, error) {
	// The /sys/devices/system/node/nodeX directory contains a subdirectory
	// called 'cpuX' for each logical processor assigned to the node. Each of
	// those subdirectories contains a topology subdirectory which has a
	// core_id file that indicates the 0-based identifier of the physical core
	// the logical processor (hardware thread) is on.
	paths := linuxpath.New(ctx)
	path := filepath.Join(
		paths.SysDevicesSystemNode,
		fmt.Sprintf("node%d", nodeID),
	)
	cores := make([]*ProcessorCore, 0)

	findCoreByID := func(coreID int) *ProcessorCore {
		for _, c := range cores {
			if c.ID == coreID {
				return c
			}
		}

		c := &ProcessorCore{
			ID:                coreID,
			LogicalProcessors: []int{},
		}
		cores = append(cores, c)
		return c
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, "cpu") {
			continue
		}
		if filename == "cpumap" || filename == "cpulist" {
			// There are two files in the node directory that start with 'cpu'
			// but are not subdirectories ('cpulist' and 'cpumap'). Ignore
			// these files.
			continue
		}
		// Grab the logical processor ID by cutting the integer from the
		// /sys/devices/system/node/nodeX/cpuX filename
		cpuPath := filepath.Join(path, filename)
		procID, err := strconv.Atoi(filename[3:])
		if err != nil {
			ctx.Warn(
				"failed to determine procID from %s. Expected integer after 3rd char.",
				filename,
			)
			continue
		}
		coreIDPath := filepath.Join(cpuPath, "topology", "core_id")
		coreID := util.SafeIntFromFile(ctx, coreIDPath)
		core := findCoreByID(coreID)
		core.LogicalProcessors = append(
			core.LogicalProcessors,
			procID,
		)
	}

	for _, c := range cores {
		c.NumThreads = uint32(len(c.LogicalProcessors))
	}

	return cores, nil
}

// logicalProcessor contains information about a single logical processor
// on the host.
type logicalProcessor struct {
	// This is the logical processor ID assigned by the host. In /proc/cpuinfo,
	// this is the zero-based index of the logical processor as it appears in
	// the /proc/cpuinfo file and matches the "processor" attribute. In
	// /sys/devices/system/cpu/cpu{N} pseudodir entries, this is the N value.
	ID int
	// The entire collection of string attribute name/value pairs for the
	// logical processor.
	Attrs map[string]string
}

// logicalProcessorsFromProcCPUInfo reads the `/proc/cpuinfo` pseudofile and
// returns a map, keyed by logical processor ID, of logical processor structs.
//
// `/proc/cpuinfo` files look like the following:
//
// ```
// processor	: 0
// vendor_id	: AuthenticAMD
// cpu family	: 23
// model		: 8
// model name	: AMD Ryzen 7 2700X Eight-Core Processor
// stepping	: 2
// microcode	: 0x800820d
// cpu MHz		: 2200.000
// cache size	: 512 KB
// physical id	: 0
// siblings	: 16
// core id		: 0
// cpu cores	: 8
// apicid		: 0
// initial apicid	: 0
// fpu		: yes
// fpu_exception	: yes
// cpuid level	: 13
// wp		: yes
// flags		: fpu vme de pse tsc msr pae mce <snip...>
// bugs		: sysret_ss_attrs null_seg spectre_v1 spectre_v2 spec_store_bypass retbleed smt_rsb
// bogomips	: 7386.41
// TLB size	: 2560 4K pages
// clflush size	: 64
// cache_alignment	: 64
// address sizes	: 43 bits physical, 48 bits virtual
// power management: ts ttp tm hwpstate cpb eff_freq_ro [13] [14]
//
// processor	: 1
// vendor_id	: AuthenticAMD
// cpu family	: 23
// model		: 8
// model name	: AMD Ryzen 7 2700X Eight-Core Processor
// stepping	: 2
// microcode	: 0x800820d
// cpu MHz		: 1885.364
// cache size	: 512 KB
// physical id	: 0
// siblings	: 16
// core id		: 1
// cpu cores	: 8
// apicid		: 2
// initial apicid	: 2
// fpu		: yes
// fpu_exception	: yes
// cpuid level	: 13
// wp		: yes
// flags		: fpu vme de pse tsc msr pae mce <snip...>
// bugs		: sysret_ss_attrs null_seg spectre_v1 spectre_v2 spec_store_bypass retbleed smt_rsb
// bogomips	: 7386.41
// TLB size	: 2560 4K pages
// clflush size	: 64
// cache_alignment	: 64
// address sizes	: 43 bits physical, 48 bits virtual
// power management: ts ttp tm hwpstate cpb eff_freq_ro [13] [14]
// ```
//
// with blank line-separated blocks of colon-delimited attribute name/value
// pairs for a specific logical processor on the host.
func logicalProcessorsFromProcCPUInfo(
	ctx *context.Context,
) map[int]*logicalProcessor {
	paths := linuxpath.New(ctx)
	r, err := os.Open(paths.ProcCpuinfo)
	if err != nil {
		return nil
	}
	defer util.SafeClose(r)

	lps := map[int]*logicalProcessor{}

	// A map of attributes describing the logical processor
	lpAttrs := map[string]string{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			// Output of /proc/cpuinfo has a blank newline to separate logical
			// processors, so here we collect up all the attributes we've
			// collected for this logical processor block
			lpIDstr, ok := lpAttrs["processor"]
			if !ok {
				ctx.Warn("expected to find 'processor' key in /proc/cpuinfo attributes")
				continue
			}
			lpID, _ := strconv.Atoi(lpIDstr)
			lp := &logicalProcessor{
				ID:    lpID,
				Attrs: lpAttrs,
			}
			lps[lpID] = lp
			// Reset the current set of processor attributes...
			lpAttrs = map[string]string{}
			continue
		}
		parts := strings.Split(line, ":")
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		lpAttrs[key] = value
	}
	return lps
}

//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package memory

import (
	"github.com/StackExchange/wmi"

	"github.com/jaypipes/ghw/pkg/unitutil"
)

const wqlOperatingSystem = "SELECT TotalVisibleMemorySize FROM Win32_OperatingSystem"

type win32OperatingSystem struct {
	TotalVisibleMemorySize *uint64
}

const wqlPhysicalMemory = "SELECT BankLabel, Capacity, DataWidth, Description, DeviceLocator, Manufacturer, Model, Name, PartNumber, PositionInRow, SerialNumber, Speed, Tag, TotalWidth FROM Win32_PhysicalMemory"

type win32PhysicalMemory struct {
	BankLabel     *string
	Capacity      *uint64
	DataWidth     *uint16
	Description   *string
	DeviceLocator *string
	Manufacturer  *string
	Model         *string
	Name          *string
	PartNumber    *string
	PositionInRow *uint32
	SerialNumber  *string
	Speed         *uint32
	Tag           *string
	TotalWidth    *uint16
}

func (i *Info) load() error {
	// Getting info from WMI
	var win32OSDescriptions []win32OperatingSystem
	if err := wmi.Query(wqlOperatingSystem, &win32OSDescriptions); err != nil {
		return err
	}
	var win32MemDescriptions []win32PhysicalMemory
	if err := wmi.Query(wqlPhysicalMemory, &win32MemDescriptions); err != nil {
		return err
	}
	// We calculate total physical memory size by summing the DIMM sizes
	var totalPhysicalBytes uint64
	i.Modules = make([]*Module, 0, len(win32MemDescriptions))
	for _, description := range win32MemDescriptions {
		totalPhysicalBytes += *description.Capacity
		i.Modules = append(i.Modules, &Module{
			Label:        *description.BankLabel,
			Location:     *description.DeviceLocator,
			SerialNumber: *description.SerialNumber,
			SizeBytes:    int64(*description.Capacity),
			Vendor:       *description.Manufacturer,
		})
	}
	var totalUsableBytes uint64
	for _, description := range win32OSDescriptions {
		// TotalVisibleMemorySize is the amount of memory available for us by
		// the operating system **in Kilobytes**
		totalUsableBytes += *description.TotalVisibleMemorySize * uint64(unitutil.KB)
	}
	i.TotalUsableBytes = int64(totalUsableBytes)
	i.TotalPhysicalBytes = int64(totalPhysicalBytes)
	return nil
}

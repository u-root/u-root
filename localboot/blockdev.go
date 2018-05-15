package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rekby/gpt"
)

type BlockDev struct {
	Name string
	Stat BlockStat
}

// https://www.kernel.org/doc/Documentation/block/stat.txt
func (b BlockDev) Summary() string {
	return fmt.Sprintf(`BlockStat{
    Name: %v
    ReadIOs: %v
    ReadMerges: %v
    ReadSectors: %v
    ReadTicks: %v
    WriteIOs: %v
    WriteMerges: %v
    WriteSectors: %v
    WriteTicks: %v
    InFlight: %v
    IOTicks: %v
    TimeInQueue: %v
}`,
		b.Name,
		b.Stat.ReadIOs,
		b.Stat.ReadMerges,
		b.Stat.ReadSectors,
		b.Stat.ReadTicks,
		b.Stat.WriteIOs,
		b.Stat.WriteMerges,
		b.Stat.WriteSectors,
		b.Stat.WriteTicks,
		b.Stat.InFlight,
		b.Stat.IOTicks,
		b.Stat.TimeInQueue,
	)
}

type BlockStat struct {
	ReadIOs      uint64
	ReadMerges   uint64
	ReadSectors  uint64
	ReadTicks    uint64
	WriteIOs     uint64
	WriteMerges  uint64
	WriteSectors uint64
	WriteTicks   uint64
	InFlight     uint64
	IOTicks      uint64
	TimeInQueue  uint64
}

// TODO use a better GUID library. Endianess, comparison and so on.
// System partitions have GUID == C12A7328-F81F-11D2-BA4B-00A0C93EC93B
var SystemPartitionGUID = gpt.Guid([...]byte{
	0x28, 0x73, 0x2a, 0xc1,
	0x1f, 0xf8,
	0xd2, 0x11,
	0xba, 0x4b,
	0x00, 0xa0, 0xc9, 0x3e, 0xc9, 0x3b,
})

func BlockStatFromBytes(buf []byte) (*BlockStat, error) {
	fields := strings.Fields(string(buf))
	// BlockStat has 11 fields
	if len(fields) != 11 {
		return nil, errors.New("Invalid number of fields")
	}
	intfields := make([]uint64, 0)
	for _, field := range fields {
		v, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return nil, err
		}
		intfields = append(intfields, v)
	}
	return &BlockStat{
		ReadIOs:      intfields[0],
		ReadMerges:   intfields[1],
		ReadSectors:  intfields[2],
		ReadTicks:    intfields[3],
		WriteIOs:     intfields[4],
		WriteMerges:  intfields[5],
		WriteSectors: intfields[6],
		WriteTicks:   intfields[7],
		InFlight:     intfields[8],
		IOTicks:      intfields[9],
		TimeInQueue:  intfields[10],
	}, nil
}

func GetBlockStats() ([]BlockDev, error) {
	blockdevs := make([]BlockDev, 0)
	devnames := make([]string, 0)
	root := "/sys/class/block"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		devnames = append(devnames, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, devname := range devnames {
		fd, err := os.Open(fmt.Sprintf("%s/%s/stat", root, devname))
		if err != nil {
			return nil, err
		}
		defer fd.Close()
		buf, err := ioutil.ReadAll(fd)
		if err != nil {
			return nil, err
		}
		bstat, err := BlockStatFromBytes(buf)
		if err != nil {
			return nil, err
		}
		blockdevs = append(blockdevs, BlockDev{Name: devname, Stat: *bstat})
	}
	return blockdevs, nil
}

func GetGPTTable(device BlockDev) (*gpt.Table, error) {
	fd, err := os.Open(fmt.Sprintf("/dev/%s", device.Name))
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	if _, err := fd.Seek(512, os.SEEK_SET); err != nil {
		return nil, err
	}
	table, err := gpt.ReadTable(fd, 512)
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func FilterEFISystemPartitions(devices []BlockDev) ([]BlockDev, error) {
	esps := make([]BlockDev, 0)
	for _, device := range devices {
		table, err := GetGPTTable(device)
		if err != nil {
			debug("Skipping %s: %v", device.Name, err)
			continue
		}
		for _, part := range table.Partitions {
			if part.IsEmpty() {
				continue
			}
			if part.Type.String() == SystemPartitionGUID.String() {
				esps = append(esps, device)
			}
		}
	}
	return esps, nil
}

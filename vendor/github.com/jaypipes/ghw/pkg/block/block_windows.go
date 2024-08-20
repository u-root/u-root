// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package block

import (
	"strconv"
	"strings"

	"github.com/StackExchange/wmi"

	"github.com/jaypipes/ghw/pkg/util"
)

type physicalDiskMediaType int

const (
	PHYSICAL_DISK_MEDIA_TYPE_UNSPECIFIED physicalDiskMediaType = 0
	PHYSICAL_DISK_MEDIA_TYPE_HDD         physicalDiskMediaType = 3
	PHYSICAL_DISK_MEDIA_TYPE_SSD         physicalDiskMediaType = 4
	PHYSICAL_DISK_MEDIA_TYPE_SCM         physicalDiskMediaType = 5
)

func (dt physicalDiskMediaType) ToDriveType() DriveType {
	switch dt {
	case PHYSICAL_DISK_MEDIA_TYPE_UNSPECIFIED:
		return DRIVE_TYPE_UNKNOWN
	case PHYSICAL_DISK_MEDIA_TYPE_HDD:
		return DRIVE_TYPE_HDD
	case PHYSICAL_DISK_MEDIA_TYPE_SSD:
		return DRIVE_TYPE_SSD
	case PHYSICAL_DISK_MEDIA_TYPE_SCM:
		return DRIVE_TYPE_UNKNOWN
	}
	return DRIVE_TYPE_UNKNOWN
}

const wqlDiskDrive = "SELECT Caption, CreationClassName, DefaultBlockSize, Description, DeviceID, Index, InterfaceType, Manufacturer, MediaType, Model, Name, Partitions, SerialNumber, Size, TotalCylinders, TotalHeads, TotalSectors, TotalTracks, TracksPerCylinder FROM Win32_DiskDrive"

type win32DiskDrive struct {
	Caption           *string
	CreationClassName *string
	DefaultBlockSize  *uint64
	Description       *string
	DeviceID          *string
	Index             *uint32 // Used to link with partition
	InterfaceType     *string
	Manufacturer      *string
	MediaType         *string
	Model             *string
	Name              *string
	Partitions        *int32
	SerialNumber      *string
	Size              *uint64
	TotalCylinders    *int64
	TotalHeads        *int32
	TotalSectors      *int64
	TotalTracks       *int64
	TracksPerCylinder *int32
}

const wqlDiskPartition = "SELECT Access, BlockSize, Caption, CreationClassName, Description, DeviceID, DiskIndex, Index, Name, Size, SystemName, Type FROM Win32_DiskPartition"

type win32DiskPartition struct {
	Access            *uint16
	BlockSize         *uint64
	Caption           *string
	CreationClassName *string
	Description       *string
	DeviceID          *string
	DiskIndex         *uint32 // Used to link with Disk Drive
	Index             *uint32
	Name              *string
	Size              *int64
	SystemName        *string
	Type              *string
}

const wqlLogicalDiskToPartition = "SELECT Antecedent, Dependent FROM Win32_LogicalDiskToPartition"

type win32LogicalDiskToPartition struct {
	Antecedent *string
	Dependent  *string
}

const wqlLogicalDisk = "SELECT Caption, CreationClassName, Description, DeviceID, FileSystem, FreeSpace, Name, Size, SystemName FROM Win32_LogicalDisk"

type win32LogicalDisk struct {
	Caption           *string
	CreationClassName *string
	Description       *string
	DeviceID          *string
	FileSystem        *string
	FreeSpace         *uint64
	Name              *string
	Size              *uint64
	SystemName        *string
}

const wqlPhysicalDisk = "SELECT DeviceId, MediaType FROM MSFT_PhysicalDisk"

type win32PhysicalDisk struct {
	DeviceId  string
	MediaType physicalDiskMediaType
}

func (i *Info) load() error {
	win32DiskDriveDescriptions, err := getDiskDrives()
	if err != nil {
		return err
	}

	win32DiskPartitionDescriptions, err := getDiskPartitions()
	if err != nil {
		return err
	}

	win32LogicalDiskToPartitionDescriptions, err := getLogicalDisksToPartitions()
	if err != nil {
		return err
	}

	win32LogicalDiskDescriptions, err := getLogicalDisks()
	if err != nil {
		return err
	}

	win32PhysicalDisks, err := getPhysicalDisks()
	if err != nil {
		return err
	}

	// Converting into standard structures
	disks := make([]*Disk, 0)
	for _, diskdrive := range win32DiskDriveDescriptions {
		var physicalDiskMediaType physicalDiskMediaType
		for _, physicalDisk := range win32PhysicalDisks {
			if id, err := strconv.Atoi(physicalDisk.DeviceId); err != nil {
				return err
			} else if uint32(id) == *diskdrive.Index {
				physicalDiskMediaType = physicalDisk.MediaType
			}
		}
		disk := &Disk{
			Name:                   strings.TrimSpace(*diskdrive.DeviceID),
			SizeBytes:              *diskdrive.Size,
			PhysicalBlockSizeBytes: *diskdrive.DefaultBlockSize,
			DriveType:              toDriveType(physicalDiskMediaType, *diskdrive.MediaType, *diskdrive.Caption),
			StorageController:      toStorageController(*diskdrive.InterfaceType),
			BusPath:                util.UNKNOWN, // TODO: add information
			NUMANodeID:             -1,
			Vendor:                 strings.TrimSpace(*diskdrive.Manufacturer),
			Model:                  strings.TrimSpace(*diskdrive.Caption),
			SerialNumber:           strings.TrimSpace(*diskdrive.SerialNumber),
			WWN:                    util.UNKNOWN, // TODO: add information
			Partitions:             make([]*Partition, 0),
		}
		for _, diskpartition := range win32DiskPartitionDescriptions {
			// Finding disk partition linked to current disk drive
			if diskdrive.Index == nil || diskpartition.DiskIndex == nil {
				continue
			}
			if *diskdrive.Index == *diskpartition.DiskIndex {
				disk.PhysicalBlockSizeBytes = *diskpartition.BlockSize
				// Finding logical partition linked to current disk partition
				for _, logicaldisk := range win32LogicalDiskDescriptions {
					for _, logicaldisktodiskpartition := range win32LogicalDiskToPartitionDescriptions {
						var desiredAntecedent = "\\\\" + *diskpartition.SystemName + "\\root\\cimv2:" + *diskpartition.CreationClassName + ".DeviceID=\"" + *diskpartition.DeviceID + "\""
						var desiredDependent = "\\\\" + *logicaldisk.SystemName + "\\root\\cimv2:" + *logicaldisk.CreationClassName + ".DeviceID=\"" + *logicaldisk.DeviceID + "\""
						if *logicaldisktodiskpartition.Antecedent == desiredAntecedent && *logicaldisktodiskpartition.Dependent == desiredDependent {
							// Appending Partition
							p := &Partition{
								Name:       strings.TrimSpace(*logicaldisk.Caption),
								Label:      strings.TrimSpace(*logicaldisk.Caption),
								SizeBytes:  *logicaldisk.Size,
								MountPoint: *logicaldisk.DeviceID,
								Type:       *diskpartition.Type,
								IsReadOnly: toReadOnly(*diskpartition.Access),
								UUID:       "",
							}
							disk.Partitions = append(disk.Partitions, p)
							break
						}
					}
				}
			}
		}
		disks = append(disks, disk)
	}

	i.Disks = disks
	var tpb uint64
	for _, d := range i.Disks {
		tpb += d.SizeBytes
	}
	i.TotalPhysicalBytes = tpb
	return nil
}

func getDiskDrives() ([]win32DiskDrive, error) {
	// Getting disks drives data from WMI
	var win3232DiskDriveDescriptions []win32DiskDrive
	if err := wmi.Query(wqlDiskDrive, &win3232DiskDriveDescriptions); err != nil {
		return nil, err
	}
	return win3232DiskDriveDescriptions, nil
}

func getDiskPartitions() ([]win32DiskPartition, error) {
	// Getting disk partitions from WMI
	var win32DiskPartitionDescriptions []win32DiskPartition
	if err := wmi.Query(wqlDiskPartition, &win32DiskPartitionDescriptions); err != nil {
		return nil, err
	}
	return win32DiskPartitionDescriptions, nil
}

func getLogicalDisksToPartitions() ([]win32LogicalDiskToPartition, error) {
	// Getting links between logical disks and partitions from WMI
	var win32LogicalDiskToPartitionDescriptions []win32LogicalDiskToPartition
	if err := wmi.Query(wqlLogicalDiskToPartition, &win32LogicalDiskToPartitionDescriptions); err != nil {
		return nil, err
	}
	return win32LogicalDiskToPartitionDescriptions, nil
}

func getLogicalDisks() ([]win32LogicalDisk, error) {
	// Getting logical disks from WMI
	var win32LogicalDiskDescriptions []win32LogicalDisk
	if err := wmi.Query(wqlLogicalDisk, &win32LogicalDiskDescriptions); err != nil {
		return nil, err
	}
	return win32LogicalDiskDescriptions, nil
}

func getPhysicalDisks() ([]win32PhysicalDisk, error) {
	// Getting physical disks from WMI
	var win32PhysicalDisks []win32PhysicalDisk
	if err := wmi.QueryNamespace(wqlPhysicalDisk, &win32PhysicalDisks, "root\\Microsoft\\Windows\\Storage"); err != nil {
		return nil, err
	}
	return win32PhysicalDisks, nil
}

func toDriveType(physicalDiskMediaType physicalDiskMediaType, mediaType string, caption string) DriveType {
	if driveType := physicalDiskMediaType.ToDriveType(); driveType != DRIVE_TYPE_UNKNOWN {
		return driveType
	}

	mediaType = strings.ToLower(mediaType)
	caption = strings.ToLower(caption)
	if strings.Contains(mediaType, "fixed") || strings.Contains(mediaType, "ssd") || strings.Contains(caption, "ssd") {
		return DRIVE_TYPE_SSD
	} else if strings.ContainsAny(mediaType, "hdd") {
		return DRIVE_TYPE_HDD
	}
	return DRIVE_TYPE_UNKNOWN
}

// TODO: improve
func toStorageController(interfaceType string) StorageController {
	var storageController StorageController
	switch interfaceType {
	case "SCSI":
		storageController = STORAGE_CONTROLLER_SCSI
	case "IDE":
		storageController = STORAGE_CONTROLLER_IDE
	default:
		storageController = STORAGE_CONTROLLER_UNKNOWN
	}
	return storageController
}

// TODO: improve
func toReadOnly(access uint16) bool {
	// See Access property from: https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-diskpartition
	return access == 0x1
}

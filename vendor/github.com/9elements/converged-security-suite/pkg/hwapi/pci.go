package hwapi

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

//PCIReadConfigSpace reads from PCI config space into buf
func (t TxtAPI) PCIReadConfigSpace(bus int, device int, devFn int, off int, buf interface{}) error {
	var path string
	path = fmt.Sprintf("/sys/bus/pci/devices/0000:%02x:%02x.%1x/config", bus, device, devFn)

	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(int64(off), io.SeekStart); err != nil {
		return err
	}
	return binary.Read(f, binary.LittleEndian, buf)
}

//PCIReadConfig16 reads 16bits from PCI config space
func (t TxtAPI) PCIReadConfig16(bus int, device int, devFn int, off int) (uint16, error) {
	var reg16 uint16

	err := t.PCIReadConfigSpace(bus, device, devFn, off, &reg16)
	if err != nil {
		return 0, err
	}

	return reg16, nil
}

//PCIReadConfig32 reads 32bits from PCI config space
func (t TxtAPI) PCIReadConfig32(bus int, device int, devFn int, off int) (uint32, error) {
	var reg32 uint32

	err := t.PCIReadConfigSpace(bus, device, devFn, off, &reg32)
	if err != nil {
		return 0, err
	}

	return reg32, nil
}

//PCIReadVendorID reads the device vendor ID from PCI config space
func (t TxtAPI) PCIReadVendorID(bus int, device int, devFn int) (uint16, error) {
	id, err := t.PCIReadConfig16(bus, device, devFn, 0)
	if err != nil {
		return 0, err
	}

	return id, nil
}

//PCIReadDeviceID reads the device ID from PCI config space
func (t TxtAPI) PCIReadDeviceID(bus int, device int, devFn int) (uint16, error) {
	id, err := t.PCIReadConfig16(bus, device, devFn, 2)
	if err != nil {
		return 0, err
	}

	return id, nil
}

const (
	// TsegPCIRegSandyAndNewer is the offset withing the MCH PCI config space since SandyBridge
	TsegPCIRegSandyAndNewer = 0xb8
	// TSEGPCIBroadwellde is the offset withing the MCH PCI config space
	TSEGPCIBroadwellde = 0xa8

	// DPRPCIRegSandyAndNewer is the offset withing the MCH PCI config space since SandyBridge
	DPRPCIRegSandyAndNewer = 0x5c
	// DPRPciRegBroadwellDE offset withing the VTd PCI config space
	DPRPciRegBroadwellDE = 0x290
)

var (
	// FIXME: Baytrail and Braswell have TSEG in IOSF BUNIT

	// HostbridgeIDsBroadwellDE lookup table is special...
	HostbridgeIDsBroadwellDE = []uint16{
		0x2F00,
		0x6F00,
	}

	// HostbridgeIDsSandyCompatible lookup table for most stuff that seems compatible with Sandy Bridge
	HostbridgeIDsSandyCompatible = []uint16{
		/* Sandy bridge */
		0x0100,
		0x0104,
		/* Ivy bridge */
		0x0150,
		0x0154,
		0x0158,
		/* Haswell */
		0x0c00,
		0x0c04,
		0x0a04,
		0x0c08,
		/* Denverton NS */
		0x1980,
		0x1995,
		/* Broadwell */
		0x1604,
		0x1610,
		0x1614,
		/* Apollolake */
		0x5af0,
		/* Gemini Lake */
		0x31f0,
		/* Skylake */
		0x1900,
		0x1904,
		0x190c,
		0x190f,
		0x1910,
		0x1918,
		0x191f,
		0x1924,
		/* Kabylake */
		0x5904,
		0x590c,
		0x590f,
		0x5910,
		0x5914,
		0x5918,
		0x591f,
		/* Cannonlake */
		0x5a04,
		0x5a02,
		/* Whiskylake */
		0x3E34,
		0x3E35,
		/* Coffeelake */
		0x3ED0,
		0x3ec4,
		0x3e20,
		0x3ec2,
		0x3e30,
		0x3e31,
		/* Icelake */
		0x8A12,
		0x8A02,
		0x8A10,
		0x8A00,
		/* Cometlake */
		0x9B61,
		0x9B71,
		0x9B51,
		0x9B60,
		0x9B55,
		0x9B35,
		0x9B54,
		0x9B44,
	}
)

//ReadHostBridgeTseg returns TSEG base and TSEG limit
func (t TxtAPI) ReadHostBridgeTseg() (uint32, uint32, error) {
	var tsegBaseOff int
	var tsegLimitOff int
	var tsegBroadwellDEfix bool
	var devicenum int

	vendorid, err := t.PCIReadVendorID(0, 0, 0)
	if err != nil {
		return 0, 0, err
	}
	if vendorid != 0x8086 {
		return 0, 0, fmt.Errorf("Hostbridge is not made by Intel")
	}
	deviceid, err := t.PCIReadDeviceID(0, 0, 0)

	var found bool
	for _, id := range HostbridgeIDsSandyCompatible {
		if id == deviceid {
			found = true
			tsegBaseOff = TsegPCIRegSandyAndNewer
			tsegLimitOff = TsegPCIRegSandyAndNewer + 4
			devicenum = 0
			break
		}
	}
	if !found {
		for _, id := range HostbridgeIDsBroadwellDE {
			if id == deviceid {
				found = true
				tsegBroadwellDEfix = true
				tsegBaseOff = TSEGPCIBroadwellde
				tsegLimitOff = TSEGPCIBroadwellde + 4
				devicenum = 5
				break
			}
		}
	}

	if !found {
		return 0, 0, fmt.Errorf("Hostbridge is unsupported")
	}

	var tsegbase uint32
	var tseglimit uint32

	tsegbase, err = t.PCIReadConfig32(0, devicenum, 0, tsegBaseOff)
	if err != nil {
		return 0, 0, err
	}

	tseglimit, err = t.PCIReadConfig32(0, devicenum, 0, tsegLimitOff)
	if err != nil {
		return 0, 0, err
	}

	if tsegBroadwellDEfix {
		// On BroadwellDe TSEG limit lower 19bits are don't care, thus add 1 MiB.
		tseglimit += 1024 * 1024
	}

	return tsegbase, tseglimit, nil
}

//DMAProtectedRange encodes the DPR register
type DMAProtectedRange struct {
	Lock bool
	// Reserved 1-3
	Size uint8
	// Reserved 12-19
	Top uint16
}

//ReadHostBridgeDPR reads the DPR register from PCI config space
func (t TxtAPI) ReadHostBridgeDPR() (DMAProtectedRange, error) {
	var dprOff int
	var devicenum int
	var ret DMAProtectedRange

	vendorid, err := t.PCIReadVendorID(0, 0, 0)
	if err != nil {
		return ret, err
	}
	if vendorid != 0x8086 {
		return ret, fmt.Errorf("Hostbridge is not made by Intel")
	}
	deviceid, err := t.PCIReadDeviceID(0, 0, 0)

	var found bool
	for _, id := range HostbridgeIDsSandyCompatible {
		if id == deviceid {
			found = true
			dprOff = DPRPCIRegSandyAndNewer
			devicenum = 0
			break
		}
	}
	if !found {
		for _, id := range HostbridgeIDsBroadwellDE {
			if id == deviceid {
				found = true
				dprOff = DPRPciRegBroadwellDE
				devicenum = 5
				break
			}
		}
	}

	if !found {
		return ret, fmt.Errorf("Hostbridge is unsupported")
	}

	var u32 uint32

	u32, err = t.PCIReadConfig32(0, devicenum, 0, dprOff)
	if err != nil {
		return ret, err
	}

	ret.Lock = u32&1 != 0
	ret.Size = uint8((u32 >> 4) & 0xff)   // 11:4
	ret.Top = uint16((u32 >> 20) & 0xfff) // 31:20

	return ret, nil
}

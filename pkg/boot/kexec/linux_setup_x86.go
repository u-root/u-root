package kexec

import (
	"fmt"
	"io"

	"github.com/u-root/u-root/pkg/uio"
)

const (
	INITRD_BASE = 0x1000000 /* 16MB */
)

func SetupLinuxBootloaderParameters(lph *LinuxParamHeader, kmem *Memory, ramfs io.ReaderAt, cmdLineOffset, setupBase uint, cmdline string) error {
	var initrdBase, initrdAddrMax, initrdSize uint

	/* Say I am a boot loader */
	lph.LoaderType = LOADER_TYPE_KEXEC << 4

	/* No loader flags */
	lph.LoaderFlagzaAq = 0

	/* Find the maximum initial ramdisk address */
	initrdAddrMax = DEFAULT_INITRD_ADDR_MAX
	if lph.ProtocolVersion >= 0x0203 {
		initrdAddrMax = uint(lph.InitrdAddrMax)
	}

	/* Load the initrd if we have one */
	var ramfsBytes []byte
	if ramfs != nil {
		r, err := uio.ReadAll(ramfs)
		if err == nil {
			ramfsBytes = r
		}
	}
	if len(ramfsBytes) > 0 { // TODO(10000TB): more sanity checks ?
		ramfsRange, err := kmem.AddPhysSegment(
			ramfsBytes,
			RangeFromInterval(
				uintptr(INITRD_BASE),
				uintptr(initrdAddrMax),
			),
			// TODO(10000TB): this by default align up by page size.
			// Add support alignement of 4096.
		)

		if err != nil {
			return fmt.Errorf("add ramfs bytes to kexec segments: %v", err)
		}

		initrdBase = uint(ramfsRange.Start)
		initrdSize = uint(len(ramfsBytes))
	} else {
		initrdBase, initrdSize = 0, 0
	}

	/* Ramdisk address and size */
	lph.InitrdStart = uint32(initrdBase)
	lph.InitrdSize = uint32(initrdSize)

	if lph.ProtocolVersion >= 0x020c && uint(uint32(initrdBase)) != initrdBase {
		lph.ExtRamdiskImage = uint32(initrdBase >> 32)
	}

	if lph.ProtocolVersion >= 0x020c && uint(uint32(initrdSize)) != initrdSize {
		lph.ExtRamdiskSize = uint32(initrdSize >> 32)
	}

	/* The location of the command line */
	lph.ClMagic = CL_MAGIC_VALUE
	lph.ClOffset = uint16(cmdLineOffset)

	if lph.ProtocolVersion >= 0x0202 {
		cmdLinePtr := setupBase + cmdLineOffset
		lph.CmdLinePtr = uint32(cmdLinePtr)

		if lph.ProtocolVersion >= 0x020c && uint(uint32(cmdLinePtr)) != cmdLinePtr {
			lph.ExtCmdLinePtr = uint32(cmdLinePtr >> 32)
		}
	}

	/* Fill in the command line */

	// TODO(10000TB): add support to fill the cmdline.

	return nil
}

func getBootParams(off uint) {

}

func setupSubarch(lph *LinuxParamHeader) error {
	//offset := 0x23C
	// TODO(10000TB): add support for subarch param setting.
	return nil
}

func setupLinuxVesaFb(lph *LinuxParamHeader) error {
	// TODO(10000TB): impl this.
	return nil
}

func setupE820(lph *LinuxParamHeader) error {
	// TODO(10000TB): impl this.
	return nil
}

func setupEddInfo(lph *LinuxParamHeader) error {
	// TODO(10000TB): impl this.
	return nil
}

func getACPIRsdp() uint64 {
	var acpiRsdp uint64 = 0

	//acpiRsdp = bootparam_get_acpi_rsdp()

	if acpiRsdp == 0 {
		// TODO(10000TB): get efi acpi rsdp.
		//
		// acpiRsdp = efi_get_acpi_rsdp()
	}

	return acpiRsdp
}

func SetupLinuxSystemParameters(lph *LinuxParamHeader) error {
	/* get subarch from running kernel */
	if err := setupSubarch(lph); err != nil {
		return err
	}

	/* Default screen size.
	 *
	 * Probably not needed for Linuxboot's non GUI boot in most cases, but
	 * just in case.
	 */
	lph.OrigX = 0
	lph.OrigY = 0
	lph.OrigVideoPage = 0
	lph.OrigVideoMode = 0
	lph.OrigVideoCols = 80
	lph.OrigVideoLines = 25
	lph.OrigVideoEgaBx = 0
	lph.OrigVideoIsVGA = 1
	lph.OrigVideoPoints = 16

	/* setup vesa fb if possible, or just use original screen_info */
	if err := setupLinuxVesaFb(lph); err != nil {
		/* save and restore the old cmdline param if needed */
		// TODO(10000TB): restore old cmdline params.
	}

	/* Fill in the memsize later */
	lph.ExtMemK = 0
	lph.AltMemK = 0
	lph.E820MapNr = 0

	/* Default APM info */
	lph.ApmBIOSInfo = ApmBIOSInfo{}
	/* Default drive info */
	lph.DriveInfo = DriveInfo{}
	/* Default sysdesc table */
	lph.SysDescTable = SysDescTable{}

	/* Default to yes.
	 *
	 * One can overwrite from cmdline.
	 */
	lph.MountRootRdonly = 0xFFFF

	/* Default to /dev/hda.
	 *
	 * One can overwrite from cmdline.
	 */
	lph.RootDev = (0x3 << 8) | 0

	/* another safe default */
	lph.AuxDeviceInfo = 0

	if err := setupE820(lph); err != nil {
		return fmt.Errorf("setup e820: %v", err)
	}

	if err := setupEddInfo(lph); err != nil {
		return fmt.Errorf("setup Edd info: %v", err)
	}

	/* Always try to fill acpi_rsdp_addr */
	lph.ACPIRsdpAddr = getACPIRsdp()

	return nil
}

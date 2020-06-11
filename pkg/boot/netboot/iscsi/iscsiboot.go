package iscsi

import (
	"context"
	"fmt"
	"net"

	"github.com/u-root/iscsinl"
	"github.com/u-root/u-root/pkg/boot/esxi"
	"github.com/u-root/u-root/pkg/boot/ibft"
	"github.com/u-root/u-root/pkg/boot/menu"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/u-root/u-root/pkg/ulog"
)

type ESXIBoot struct{}

func (ESXIBoot) Name() string { return "ESXi boot" }

func (ESXIBoot) Parse(ctx context.Context, ibft *ibft.IBFT, partitionedDevice string) ([]menu.Entry, error) {
	imgs, err := esxi.LoadDisk(partitionedDevice)
	if err != nil {
		return nil, err
	}

	var es []menu.Entry
	for _, img := range imgs {
		if ibft != nil {
			// Insert iBFT.
			img.IBFT = ibft
			img.Name = fmt.Sprintf("%s from iSCSI target %s", img.Name, ibft.Target0.Target)
		}
		es = append(es, menu.OSImages(true, img)...)
	}
	return es, nil
}

type DiskParser interface {
	Name() string
	Parse(ctx context.Context, ibft *ibft.IBFT, partitionedDevice string) ([]menu.Entry, error)
}

type ISCSIBoot struct {
	Log         ulog.Logger
	CreateIBFT  bool
	DiskParsers []DiskParser
}

func (ISCSIBoot) Name() string { return "iSCSI boot" }

// Parse attempts to boot from an iSCSI volume.
func (ib *ISCSIBoot) Parse(ctx context.Context, lease dhclient.Lease) ([]menu.Entry, error) {
	target, volume, err := lease.ISCSIBoot()
	if err != nil {
		return nil, fmt.Errorf("no usable iSCSI path in lease: %v", err)
	}

	devices, err := iscsinl.MountIscsi(
		iscsinl.WithTarget(target.String(), volume),
		iscsinl.WithScheduler("noop"),
		iscsinl.WithCmdsMax(128),
		iscsinl.WithQueueDepth(16),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to mount iSCSI for target %s @ volume %s: %v", target, volume, err)
	}

	var ibft *ibft.IBFT
	if ib.CreateIBFT {
		ibft = CreateIBFT(lease.Link().Attrs().HardwareAddr, target, volume)
	}

	var es []menu.Entry
	// There really should only be one device, I think?
	for _, dev := range devices {
		device := fmt.Sprintf("/dev/%s", dev)

		for _, parser := range ib.DiskParsers {
			entries, err := parser.Parse(ctx, ibft, device)
			if err != nil {
				ib.Log.Printf("Could not interpret iSCSI disk contents as %s: %v", parser.Name(), err)
			}
			es = append(es, entries...)
		}
	}
	return es, nil
}

// CreateIBFT makes an iBFT. It's only guaranteed to work with IPv4 ESXi right now.
func CreateIBFT(mac net.HardwareAddr, target *net.TCPAddr, volume string) *ibft.IBFT {
	return &ibft.IBFT{
		Initiator: ibft.Initiator{
			Valid: true,
			Boot:  true,
			Name:  "NERF",
		},
		NIC0: ibft.NIC{
			Valid:      true,
			Boot:       true,
			Global:     true,
			MACAddress: mac,

			// ESXi can live without this information. It'll do a
			// DHCP request to get it.
			//
			// This makes us trivially compatible with IPv6 as
			// well, since we don't have to query the kernel for
			// the RA information that would fill these fields in
			// addition to DHCPv6. But, note: IPv6 is untested.
			//
			// Plus, routes are gonna be more complicated than just
			// a "gateway". I'm curious what iPXE here.
			IPNet:   nil,
			Gateway: nil,
		},
		Target0: ibft.Target{
			Valid:          true,
			Boot:           true,
			NICAssociation: 0,
			Target:         target,
			TargetName:     volume,
		},
	}
}

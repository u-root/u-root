package virtio

const (
	// The number of free descriptors in virt queue must exceed
	// MAX_SKB_FRAGS (16). Otherwise, packet transmission from
	// the guest to the host will be stopped.
	//
	// refs https://github.com/torvalds/linux/blob/5859a2b/drivers/net/virtio_net.c#L1754
	QueueSize = 32
)

type IRQInjector interface {
	InjectVirtioNetIRQ() error
	InjectVirtioBlkIRQ() error
}

type commonHeader struct {
	_        uint32 // hostFeatures
	_        uint32 // guestFeatures
	_        uint32 // queuePFN
	queueNUM uint16
	queueSEL uint16
	_        uint16 // queueNotify
	_        uint8  // status
	isr      uint8
}

// refs: https://wiki.osdev.org/Virtio#Virtual_Queue_Descriptor
type VirtQueue struct {
	DescTable [QueueSize]struct {
		Addr  uint64
		Len   uint32
		Flags uint16
		Next  uint16
	}

	AvailRing struct {
		Flags     uint16
		Idx       uint16
		Ring      [QueueSize]uint16
		UsedEvent uint16
	}

	// padding for 4096 byte alignment
	_ [4096 - ((16*QueueSize + 6 + 2*QueueSize) % 4096)]uint8

	UsedRing struct {
		Flags uint16
		Idx   uint16
		Ring  [QueueSize]struct {
			Idx uint32
			Len uint32
		}
		availEvent uint16
	}
}

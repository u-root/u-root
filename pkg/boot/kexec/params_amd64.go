package kexec

// Params communicates boot information to
// the purgatory, and possibly the kernel.
type Params struct {
	Entry  uint64
	Params uint64
	_      [5]uint64
}

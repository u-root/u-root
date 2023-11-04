package iodev

type Noop struct {
	Port  uint64
	Psize uint64
}

func (n *Noop) Read(port uint64, data []byte) error {
	return nil
}

func (n *Noop) Write(port uint64, data []byte) error {
	return nil
}

func (n *Noop) IOPort() uint64 {
	return n.Port
}

func (n *Noop) Size() uint64 {
	return n.Psize
}

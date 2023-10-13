package serial

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	COM1Addr = 0x03f8
)

// Note that this identical interface is defined across
// multiple packages. It should be defined by the machine.

type IRQInjector interface {
	InjectSerialIRQ() error
}

type Serial struct {
	IER byte
	LCR byte

	inputChan chan byte

	irqInjector IRQInjector
}

func New(irqInjector IRQInjector) (*Serial, error) {
	s := &Serial{
		IER: 0, LCR: 0,
		inputChan:   make(chan byte, 10000),
		irqInjector: irqInjector,
	}

	return s, nil
}

func (s *Serial) GetInputChan() chan<- byte {
	return s.inputChan
}

func (s *Serial) dlab() bool {
	return s.LCR&0x80 != 0
}

func (s *Serial) In(port uint64, values []byte) error {
	port -= COM1Addr

	switch {
	case port == 0 && !s.dlab():
		// RBR
		if len(s.inputChan) > 0 {
			values[0] = <-s.inputChan
		}
	case port == 0 && s.dlab():
		// DLL
		values[0] = 0xc // baud rate 9600
	case port == 1 && !s.dlab():
		// IER
		values[0] = s.IER
	case port == 1 && s.dlab():
		// DLM
		values[0] = 0x0 // baud rate 9600
	case port == 2:
		// IIR
	case port == 3:
		// LCR
	case port == 4:
		// MCR
	case port == 5:
		// LSR
		values[0] |= 0x20 // Empty Transmitter Holding Register
		values[0] |= 0x40 // Empty Data Holding Registers

		if len(s.inputChan) > 0 {
			values[0] |= 0x1 // Data Ready
		}
	case port == 6:
		// MSR
		break
	}

	return nil
}

func (s *Serial) Out(port uint64, values []byte) error {
	port -= COM1Addr

	var err error

	switch {
	case port == 0 && !s.dlab():
		// THR
		fmt.Printf("%c", values[0])
	case port == 0 && s.dlab():
		// DLL
	case port == 1 && !s.dlab():
		// IER
		s.IER = values[0]
		if s.IER != 0 {
			err = s.irqInjector.InjectSerialIRQ()
		}
	case port == 1 && s.dlab():
		// DLM
	case port == 2:
		// FCR
	case port == 3:
		// LCR
		s.LCR = values[0]
	case port == 4:
		// MCR
	default:
		// factory test or not used
		break
	}

	return err
}

func (s *Serial) StartSerial(in bufio.Reader, restoreMode func(), irqInject func() error) {
	var before byte = 0

	go func() {
		for {
			b, err := in.ReadByte()
			if err != nil {
				log.Printf("%v", err)

				break
			}
			s.GetInputChan() <- b

			if len(s.GetInputChan()) > 0 {
				if err := irqInject(); err != nil {
					log.Printf("InjectSerialIRQ: %v", err)
				}
			}

			if before == 0x1 && b == 'x' {
				restoreMode()
				os.Exit(0)
			}

			before = b
		}
	}()
}

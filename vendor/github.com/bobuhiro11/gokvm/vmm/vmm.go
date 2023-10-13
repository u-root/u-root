package vmm

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/bobuhiro11/gokvm/machine"
	"github.com/bobuhiro11/gokvm/pvh"
	"github.com/bobuhiro11/gokvm/term"
)

// Config defines the configuration of the
// virtual machine, as determined by flags.
type Config struct {
	Debug      bool
	Dev        string
	Kernel     string
	Initrd     string
	Params     string
	TapIfName  string
	Disk       string
	NCPUs      int
	MemSize    int
	TraceCount int
}

type VMM struct {
	*machine.Machine
	Config
}

func New(c Config) *VMM {
	return &VMM{
		Machine: nil,
		Config:  c,
	}
}

// Init instantiates a machine.
func (v *VMM) Init() error {
	m, err := machine.New(v.Dev, v.NCPUs, v.MemSize)
	if err != nil {
		return err
	}

	if len(v.TapIfName) > 0 {
		if err := m.AddTapIf(v.TapIfName); err != nil {
			return err
		}
	}

	if len(v.Disk) > 0 {
		if err := m.AddDisk(v.Disk); err != nil {
			return err
		}
	}

	v.Machine = m

	return nil
}

func (v *VMM) Setup() error {
	var initrd *os.File
	// Kernel arg required to load kernel or firmware image
	kern, err := os.Open(v.Kernel)
	if err != nil {
		return err
	}

	isPVH, err := pvh.CheckPVH(kern)
	if err != nil {
		return err
	}

	if v.Initrd != "" {
		initrd, err = os.Open(v.Initrd)
		if err != nil {
			return err
		}
	}

	if isPVH {
		if err := v.Machine.LoadPVH(kern, initrd, v.Params); err != nil {
			return err
		}
	} else {
		if err := v.Machine.LoadLinux(kern, initrd, v.Params); err != nil {
			return err
		}
	}

	return nil
}

func (v *VMM) Boot() error {
	var err error

	var wg sync.WaitGroup

	trace := v.TraceCount > 0
	if err := v.SingleStep(trace); err != nil {
		return fmt.Errorf("setting trace to %v:%w", trace, err)
	}

	for cpu := 0; cpu < v.NCPUs; cpu++ {
		fmt.Printf("Start CPU %d of %d\r\n", cpu, v.NCPUs)
		v.StartVCPU(cpu, v.TraceCount, &wg)
		wg.Add(1)
	}

	if !term.IsTerminal() {
		fmt.Fprintln(os.Stderr, "this is not terminal and does not accept input")
		select {}
	}

	restoreMode, err := term.SetRawMode()
	if err != nil {
		return err
	}

	defer restoreMode()

	if err := v.SingleStep(trace); err != nil {
		log.Printf("SingleStep(%v): %v", trace, err)

		return err
	}

	in := bufio.NewReader(os.Stdin)

	v.GetSerial().StartSerial(*in, restoreMode, v.InjectSerialIRQ)

	fmt.Printf("Waiting for CPUs to exit\r\n")
	wg.Wait()
	fmt.Printf("All cpus done\n\r")

	return nil
}

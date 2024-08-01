// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ip manipulates network addresses, interfaces, routing, and other config.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/gopacket/gopacket"
	pcap "github.com/packetcap/go-pcap"

	"github.com/gopacket/gopacket/layers"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type flags struct {
	help      bool
	countPkg  int
	filter    string
	snapLen   int
	device    string
	noPromisc bool
}

const tcpdumpHelp = `Usage: tcpdump [ -h ][ -c count ] [ -i interface ][ -s snaplen ] [ -no-promiscuous-mode ] [ expression ]`

func parseFlags(args []string, out io.Writer) (cmd, error) {
	opts := flags{}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.IntVar(&opts.countPkg, "c", 0, "Exit after receiving count packets")
	fs.BoolVar(&opts.help, "help", false, "Print help message")
	fs.BoolVar(&opts.help, "h", false, "Print help message")
	fs.StringVar(&opts.device, "i", "", "Listen on interface")
	fs.StringVar(&opts.device, "interface", "", "Listen on interface")
	fs.IntVar(&opts.snapLen, "s", 262144, "narf snaplen bytes of data from each packet rather than the default of 262144 bytes")
	fs.IntVar(&opts.snapLen, "snapshot-length", 262144, "narf snaplen bytes of data from each packet rather than the default of 262144 bytes")
	fs.BoolVar(&opts.noPromisc, "p", true, "Set non-promiscuous mode")
	fs.BoolVar(&opts.noPromisc, "no-promiscuous-mode", true, "Set non-promiscuous mode")

	fs.Usage = func() {
		fmt.Fprintf(out, "%s\n\n", tcpdumpHelp)

		fs.PrintDefaults()
	}

	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))

	filter := ""
	if fs.NArg() > 0 {
		for _, arg := range fs.Args() {
			filter += arg + " "
		}
	}

	opts.filter = filter

	return cmd{Opts: opts, out: out}, nil
}

type cmd struct {
	out  io.Writer
	Opts flags
}

func (cmd *cmd) run() error {
	var (
		src *pcap.Handle
		err error
	)

	if cmd.Opts.help {
		fmt.Println(tcpdumpHelp)

		return nil
	}

	if cmd.Opts.device == "" {
		return fmt.Errorf("no device specified")
	}

	sigChan := make(chan os.Signal, 1)
	doneChan := make(chan bool, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		doneChan <- true
	}()

	if src, err = pcap.OpenLive(cmd.Opts.device, int32(cmd.Opts.snapLen), !cmd.Opts.noPromisc, 0, false); err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			return fmt.Errorf("you don't have permission to capture on that/these device(s)")
		}

		return err

	}
	defer src.Close()

	if err := src.SetBPFFilter(cmd.Opts.filter); err != nil {
		return err
	}

	packetSource := gopacket.NewPacketSource(src, layers.LinkTypeEthernet)
	packetSource.NoCopy = true

	fmt.Fprintf(cmd.out, "tcpdump: listening on %s\n", cmd.Opts.device)

	capturedPackets := 0

	for {
		select {
		case packet := <-packetSource.Packets():
			capturedPackets++
			if cmd.Opts.countPkg > 0 && capturedPackets >= cmd.Opts.countPkg {
				doneChan <- true
			}
			cmd.processPacket(packet, capturedPackets)
		case <-doneChan:
			fmt.Fprintf(cmd.out, "\n%d packets captured\n", capturedPackets)

			return nil
		}
	}
}

func (cmd *cmd) processPacket(packet gopacket.Packet, num int) {
	if packet == nil || packet.NetworkLayer() == nil || packet.TransportLayer() == nil {
		return
	}

	if err := packet.ErrorLayer(); err != nil {
		fmt.Fprintf(cmd.out, "skipping packet no. %d: %v", num, err)

		return
	}

	networkSrc, networkDst := packet.NetworkLayer().NetworkFlow().Endpoints()
	transportSrc, transportDst := packet.TransportLayer().TransportFlow().Endpoints()

	srcIP, dstIP := networkSrc.String(), networkDst.String()
	srcPort, dstPort := transportSrc.String(), transportDst.String()

	if dstPort == "53" {
		dstPort = "domain"
	}

	srcHostnames, err := net.LookupAddr(srcIP)
	if err != nil || len(srcHostnames) == 0 {
		srcHostnames = []string{srcIP}
	}

	dstHostnames, err := net.LookupAddr(dstIP)
	if err != nil || len(dstHostnames) == 0 {
		dstHostnames = []string{dstIP}
	}

	applicationData := ""
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		switch layer := applicationLayer.(type) {
		case *layers.DNS:
			applicationData = fmt.Sprintf("%d", layer.ID)

			if len(layer.Answers)+len(layer.Authorities)+len(layer.Additionals) >= 1 {
				applicationData += fmt.Sprintf(" %d/%d/%d ", len(layer.Answers), len(layer.Authorities), len(layer.Additionals))
			}

			for _, question := range layer.Questions {
				applicationData += fmt.Sprintf(" %s %s,", question.Type.String(), question.Name)
			}

			for _, answer := range layer.Answers {
				applicationData += answer.String() + ","
			}

			applicationData = strings.TrimRight(applicationData, ",")

			applicationData += fmt.Sprintf((" (%d)"), len(layer.Contents))
		case *layers.TLS:
			applicationData = "TLS"
		default:
		}
	}

	fmt.Fprintf(cmd.out, "%s %s %s %s%s > %s%s: %s\n",
		packet.Metadata().Timestamp.Format("15:04:05.000000"),
		packet.NetworkLayer().NetworkFlow().EndpointType(),
		cmd.Opts.device,
		srcHostnames[0],
		srcPort,
		dstHostnames[0],
		dstPort,
		applicationData)
}

func main() {
	// Disable logrus logging
	logrus.SetLevel(logrus.PanicLevel)

	cmd, err := parseFlags(os.Args, os.Stdout)
	if err != nil {
		log.Fatalf("tcpdump: %v", err)
	}

	err = cmd.run()
	if err != nil {
		log.Fatalf("tcpdump: %v", err)
	}
}

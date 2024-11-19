// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build !tinygo || tinygo.enable

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gopacket/gopacket"
	pcap "github.com/packetcap/go-pcap"

	"github.com/gopacket/gopacket/layers"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

type flags struct {
	Help                   bool
	CountPkg               int
	Filter                 string
	SnapshotLength         int
	Device                 string
	NoPromisc              bool
	Count                  bool
	ListDevices            bool
	Numerical              bool
	Number                 bool
	T                      bool
	TT                     bool
	TTT                    bool
	TTTT                   bool
	TTTTT                  bool
	FirstPacketProcessed   bool
	Verbose                bool
	Data                   bool
	DataWithHeader         bool
	Quiet                  bool
	ASCII                  bool
	Ether                  bool
	FilterFile             string
	TimeStampInNanoSeconds bool
	IcmpOnly               bool
}

const tcpdumpHelp = `       tcpdump [ -ADehnpqtvx# ] [ -icmp ]
                [ -c count ] [ --count ] [ -F file ][ -i interface ]
			    [ --number ] [ --print ] [ -s snaplen ] [ --nano ] 
				[ EXPRESSION ]
	EXPRESSION := [ EXPRESSION ] [ and ] [ or ] [ not ] 
				  [ gateway host ] [ proto protocol ] [ ether type ] [ src host ]
				  [ dst host ] [ net net ] [ port port ] [ portrange X-Y ] 
				  [ ip host ] [ ip4 ] [ ip6 ] [ tcp ] [ udp ]`

// parseFlags parses the flags and returns the cmd.
func parseFlags(args []string, out io.Writer) (cmd, error) {
	opts := flags{}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.IntVar(&opts.CountPkg, "c", 0, "Exit after receiving count packets")
	fs.BoolVar(&opts.Help, "help", false, "Print help message")
	fs.BoolVar(&opts.Help, "h", false, "Print help message")
	fs.StringVar(&opts.Device, "i", "", "Listen on interface")
	fs.StringVar(&opts.Device, "interface", "", "Listen on interface")
	fs.IntVar(&opts.SnapshotLength, "s", 262144, "snarf snaplen bytes of data from each packet rather than the default of 262144 bytes")
	fs.IntVar(&opts.SnapshotLength, "snapshot-length", 262144, "narf snaplen bytes of data from each packet rather than the default of 262144 bytes")
	fs.BoolVar(&opts.NoPromisc, "p", false, "Set non-promiscuous mode")
	fs.BoolVar(&opts.NoPromisc, "no-promiscuous-mode", false, "Set non-promiscuous mode")
	fs.BoolVar(&opts.Count, "count", false, "Print only the number of packets captured")
	fs.BoolVar(&opts.ListDevices, "D", false, "Print  the  list of the network interfaces available on the system and on which tcpdump can capture packets")
	fs.BoolVar(&opts.ListDevices, "list-interfaces", false, "Print  the  list of the network interfaces available on the system and on which tcpdump can capture packets")
	fs.BoolVar(&opts.Numerical, "n", false, "Don't convert addresses (i.e., host addresses, port numbers, etc.) to names")
	fs.BoolVar(&opts.Number, "#", false, " Print an optional packet number at the beginning of the line")
	fs.BoolVar(&opts.Number, "number", false, " Print an optional packet number at the beginning of the line")
	fs.BoolVar(&opts.IcmpOnly, "icmp", false, "Only capture ICMP packets")
	fs.BoolVar(&opts.Ether, "e", false, "Print the link-level header on each dump line.  This can be used, for example, to print MAC layer addresses for protocols such as Ethernet and IEEE 802.11.")
	fs.BoolVar(&opts.T, "t", false, "Don't print a timestamp on each dump line")
	fs.BoolVar(&opts.TT, "tt", false, "Print the timestamp, as seconds since January 1, 1970, 00:00:00, UTC, and fractions of a second since that time, on each dump line")
	fs.BoolVar(&opts.TTT, "ttt", false, "Print a delta (microsecond or nanosecond resolution depending on the --time-stamp-precision option) between current and previous line on each dump line.  The default is microsecond resolution")
	fs.BoolVar(&opts.TTTT, "tttt", false, "Print a timestamp, as hours, minutes, seconds, and fractions of a second since midnight, preceded by the date, on each dump line")
	fs.BoolVar(&opts.TTTTT, "ttttt", false, "Print  a delta (microsecond or nanosecond resolution depending on the --time-stamp-precision option) between current and first line on each dump line.  The default is microsecond resolution")
	fs.BoolVar(&opts.TimeStampInNanoSeconds, "nano", false, "Print the timestamp in nanosecond resolution (instead of microseconds)")
	fs.BoolVar(&opts.Data, "x", false, "When parsing and printing, in addition to printing the headers of each packet, print the data of each packet (minus its link level header) in hex")
	fs.BoolVar(&opts.DataWithHeader, "xx", false, "When parsing and printing, in addition to printing the headers of each packet, print the data of each packet (including its link level header) in hex")
	fs.StringVar(&opts.FilterFile, "F", "", "Use file as input for the filter expression.  An additional expression given on the command line is ignored.")
	fs.BoolVar(&opts.ASCII, "A", false, "Print each packet (minus its link level header) in ASCII.  Handy for capturing web pages")
	fs.BoolVar(&opts.Quiet, "q", false, "Quiet output. Print less protocol information so output lines are shorter")
	fs.BoolVar(&opts.Verbose, "v", false, "When parsing and printing, produce (slightly more) verbose output.  For example, the time to live, identification, total length and options in an IP packet are printed.  Also enables additional packet integrity checks such as verifying the IP and ICMP header checksum")
	fs.BoolVar(&opts.Verbose, "verbose", false, "When parsing and printing, produce (slightly more) verbose output.  For example, the time to live, identification, total length and options in an IP packet are printed.  Also enables additional packet integrity checks such as verifying the IP and ICMP header checksum")

	fs.Usage = func() {
		fmt.Fprintf(out, "%s\n\n", tcpdumpHelp)

		fs.SetOutput(out)
		fs.PrintDefaults()
	}

	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))

	if opts.Verbose && opts.Quiet {
		return cmd{}, fmt.Errorf("cannot use both -v and -q flags")
	}

	filter := ""
	if fs.NArg() > 0 {
		for _, arg := range fs.Args() {
			filter += arg + " "
		}
	}

	opts.Filter = filter

	if opts.FilterFile != "" {
		if data, err := os.ReadFile(opts.FilterFile); err == nil {
			opts.Filter = string(data)
		} else {
			return cmd{}, fmt.Errorf("failed to read filter file: %w", err)
		}
	}

	return cmd{Opts: opts, Out: out, usage: fs.Usage}, nil
}

type cmd struct {
	Out  io.Writer
	Opts flags

	usage func()
}

func (cmd *cmd) run() error {
	var (
		src *pcap.Handle
		err error
	)

	if cmd.Opts.Help {
		cmd.usage()

		return nil
	}

	if cmd.Opts.ListDevices {
		return listDevices()
	}

	if cmd.Opts.Device == "" {
		return fmt.Errorf("no device specified")
	}

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-sigChan
		cancel()
	}()

	if src, err = pcap.OpenLive(cmd.Opts.Device, int32(cmd.Opts.SnapshotLength), !cmd.Opts.NoPromisc, 0, false); err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			return fmt.Errorf("you don't have permission to capture on that/these device(s)")
		}

		return err

	}
	defer src.Close()

	if err := src.SetBPFFilter(cmd.Opts.Filter); err != nil {
		return err
	}

	packetSource := gopacket.NewPacketSource(src, layers.LinkTypeEthernet)
	packetSource.NoCopy = true

	fmt.Fprintf(cmd.Out, "tcpdump: verbose output suppressed, use -v for full protocol decode\nlistening on %s, link-type %d, snapshot length %d bytes\n", cmd.Opts.Device, src.LinkType(), cmd.Opts.SnapshotLength)

	var (
		capturedPackets int
		timeStamp       time.Time
	)
	if cmd.Opts.TTTTT {
		timeStamp = time.Now()
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintf(cmd.Out, "\n%d packets captured\n", capturedPackets)

			return nil
		case packet := <-packetSource.PacketsCtx(ctx):
			capturedPackets++
			if cmd.Opts.CountPkg > 0 && capturedPackets >= cmd.Opts.CountPkg {
				return nil
			}

			if !cmd.Opts.Count {
				pkgTime := cmd.processPacket(packet, capturedPackets, timeStamp)

				if cmd.Opts.TTT {
					timeStamp = pkgTime
				}
			}

		}
	}
}

// processPacket processes a packet and prints the output to the output writer.
// A timestamp of the packet is returned.
func (cmd *cmd) processPacket(packet gopacket.Packet, num int, lastPkgTimeStamp time.Time) time.Time {
	var (
		no        string
		srcAddr   string
		srcPort   string
		dstAddr   string
		dstPort   string
		timeStamp string
	)

	if cmd.Opts.Number {
		no = fmt.Sprintf("%d  ", num)
	}

	if packet == nil {
		return lastPkgTimeStamp
	}

	if err := packet.ErrorLayer(); err != nil {
		fmt.Fprintf(cmd.Out, "skipping packet no. %d: %v\n", num, err)

		return lastPkgTimeStamp
	}

	ethernetLayer := packet.LinkLayer()
	if ethernetLayer == nil {
		return lastPkgTimeStamp
	}

	networkLayer := packet.NetworkLayer()

	if networkLayer == nil {
		return lastPkgTimeStamp
	}

	etherInfo := cmd.ethernetInfo(ethernetLayer, networkLayer)

	if cmd.Opts.Verbose {
		switch layer := networkLayer.(type) {
		case *layers.IPv4:
			etherInfo += fmt.Sprintf(" (tos 0x%x, ttl %d, id %d, offset %d, flags [%s], proto %s (%d), length %d)\n", layer.TOS, layer.TTL, layer.Id, layer.FragOffset, layer.Flags, layer.Protocol, layer.Protocol, len(layer.Contents)+len(layer.Payload))
		case *layers.IPv6:
			etherInfo += fmt.Sprintf(" (flowlabel 0x%x, hlim %d, next-header %s (%d), payload length: %d)\n", layer.FlowLabel, layer.HopLimit, layer.NextHeader, layer.NextHeader, len(layer.Payload))
		}
	}

	networkSrc, networkDst := networkLayer.NetworkFlow().Endpoints()

	srcAddr, dstAddr = networkSrc.String(), networkDst.String()

	if srcHostNames, err := net.LookupAddr(srcAddr); err == nil && len(srcHostNames) > 0 && !cmd.Opts.Numerical {
		srcAddr = srcHostNames[0]
	}

	if dstHostNames, err := net.LookupAddr(dstAddr); err == nil && len(dstHostNames) > 0 && !cmd.Opts.Numerical {
		dstAddr = dstHostNames[0]
	}

	// Append a dot to the end of the addresses if it doesn't have one
	if !strings.HasSuffix(srcAddr, ".") {
		srcAddr += "."
	}
	if !strings.HasSuffix(dstAddr, ".") {
		dstAddr += "."
	}

	data := parseICMP(packet)

	if cmd.Opts.IcmpOnly && data == "" {
		return lastPkgTimeStamp
	}

	transportLayer := packet.TransportLayer()

	// Set the source and destination ports, if a transport layer is present
	if transportLayer != nil {
		transportSrc, transportDst := transportLayer.TransportFlow().Endpoints()

		srcPort, dstPort = transportSrc.String(), cmd.wellKnownPorts(transportDst.String())
	}

	// parse the application layer
	applicationLayer := packet.ApplicationLayer()

	if applicationLayer != nil && !cmd.Opts.Quiet {
		switch layer := applicationLayer.(type) {
		case *layers.DNS:
			data = dnsData(layer)
		}
	}

	if data == "" {
		var length int

		if applicationLayer != nil {
			length = len(applicationLayer.LayerContents())
		}

		switch layer := transportLayer.(type) {
		case *layers.TCP:
			data = tcpData(layer, length, cmd.Opts.Verbose, cmd.Opts.Quiet)
		case *layers.UDP:
			data = fmt.Sprintf("UDP, length %d", length)
		case *layers.UDPLite:
			data = fmt.Sprintf("UDPLite, length %d", length)
		case nil:
			data = fmt.Sprintf("unknown transport layer, length %d", length)
		default:
			data = fmt.Sprintf("%s, length %d", layer.LayerType(), length)
		}
	}

	pkgTimeStamp := packet.Metadata().Timestamp

	timeStamp = cmd.parseTimeStamp(pkgTimeStamp, lastPkgTimeStamp)

	fmt.Fprintf(cmd.Out, "%s%s %s %s%s > %s%s: %s\n",
		no,
		timeStamp,
		etherInfo,
		srcAddr,
		srcPort,
		dstAddr,
		dstPort,
		data)

	switch {
	case cmd.Opts.ASCII:
		content := []byte("")
		if applicationLayer != nil {
			content = applicationLayer.LayerContents()
		}
		fmt.Fprintf(cmd.Out, "%s\n", content)
	case cmd.Opts.Data:
		fmt.Fprintf(cmd.Out, "%s\n", formatPacketData(packet.Data()[14:]))
	case cmd.Opts.DataWithHeader:
		fmt.Fprintf(cmd.Out, "%s\n", formatPacketData(packet.Data()))
	}

	return pkgTimeStamp
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

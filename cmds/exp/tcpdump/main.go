// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
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
	help                   bool
	countPkg               int
	filter                 string
	snapshotLength         int
	device                 string
	noPromisc              bool
	count                  bool
	listDevices            bool
	numerical              bool
	number                 bool
	t                      bool
	tt                     bool
	ttt                    bool
	tttt                   bool
	ttttt                  bool
	firstPacketProcessed   bool
	verbose                bool
	data                   bool
	dataWithHeader         bool
	quiet                  bool
	ascii                  bool
	ether                  bool
	filterFile             string
	timeStampInNanoSeconds bool
	icmpOnly               bool
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
	fs.IntVar(&opts.countPkg, "c", 0, "Exit after receiving count packets")
	fs.BoolVar(&opts.help, "help", false, "Print help message")
	fs.BoolVar(&opts.help, "h", false, "Print help message")
	fs.StringVar(&opts.device, "i", "", "Listen on interface")
	fs.StringVar(&opts.device, "interface", "", "Listen on interface")
	fs.IntVar(&opts.snapshotLength, "s", 262144, "snarf snaplen bytes of data from each packet rather than the default of 262144 bytes")
	fs.IntVar(&opts.snapshotLength, "snapshot-length", 262144, "narf snaplen bytes of data from each packet rather than the default of 262144 bytes")
	fs.BoolVar(&opts.noPromisc, "p", false, "Set non-promiscuous mode")
	fs.BoolVar(&opts.noPromisc, "no-promiscuous-mode", false, "Set non-promiscuous mode")
	fs.BoolVar(&opts.count, "count", false, "Print only the number of packets captured")
	fs.BoolVar(&opts.listDevices, "D", false, "Print  the  list of the network interfaces available on the system and on which tcpdump can capture packets")
	fs.BoolVar(&opts.listDevices, "list-interfaces", false, "Print  the  list of the network interfaces available on the system and on which tcpdump can capture packets")
	fs.BoolVar(&opts.numerical, "n", false, "Don't convert addresses (i.e., host addresses, port numbers, etc.) to names")
	fs.BoolVar(&opts.number, "#", false, " Print an optional packet number at the beginning of the line")
	fs.BoolVar(&opts.number, "number", false, " Print an optional packet number at the beginning of the line")
	fs.BoolVar(&opts.icmpOnly, "icmp", false, "Only capture ICMP packets")
	fs.BoolVar(&opts.ether, "e", false, "Print the link-level header on each dump line.  This can be used, for example, to print MAC layer addresses for protocols such as Ethernet and IEEE 802.11.")
	fs.BoolVar(&opts.t, "t", false, "Don't print a timestamp on each dump line")
	fs.BoolVar(&opts.tt, "tt", false, "Print the timestamp, as seconds since January 1, 1970, 00:00:00, UTC, and fractions of a second since that time, on each dump line")
	fs.BoolVar(&opts.ttt, "ttt", false, "Print a delta (microsecond or nanosecond resolution depending on the --time-stamp-precision option) between current and previous line on each dump line.  The default is microsecond resolution")
	fs.BoolVar(&opts.tttt, "tttt", false, "Print a timestamp, as hours, minutes, seconds, and fractions of a second since midnight, preceded by the date, on each dump line")
	fs.BoolVar(&opts.ttttt, "ttttt", false, "Print  a delta (microsecond or nanosecond resolution depending on the --time-stamp-precision option) between current and first line on each dump line.  The default is microsecond resolution")
	fs.BoolVar(&opts.timeStampInNanoSeconds, "nano", false, "Print the timestamp in nanosecond resolution (instead of microseconds)")
	fs.BoolVar(&opts.data, "x", false, "When parsing and printing, in addition to printing the headers of each packet, print the data of each packet (minus its link level header) in hex")
	fs.BoolVar(&opts.dataWithHeader, "xx", false, "When parsing and printing, in addition to printing the headers of each packet, print the data of each packet (including its link level header) in hex")
	fs.StringVar(&opts.filterFile, "F", "", "Use file as input for the filter expression.  An additional expression given on the command line is ignored.")
	fs.BoolVar(&opts.ascii, "A", false, "Print each packet (minus its link level header) in ASCII.  Handy for capturing web pages")
	fs.BoolVar(&opts.quiet, "q", false, "Quiet output. Print less protocol information so output lines are shorter")
	fs.BoolVar(&opts.verbose, "v", false, "When parsing and printing, produce (slightly more) verbose output.  For example, the time to live, identification, total length and options in an IP packet are printed.  Also enables additional packet integrity checks such as verifying the IP and ICMP header checksum")
	fs.BoolVar(&opts.verbose, "verbose", false, "When parsing and printing, produce (slightly more) verbose output.  For example, the time to live, identification, total length and options in an IP packet are printed.  Also enables additional packet integrity checks such as verifying the IP and ICMP header checksum")

	fs.Usage = func() {
		fmt.Fprintf(out, "%s\n\n", tcpdumpHelp)

		fs.PrintDefaults()
	}

	fs.Parse(unixflag.ArgsToGoArgs(args[1:]))

	if opts.verbose && opts.quiet {
		return cmd{}, fmt.Errorf("cannot use both -v and -q flags")
	}

	filter := ""
	if fs.NArg() > 0 {
		for _, arg := range fs.Args() {
			filter += arg + " "
		}
	}

	opts.filter = filter

	if opts.filterFile != "" {
		if data, err := os.ReadFile(opts.filterFile); err == nil {
			opts.filter = string(data)
		} else {
			return cmd{}, fmt.Errorf("failed to read filter file: %v", err)
		}
	}

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

	if cmd.Opts.listDevices {
		return listDevices()
	}

	if cmd.Opts.device == "" {
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

	if src, err = pcap.OpenLive(cmd.Opts.device, int32(cmd.Opts.snapshotLength), !cmd.Opts.noPromisc, 0, false); err != nil {
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

	fmt.Fprintf(cmd.out, "tcpdump: verbose output suppressed, use -v for full protocol decode\nlistening on %s, link-type %d, snapshot length %d bytes\n", cmd.Opts.device, src.LinkType(), cmd.Opts.snapshotLength)

	var (
		capturedPackets int
		timeStamp       time.Time
	)
	if cmd.Opts.ttttt {
		timeStamp = time.Now()
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintf(cmd.out, "\n%d packets captured\n", capturedPackets)

			return nil
		case packet := <-packetSource.PacketsCtx(ctx):
			capturedPackets++
			if cmd.Opts.countPkg > 0 && capturedPackets >= cmd.Opts.countPkg {
				return nil
			}

			if !cmd.Opts.count {
				pkgTime := cmd.processPacket(packet, capturedPackets, timeStamp)

				if cmd.Opts.ttt {
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

	if cmd.Opts.number {
		no = fmt.Sprintf("%d  ", num)
	}

	if packet == nil {
		return lastPkgTimeStamp
	}

	if err := packet.ErrorLayer(); err != nil {
		fmt.Fprintf(cmd.out, "skipping packet no. %d: %v\n", num, err)

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

	if cmd.Opts.verbose {
		switch layer := networkLayer.(type) {
		case *layers.IPv4:
			etherInfo += fmt.Sprintf(" (tos 0x%x, ttl %d, id %d, offset %d, flags [%s], proto %s (%d), length %d)\n", layer.TOS, layer.TTL, layer.Id, layer.FragOffset, layer.Flags, layer.Protocol, layer.Protocol, len(layer.Contents)+len(layer.Payload))
		case *layers.IPv6:
			etherInfo += fmt.Sprintf(" (flowlabel 0x%x, hlim %d, next-header %s (%d), payload length: %d)\n", layer.FlowLabel, layer.HopLimit, layer.NextHeader, layer.NextHeader, len(layer.Payload))
		}
	}

	networkSrc, networkDst := networkLayer.NetworkFlow().Endpoints()

	srcAddr, dstAddr = networkSrc.String(), networkDst.String()

	if srcHostNames, err := net.LookupAddr(srcAddr); err == nil && len(srcHostNames) > 0 && !cmd.Opts.numerical {
		srcAddr = srcHostNames[0]
	}

	if dstHostNames, err := net.LookupAddr(dstAddr); err == nil && len(dstHostNames) > 0 && !cmd.Opts.numerical {
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

	if cmd.Opts.icmpOnly && data == "" {
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

	if applicationLayer != nil && !cmd.Opts.quiet {
		switch layer := applicationLayer.(type) {
		case *layers.DNS:
			data = dnsData(layer)
		}
	}

	if data == "" {
		var length int

		if applicationLayer != nil {
			length = len(applicationLayer.LayerContents())
		} else {
			length = 0
		}

		switch layer := transportLayer.(type) {
		case *layers.TCP:
			data = tcpData(layer, length, cmd.Opts.verbose, cmd.Opts.quiet)
		case *layers.UDP:
			data = fmt.Sprintf("UDP, length %d", length)
		case *layers.UDPLite:
			data = fmt.Sprintf("UDPLite, length %d", length)
		default:
			data = fmt.Sprintf("%s, length %d", layer.LayerType(), length)
		}
	}

	pkgTimeStamp := packet.Metadata().Timestamp

	timeStamp = cmd.parseTimeStamp(pkgTimeStamp, lastPkgTimeStamp)

	fmt.Fprintf(cmd.out, "%s%s %s %s%s > %s%s: %s\n",
		no,
		timeStamp,
		etherInfo,
		srcAddr,
		srcPort,
		dstAddr,
		dstPort,
		data)

	switch {
	case cmd.Opts.ascii:
		fmt.Fprintf(cmd.out, "%s\n", applicationLayer.LayerContents())
	case cmd.Opts.data:
		fmt.Fprintf(cmd.out, "%s\n", formatPacketData(packet.Data()[14:]))
	case cmd.Opts.dataWithHeader:
		fmt.Fprintf(cmd.out, "%s\n", formatPacketData(packet.Data()))
	}

	return pkgTimeStamp
}

func formatPacketData(data []byte) string {
	var result string
	for i := 0; i < len(data); i += 16 {
		// Print the offset
		result += fmt.Sprintf("0x%04x:  ", i)

		// Print the hex values
		for j := 0; j < 16; j++ {
			if i+j < len(data) {
				result += fmt.Sprintf("%02x", data[i+j])
			} else {
				result += "  "
			}
			if j%2 == 1 {
				result += " "
			}
		}
		result += "\n"
	}
	return result
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

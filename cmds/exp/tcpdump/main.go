// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ip manipulates network addresses, interfaces, routing, and other config.
package main

import (
	"context"
	"encoding/binary"
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
	"github.com/vishvananda/netlink"

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
	verbose                bool
	writeToFile            string
	readFromFile           string
	data                   bool
	dataWithHeader         bool
	absoluteTCPSeq         bool
	quiet                  bool
	direction              bool
	alwaysPrint            bool
	ascii                  bool
	ether                  bool
	filterFile             string
	timeStampInNanoSeconds bool
}

const tcpdumpHelp = `       tcpdump [ -ADehnpqStvx# ] 
                [ -c count ] [ --count ] [ -F file ][ -i interface ]
			    [ --number ] [ --print ] [ -Q in|out|inout ] [ -r file ] 
				[ -s snaplen ][ -w file ] [ --nano ] [ expression ]`

func parseFlags(args []string, out io.Writer) (cmd, error) {
	opts := flags{}

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.IntVar(&opts.countPkg, "c", 0, "Exit after receiving count packets")
	fs.BoolVar(&opts.help, "help", false, "Print help message")
	fs.BoolVar(&opts.help, "h", false, "Print help message")
	fs.StringVar(&opts.device, "i", "", "Listen on interface")
	fs.StringVar(&opts.device, "interface", "", "Listen on interface")
	fs.IntVar(&opts.snapshotLength, "s", 262144, "narf snaplen bytes of data from each packet rather than the default of 262144 bytes")
	fs.IntVar(&opts.snapshotLength, "snapshot-length", 262144, "narf snaplen bytes of data from each packet rather than the default of 262144 bytes")
	fs.BoolVar(&opts.noPromisc, "p", true, "Set non-promiscuous mode")
	fs.BoolVar(&opts.noPromisc, "no-promiscuous-mode", true, "Set non-promiscuous mode")
	fs.BoolVar(&opts.count, "count", false, "Print only the number of packets captured")
	fs.BoolVar(&opts.listDevices, "D", false, "Print  the  list of the network interfaces available on the system and on which tcpdump can capture packets")
	fs.BoolVar(&opts.listDevices, "list-interfaces", false, "Print  the  list of the network interfaces available on the system and on which tcpdump can capture packets")
	fs.BoolVar(&opts.numerical, "n", false, "Don't convert addresses (i.e., host addresses, port numbers, etc.) to names")
	fs.BoolVar(&opts.number, "#", false, " Print an optional packet number at the beginning of the line")
	fs.BoolVar(&opts.number, "number", false, " Print an optional packet number at the beginning of the line")
	// TODO: Implement remaining flags
	fs.BoolVar(&opts.t, "t", false, "Don't print a timestamp on each dump line")
	fs.BoolVar(&opts.tt, "tt", false, "Print the timestamp, as seconds since January 1, 1970, 00:00:00, UTC, and fractions of a second since that time, on each dump line")
	fs.BoolVar(&opts.ttt, "ttt", false, "Print a delta (microsecond or nanosecond resolution depending on the --time-stamp-precision option) between current and previous line on each dump line.  The default is microsecond resolution")
	fs.BoolVar(&opts.tttt, "tttt", false, "Print a timestamp, as hours, minutes, seconds, and fractions of a second since midnight, preceded by the date, on each dump line")
	fs.BoolVar(&opts.ttttt, "ttttt", false, "Print  a delta (microsecond or nanosecond resolution depending on the --time-stamp-precision option) between current and first line on each dump line.  The default is microsecond resolution")
	fs.BoolVar(&opts.verbose, "v", false, "When parsing and printing, produce (slightly more) verbose output.  For example, the time to live, identification, total length and options in an IP packet are printed.  Also enables additional packet integrity checks such as verifying the IP and ICMP header checksum")
	fs.BoolVar(&opts.verbose, "verbose", false, "When parsing and printing, produce (slightly more) verbose output.  For example, the time to live, identification, total length and options in an IP packet are printed.  Also enables additional packet integrity checks such as verifying the IP and ICMP header checksum")
	fs.StringVar(&opts.writeToFile, "w", "", "Write the raw packets to file rather than parsing and printing them out.  They can later be printed with the -r option.  Standard output is used if file is ``-''")
	fs.BoolVar(&opts.data, "x", false, "When parsing and printing, in addition to printing the headers of each packet, print the data of each packet (minus its link level header) in hex")
	fs.BoolVar(&opts.dataWithHeader, "xx", false, "When parsing and printing, in addition to printing the headers of each packet, print the data of each packet (including its link level header) in hex")
	fs.BoolVar(&opts.absoluteTCPSeq, "S", false, "Print absolute, rather than relative, TCP sequence numbers")
	fs.BoolVar(&opts.absoluteTCPSeq, "absolute-tcp-sequence-numbers", false, "Print absolute, rather than relative, TCP sequence numbers")
	fs.StringVar(&opts.readFromFile, "r", "", "Read packets from file (which was created with the -w option) rather than from a network interface")
	fs.BoolVar(&opts.quiet, "q", false, "Quiet output. Print less protocol information so output lines are shorter")
	fs.BoolVar(&opts.direction, "Q", false, "Choose send/receive direction direction for which packets should be captured. Possible values are `in', `out' and `inout'")
	fs.BoolVar(&opts.direction, "direction", false, "Choose send/receive direction direction for which packets should be captured. Possible values are `in', `out' and `inout'")
	fs.BoolVar(&opts.alwaysPrint, "print", false, "Print parsed packet output, even if the raw packets are being saved to a file with the -w flag.")
	fs.BoolVar(&opts.ascii, "A", false, "Print each packet (minus its link level header) in ASCII.  Handy for capturing web pages")
	fs.BoolVar(&opts.ether, "e", false, "Print the link-level header on each dump line.  This can be used, for example, to print MAC layer addresses for protocols such as Ethernet and IEEE 802.11.")
	fs.StringVar(&opts.filterFile, "F", "", "Use file as input for the filter expression.  An additional expression given on the command line is ignored.")
	fs.BoolVar(&opts.timeStampInNanoSeconds, "nano", false, "Print the timestamp in nanosecond resolution (instead of microseconds)")

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

	fmt.Fprintf(cmd.out, "tcpdump: listening on %s, link-type %d, snapshot length %d bytes\n", cmd.Opts.device, src.LinkType(), cmd.Opts.snapshotLength)

	capturedPackets := 0

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
				cmd.processPacket(packet, capturedPackets)
			}

		}
	}
}

func listDevices() error {
	links, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for idx, link := range links {
		fmt.Printf("%d.%s [%s]\n", idx, link.Attrs().Name, link.Attrs().OperState)
	}

	return nil
}

func (cmd *cmd) processPacket(packet gopacket.Packet, num int) {
	if packet == nil {
		return
	}

	if err := packet.ErrorLayer(); err != nil {
		fmt.Fprintf(cmd.out, "skipping packet no. %d: %v\n", num, err)

		return
	}

	if packet.NetworkLayer() == nil || packet.TransportLayer() == nil {
		return
	}

	networkLayer := packet.NetworkLayer()
	transportLayer := packet.TransportLayer()
	applicationLayer := packet.ApplicationLayer()

	networkSrc, networkDst := networkLayer.NetworkFlow().Endpoints()
	transportSrc, transportDst := transportLayer.TransportFlow().Endpoints()

	srcIP, dstIP := networkSrc.String(), cmd.wellKnownPorts(networkDst.String())
	srcPort, dstPort := transportSrc.String(), cmd.wellKnownPorts(transportDst.String())

	srcHostnames, err := net.LookupAddr(srcIP)
	if err != nil || len(srcHostnames) == 0 || cmd.Opts.numerical {
		srcHostnames = []string{srcIP}
	}

	dstHostnames, err := net.LookupAddr(dstIP)
	if err != nil || len(dstHostnames) == 0 || cmd.Opts.numerical {
		dstHostnames = []string{dstIP}
	}

	// Check if the hostname does not end with a dot and add one if necessary
	if !strings.HasSuffix(srcHostnames[0], ".") {
		srcHostnames[0] += "."
	}

	if !strings.HasSuffix(dstHostnames[0], ".") {
		dstHostnames[0] += "."
	}

	data := ""

	if applicationLayer != nil {
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
			data = tcpData(layer, length)
		case *layers.UDP:
			data = fmt.Sprintf("UDP, length %d", length)
		case *layers.UDPLite:
			data = fmt.Sprintf("UDPLite, length %d", length)
		default:
			data = fmt.Sprintf("%s, length %d", layer.LayerType(), length)
		}
	}

	no := ""
	if cmd.Opts.number {
		no = fmt.Sprintf("%d  ", num)
	}

	fmt.Fprintf(cmd.out, "%s%s %s %s %s%s > %s%s: %s\n",
		no,
		packet.Metadata().Timestamp.Format("15:04:05.000000"),
		networkLayer.NetworkFlow().EndpointType(),
		cmd.Opts.device,
		srcHostnames[0],
		srcPort,
		dstHostnames[0],
		dstPort,
		data)
}

func tcpData(layer *layers.TCP, length int) string {
	var data string

	flags := tcpFlags(*layer)
	opts := tcpOptions(layer.Options)

	data = fmt.Sprintf("Flags [%s], seq %d, ack %d, win %d, options [%s], length %d", flags, layer.Seq, layer.Ack, layer.Window, opts, length)

	return data
}

func tcpFlags(layer layers.TCP) string {
	var flags string
	if layer.PSH {
		flags += "P"
	}
	if layer.FIN {
		flags += "F"
	}
	if layer.SYN {
		flags += "S"
	}
	if layer.RST {
		flags += "R"
	}
	if layer.URG {
		flags += "U"
	}
	if layer.ECE {
		flags += "E"
	}
	if layer.CWR {
		flags += "C"
	}
	if layer.NS {
		flags += "N"
	}
	if layer.ACK {
		flags += "."
	}

	return flags
}

func tcpOptions(options []layers.TCPOption) string {
	var opts string

	for _, opt := range options {
		opts += tcpOptionToString(opt) + ","
	}

	return strings.TrimRight(opts, ",")
}

func tcpOptionToString(opt layers.TCPOption) string {
	switch opt.OptionType {
	case layers.TCPOptionKindMSS:
		if len(opt.OptionData) >= 2 {
			return fmt.Sprintf("%s val %v",
				opt.OptionType,
				binary.BigEndian.Uint16(opt.OptionData))
		}

	case layers.TCPOptionKindTimestamps:
		if len(opt.OptionData) == 8 {
			return fmt.Sprintf("%s val %v",
				opt.OptionType,
				binary.BigEndian.Uint32(opt.OptionData[:4]))
		}
	}

	return fmt.Sprintf("%s", opt.OptionType)
}

func dnsData(layer *layers.DNS) string {
	applicationData := fmt.Sprintf("%d", layer.ID)

	if layer.ResponseCode != layers.DNSResponseCodeNoErr {
		applicationData += fmt.Sprintf(" %s", layer.ResponseCode.String())
	}

	if layer.AA {
		applicationData += "*"
	} else if layer.RD {
		applicationData += "+"
	}

	if len(layer.Answers)+len(layer.Authorities)+len(layer.Additionals) > 1 {
		applicationData += fmt.Sprintf(" %d/%d/%d ", len(layer.Answers), len(layer.Authorities), len(layer.Additionals))
	}

	if layer.ResponseCode == layers.DNSResponseCodeNoErr {

		if len(layer.Answers) == 0 {
			for _, question := range layer.Questions {
				applicationData += fmt.Sprintf(" %s? %s, ", question.Type.String(), question.Name)
			}
		}

		for _, answer := range layer.Answers {
			applicationData += answer.String() + ", "
		}

		applicationData = strings.TrimRight(applicationData, ", ")
	}

	applicationData += fmt.Sprintf((" (%d)"), len(layer.Contents))

	return applicationData
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

var wellKnownPortsMap = map[string]string{
	"20":  "ftp-data",
	"21":  "ftp",
	"22":  "ssh-scp",
	"23":  "telnet",
	"25":  "smtp",
	"53":  "domain",
	"80":  "http",
	"88":  "kerberos",
	"110": "pop3",
	"119": "nntp",
	"123": "ntp",
	"143": "imap",
	"443": "https",
	"465": "smtps",
	"563": "nntps",
	"989": "ftps-data",
	"990": "ftps",
	"993": "imaps",
	"995": "pop3s",
}

func (cmd cmd) wellKnownPorts(port string) string {
	if name, ok := wellKnownPortsMap[port]; ok && !cmd.Opts.numerical {
		return name
	}

	return port
}

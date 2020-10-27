package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/u-root/u-root/pkg/dhclient"
	"github.com/vishvananda/netlink"
	"github.com/xi2/xz"
)

var url = "https://blobs.9esec.io/os/nightly/Fedora-HWT-disk-31-buildserver-MBR.img.xz"

var (
	doDebug = flag.Bool("d", true, "Print debug output")
)

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

var banner = `

 _________________________________
< Booting with 9e is even hotter! >
 ---------------------------------
                                /
                     __       //
                     -\= \=\ //
                   --=_\=---//=--
                 -_==/  \/ //\/--
                  ==/   /O   O\==--
     _ _ _ _     /_/    \  ]  /--
    /\ ( (- \    /       ] ] ]==-
   (\ _\_\_\-\__/     \  (,_,)--
  (\_/                 \     \-
  \/      /       (   ( \  ] /)
  /      (         \   \_ \./ )
  (       \         \      )  \
  (       /\_ _ _ _ /---/ /\_  \
   \     / \     / ____/ /   \  \
    (   /   )   / /  /__ )   (  )
    (  )   / __/ '---'       / /
    \  /   \ \             _/ /
    ] ]     )_\_         /__\/
    /_\     ]___\
   (___)
`
var debug = func(string, ...interface{}) {}

func configureDHCPNetwork() error {
	log.Printf("Trying to configure network configuration dynamically...")

	link, err := findNetworkInterface()
	if err != nil {
		return err
	}

	var links []netlink.Link
	links = append(links, link)

	var level dhclient.LogLevel

	config := dhclient.Config{
		Timeout:  6 * time.Second,
		Retries:  4,
		LogLevel: level,
	}

	r := dhclient.SendRequests(context.TODO(), links, true, false, config, 30*time.Second)
	for result := range r {
		if result.Err == nil {
			return result.Lease.Configure()
		} else {
			log.Printf("dhcp response error: %v", result.Err)
		}
	}
	return errors.New("no valid DHCP configuration recieved")
}

func findNetworkInterface() (netlink.Link, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	if len(ifaces) == 0 {
		return nil, errors.New("No network interface found")
	}

	var ifnames []string
	for _, iface := range ifaces {
		ifnames = append(ifnames, iface.Name)
		// skip loopback
		if iface.Flags&net.FlagLoopback != 0 || iface.HardwareAddr.String() == "" {
			continue
		}
		log.Printf("Try using %s", iface.Name)
		link, err := netlink.LinkByName(iface.Name)
		if err == nil {
			return link, nil
		}
		log.Print(err)
	}

	return nil, fmt.Errorf("Could not find a non-loopback network interface with hardware address in any of %v", ifnames)
}

func main() {
	if *doDebug {
		debug = log.Printf
	}
	log.Print(banner)

	configureDHCPNetwork()

	log.Printf("Downloading file from %s\n", url)

	err := DownloadFile("/tmp/fedora.img.xz", url)
	if err != nil {
		fmt.Printf("Downloading file at %s failed with %v\n", url, err)
	}

	data, err := ioutil.ReadFile("/tmp/fedora.img.xz")

	out, err := os.Create("/tmp/fedora.img")

	log.Printf("Decompressing now..")
	// create an xz.Reader to decompress the data
	r, err := xz.NewReader(bytes.NewReader(data), 0)
	if err != nil {
		log.Fatal(err)
		return
	}
	// write the decompressed data to os.Stdout
	_, err = io.Copy(out, r)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = os.Remove("/tmp/fedora.img.xz")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Execute dd if=/tmp/fedora.img of=/dev/nvme0n1 bs=4096 now
	cmd := exec.Command("dd if=/tmp/fedora.img of=/dev/nvme0n1 bs=4096")
	log.Printf("Running %s", strings.Join(cmd.Args, " "))
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error calling dd: %v", err)
		return
	}
	log.Printf("Done.\n")

}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory. We pass an io.TeeReader
// into Copy() to report progress on the download.
func DownloadFile(filepath string, url string) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Get the data
	resp, err := client.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	return nil
}

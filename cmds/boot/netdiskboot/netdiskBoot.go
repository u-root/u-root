package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/u-root/u-root/pkg/boot/kexec"
	"github.com/u-root/u-root/pkg/mount"

	"github.com/dustin/go-humanize"
	"github.com/xi2/xz"
)

var (
	// uroot.uinitargs='-dryRun' as kernel cmd args
	dryRun       = flag.Bool("dryrun", false, "Do everything except downloading and unpacking image - Caution: Only for QEMU usage")
	cfgfile      = "config.json"
	cfg          *netdiskBootConfig
	mountPath    = "/mnt/drive"
	bootfilePath = "/mnt/drive/boot"
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

// PrintProgress prints progess of an ongoing process. In this case the downloading of a given image file.
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

func main() {
	if *dryRun {
		log.Print("Dryrun mode!")
	}

	cfg, err := loadConfig(cfgfile)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("\n%v", cfg.String())

	log.Print(banner)

	err = configureDHCPNetwork()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Only execute this on real hardware with a hdd/ssd available. If you want to play around - Use QEMU
	if !*dryRun {
		log.Printf("Downloading file from %s\n", cfg.ImgURL)

		err = DownloadFile("/tmp/Image.img.xz", cfg.ImgURL)
		if err != nil {
			fmt.Printf("Downloading file at %s failed with %v\n", cfg.ImgURL, err)
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
		// write the decompressed data into fedora.img
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
		cmd := exec.Command("dd if=/tmp/fedora.img of=" + cfg.Device + " bs=4096")
		log.Printf("Running %s", strings.Join(cmd.Args, " "))
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error calling dd: %v", err)
			return
		}
	}

	// Mount cfg.Device as "new" filesystem
	if _, err := mount.Mount(cfg.Device, mountPath, "ext4", "", 0); err != nil {
		log.Fatal(err)
		return
	}
	defer mount.Unmount(mountPath, true, false)

	// Find kernel and initramfs on cfg.Device
	kernel, initramfs, err := getBootFiles(cfg.KernelPrefix, cfg.InitramPrefix)
	if err != nil {
		log.Printf("Kernel: %v - Initramfs: %v - Error: %v", kernel, initramfs, err)
	}

	// kexec into the system
	if err := kexec.FileLoad(kernel, initramfs, cfg.Args); err != nil {
		log.Fatal(err)
		return
	}

	if err = kexec.Reboot(); err != nil {
		log.Fatal(err)
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

func getBootFiles(kernelprefix, initramprefix string) (*os.File, *os.File, error) {
	var kernel, initramfs *os.File
	files, err := ioutil.ReadDir(bootfilePath)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range files {
		if strings.HasPrefix(item.Name(), kernelprefix) {
			kernel, err = os.Open(bootfilePath + "/" + item.Name())
			if err != nil {
				return nil, nil, err
			}
		}
		if strings.HasPrefix(item.Name(), initramprefix) {
			initramfs, err = os.Open(bootfilePath + "/" + item.Name())
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return kernel, initramfs, nil
}

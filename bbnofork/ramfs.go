// bbramfs builds a simple initramfs given an existing built bb; see src/bb.go
// You have to run bb first, which creates src/cmds/bb/bbsh. cd to that directory,
// and run bbramfs, and you have a single binary which does all u-root commands.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type copyfiles struct {
	dir  string
	spec string
}

const (
	bbList = `{{.Uroot}}/src/bb/bbsh
init`
)

func lsr(n string, w *os.File) error {
	n = n + "/"
	err := filepath.Walk(n, func(name string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		cn := strings.TrimPrefix(name, n)
		fmt.Fprintf(w, "%v\n", cn)
		return nil
	})
	return err
}

// we'll keep using cpio and hope the kernel gets fixed some day.
func cpiop(c string) error {

	t := template.Must(template.New("filelist").Parse(c))
	var b bytes.Buffer
	if err := t.Execute(&b, config); err != nil {
		log.Fatalf("spec %v: %v\n", c, err)
	}

	n := strings.Split(b.String(), "\n")
	if config.Debug {
		fmt.Fprintf(os.Stderr, "Strings :%v:\n", n)
	}

	r, w, err := os.Pipe()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	cmd := exec.Command("cpio", "--make-directories", "-p", config.TempDir)
	cmd.Dir = n[0]
	cmd.Stdin = r
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if config.Debug {
		log.Printf("Run %v @ %v", cmd, cmd.Dir)
	}
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	for _, v := range n[1:] {
		if config.Debug {
			fmt.Fprintf(os.Stderr, "%v\n", v)
		}
		err := filepath.Walk(path.Join(n[0], v), func(name string, fi os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf(" WALK FAIL%v: %v\n", name, err)
				// That's ok, sometimes things are not there.
				return filepath.SkipDir
			}
			cn := strings.TrimPrefix(name, n[0]+"/")
			if cn == ".git" {
				return filepath.SkipDir
			}
			fmt.Fprintf(w, "%v\n", cn)
			fmt.Printf("c.dir %v %v %v\n", n[0], name, cn)
			return nil
		})
		if err != nil {
			fmt.Printf("%s: %v\n", v, err)
		}
	}
	w.Close()
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	return nil
}

func sanity() {
	goBinGo := path.Join(config.Goroot, "bin/go")
	_, err := os.Stat(goBinGo)
	if err == nil {
		config.Go = goBinGo
	}
	// but does the one in go/bin/OS_ARCH exist too?
	goBinGo = path.Join(config.Goroot, fmt.Sprintf("bin/%s_%s/go", config.Goos, config.Arch))
	_, err = os.Stat(goBinGo)
	if err == nil {
		config.Go = goBinGo
	}
	if config.Go == "" {
		log.Fatalf("Can't find a go binary! Is GOROOT set correctly?")
	}
}

// sad news. If I concat the Go cpio with the other cpios, for reasons I don't understand,
// the kernel can't unpack it. Don't know why, don't care. Need to create one giant cpio and unpack that.
// It's not size related: if the go archive is first or in the middle it still fails.
func ramfs() {
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	// First create the archive and put the device cpio in it.
	dev, err := ioutil.ReadFile(path.Join(config.Uroot, "scripts", "dev.cpio"))
	if err != nil {
		log.Fatalf("%v %v\n", dev, err)
	}

	oname := fmt.Sprintf("/tmp/initramfs.%v_%v.cpio", config.Goos, config.Arch)
	if err := ioutil.WriteFile(oname, dev, 0600); err != nil {
		log.Fatalf("%v\n", err)
	}

	// Now use the append option for cpio to append to it.
	// That way we get one cpio.
	cmd := exec.Command("cpio", "-H", "newc", "-o", "-A", "-F", oname)
	cmd.Dir = path.Join(config.Uroot, "src/bb/bbsh")
	cmd.Stdin = r
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if config.Debug {
		fmt.Fprintf(os.Stderr, "Run %v @ %v", cmd, cmd.Dir)
	}
	err = cmd.Start()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	w.Write([]byte("init\n"))
	w.Close()
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	fmt.Printf("Output file is in %v\n", oname)
}

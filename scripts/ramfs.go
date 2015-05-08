package main

import (
	"bytes"
	"flag"
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
	dir string
	spec string
}

const (
	goList = `{{.Gosrcroot}}
go/bin/go
go/pkg/include
go/src
go/VERSION.cache
go/misc
go/bin/{{.Goos}}_{{.Arch}}/go
go/pkg/tool/{{.Goos}}_{{.Arch}}/{{.Letter}}g
go/pkg/tool/{{.Goos}}_{{.Arch}}/{{.Letter}}l
go/pkg/tool/{{.Goos}}_{{.Arch}}/asm
go/pkg/tool/{{.Goos}}_{{.Arch}}/old{{.Letter}}a`
	initList=`{{.Gopath}}/src/cmds/init
init`
	urootList=`{{.Gopath}}
src`
)

var (
	config struct {
		Goroot string
		Gosrcroot string
		Arch string
		Goos string
		Letter string
		Gopath string
		TempDir string
	}
	letter = map[string]string{
		"amd64": "6",
		"arm": "5",
		"ppc": "9",
		}
)

func getenv(e, d string) string {
	v := os.Getenv(e)
	if v == "" {
		v = d
	}
	return v
}

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
	fmt.Fprintf(os.Stderr, "Strings :%v:\n", n)
	
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	cmd := exec.Command("cpio", "--make-directories", "-p", config.TempDir)
	cmd.Dir = n[0]
	cmd.Stdin = r
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Printf("Run %v @ %v", cmd, cmd.Dir)
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	
	for _, v := range n[1:] {
		fmt.Fprintf(os.Stderr, "%v\n", v)
		err := filepath.Walk(path.Join(n[0],v), func(name string, fi os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf(" WALK FAIL%v: %v\n", name, err)
				// That's ok, sometimes things are not there.
				return filepath.SkipDir
			}
			cn := strings.TrimPrefix(name, n[0] + "/")
			if cn == ".git" {
				return filepath.SkipDir
			}
			fmt.Fprintf(w, "%v\n", cn)
			//fmt.Printf("c.dir %v %v %v\n", n[0], name, cn)
			return nil
		})
		fmt.Printf("WALKED %v\n", v)
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
// sad news. If I concat the Go cpio with the other cpios, for reasons I don't understand,
// the kernel can't unpack it. Don't know why, don't care. Need to create one giant cpio and unpack that.
// It's not size related: if the go archive is first or in the middle it still fails.
func main() {
	flag.Parse()
	var err error
	config.Arch = getenv("GOARCH", "amd64")
	config.Goroot = getenv("GOROOT", "/")
	config.Gosrcroot = path.Dir(config.Goroot)
	config.Gopath = getenv("GOPATH", "")
	config.Goos = "linux"
	config.Letter = letter[config.Arch]
	config.TempDir, err = ioutil.TempDir("", "u-root")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// sanity checking: do /go/bin/go, and some basic source files exist?
	//sanity()
	// Build init
	cmd := exec.Command("go", "build", "init.go")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = path.Join(config.Gopath, "src/cmds/init")

	err = cmd.Run()
	if err != nil {
		log.Fatalf("%v\n", err)
		os.Exit(1)
	}

	// These produce arrays of strings, the first element being the
	// directory to walk from.
	cpio := []string{
		goList,
		urootList,
		"{{.Gopath}}/src/cmds/init\ninit",
	}
	for _, c := range cpio {
		if err := cpiop(c); err != nil {
			log.Printf("Things went south. TempDir is %v", config.TempDir)
			log.Fatalf("Bailing out near line 666")
		}
	}

	// now try to create the cpio.
	f, err := ioutil.TempFile("", "u-root")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	cmd = exec.Command("cpio", "-H", "newc", "-o")
	cmd.Dir = config.TempDir
	cmd.Stdin = r
	cmd.Stderr = os.Stderr
	cmd.Stdout = f
	fmt.Fprintf(os.Stderr, "Run %v @ %v", cmd, cmd.Dir)
	err = cmd.Start()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	if err := lsr(config.TempDir, w); err != nil {
		log.Fatal("%v\n", err)
	}
	if err := os.Rename(f.Name(), "initramfs.cpio"); err != nil {
		log.Fatal("%v\n", err)
	}
}

/*
# 1. Copy the "myinit" program (compiled above) into the
#    initramfs directory (and rename it to "init"):
#cp myinit initramfs/init
set -e
echo 'if this cp fails, run README in u-root'
cp u-root/init initramfs

# 2. Create the CPIO archive:
cd initramfs
mkdir -p lib/x86_64-linux-gnu lib64
rsync -av ../u-root/src .
rsync -av ../u-root/etc .
cpio -id < ../u-root/go.cpio
cpio -id -E ../u-root/tinycorebase/filelist < ../u-root/tinycorebase/corepure64.cpio
#cpio -id < ../u-root/tinycorebase/tinycorebase.cpio

cp ../u-root//lib/x86_64-linux-gnu/libm.so.6 lib/x86_64-linux-gnu 
cp ../u-root//lib/x86_64-linux-gnu/libc.so.6  lib/x86_64-linux-gnu 
cp ../u-root/lib64/ld-linux-x86-64.so.2  lib64

#fakeroot # this is pure magic (it allows us to pretend to be root)
chown root init
find . | cpio -H newc -o > ../initramfs.cpio # <-- this is the actual initramfs
#exit # leave the fakeroot shell
cd ..
ls -l initramfs.cpio
cp initramfs.cpio linux-3.14.17
*/

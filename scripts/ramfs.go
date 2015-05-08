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
{{.Gosrcroot}}/go/bin/go
{{.Gosrcroot}}/go/pkg/include
{{.Gosrcroot}}/go/src
{{.Gosrcroot}}/go/VERSION.cache
{{.Gosrcroot}}/go/misc
{{.Gosrcroot}}/go/bin/{{.Goos}}_{{.Arch}}/go
{{.Gosrcroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/{{.Letter}}g
{{.Gosrcroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/{{.Letter}}l
{{.Gosrcroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/asm
{{.Gosrcroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/old{{.Letter}}a
`
	initList="init"
	urootList="{{.Gopath}}/src"
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
)

func getenv(e, d string) string {
	v := os.Getenv(e)
	if v == "" {
		v = d
	}
	return v
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
		//		{config.Gopath, urootList},
		"{{.Gopath}}/src/cmds/init\ninit",
	}
	cpiop(cpio[0])
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

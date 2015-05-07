package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"text/template"
)

const (
	copylist = `{{.Goroot}}/go/src
{{.Goroot}}/go/VERSION.cache
{{.Goroot}}/go/misc
{{.Goroot}}/go/pkg/include
{{.Goroot}}/go/bin/{{.Goos}}_{{.Arch}}/go
{{.Goroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/{{.Letter}}g
{{.Goroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/{{.Letter}}l
{{.Goroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/asm
{{.Goroot}}/go/pkg/tool/{{.Goos}}_{{.Arch}}/old{{.Letter}}a
`
)
var (
	fail int
	t = template.Must(template.New("filelist").Parse(copylist))
	letter = map[string]string{
		"amd64": "6",
		"arm": "5",
		"ppc": "9",
		}
)

func cp(in, out string) error {
	b, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", in, err)
		fail++
		return err
	}
	err = ioutil.WriteFile(out, b, 0444)
	if err != nil {
	     	fmt.Fprintf(os.Stderr, "%v: %v\n", out, err)
		fail++
	}
	return nil
}

func getenv(e, d string) string {
	v := os.Getenv(e)
	if v == "" {
		v = d
	}
	return v
}

func main() {

	
	type config struct {
		Goroot string
		Arch string
		Goos string
		Letter string
		Gopath string
	}
	var a config
	flag.Parse()
	var err error
	a.Arch = getenv("GOARCH", "amd64")
	a.Goroot = getenv("GOROOT", "/")
	a.Gopath = getenv("GOPATH", "")
	a.Goos = "linux"

	// Build init
	cmd := exec.Command("go", "build", ".")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = path.Join(a.Gopath, "src/cmds/init")

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	cmd = exec.Command("cpio", "-H", "newc", "--verbose", "-o")
	cmd.Dir = path.Join(a.Gopath, "src/cmds/init")
	cmd.Stdout, err = os.Create(path.Join(a.Gopath, fmt.Sprintf("%v_%vinit.cpio", a.Goos, a.Arch)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	w, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	fmt.Fprintf(w, "init\n")
	w.Close()
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	cmd = exec.Command("cpio", "-H", "newc", "--verbose", "-o")
	cmd.Dir = a.Gopath
	cmd.Stdout, err = os.Create(path.Join(a.Gopath, fmt.Sprintf("%v_%vuroot.cpio", a.Goos, a.Arch)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	w, err = cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	
	err = cmd.Start()
	err = filepath.Walk(path.Join(a.Gopath,"src"), func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("%v: %v\n", path, err)
			return err
		}
		fmt.Fprintf(w, "%v\n", path)
		fmt.Printf("%v\n", path)
		return err
	})
	if err != nil {
		fmt.Printf("%s: %v\n", a.Gopath, err)
	}
	w.Close()
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
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

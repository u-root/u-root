package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var urpath = "/go/bin:/bin:/buildbin:/usr/local/bin:"

func main() {
	/* e.g. (GOBIN=`pwd`/bin go install date) */
	p := os.Getenv("PATH")
	log.Printf("PATH in installcommand: %s\n", p)
	cleanPath := path.Clean(os.Args[0])
	log.Printf("cleanPath %v\n", cleanPath)
	binDir, commandName := path.Split(cleanPath)
	log.Printf("bindri, commandname %v %v\n", binDir, commandName)
	myRoot := path.Base(binDir)
	log.Printf("Myroot i %v\n", myRoot)
	destDir := path.Join(/*myRoot*/"/", "/bin")
	destFile := path.Join(destDir, commandName)
	// Still not totally sure if we want to do all this here
	// or in some hypothetical 'init' command.
	// Pros for doing it here are if the user screws up
	// the environment, it will break installcommand.
	e := os.Environ()
	// sudo sensibly doesn't inherit the path if you are root.
	// and probably doesn't in general so we need to do this much.
	// interestingly, on arch, I did not need to do this. Sounds bad.
	np := strings.NewReplacer("PATH=", "PATH=/go/bin:/bin:/buildbin:")
	for i := range e {
		e[i] = np.Replace(e[i])
	}
	e = append(e, "GOROOT=/go")
	e = append(e, "GOPATH=/")
	e = append(e, "GOBIN=/bin")
	os.Setenv("GOROOT", "/go")
	os.Setenv("GOPATH", "/")
	os.Setenv("GOBIN", "/bin")
	// oh, and, Go looks in the environment, NOT the env in the cmd.
	// Add the prefix if we don't have it.
	if !strings.HasPrefix(p, urpath) {
		if err := os.Setenv("PATH", urpath+p); err != nil {
			fmt.Printf("Couldn't set path; %v\n", err)
			os.Exit(1)
		}
	}
	p = os.Getenv("PATH")
	log.Printf("PATH in installcommand after change: %s\n", p)
	p = os.Getenv("LD_LIBRARY_PATH")
	// tinycore requires /usr/local/lib; make it always last.
	if !strings.HasSuffix(p, ":/usr/local/lib") {
		if err := os.Setenv("LD_LIBRARY_PATH", p+":/usr/local/lib"); err != nil {
			fmt.Printf("Couldn't set LD_LIBRARY_PATH; %v\n", err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("go", "install", "-x", commandName)
	cmd.Env = append(os.Environ(),
		"GOBIN="+path.Join(myRoot, "bin"),
		"CGO_ENABLED=0")

	cmd.Dir = myRoot

	log.Printf("Run %v", cmd)
	out, err := cmd.CombinedOutput()
	log.Printf("installcommand: go build returned")

	if err != nil {
		p := os.Getenv("PATH")
		log.Fatalf("installcommand: trying to build cleanPath: %v, PATH %s, err %v, out %s", cleanPath, p, err, out)
	}

	if false {
		log.Printf(string(out))
	}

	log.Printf("WE BUILT IT? %v\n", destFile)
	cmd = exec.Command(destFile)

	cmd.Args = append([]string{commandName}, os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

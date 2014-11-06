package main

import (
	"log"
	"os"
	"os/exec"
	"path"
)

var urpath = "/go/bin:/bin:/buildbin:/usr/local/bin:"

func main() {
	/* e.g. (GOBIN=`pwd`/bin go install date) */

	cleanPath := path.Clean(os.Args[0])
	log.Printf("cleanPath %v\n", cleanPath)
	binDir, commandName := path.Split(cleanPath)
	log.Printf("bindir, commandname %v %v\n", binDir, commandName)
	destDir := "/bin"
	destFile := path.Join(destDir, commandName)

	cmd := exec.Command("go", "install", "-x", commandName)

	cmd.Dir = "/"

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

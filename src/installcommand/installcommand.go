package main

import (
	"log"
	"os"
	"os/exec"
	"path"
)

func main() {
	/* e.g. (GOBIN=`pwd`/bin go install date) */
	p := os.Getenv("PATH")
	log.Printf("PATH in installcommand: %s\n", p)
	cleanPath := path.Clean(os.Args[0])
	binDir, commandName := path.Split(cleanPath)
	myRoot := path.Base(binDir)
	destDir := path.Join(myRoot, "/bin")
	destFile := path.Join(destDir, commandName)

	cmd := exec.Command("go", "install" /*"-x", */, commandName)
	cmd.Env = append(os.Environ(),
		"GOBIN="+path.Join(myRoot, "bin"),
		"CGO_ENABLED=0")

	cmd.Dir = myRoot

	out, err := cmd.CombinedOutput()

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

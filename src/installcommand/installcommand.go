package main

import (
	"path"
	"os"
	"os/exec"
	"log"
)

func main(){
	/* e.g. (GOBIN=`pwd`/bin go install date) */
	cleanPath := path.Clean(os.Args[0])
	binDir, commandName := path.Split(cleanPath)
	myRoot := path.Base(binDir)
	destDir := path.Join(myRoot, "/bin")
	destFile := path.Join(destDir, commandName)

	e := os.Environ()
	e = append(e, "GOBIN="+path.Join(myRoot,"bin"))
	e = append(e, "CGO_ENABLED=0")

	cmd := exec.Command("go", "install", commandName)
	cmd.Dir = myRoot
	
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err, string(out))
	}

	cmd = exec.Command(destFile)
	cmd.Args = os.Args[1:]
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"path"
	"os"
	"os/exec"
	"log"
)

func main(){
	buildPath := os.Getenv("BUILDPATH")

	cleanPath := path.Clean(os.Args[0])
	arg0 := path.Base(cleanPath)

	cmd := exec.Command("go", "install")
	cmd.Dir = path.Join(buildPath, arg0)
	
	out, err := cmd.CombinedOutput()

	if err != nil {
		log.Fatal(err, string(out))
	}

	cmd = exec.Command(path.Join(buildPath, "bin", arg0))
	cmd.Args = os.Args[1:]
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

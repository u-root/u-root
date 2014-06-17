package main

import (
	"path"
	"os"
	"os/exec"
	"log"
)

func main(){
	myRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cleanPath := path.Clean(os.Args[0])
	commandName := path.Base(cleanPath)
	buildingBin := path.Join(myRoot, "/buildbin")
	destDir := path.Join(myRoot, "/bin")
	destFile := path.Join(destDir, commandName)
	sourceDir := path.Join(myRoot, commandName)
	log.Println("%v, %v, %v, %v, %v %v\n", 
	cleanPath,
	commandName,
	buildingBin,
	myRoot,
	destDir, destFile)

	cmd := exec.Command("go", "install")
	cmd.Dir = sourceDir
	log.Println("RUn %v\n", cmd)
	
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

package main

import (
	"path"
	"os"
	"os/exec"
	"log"
)

func main(){
	/* e.g. (GOBIN=`pwd`/bin go install date) */
 	/*
	myRoot, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
 */
	cleanPath := path.Clean(os.Args[0])
	binDir, commandName := path.Split(cleanPath)
	myRoot := path.Base(binDir)
	buildingBin := path.Join(myRoot, "/buildbin")
	destDir := path.Join(myRoot, "/bin")
	destFile := path.Join(destDir, commandName)
	//sourceDir := path.Join(myRoot, commandName)
	log.Println("%v, %v, %v, %v, %v %v\n", 
	cleanPath,
	commandName,
	buildingBin,
	myRoot,
	destDir, destFile)

	e := os.Environ()
	e = append(e, "GOBIN="+path.Join(myRoot,"bin"))

	cmd := exec.Command("go", "install", commandName)
	cmd.Dir = myRoot
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

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"gopkg.in/yaml.v2"
)

var yamlFile = flag.String("yaml", ".github/workflows/go.yml", "YAML CI file")

func main() {

	flag.Parse()

	b, err := ioutil.ReadFile(*yamlFile)
	if err != nil {
		fmt.Printf("yamlFile.Get err #%v ", err)
	}

	var ci CI
	err = yaml.Unmarshal(b, &ci)
	if err != nil {
		fmt.Printf("Unmarshal: %v", err)
	}

	l := ci.Jobs.Linters
	for _, s := range l.Steps {

		if out, err := exec.Command("bash", "-c", s.Run).CombinedOutput(); err != nil {
			log.Printf("Linter %s failed:%s:%v", s, string(out), err)
		}
		log.Printf("Linter %s:OK", s)

	}
}

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

	for _, s := range ci.Jobs.Linters.Steps {
		log.Printf("Linter %s: run %s", s.Name, s.Run)
		if out, err := exec.Command("bash", "-c", s.Run).CombinedOutput(); err != nil {
			log.Printf("Linter %s failed:%s:%v", s.Name, string(out), err)
		}
		log.Printf("Linter %s:OK", s.Name)
	}

	for _, s := range ci.Jobs.Build.Steps {
		log.Printf("Build %s: run %s", s.Name, s.Run)
		if out, err := exec.Command("bash", "-c", s.Run).CombinedOutput(); err != nil {
			log.Printf("Build %s failed:%s:%v", s.Name, string(out), err)
		}
		log.Printf("Build %s:OK", s.Name)
	}

	for _, s := range ci.Jobs.Badbuild.Steps {
		log.Printf("Badbuild %s: run %s", s.Name, s.Run)
		if out, err := exec.Command("bash", "-c", s.Run).CombinedOutput(); err == nil {
			log.Printf("Badbuild %s did not fail though it should:%s", s.Name, string(out))
		}
		log.Printf("Badbuild %s:OK", s.Name)
	}

	//
	m := ci.Jobs.MultiOsArch
	for _, os := range m.Strategy.Matrix.Os {
		for _, arch := range m.Strategy.Matrix.Arch {
			for _, s := range m.Steps {
				log.Printf("MultiOsArch:%s(%s,%s)", s.Name, os, arch)
				// todo: GoVersion
				// The Run uses a template but not that will work so well with Go
				// templates.
				c := exec.Command("bash", "-c", "go build .")
				c.Env = append(c.Env, "GOOS="+os, "GOARCH="+arch)
				if out, err := c.CombinedOutput(); err != nil {
					log.Printf("MultiOsArch:%s(%s)(%s,%s) failed:(%s,%v)", s.Name, s.Run, os, arch, string(out), err)
				}
				log.Printf("MultiOsArch:%s(%s,%s):OK", s.Name, os, arch)
			}
		}
	}
}

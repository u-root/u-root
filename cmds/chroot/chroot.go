// +build !windows !plan9

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type userSpec struct {
	uid uint32
	gid uint32
}

func (u *userSpec) Set(s string) error {
	var err error
	userspecSplit := strings.Split(s, ":")
	if len(userspecSplit) != 2 || userspecSplit[1] == "" {
		return fmt.Errorf("Expected user spec flag to be \":\" separated values received %s", s)
	}

	u.uid, err = stringToUint32(userspecSplit[0])
	if err != nil {
		return err
	}

	u.gid, err = stringToUint32(userspecSplit[1])
	if err != nil {
		return err
	}

	return nil
}

func (u *userSpec) Get() interface{} {
	return *u
}

func (u *userSpec) String() string {
	return fmt.Sprintf("%d:%d", u.uid, u.gid)
}

func defaultUser() userSpec {
	return userSpec{
		uid: uint32(os.Getuid()),
		gid: uint32(os.Getgid()),
	}
}

type groupsSpec struct {
	groups []uint32
}

func (g *groupsSpec) Set(s string) error {
	groupStrs := strings.Split(s, ",")
	g.groups = make([]uint32, len(groupStrs))

	for index, group := range groupStrs {

		gid, err := stringToUint32(group)
		if err != nil {
			return err
		}

		g.groups[index] = gid
	}

	return nil
}

func (g *groupsSpec) Get() interface{} {
	return *g
}

func (g *groupsSpec) String() string {
	var buffer bytes.Buffer

	for index, gid := range g.groups {
		buffer.WriteString(fmt.Sprint(gid))
		if index < len(g.groups)-1 {
			buffer.WriteString(",")
		}
	}

	return buffer.String()
}

var (
	skipchdirFlag bool
	user          userSpec   = defaultUser()
	groups        groupsSpec = groupsSpec{}
)

func init() {
	flag.Var(&user, "u", "specify user and group (ID only) as USER:GROUP")
	flag.Var(&groups, "g", "specify supplementary group ids as g1,g2,..,gN")
	flag.BoolVar(&skipchdirFlag, "s", false, fmt.Sprint("Use this option to not change",
		"the working directory to / after changing the root directory to newroot, i.e., ",
		"inside the chroot. This option is only permitted when newroot is the old / directory."))
}

func stringToUint32(str string) (uint32, error) {
	ret, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(ret), nil
}

func parseCommand(args []string) []string {
	if len(args) > 1 {
		return args[1:]
	}
	return []string{"/bin/sh", "-i"}
}

func parseRoot(args []string) (root string, err error) {
	if len(args) < 1 {
		return "", fmt.Errorf("Missing operand")
	}

	return filepath.Abs(args[0])
}

func isRoot(dir string) (bool, error) {
	realPath, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return false, err
	}
	absolutePath, err := filepath.Abs(realPath)
	if err != nil {
		return false, err
	}
	if absolutePath == "/" {
		return true, nil
	}
	return false, nil
}

func main() {
	var (
		newRoot   string
		isOldroot bool
		err       error
	)

	flag.Parse()

	if flag.NFlag() == 0 && flag.NArg() == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	newRoot, err = parseRoot(flag.Args())
	if err != nil {
		log.Fatal(err)
	}
	isOldroot, err = isRoot(newRoot)
	if err != nil {
		log.Fatal(err)
	}

	if !skipchdirFlag {
		err = os.Chdir(newRoot)
		if err != nil {
			log.Fatal(err)
		}
	} else if !isOldroot {
		log.Fatal("The -s option is only permitted when newroot is the old / directory")
	}

	argv := parseCommand(flag.Args())

	cmd := exec.Command(argv[0], argv[1:]...)

	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid:    user.uid,
			Gid:    user.gid,
			Groups: groups.groups,
		},
		Chroot: newRoot,
	}

	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

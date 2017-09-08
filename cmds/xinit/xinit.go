package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var (
	tczPackages  = []string{"aterm", "fltk-1.3", "flwm", "freetype", "glib2", "harfbuzz", "imlib2-bin", "imlib2", "libffi", "libfontenc", "libICE", "libjpeg-turbo", "libpng", "libSM", "libX11", "libXau", "libxcb", "libXdmcp", "libXext", "libXfont", "libXi", "libXmu", "libXpm", "libXrandr", "libXrender", "libXt", "pcre", "wbar", "Xfbdev", "Xlibs", "Xorg-fonts", "Xprogs", "Xorg-7.7", "links"}
	sshdCommands = []string{"Protocol 2", "AcceptEnv LANG LC_*", "UsePAM no", "ChallengeResponseAuthentication no", "passwordauthentication no", "AuthorizedKeysFile ~/.ssh/authorized_keys", "PermitRootLogin without-password", "X11Forwarding yes", "RSAAuthentication yes", "PubkeyAuthentication yes", "X11DisplayOffset 10", "X11UseLocalhost yes"}
)

func setup() error {
	if false {
		return syscall.Mount("/tmp", "/tmp", "tmpfs", syscall.MS_MGC_VAL, "")
	}
	cmd := exec.Command("dhclient", "&")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf(" error: %v. Continuing...", err)
	}
	fmt.Printf("Ip link below: \n")
	cmd = exec.Command("ip", "link", "show")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf(" error: %v. Continuing...", err)
	}
	return nil
}

func sshSetup() error {
	sshdLoc := "/etc/ssh/sshd_config"
	cmd := exec.Command("tcz", "openssh")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	//cp /etc/ssh/sshd_config /etc/ssh/sshd_config_backup
	if _, err := os.Stat(sshdLoc); err != nil {
		return err
	}
	fileContent, err := ioutil.ReadFile(sshdLoc)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(sshdLoc+"_backup", fileContent, 0777); err != nil {
		return nil
	}
	var buffer bytes.Buffer
	byteFile := []byte{}
	for arg := range sshdCommands {
		buffer.WriteString(fmt.Sprintf("%v\n", arg))
	}
	fmt.Printf("SSHD file is : %s", string(byteFile))
	if err := ioutil.WriteFile(sshdLoc, byteFile, 0777); err != nil {
		return nil
	}
	fmt.Printf("not functional \n")
	//scp ananyajoshi@100.96.221.137:~/.ssh/id_rsa.pub /root/.ssh
	//cat /root/.ssh/id_rsa.pub >> /root/.ssh/authorized_keys
	///usr/sbin/sshd -d -d -d -D -e*
	return nil
}

func tczSetup() error {
	get := []string{"tcz", "-v", "8.x"}
	get = append(get, tczPackages...)
	cmd := exec.Command("sudo", get...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	if false {
		setup()
		sshSetup()
	}
	tczSetup()
}

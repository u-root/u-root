package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	tczPackages = []string{
		"aterm",
		"bash",
		"fltk-1.3",
		"flwm",
		"freetype",
		"glib2",
		"harfbuzz",
		"imlib2-bin",
		"imlib2",
		"libffi",
		"libfontenc",
		"libICE",
		"libjpeg-turbo",
		"libpng",
		"libSM",
		"libX11",
		"libXau",
		"libxcb",
		"libXdmcp",
		"libXext",
		"libXfont",
		"libXi",
		"libXmu",
		"libXpm",
		"libXrandr",
		"libXrender",
		"libXt",
		"pcre",
		"wbar",
		"Xfbdev",
		"Xlibs",
		"Xorg-fonts",
		"Xprogs",
		"Xorg-7.7",
		"links",
		"opera-12",
	}
	sshdCommands = []string{
		"Protocol 2",
		"AcceptEnv LANG LC_*",
		"UsePAM no",
		"ChallengeResponseAuthentication no",
		"passwordauthentication no",
		"AuthorizedKeysFile ~/.ssh/authorized_keys",
		"PermitRootLogin without-password",
		"X11Forwarding yes",
		"RSAAuthentication yes",
		"PubkeyAuthentication yes",
		"X11DisplayOffset 10",
		"X11UseLocalhost yes",
	}
	ssh = flag.Bool("ssh", false, "Ssh default set to false")
)

func setup() error {
	var errors error
	if err := os.Symlink("/usr/local/bin/bash", "/bin/bash"); err != nil {
		errors = fmt.Errorf("bash symlink: %v", err)
	}
	if err := os.Symlink("/lib/ld-linux-x86-64.so.2", "/lib64/ld-linux-x86-64.so.2"); err != nil {
		errors = fmt.Errorf("%v. ld-linux symlink: %v", errors, err)
	}
	return errors
}

func sshSetup() error {
	sshdLoc := "/etc/ssh/sshd_config"
	cmd := exec.Command("tcz", "openssh")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	//This block of code does cp /etc/ssh/sshd_config /etc/ssh/sshd_config_backup
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
	// This writes the new sshd file
	var buffer bytes.Buffer
	byteFile := []byte{}
	for arg := range sshdCommands {
		buffer.WriteString(fmt.Sprintf("%v\n", arg))
	}
	fmt.Printf("SSHD file is : %s", string(byteFile))
	if err := ioutil.WriteFile(sshdLoc, byteFile, 0777); err != nil {
		return nil
	}

	if err := ioutil.WriteFile("/etc/passwd", []byte("root:x:0:0:root:/root:/bin/bash\nnobody::27:27:nobody privsep:/var/empty:/sbin/nologin"), 0777); err != nil {
		return err
	}
	if err := ioutil.WriteFile("/etc/group", []byte("nobody::27:"), 0777); err != nil {
		return err
	}

	//this does scp yourname@iplink:~/.ssh/id_rsa.pub /root/.ssh
	locationKeys := "yourname@iplink:~/.ssh/id_rsa.pub"
	for true {
		fmt.Printf("Where are your public keys located? ex:%s They will be copied to /root/.ssh and be called id_rsa.pub. \n", locationKeys)
		_, err := fmt.Scanf("%s", &locationKeys)
		if err != nil {
			return err
		}
		cmd := exec.Command("scp", locationKeys, "/root/.ssh/id_rsa.pub")
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("tried to scp your public key. Failed because error: %v", err)
		}
		break
	}
	//this does cat /root/.ssh/id_rsa.pub >> /root/.ssh/authorized_keys
	rsaKey, err := ioutil.ReadFile("/root/.ssh/id_rsa.pub")
	if err != nil {
		return err
	}
	ioutil.WriteFile("/root/.ssh/authorized_keys", rsaKey, os.ModeAppend)

	cmd = exec.Command("usr/sbin/sshd", "-d", "-d", "-d", "-D", "-e")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func tczSetup() error {
	get := []string{"-v", "8.x"}
	get = append(get, tczPackages...)
	cmd := exec.Command("tcz", get...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func xSetup() error {
	go func() {
		cmd := exec.Command("Xfbdev")
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("X11 startup: %v", err)
		}
	}()
	time.Sleep(5 * time.Second)
	for _, f := range []string{"flwm", "opera-12"} {
		log.Printf("Run %v", f)
		go func(n string) {
			cmd := exec.Command(n, "-display", ":0")
			cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
			if err := cmd.Run(); err != nil {
				log.Printf("%v: %v", n, err)
			}
		}(f)
	}
	cmd := exec.Command("aterm", "-display", ":0")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := setup(); err != nil {
		log.Printf("Error is %v", err)
	}
	if err := tczSetup(); err != nil {
		log.Printf("Error is %v", err)
	}
	if false {
		if err := sshSetup(); err != nil {
			log.Printf("Error is %v", err)
		}
	}
	if err := xSetup(); err != nil {
		log.Printf("Error is %v", err)
	}
}

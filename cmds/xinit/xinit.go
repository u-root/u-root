package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
	"flag"
)

var (
	tczPackages  = []string{"Xorg-7.7-bin", "Xorg-7.7-dev", "Xorg-7.7-lib-dev", "Xorg-7.7-lib", "Xorg-7.7", "Xorg-fonts", "xorg-proto", "xorg-server-dev",
"xorg-server","Xprogs","Xlibs","xpad-locale","xpad","Xlibs","Xfbdev","xvid","xkeyboard-config","xz","gdb","libX11-dev","libX11","libXfixes","i2c-4.8.17-tinycore64","graphics-4.8.17-tinycore64","libXvmc","libdrm","pixman","libXinerama","libXrandr","libXdamage","libXcursor","libXtst","libxshmfence","xf86-video-intel","openssh","aterm", "bash", "strace", "openssh" , "opera-12"}
	sshdCommands = []string{"Protocol 2", "AcceptEnv LANG LC_*", "UsePAM no", "ChallengeResponseAuthentication no", "passwordauthentication no", "AuthorizedKeysFile ~/.ssh/authorized_keys", "PermitRootLogin without-password", "X11Forwarding yes", "RSAAuthentication yes", "PubkeyAuthentication yes", "X11DisplayOffset 10", "X11UseLocalhost yes"}
	ssh         = flag.Bool("ssh", false, "Ssh default set to false")
)

func setup() error {
	cmd := exec.Command("dhclient", "-ipv4=true", "-ipv6=false", "-verbose")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf(" error: %v. Continuing...", err)
	}
	if false {
		cmd := exec.Command("bash", "-c", "OPERA_DIR=/tmp/tcloop/opera-12/usr/local/share/opera-12 /tmp/tcloop/opera-12/usr/local/lib/opera-12/opera-12")
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf(" error: %v. Continuing...", err)
		}	
	}
	

	fmt.Printf("Ip link below: \n")
	cmd = exec.Command("ip", "link", "show")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf(" error: %v. Continuing...", err)
	}
	if err = os.Symlink("/usr/local/bin/bash","/bin/bash"); err != nil{
		return err	
	}
	if err = os.Symlink("/lib/ld-linux-x86-64.so.2","/lib64"); err != nil{
		return err	
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
	for true{
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

func xSetup() error{
	cmd := exec.Command("Xfbdev")
	cmd.SysProcAttr=& syscall.SysProcAttr{Foreground:false}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	time.Sleep(5)
	cmd = exec.Command("aterm",  "-display",":0")
	cmd.SysProcAttr=& syscall.SysProcAttr{Foreground:false}
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := setup(); err != nil{
		log.Printf("Error is %v", err)
	}	
	if err := tczSetup(); err != nil {
		log.Printf("Error is %v", err)
	}
	if false {
		if err := sshSetup(); err != nil{
			log.Printf("Error is %v", err)
		}
	}
	if err := xSetup(); err != nil{
		log.Printf("Error is %v", err)	
	}
}

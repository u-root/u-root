package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"

	"github.com/mholt/archiver"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/interfaces"
	"github.com/insomniacslk/dhcp/netboot"
	flag "github.com/spf13/pflag"
)

var (
	alpineMirror       = flag.String("m", "http://ftp.halifax.rwth-aachen.de/alpine", "Alpine linux mirror url")
	alpineRelease      = flag.String("r", "v3.8", "Alpine linux release")
	alpineToolsVersion = flag.String("v", "2.10.1-r0", "Alpine linux installer version")
	setupNetwork       = flag.Bool("n", true, "Setup network on all interfaces")
	startChrome        = flag.Bool("c", false, "Run chrome on Alpine Linux")
	machineArch        string
)

const (
	interfaceUpTimeout = 1 * time.Second
	retries            = 3
	alpineRootPath     = "/tmp/alpine"
)

func runCmd(executable string, subCommands ...string) {
	baseCmd, err := exec.LookPath(executable)
	if err != nil {
		log.Fatalf("Couldn't find executable: %v", err)
	}

	cmdLoad := exec.Command(baseCmd, subCommands...)
	if err := cmdLoad.Run(); err != nil {
		//log.Fatalf("Couldn't execute command: %s ", baseCmd)
	}
}

func machine() (string, error) {
	u := syscall.Utsname{}
	err := syscall.Uname(&u)
	if err != nil {
		return "", err
	}

	var m string
	for _, val := range u.Machine {
		m += string(int(val))
	}

	return m, nil
}

func setupChroot(chrootPath string) error {
	if _, err := os.Stat(path.Join(chrootPath, "etc", "resolv.conf")); !os.IsNotExist(err) {
		return nil
	}
	err := syscall.Mount("/proc", path.Join(chrootPath, "proc"), "",
		syscall.MS_BIND, "")
	if err != nil {
		return err
	}
	err = syscall.Mount("/sys", path.Join(chrootPath, "sys"), "",
		syscall.MS_BIND, "")
	if err != nil {
		return err
	}
	err = syscall.Mount("/dev", path.Join(chrootPath, "dev"), "",
		syscall.MS_BIND, "")
	if err != nil {
		return err
	}
	resolvConf, err := ioutil.ReadFile("/etc/resolv.conf")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(chrootPath, "etc", "resolv.conf"), resolvConf, 0644)
	if err != nil {
		return err
	}

	return nil
}

func userChroot(chrootPath string) error {
	setupChroot(alpineRootPath)

	os.Chdir(alpineRootPath)
	cmd := exec.Command("/bin/sh", "-i")

	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(os.Getuid()),
			Gid: uint32(os.Getgid()),
		},
		Chroot: alpineRootPath,
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func sysChroot(chrootPath string) (func() error, error) {
	setupChroot(alpineRootPath)

	root, err := os.Open("/")
	if err != nil {
		return nil, err
	}

	if err := syscall.Chroot(chrootPath); err != nil {
		root.Close()
		return nil, err
	}

	return func() error {
		defer root.Close()
		if err := root.Chdir(); err != nil {
			return err
		}
		return syscall.Chroot(".")
	}, nil
}

func configureNetwork(ifname string, attempts int) (*netboot.NetConf, error) {
	var (
		conv []*dhcpv4.DHCPv4
		err  error
	)
	_, err = netboot.IfUp(ifname, interfaceUpTimeout)
	if err != nil {
		return nil, fmt.Errorf("Ifup failed: %v", err)
	}
	if attempts < 1 {
		attempts = 1
	}
	client := dhcpv4.NewClient()
	for attempt := 0; attempt < attempts; attempt++ {
		log.Printf("Attempt to get DHCP lease %d of %d for interface %s", attempt+1, attempts, ifname)
		conv, err = client.Exchange(ifname)
		if err != nil && attempt < attempts {
			log.Printf("Error: %v", err)
			continue
		}
		break
	}
	log.Printf("Interface %s configured!", ifname)

	netconf, _, err := netboot.ConversationToNetconfv4(conv)
	return netconf, err
}

func runDhcp() error {
	iflist, err := interfaces.GetNonLoopbackInterfaces()
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, iface := range iflist {
		netconf, err := configureNetwork(iface.Name, retries+1)
		if err != nil {
			log.Printf("%v", err)
			continue
		}

		if err := netboot.ConfigureInterface(iface.Name, netconf); err != nil {
			return err
		}
	}

	return nil
}

func installAlpine() error {
	if _, err := os.Stat(alpineRootPath); !os.IsNotExist(err) {
		return nil
	}
	os.MkdirAll(alpineRootPath, os.ModePerm)

	urlString := "http://ftp.halifax.rwth-aachen.de/alpine/v3.8/main/x86_64/apk-tools-static-2.10.1-r0.apk"
	//urlString := *alpineMirror + "/" + *alpineRelease + "/main/" + machineArch + "/" + "apk-tools-static-" + *alpineToolsVersion + ".apk"
	log.Printf("%s", urlString)
	resp, err := http.Get(urlString)
	if err != nil {
		log.Fatalf("Couldn't download alpine tools with URL: %s", urlString)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("%s for URL: %s", resp.Status, urlString)
	}

	out, err := ioutil.TempFile("/tmp", "alpine-linux-")
	if err != nil {
		log.Fatalf("Couldn't create tempfile: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Couldn't write data from url to tempfile: %v", err)
	}

	os.RemoveAll("/tmp/apk/")
	err = archiver.DefaultTarGz.Unarchive(out.Name(), "/tmp/apk")
	if err != nil {
		log.Fatalf("Couldn't unpack data from tempfile: %v", err)
	}

	var runCommands []string
	mainRepo := *alpineMirror + "/" + *alpineRelease + "/main"
	runCommands = append(runCommands, "--repository")
	runCommands = append(runCommands, mainRepo)
	runCommands = append(runCommands, "--update-cache")
	runCommands = append(runCommands, "--allow-untrusted")
	runCommands = append(runCommands, "--root")
	runCommands = append(runCommands, alpineRootPath)
	runCommands = append(runCommands, "--initdb")
	runCommands = append(runCommands, "add")
	runCommands = append(runCommands, "alpine-base")

	runCmd("/tmp/apk/sbin/apk.static", runCommands...)

	communityRepo := *alpineMirror + "/" + *alpineRelease + "/community"
	repositoryFilePath := alpineRootPath + "/etc/apk/repositories"
	f, err := os.Create(repositoryFilePath)
	if err != nil {
		log.Fatalf("Couldn't create Alpine repository file: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(mainRepo + "\n")
	if err != nil {
		log.Fatalf("Couldn't write Alpine repository data: %v", err)
	}

	_, err = f.WriteString(communityRepo + "\n")
	if err != nil {
		log.Fatalf("Couldn't write Alpine repository data: %v", err)
	}

	return nil
}

func setupAlpine() {
	exit, err := sysChroot(alpineRootPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	runCmd("apk", "update")
	runCmd("apk", "upgrade")

	if err := exit(); err != nil {
		panic(err)
	}
}

func runChrome() {
	exit, err := sysChroot(alpineRootPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	runCmd("setup-xorg-base")
	runCmd("apk", "add", "chromium", "xinit", "xrandr")
	runCmd("X", "-configure")

	if err := exit(); err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	var err error

	machineArch, err = machine()
	if err != nil {
		log.Fatalf("Could not obtain machine architecture: %v", err)
	}

	if *setupNetwork {
		runDhcp()
	}

	err = installAlpine()
	if err != nil {
		log.Fatalf("Could not obtain chroot path: %v", err)
	}

	setupAlpine()

	if *startChrome {
		runChrome()
	} else {
		userChroot(alpineRootPath)
	}
}

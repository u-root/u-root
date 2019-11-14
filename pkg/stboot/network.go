package stboot

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/client4"
	"github.com/insomniacslk/dhcp/netboot"
)

// SetupIOFromNetVars sets up your eth interface from netvars.json
func ConfigureStaticNetwork(vars NetVars, doDebug bool) error {
	//setup ip
	debug("Setup network configuration with IP: " + vars.HostIP)
	cmd := exec.Command("ip", "addr", "add", vars.HostIP, "dev", eth)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
	}
	cmd = exec.Command("ip", "link", "set", eth, "up")
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
	}
	cmd = exec.Command("ip", "route", "add", "default", "via", vars.DefaultGateway, "dev", eth)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
	}

	if doDebug {
		cmd = exec.Command("ip", "addr")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error executing %v: %v", cmd, err)
		}
		cmd = exec.Command("ip", "route")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error executing %v: %v", cmd, err)
		}
	}

	return nil
}

// ConfigureDHCPNetwork configures DHCP on eth0
func ConfigureDHCPNetwork() error {

	debug("Trying to configure network configuration dynamically..")
	attempts := 10
	var conversation []*dhcpv4.DHCPv4

	_, err := netboot.IfUp(eth, interfaceUpTimeout)
	if err != nil {
		log.Println("Enabling eth0 failed.")
		return fmt.Errorf("Ifup failed: %v", err)
	}
	if attempts < 1 {
		attempts = 1
	}

	client := client4.NewClient()
	for attempt := 0; attempt < attempts; attempt++ {
		debug("Attempt to get DHCP lease %d of %d for interface %s", attempt+1, attempts, eth)
		conversation, err = client.Exchange(eth)

		if err != nil && attempt < attempts {
			log.Printf("Error: %v", err)
			continue
		}
		break
	}

	if conversation[3] == nil {
		return fmt.Errorf("Gateway is null")
	}
	netbootConfig, err := netboot.GetNetConfFromPacketv4(conversation[3])

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	err = netboot.ConfigureInterface(eth, netbootConfig)

	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	// Some manual shit - for now
	cmd := exec.Command("ip", "route", "add", "default", "via", netbootConfig.Routers[0].String()+"/24", "dev", eth)
	if err := cmd.Run(); err != nil {
		log.Printf("Error executing %v: %v", cmd, err)
		return err
	}

	return nil
}

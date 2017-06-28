package dhcp6client

import (
	"net"

	"github.com/mdlayher/dhcp6"
)

/***
* DHCP packet
 */
func newSolicitOptions(mac *net.HardwareAddr) (dhcp6.Options, error) {
	// make options: iana
	var id = [4]byte{'r', 'o', 'o', 't'}
	options := make(dhcp6.Options)
	if err := options.Add(dhcp6.OptionIANA, dhcp6.NewIANA(id, 0, 0, nil)); err != nil {
		return nil, err
	}
	// make options: rapid commit
	if err := options.Add(dhcp6.OptionRapidCommit, nil); err != nil {
		return nil, err
	}
	// make options: elapsed time
	var et dhcp6.ElapsedTime
	et.UnmarshalBinary([]byte{0x00, 0x00})
	if err := options.Add(dhcp6.OptionElapsedTime, et); err != nil {
		return nil, err
	}
	// make options: option request option
	oro := make(dhcp6.OptionRequestOption, 4)
	oro.UnmarshalBinary([]byte{0x00, 0x17, 0x00, 0x18})
	if err := options.Add(dhcp6.OptionORO, oro); err != nil {
		return nil, err
	}
	// make options: duid with mac address
	duid := dhcp6.NewDUIDLL(6, *mac)
	db, err := duid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	// add row
	options[dhcp6.OptionClientID] = append(options[dhcp6.OptionClientID], db)

	return options, nil
}

func newSolicitPacket(mac *net.HardwareAddr) (*dhcp6.Packet, error) {
	options, err := newSolicitOptions(mac)
	if err != nil {
		return nil, err
	}

	return &dhcp6.Packet{
		MessageType:   dhcp6.MessageTypeSolicit,
		TransactionID: [3]byte{0x00, 0x01, 0x02},
		Options:       options,
	}, nil
}

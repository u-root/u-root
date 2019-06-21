package netboot

import "github.com/jsimonetti/rtnetlink"

// getOperState returns the operational state for the given interface index.
func getOperState(iface int) (rtnetlink.OperationalState, error) {
	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return 0, err
	}
	msg, err := conn.Link.Get(uint32(iface))
	if err != nil {
		return 0, err
	}
	return msg.Attributes.OperationalState, nil
}

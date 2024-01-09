// Package rtnl provides a convenient API on top of the rtnetlink library.
package rtnl

import (
	"github.com/jsimonetti/rtnetlink"
	"github.com/mdlayher/netlink"
)

// Conn represents the underlying netlink connection
type Conn struct {
	Conn *rtnetlink.Conn // a route netlink connection
}

// Dial the netlink socket. Establishes a new connection. The typical initialisation is:
// 	conn, err := rtnl.Dial(nil)
//	if err != nil {
//		log.Fatal("can't establish netlink connection: ", err)
//	}
//	defer conn.Close()
//	// use conn for your calls
//
func Dial(cfg *netlink.Config) (*Conn, error) {
	conn, err := rtnetlink.Dial(cfg)
	if err != nil {
		return nil, err
	}
	return &Conn{Conn: conn}, nil
}

// Close the connection.
func (c *Conn) Close() error {
	return c.Conn.Close()
}

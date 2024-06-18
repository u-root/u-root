// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceroute

type Collection struct {
	Num     int
	Entries []Entry
}

type Entry struct {
	id    int
	Probe *Probe
}

func (c *Collection) AddProbe(pb *Probe) {
	id := c.Num

	e := Entry{
		id:    id,
		Probe: pb,
	}
	c.Entries = append(c.Entries, e)
}

func (c *Collection) ModifyProbe(pb *Probe) {
	for _, e := range c.Entries {
		if e.Probe.id == pb.id {
			e.Probe.saddr = pb.saddr
			e.Probe.recvTime = pb.sendtime
		}
	}
}

func (c *Collection) GetByTTL(ttl int) []*Probe {
	pbs := make([]*Probe, 0)
	for _, e := range c.Entries {
		if e.Probe.ttl == ttl {
			pbs = append(pbs, e.Probe)
		}
	}
	return pbs
}

func (c *Collection) MaxTTL() int {
	ttl := 1
	for _, e := range c.Entries {
		if ttl < e.Probe.ttl {
			ttl = e.Probe.ttl
		}
	}

	return ttl
}

// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brctl

const (
	BRCTL_SYS_NET    = "/sys/class/net"
	BRCTL_SYS_SUFFIX = 0x0a
)

const (
	BRCTL_ADD_BRIDGE          = 2
	BRCTL_DEL_BRIDGE          = 3
	BRCTL_ADD_I               = 4
	BRCTL_DEL_I               = 5
	BRCTL_SET_AEGING_TIME     = 11
	BRCTL_SET_BRIDGE_PRIORITY = 15
	BRCTL_SET_PORT_PRIORITY   = 16
	BRCTL_SET_PATH_COST       = 17
)

const (
	BRCTL_DESIGNATED_ROOT          = "root_id"
	BRCTL_ROOT_PATH_COST           = "root_path_cost"
	BRCTL_AGEING_TIME              = "ageing_time"
	BRCTL_STP_STATE                = "stp_state"
	BRCTL_BRIDGE_PRIO              = "priority"
	BRCTL_FORWARD_DELAY            = "forward_delay"
	BRCTL_HELLO_TIME               = "hello_time"
	BRCTL_HELLO_TIMER_VALUE        = "hello_timer"
	BRCTL_TCN_TIMER                = "tcn_timer"
	BRCTL_TOPOLOGY_CHANGE          = "topology_change"
	BRCTL_TOPOLOGY_CHANGE_DETECTED = "topology_change_detected"
	BRCTL_TOPOLOGY_CHANGE_TIMER    = "topology_change_timer"
	BRCTL_GC_TIMER_VALUE           = "gc_timer"
	BRCTL_ROOT_PORT                = "root_port"
	BRCTL_TRILL_ENABLED            = "trill_state"
	BRCTL_MAX_AGE                  = "max_age"
	BRCTL_PATH_COST                = "path_cost"
	BRCTL_PRIORITY                 = "priority"
	BRCTL_HAIRPIN                  = "hairpin_mode"
	BRCTL_BRFORWARD                = "brforward"
	BRCTL_BRIDGEID                 = "bridge_id"
	BRCTL_BRIDGE_DIR               = "bridge"
	BRCTL_BRIDGE_INTERFACE_DIR     = "brif"
)

//+build linux

package netlink

import (
	"syscall"
	"testing"
	"unsafe"
)

func TestHeaderMemoryLayoutLinux(t *testing.T) {
	var nh Header
	var sh syscall.NlMsghdr

	if want, got := unsafe.Sizeof(sh), unsafe.Sizeof(nh); want != got {
		t.Fatalf("unexpected structure sizes:\n- want: %v\n-  got: %v",
			want, got)
	}

	sh = syscall.NlMsghdr{
		Len:   0x10101010,
		Type:  0x2020,
		Flags: 0x3030,
		Seq:   0x40404040,
		Pid:   0x50505050,
	}
	nh = sysToHeader(sh)

	if want, got := sh.Len, nh.Length; want != got {
		t.Fatalf("unexpected header length:\n- want: %v\n-  got: %v",
			want, got)
	}
	if want, got := sh.Type, uint16(nh.Type); want != got {
		t.Fatalf("unexpected header type:\n- want: %v\n-  got: %v",
			want, got)
	}
	if want, got := sh.Flags, uint16(nh.Flags); want != got {
		t.Fatalf("unexpected header flags:\n- want: %v\n-  got: %v",
			want, got)
	}
	if want, got := sh.Seq, nh.Sequence; want != got {
		t.Fatalf("unexpected header sequence:\n- want: %v\n-  got: %v",
			want, got)
	}
	if want, got := sh.Pid, nh.PID; want != got {
		t.Fatalf("unexpected header PID:\n- want: %v\n-  got: %v",
			want, got)
	}
}

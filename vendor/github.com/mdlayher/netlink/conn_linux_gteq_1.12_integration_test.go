//+build go1.12,linux

package netlink_test

import (
	"net"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jsimonetti/rtnetlink/rtnl"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

func TestIntegrationConnTimeout(t *testing.T) {
	t.Parallel()

	conn, err := netlink.Dial(unix.NETLINK_GENERIC, nil)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	timeout := 1 * time.Millisecond
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		t.Fatalf("failed to set deadline: %v", err)
	}

	errC := make(chan error)
	go func() {
		_, err := conn.Receive()
		errC <- err
	}()

	select {
	case err := <-errC:
		mustBeTimeoutNetError(t, err)
	case <-time.After(timeout + 100*time.Millisecond):
		t.Fatalf("timeout did not fire")
	}
}

func TestIntegrationConnExecuteAfterReadDeadline(t *testing.T) {
	t.Parallel()

	conn, err := netlink.Dial(unix.NETLINK_GENERIC, nil)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	timeout := 1 * time.Millisecond
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		t.Fatalf("failed to set deadline: %v", err)
	}
	time.Sleep(2 * timeout)

	req := netlink.Message{
		Header: netlink.Header{
			Flags:    netlink.Request | netlink.Acknowledge,
			Sequence: 1,
		},
	}
	got, err := conn.Execute(req)
	if err == nil {
		t.Fatalf("Execute succeeded: got %v", got)
	}
	mustBeTimeoutNetError(t, err)
}

func TestIntegrationConnNetNSExplicit(t *testing.T) {
	t.Parallel()

	skipUnprivileged(t)

	// Create a network namespace for use within this test.
	const ns = "nltest0"
	shell(t, "ip", "netns", "add", ns)
	defer shell(t, "ip", "netns", "del", ns)

	f, err := os.Open("/var/run/netns/" + ns)
	if err != nil {
		t.Fatalf("failed to open namespace file: %v", err)
	}
	defer f.Close()

	// Create a connection in each the host namespace and the new network
	// namespace. We will use these to validate that a namespace was entered
	// and that an interface creation notification was only visible to the
	// connection within the namespace.
	hostC, hostDone := rtnlDial(t, 0)
	defer hostDone()

	nsC, nsDone := rtnlDial(t, int(f.Fd()))
	defer nsDone()

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	go func() {
		defer wg.Done()

		_, err := hostC.Receive()
		if err == nil {
			panic("received netlink message in host namespace")
		}

		// Timeout means we were interrupted, so return.
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			return
		}

		panicf("failed to receive in host namespace: %v", err)
	}()

	// Create a temporary interface within the new network namespace.
	const ifName = "nltestns0"
	defer shell(t, "ip", "netns", "exec", ns, "ip", "link", "del", ifName)

	ifi := rtnlReceive(t, nsC, func() {
		// Trigger a notification in the new namespace.
		shell(t, "ip", "netns", "exec", ns, "ip", "tuntap", "add", ifName, "mode", "tun")
	})

	// And finally interrupt the host connection so it can exit its
	// receive goroutine.
	if err := hostC.SetDeadline(time.Unix(1, 0)); err != nil {
		t.Fatalf("failed to interrupt host connection: %v", err)
	}

	if diff := cmp.Diff(ifName, ifi); diff != "" {
		t.Fatalf("unexpected interface name (-want +got):\n%s", diff)
	}
}

func TestIntegrationConnNetNSImplicit(t *testing.T) {
	t.Parallel()

	skipUnprivileged(t)

	// Create a network namespace for use within this test.
	const ns = "nltest0"
	shell(t, "ip", "netns", "add", ns)
	defer shell(t, "ip", "netns", "del", ns)

	f, err := os.Open("/var/run/netns/" + ns)
	if err != nil {
		t.Fatalf("failed to open namespace file: %v", err)
	}
	defer f.Close()

	// Create an interface in the new namespace. We will attempt to find it later.
	const ifName = "nltestns0"
	shell(t, "ip", "netns", "exec", ns, "ip", "tuntap", "add", ifName, "mode", "tun")
	defer shell(t, "ip", "netns", "exec", ns, "ip", "link", "del", ifName)

	// We're going to manipulate the network namespace of this thread, so we
	// must lock OS thread and keep track of the original namespace for later.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	threadNS, err := netlink.ThreadNetNS()
	if err != nil {
		t.Fatalf("failed to get current network namespace: %v", err)
	}
	defer threadNS.Close()

	if err := threadNS.Set(int(f.Fd())); err != nil {
		t.Fatalf("failed to enter new network namespace: %v", err)
	}

	// A newly created netlink connection should enter the new network namespace
	// associated with this thread automatically.
	if !findLink(t, ifName) {
		t.Fatalf("did not find interface %q in namespace %q", ifName, ns)
	}

	// Return to the default namespace.
	//
	// A newly created netlink connection should NOT find the link because it
	// is now in the default namespace.
	if err := threadNS.Restore(); err != nil {
		t.Fatalf("failed to restore original network namespace: %v", err)
	}

	if findLink(t, ifName) {
		t.Fatalf("found interface %q in default namespace", ifName)
	}
}

func findLink(t *testing.T, name string) bool {
	t.Helper()

	c, err := rtnl.Dial(nil)
	if err != nil {
		t.Fatalf("failed to dial rtnetlink: %v", err)
	}
	defer c.Close()

	ifis, err := c.Links()
	if err != nil {
		t.Fatalf("failed to list links: %v", err)
	}

	var found bool
	for _, ifi := range ifis {
		if ifi.Name == name {
			found = true
			break
		}
	}

	return found
}

func mustBeTimeoutNetError(t *testing.T, err error) {
	t.Helper()
	ne, ok := err.(net.Error)
	if !ok {
		t.Fatalf("didn't get a net.Error: got a %T instead", err)
	}
	if !ne.Timeout() {
		t.Fatalf("didn't get a timeout")
	}
}

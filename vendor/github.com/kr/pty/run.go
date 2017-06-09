// +build !windows

package pty

import (
	"fmt"
	"os"
	"os/exec"
//	"syscall"
)

// Start assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
func Start(c *exec.Cmd) (pty *os.File, err error) {
	cc := exec.Command("/ubin/date")
	fmt.Printf("UDEATE TRY TO START %v\n\n", cc)
	a, b := cc.CombinedOutput()
	fmt.Printf("run date: %v %v\n\n", a, b)
	pty, tty, err := Open()
	if err != nil {
		return nil, fmt.Errorf("pty Start: open: %v", err)
	}
	defer tty.Close()
//	c.Stdout = tty
	//c.Stdin = tty
	//c.Stderr = tty
	//if c.SysProcAttr == nil {
		//c.SysProcAttr = &syscall.SysProcAttr{}
	//}
//	c.SysProcAttr.Setctty = true
	//c.SysProcAttr.Setsid = true
//	c.SysProcAttr.Foreground = true
	fmt.Printf("TRY TO START %v\nSysprocattr %v\n\n", c, c.SysProcAttr)
	err = c.Start()
	if err != nil {
		fmt.Printf("THAT WENT BADLY: %v", err)
		pty.Close()
		return nil, fmt.Errorf("pty Start: Start command: %v", err)
	}
	return pty, err
}

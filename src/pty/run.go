package pty

import (
	"os"
	"os/exec"
)

// Start assigns a pseudo-terminal tty os.File to c.Stdin, c.Stdout,
// and c.Stderr, calls c.Start, and returns the File of the tty's
// corresponding pty.
func Start(c *exec.Cmd) (pty *os.File, err error) {
	ptm, pts, _, err := Open()
	if err != nil {
		return nil, err
	}
	defer ptm.Close()
	c.Stdout = pts
	c.Stdin = pts
	c.Stderr = pts
	c.SysProcAttr.Setctty = true
	c.SysProcAttr.Setsid = true
	c.SysProcAttr.Ctty = 0
	err = c.Start()
	if err != nil {
		pty.Close()
		return nil, err
	}
	return ptm, err
}

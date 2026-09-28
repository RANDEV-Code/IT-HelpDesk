//go:build !windows

package main

import (
	"os/exec"
	"syscall"
)

// prepareCmd memisahkan process group agar sinyal tidak merambat.
func prepareCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// sendShutdownSignal mengirim SIGTERM secara langsung (dipakai CI Linux).
func sendShutdownSignal(cmd *exec.Cmd) error {
	return cmd.Process.Signal(syscall.SIGTERM)
}

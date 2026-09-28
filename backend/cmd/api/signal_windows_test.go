//go:build windows

package main

import (
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

// prepareCmd membuat proses anak dalam process group baru agar CTRL_BREAK
// hanya terkirim ke proses tersebut (tidak ke proses tes).
func prepareCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: windows.CREATE_NEW_PROCESS_GROUP,
	}
}

// sendShutdownSignal mengirim CTRL_BREAK_EVENT; runtime Go di Windows
// memetakannya menjadi SIGTERM sehingga handler graceful shutdown aktif.
func sendShutdownSignal(cmd *exec.Cmd) error {
	return windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, uint32(cmd.Process.Pid))
}

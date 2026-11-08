//go:build linux

package pigment

import (
	"os"
	"syscall"
	"unsafe"
)

// isTerminal reports whether f refers to a terminal (tty) device.
//
// It issues a TCGETS ioctl, which only terminal drivers implement; regular
// files, pipes, sockets and character devices such as /dev/null reject it.
// This correctly distinguishes a real tty from /dev/null, which the previous
// ModeCharDevice heuristic misclassified as a terminal.
func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	var termios [128]byte // >= sizeof(struct termios) on every unix port
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(syscall.TCGETS),
		uintptr(unsafe.Pointer(&termios[0])),
		0, 0, 0,
	)
	return errno == 0
}

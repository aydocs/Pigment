//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package pigment

import (
	"os"
	"syscall"
	"unsafe"
)

// isTerminal reports whether f refers to a terminal (tty) device.
// BSD-derived systems use TIOCGETA for the same purpose Linux uses TCGETS.
func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	var termios [256]byte
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(syscall.TIOCGETA),
		uintptr(unsafe.Pointer(&termios[0])),
		0, 0, 0,
	)
	return errno == 0
}

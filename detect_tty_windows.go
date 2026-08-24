//go:build windows

package pigment

import (
	"os"
	"syscall"
	"unsafe"
)

var getConsoleMode = syscall.NewLazyDLL("kernel32.dll").NewProc("GetConsoleMode")

// isTerminal reports whether f refers to a console handle. GetConsoleMode
// succeeds only for console input/output handles; pipes, files and NUL all
// reject it.
func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	var mode uint32
	r, _, _ := getConsoleMode.Call(f.Fd(), uintptr(unsafe.Pointer(&mode)))
	return r != 0
}

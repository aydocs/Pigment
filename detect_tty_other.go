//go:build !unix && !windows

package pigment

import "os"

// isTerminal falls back to the character-device heuristic on platforms
// without a supported tty probe in the standard library (plan9, js, wasip1).
// Unlike the unix and windows implementations it cannot distinguish
// character devices such as /dev/null from a real terminal.
func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

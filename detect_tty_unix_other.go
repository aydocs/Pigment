//go:build aix || solaris

package pigment

import "os"

// isTerminal uses the portable character-device fallback on Unix platforms
// whose ioctl constants are not exposed consistently by the standard library.
func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	fi, err := f.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

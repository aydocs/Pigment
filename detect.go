package pigment

import "os"

// Enabled reports whether colorized output is active. It is computed at start
// from the environment and terminal, but can be changed with SetEnabled.
var Enabled = detectEnabled()

// SetEnabled forces colorized output on or off.
func SetEnabled(v bool) {
	Enabled = v
}

// DetectEnabled re-runs environment/terminal detection.
func DetectEnabled() {
	Enabled = detectEnabled()
}

func detectEnabled() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	if v := os.Getenv("FORCE_COLOR"); v != "" && v != "0" {
		return true
	}
	if v := os.Getenv("CLICOLOR_FORCE"); v != "" && v != "0" {
		return true
	}
	if v := os.Getenv("CLICOLOR"); v == "0" {
		return false
	}
	return isTerminal(os.Stdout)
}

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

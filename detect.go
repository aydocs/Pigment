package pigment

import (
	"os"
	"sync/atomic"
)

// enabled holds the global color state in atomic storage so that
// SetEnabled, DetectEnabled, IsEnabled and every render path can run
// concurrently without a data race.
var enabled atomic.Bool

func init() {
	enabled.Store(detectFromEnv())
}

// SetEnabled forces colorized output on or off for the whole package,
// overriding environment detection until DetectEnabled runs again.
// It is safe to call from multiple goroutines.
func SetEnabled(v bool) {
	enabled.Store(v)
}

// IsEnabled reports the current global color state without re-running
// detection. It is safe to call from multiple goroutines.
func IsEnabled() bool {
	return enabled.Load()
}

// DetectEnabled re-runs environment/terminal detection, stores the result as
// the new global state and returns it. It is safe to call from multiple
// goroutines.
func DetectEnabled() bool {
	v := detectFromEnv()
	enabled.Store(v)
	return v
}

// colorEnabled reports whether color should currently be emitted. It is the
// single read seam used by every rendering path (styles, gradients, links).
func colorEnabled() bool {
	return enabled.Load()
}

// detectFromEnv resolves whether color should be enabled, following the
// chalk convention. Precedence, highest first:
//
//  1. FORCE_COLOR, when set: "0" or "false" disables color; any other value
//     (including the empty string) forces color on, even for a non-terminal.
//  2. CLICOLOR_FORCE, when set to anything other than "0", forces color on.
//  3. NO_COLOR, when set to a non-empty value, disables color.
//  4. TERM set to "dumb" disables color.
//  5. CLICOLOR set to "0" disables color.
//  6. Otherwise color is enabled only when os.Stdout is a terminal.
func detectFromEnv() bool {
	if v, ok := os.LookupEnv("FORCE_COLOR"); ok {
		return v != "0" && v != "false"
	}
	if v, ok := os.LookupEnv("CLICOLOR_FORCE"); ok && v != "0" {
		return true
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	if os.Getenv("CLICOLOR") == "0" {
		return false
	}
	return isTerminal(os.Stdout)
}

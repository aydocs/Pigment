package pigment

import "sync"

// Theme is a named collection of styles. It is safe for concurrent use:
// Register and Style synchronize on an internal mutex, and stored styles are
// independent snapshots, so a snapshot handed out by Style can be mutated
// freely without affecting the theme or other snapshots.
type Theme struct {
	mu     sync.RWMutex
	styles map[string]*Style
}

// NewTheme creates an empty theme.
func NewTheme() *Theme {
	return &Theme{styles: make(map[string]*Style)}
}

// Register associates name with an independent snapshot of s. Later mutations
// of s do not affect the theme; registering nil stores an absent style.
// It is safe for concurrent use.
func (t *Theme) Register(name string, s *Style) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if s == nil {
		t.styles[name] = nil
		return
	}
	t.styles[name] = s.Clone()
}

// Style returns an independent snapshot of the style registered under name,
// or nil if absent. Mutating the returned style does not affect the theme.
// It is safe for concurrent use.
func (t *Theme) Style(name string) *Style {
	t.mu.RLock()
	defer t.mu.RUnlock()
	s := t.styles[name]
	if s == nil {
		return nil
	}
	return s.Clone()
}

// Paint renders text with the named style. Unknown names return text unchanged.
func (t *Theme) Paint(name, text string) string {
	s := t.Style(name)
	if s == nil {
		return text
	}
	return s.Paint(text)
}

// Names returns the registered style names.
func (t *Theme) Names() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]string, 0, len(t.styles))
	for k := range t.styles {
		out = append(out, k)
	}
	return out
}

// DefaultTheme is the global theme preloaded with common semantic styles.
// It is safe for concurrent use; Register and Lookup operate on independent
// snapshots.
var DefaultTheme = newDefaultTheme()

func newDefaultTheme() *Theme {
	t := NewTheme()
	t.Register("error", New().Fg("#ff5555").Bold())
	t.Register("err", t.Style("error"))
	t.Register("warn", New().Fg("#ffb86c").Bold())
	t.Register("warning", t.Style("warn"))
	t.Register("info", New().Fg("#8be9fd"))
	t.Register("success", New().Fg("#50fa7b").Bold())
	t.Register("ok", t.Style("success"))
	t.Register("debug", New().Fg("#6272a4"))
	t.Register("muted", New().Fg("#888888").Faint())
	t.Register("link", New().Fg("#8be9fd").Underline())
	t.Register("title", New().Fg("#bd93f9").Bold().Underline())
	t.Register("highlight", New().Fg("#f1fa8c"))
	return t
}

// Register adds a snapshot of a style to the global DefaultTheme.
func Register(name string, s *Style) { DefaultTheme.Register(name, s) }

// Lookup returns an independent snapshot of the style registered under name
// in the global DefaultTheme (nil if absent).
func Lookup(name string) *Style { return DefaultTheme.Style(name) }

// Paint renders text with a named style from the global DefaultTheme.
func Paint(name, text string) string { return DefaultTheme.Paint(name, text) }

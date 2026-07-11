package pigment

import "sync"

// Theme is a named collection of styles. Safe for concurrent use.
type Theme struct {
	mu     sync.RWMutex
	styles map[string]*Style
}

// NewTheme creates an empty theme.
func NewTheme() *Theme {
	return &Theme{styles: make(map[string]*Style)}
}

// Register associates a name with a style.
func (t *Theme) Register(name string, s *Style) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.styles[name] = s
}

// Style returns the style registered under name (nil if absent).
func (t *Theme) Style(name string) *Style {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.styles[name]
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

// Register adds a style to the global DefaultTheme.
func Register(name string, s *Style) { DefaultTheme.Register(name, s) }

// Lookup looks up a style in the global DefaultTheme.
func Lookup(name string) *Style { return DefaultTheme.Style(name) }

// Paint renders text with a named style from the global DefaultTheme.
func Paint(name, text string) string { return DefaultTheme.Paint(name, text) }

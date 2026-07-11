package pigment

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Output is the default writer used by the Print* methods.
var Output io.Writer = os.Stdout

// SetOutput redirects the Print* output.
func SetOutput(w io.Writer) {
	if w != nil {
		Output = w
	}
}

// Attr is a text style attribute (bold, underline...). Combine with OR.
type Attr uint32

const (
	Bold            Attr = 1 << iota // bright/bold
	Faint                            // dim
	Italic                           // italic
	Underline                        // underline
	DoubleUnderline                  // double underline
	BlinkSlow                        // slow blink
	BlinkRapid                       // fast blink
	Reverse                          // swap fg/bg
	Conceal                          // hidden
	CrossedOut                       // strikethrough
	Overline                         // overline
	Framed                           // frame (rarely supported)
	Encircled                        // circle (rarely supported)
)

type colorVal struct {
	isSet bool
	isPal bool
	rgb   Color
	pal   int
}

// Style describes how text is rendered: fg, bg and attributes.
type Style struct {
	fg      colorVal
	bg      colorVal
	attrs   Attr
	enabled *bool
}

// New returns an empty Style.
func New() *Style { return &Style{} }

// Clone returns a deep copy of the style.
func (s *Style) Clone() *Style {
	c := *s
	return &c
}

// SetEnabled forces this style on/off, ignoring the global flag. nil reverts.
func (s *Style) SetEnabled(v *bool) *Style {
	s.enabled = v
	return s
}

// Enable forces this style on.
func (s *Style) Enable() *Style {
	s.enabled = boolPtr(true)
	return s
}

// Disable forces this style off.
func (s *Style) Disable() *Style {
	s.enabled = boolPtr(false)
	return s
}

func (s *Style) active() bool {
	if s.enabled != nil {
		return *s.enabled
	}
	return Enabled
}

// Fg sets the foreground from a flexible value: Color, *Color, "#rrggbb",
// color name, [3]int, or a 0..255 palette index.
func (s *Style) Fg(v any) *Style {
	if cv, ok := toColorVal(v); ok {
		s.fg = cv
	}
	return s
}

// FgRGB sets the foreground to a 24-bit RGB color.
func (s *Style) FgRGB(r, g, b int) *Style {
	s.fg = colorVal{isSet: true, rgb: RGB(r, g, b)}
	return s
}

// FgHex sets the foreground from a hex string.
func (s *Style) FgHex(h string) *Style {
	if c, err := Hex(h); err == nil {
		s.fg = colorVal{isSet: true, rgb: c}
	}
	return s
}

// FgPalette sets the foreground to an xterm 256 index (0..255).
func (s *Style) FgPalette(n int) *Style {
	s.fg = colorVal{isSet: true, isPal: true, pal: n}
	return s
}

// Bg sets the background (same rules as Fg).
func (s *Style) Bg(v any) *Style {
	if cv, ok := toColorVal(v); ok {
		s.bg = cv
	}
	return s
}

// BgRGB sets the background to a 24-bit RGB color.
func (s *Style) BgRGB(r, g, b int) *Style {
	s.bg = colorVal{isSet: true, rgb: RGB(r, g, b)}
	return s
}

// BgHex sets the background from a hex string.
func (s *Style) BgHex(h string) *Style {
	if c, err := Hex(h); err == nil {
		s.bg = colorVal{isSet: true, rgb: c}
	}
	return s
}

// BgPalette sets the background to an xterm 256 index.
func (s *Style) BgPalette(n int) *Style {
	s.bg = colorVal{isSet: true, isPal: true, pal: n}
	return s
}

// Add turns the given attributes on.
func (s *Style) Add(a Attr) *Style {
	s.attrs |= a
	return s
}

// Del turns the given attributes off.
func (s *Style) Del(a Attr) *Style {
	s.attrs &^= a
	return s
}

// Bold, Faint, Italic... are short setters returning the receiver.
func (s *Style) Bold() *Style      { return s.Add(Bold) }
func (s *Style) Faint() *Style     { return s.Add(Faint) }
func (s *Style) Italic() *Style    { return s.Add(Italic) }
func (s *Style) Underline() *Style { return s.Add(Underline) }
func (s *Style) Blink() *Style     { return s.Add(BlinkSlow) }
func (s *Style) Reverse() *Style   { return s.Add(Reverse) }
func (s *Style) Strike() *Style    { return s.Add(CrossedOut) }
func (s *Style) Overline() *Style  { return s.Add(Overline) }

const (
	esc     = "\x1b["
	reset   = "\x1b[0m"
	osc8Pre = "\x1b]8;;"
	osc8Sep = "\x1b\\"
)

// open returns the SGR opening sequence (empty when inactive).
func (s *Style) open() string {
	if !s.active() {
		return ""
	}
	var parts []string
	if s.attrs&Bold != 0 {
		parts = append(parts, "1")
	}
	if s.attrs&Faint != 0 {
		parts = append(parts, "2")
	}
	if s.attrs&Italic != 0 {
		parts = append(parts, "3")
	}
	if s.attrs&Underline != 0 {
		parts = append(parts, "4")
	}
	if s.attrs&DoubleUnderline != 0 {
		parts = append(parts, "21")
	}
	if s.attrs&BlinkSlow != 0 {
		parts = append(parts, "5")
	}
	if s.attrs&BlinkRapid != 0 {
		parts = append(parts, "6")
	}
	if s.attrs&Reverse != 0 {
		parts = append(parts, "7")
	}
	if s.attrs&Conceal != 0 {
		parts = append(parts, "8")
	}
	if s.attrs&CrossedOut != 0 {
		parts = append(parts, "9")
	}
	if s.attrs&Framed != 0 {
		parts = append(parts, "51")
	}
	if s.attrs&Encircled != 0 {
		parts = append(parts, "52")
	}
	if s.attrs&Overline != 0 {
		parts = append(parts, "53")
	}
	if s.fg.isSet {
		if s.fg.isPal {
			parts = append(parts, fmt.Sprintf("38;5;%d", s.fg.pal))
		} else {
			parts = append(parts, fmt.Sprintf("38;2;%d;%d;%d", s.fg.rgb.R, s.fg.rgb.G, s.fg.rgb.B))
		}
	}
	if s.bg.isSet {
		if s.bg.isPal {
			parts = append(parts, fmt.Sprintf("48;5;%d", s.bg.pal))
		} else {
			parts = append(parts, fmt.Sprintf("48;2;%d;%d;%d", s.bg.rgb.R, s.bg.rgb.G, s.bg.rgb.B))
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return esc + strings.Join(parts, ";") + "m"
}

func (s *Style) close() string {
	if !s.active() {
		return ""
	}
	return reset
}

// Paint wraps text with the style's escape codes.
func (s *Style) Paint(text string) string {
	o := s.open()
	if o == "" {
		return text
	}
	return o + text + reset
}

// Sprint formats with fmt.Sprint and paints the result.
func (s *Style) Sprint(a ...any) string { return s.Paint(fmt.Sprint(a...)) }

// Sprintf formats with fmt.Sprintf and paints the result.
func (s *Style) Sprintf(format string, a ...any) string {
	return s.Paint(fmt.Sprintf(format, a...))
}

// Sprintln formats with fmt.Sprintln (no trailing newline) and paints.
func (s *Style) Sprintln(a ...any) string {
	return s.Paint(strings.TrimSuffix(fmt.Sprintln(a...), "\n")) + "\n"
}

// Print writes the painted text to Output.
func (s *Style) Print(a ...any) (int, error) { return fmt.Fprint(Output, s.Paint(fmt.Sprint(a...))) }

// Println writes the painted text (with newline) to Output.
func (s *Style) Println(a ...any) (int, error) {
	return fmt.Fprint(Output, s.Sprintln(a...))
}

// Printf writes the formatted, painted text to Output.
func (s *Style) Printf(format string, a ...any) (int, error) {
	return fmt.Fprint(Output, s.Paint(fmt.Sprintf(format, a...)))
}

// Fprint writes the painted text to w.
func (s *Style) Fprint(w io.Writer, a ...any) (int, error) {
	return fmt.Fprint(w, s.Paint(fmt.Sprint(a...)))
}

// Fprintln writes the painted text (with newline) to w.
func (s *Style) Fprintln(w io.Writer, a ...any) (int, error) {
	return fmt.Fprint(w, s.Sprintln(a...))
}

// Fprintf writes the formatted, painted text to w.
func (s *Style) Fprintf(w io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprint(w, s.Paint(fmt.Sprintf(format, a...)))
}

// PrintFunc returns a new function that prints arguments with this style.
func (s *Style) PrintFunc() func(a ...any) {
	return func(a ...any) { s.Print(a...) }
}

// PrintfFunc returns a new function that formats and prints with this style.
func (s *Style) PrintfFunc() func(format string, a ...any) {
	return func(format string, a ...any) { s.Printf(format, a...) }
}

// PrintlnFunc returns a new function that prints lines with this style.
func (s *Style) PrintlnFunc() func(a ...any) {
	return func(a ...any) { s.Println(a...) }
}

// SprintFunc returns a new function that returns colorized strings.
func (s *Style) SprintFunc() func(a ...any) string {
	return func(a ...any) string { return s.Sprint(a...) }
}

// SprintfFunc returns a new function that returns formatted colorized strings.
func (s *Style) SprintfFunc() func(format string, a ...any) string {
	return func(format string, a ...any) string { return s.Sprintf(format, a...) }
}

// SprintlnFunc returns a new function that returns colorized lines.
func (s *Style) SprintlnFunc() func(a ...any) string {
	return func(a ...any) string { return s.Sprintln(a...) }
}

// FprintFunc returns a new function that prints to an io.Writer with this style.
func (s *Style) FprintFunc() func(w io.Writer, a ...any) {
	return func(w io.Writer, a ...any) { s.Fprint(w, a...) }
}

// FprintfFunc returns a new function that formats and prints to an io.Writer.
func (s *Style) FprintfFunc() func(w io.Writer, format string, a ...any) {
	return func(w io.Writer, format string, a ...any) { s.Fprintf(w, format, a...) }
}

// FprintlnFunc returns a new function that prints lines to an io.Writer.
func (s *Style) FprintlnFunc() func(w io.Writer, a ...any) {
	return func(w io.Writer, a ...any) { s.Fprintln(w, a...) }
}

// Equals reports whether two styles render identically.
func (s *Style) Equals(o *Style) bool {
	if s == nil || o == nil {
		return s == o
	}
	return s.fg == o.fg && s.bg == o.bg && s.attrs == o.attrs
}

func toColorVal(v any) (colorVal, bool) {
	switch x := v.(type) {
	case Color:
		return colorVal{isSet: true, rgb: x}, true
	case *Color:
		if x == nil {
			return colorVal{}, false
		}
		return colorVal{isSet: true, rgb: *x}, true
	case string:
		if strings.HasPrefix(x, "#") {
			if c, err := Hex(x); err == nil {
				return colorVal{isSet: true, rgb: c}, true
			}
			return colorVal{}, false
		}
		if c, ok := Named(x); ok {
			return colorVal{isSet: true, rgb: c}, true
		}
		return colorVal{}, false
	case int:
		if x >= 0 && x <= 255 {
			return colorVal{isSet: true, isPal: true, pal: x}, true
		}
	case uint8:
		return colorVal{isSet: true, rgb: Color{x, x, x}}, true
	case [3]int:
		return colorVal{isSet: true, rgb: RGB(x[0], x[1], x[2])}, true
	case [3]uint8:
		return colorVal{isSet: true, rgb: Color{x[0], x[1], x[2]}}, true
	}
	return colorVal{}, false
}

func boolPtr(v bool) *bool { return &v }

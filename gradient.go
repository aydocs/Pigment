package pigment

import (
	"sort"
	"strings"
)

// GradientStop is a color anchored at a position (0..1) along the text.
type GradientStop struct {
	At    float64
	Color Color
}

// GradientOpt configures a Gradient call.
type GradientOpt func(*gradientCfg)

type gradientCfg struct {
	stops []GradientStop
}

// Stops replaces the default two-color gradient with explicit stops.
// A Stops option carrying no stops is ignored and the default gradient kept.
func Stops(stops ...GradientStop) GradientOpt {
	return func(c *gradientCfg) {
		if len(stops) > 0 {
			c.stops = stops
		}
	}
}

// Gradient paints text with a smooth color transition. Whitespace is left
// uncolored so layouts stay intact.
func Gradient(text string, start, end Color, opts ...GradientOpt) string {
	cfg := gradientCfg{stops: []GradientStop{{At: 0, Color: start}, {At: 1, Color: end}}}
	for _, o := range opts {
		o(&cfg)
	}
	stops := normalizeStops(cfg.stops)
	return walkRunes(text, stops.colorAt)
}

// walkRunes paints every non-whitespace rune of text with the color chosen by
// colorAt, called with the rune's relative position (0..1). Space, tab,
// newline and carriage return pass through uncolored so layouts stay intact.
func walkRunes(text string, colorAt func(t float64) Color) string {
	runes := []rune(text)
	n := len(runes)
	if n == 0 {
		return text
	}
	var b strings.Builder
	for i, r := range runes {
		switch r {
		case ' ', '\t', '\n', '\r':
			b.WriteRune(r)
			continue
		}
		t := 0.0
		if n > 1 {
			t = float64(i) / float64(n-1)
		}
		b.WriteString(paintRune(r, colorAt(t)))
	}
	return b.String()
}

func paintRune(r rune, c Color) string {
	if !colorEnabled() {
		return string(r)
	}
	return esc + "38;2;" + itoa(int(c.R)) + ";" + itoa(int(c.G)) + ";" + itoa(int(c.B)) + "m" + string(r) + reset
}

// RainbowOpt configures a Rainbow call.
type RainbowOpt func(*rainbowCfg)

type rainbowCfg struct {
	start  float64
	end    float64
	hasEnd bool // keeps RainbowEnd(0) distinguishable from unset
	sat    float64
	light  float64
}

// RainbowStart sets the starting hue (0..360).
func RainbowStart(h float64) RainbowOpt {
	return func(c *rainbowCfg) { c.start = h }
}

// RainbowEnd sets the ending hue (0..360). Defaults to start+360.
func RainbowEnd(h float64) RainbowOpt {
	return func(c *rainbowCfg) { c.end = h; c.hasEnd = true }
}

// RainbowSaturation sets the HSL saturation (0..1).
func RainbowSaturation(s float64) RainbowOpt {
	return func(c *rainbowCfg) { c.sat = s }
}

// RainbowLightness sets the HSL lightness (0..1).
func RainbowLightness(l float64) RainbowOpt {
	return func(c *rainbowCfg) { c.light = l }
}

// Rainbow paints text with a hue cycle for a rainbow effect.
func Rainbow(text string, opts ...RainbowOpt) string {
	cfg := rainbowCfg{start: 0, sat: 1, light: 0.5}
	for _, o := range opts {
		o(&cfg)
	}
	if !cfg.hasEnd {
		cfg.end = cfg.start + 360
	}
	return walkRunes(text, func(t float64) Color {
		return HSL(cfg.start+t*(cfg.end-cfg.start), cfg.sat, cfg.light)
	})
}

// normalizeStops sorts user stops by position and anchors the gradient at 0
// and 1. An anchor is prepended or appended only when the extreme user stop
// does not already reach that boundary, so boundary colors are never
// duplicated. A single stop yields a constant color.
func normalizeStops(in []GradientStop) gradientStops {
	stops := make([]GradientStop, len(in))
	copy(stops, in)
	sort.Slice(stops, func(i, j int) bool { return stops[i].At < stops[j].At })
	out := make([]GradientStop, 0, len(stops)+2)
	if first := stops[0]; first.At > 0 {
		out = append(out, GradientStop{At: 0, Color: first.Color})
	}
	out = append(out, stops...)
	if last := stops[len(stops)-1]; last.At < 1 {
		out = append(out, GradientStop{At: 1, Color: last.Color})
	}
	return gradientStops(out)
}

type gradientStops []GradientStop

func (gs gradientStops) colorAt(t float64) Color {
	if t <= gs[0].At {
		return gs[0].Color
	}
	if t >= gs[len(gs)-1].At {
		return gs[len(gs)-1].Color
	}
	for i := 1; i < len(gs); i++ {
		if t <= gs[i].At {
			a, b := gs[i-1], gs[i]
			span := b.At - a.At
			var lt float64
			if span == 0 {
				lt = 0
			} else {
				lt = (t - a.At) / span
			}
			return a.Color.Mix(b.Color, lt)
		}
	}
	return gs[len(gs)-1].Color
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var buf [3]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

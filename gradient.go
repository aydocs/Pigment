package pigment

import "strings"

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
func Stops(stops ...GradientStop) GradientOpt {
	return func(c *gradientCfg) { c.stops = stops }
}

// Gradient paints text with a smooth color transition. Whitespace is left
// uncolored so layouts stay intact.
func Gradient(text string, start, end Color, opts ...GradientOpt) string {
	cfg := gradientCfg{stops: []GradientStop{{At: 0, Color: start}, {At: 1, Color: end}}}
	for _, o := range opts {
		o(&cfg)
	}
	stops := normalizeStops(cfg.stops)

	runes := []rune(text)
	n := len(runes)
	if n == 0 {
		return text
	}
	if n == 1 {
		return paintRune(runes[0], stops.colorAt(0))
	}

	var b strings.Builder
	for i, r := range runes {
		if r == ' ' || r == '\t' || r == '\n' {
			b.WriteRune(r)
			continue
		}
		t := float64(i) / float64(n-1)
		b.WriteString(paintRune(r, stops.colorAt(t)))
	}
	return b.String()
}

func paintRune(r rune, c Color) string {
	if !Enabled {
		return string(r)
	}
	return esc + "38;2;" + itoa(int(c.R)) + ";" + itoa(int(c.G)) + ";" + itoa(int(c.B)) + "m" + string(r) + reset
}

// RainbowOpt configures a Rainbow call.
type RainbowOpt func(*rainbowCfg)

type rainbowCfg struct {
	start float64
	end   float64
	sat   float64
	light float64
}

// RainbowStart sets the starting hue (0..360).
func RainbowStart(h float64) RainbowOpt {
	return func(c *rainbowCfg) { c.start = h }
}

// RainbowEnd sets the ending hue (0..360). Defaults to start+360.
func RainbowEnd(h float64) RainbowOpt {
	return func(c *rainbowCfg) { c.end = h }
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
	if cfg.end == 0 {
		cfg.end = cfg.start + 360
	}

	runes := []rune(text)
	n := len(runes)
	if n == 0 {
		return text
	}

	var b strings.Builder
	for i, r := range runes {
		if r == ' ' || r == '\t' || r == '\n' {
			b.WriteRune(r)
			continue
		}
		var t float64
		if n == 1 {
			t = 0
		} else {
			t = float64(i) / float64(n-1)
		}
		hue := cfg.start + t*(cfg.end-cfg.start)
		b.WriteString(paintRune(r, HSL(hue, cfg.sat, cfg.light)))
	}
	return b.String()
}

func normalizeStops(in []GradientStop) gradientStops {
	stops := make([]GradientStop, 0, len(in))
	stops = append(stops, in...)
	stops = append(stops, GradientStop{At: 0, Color: in[0].Color})
	stops = append(stops, GradientStop{At: 1, Color: in[len(in)-1].Color})
	for i := 1; i < len(stops); i++ {
		for j := i; j > 0 && stops[j].At < stops[j-1].At; j-- {
			stops[j], stops[j-1] = stops[j-1], stops[j]
		}
	}
	return gradientStops(stops)
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

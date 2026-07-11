package pigment

import (
	"fmt"
	"strconv"
	"strings"
)

// Color is a 24-bit RGB color.
type Color struct {
	R, G, B uint8
}

// RGB builds a Color from 0-255 int components (clamped).
func RGB(r, g, b int) Color {
	return Color{clamp8(r), clamp8(g), clamp8(b)}
}

// Hex parses "#ff0000", "ff0000", "#f00" or "f00".
func Hex(s string) (Color, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "#")
	switch len(s) {
	case 3:
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	case 6:
	default:
		return Color{}, fmt.Errorf("pigment: bad hex color %q", s)
	}
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return Color{}, fmt.Errorf("pigment: bad hex color %q: %w", s, err)
	}
	return Color{uint8(v >> 16), uint8(v >> 8), uint8(v)}, nil
}

// MustHex is like Hex but panics on error.
func MustHex(s string) Color {
	c, err := Hex(s)
	if err != nil {
		panic(err)
	}
	return c
}

// Mix blends c and d at t (0..1).
func (c Color) Mix(d Color, t float64) Color {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return Color{
		R: uint8(float64(c.R) + t*float64(int(d.R)-int(c.R))),
		G: uint8(float64(c.G) + t*float64(int(d.G)-int(c.G))),
		B: uint8(float64(c.B) + t*float64(int(d.B)-int(c.B))),
	}
}

// HSL builds a Color from hue(0..360), sat(0..1), light(0..1).
func HSL(h, s, l float64) Color {
	h = mod(h, 360)
	c := (1 - absF(2*l-1)) * s
	x := c * (1 - absF(mod(h/60, 2)-1))
	m := l - c/2
	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = c, x, 0
	case h < 120:
		r, g, b = x, c, 0
	case h < 180:
		r, g, b = 0, c, x
	case h < 240:
		r, g, b = 0, x, c
	case h < 300:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}
	return Color{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
	}
}

// Palette returns the RGB for an xterm 256-color index (0..255).
func Palette(n int) Color {
	switch {
	case n < 0:
		n = 0
	case n > 255:
		n = 255
	}
	switch {
	case n < 16:
		if n < 8 {
			return systemColors[n]
		}
		return brightColors[n-8]
	case n < 232:
		i := n - 16
		return Color{colorCube[i/36], colorCube[(i/6)%6], colorCube[i%6]}
	default:
		v := 8 + (n-232)*10
		return Color{uint8(v), uint8(v), uint8(v)}
	}
}

var colorCube = [6]uint8{0, 95, 135, 175, 215, 255}

var systemColors = [8]Color{
	{0, 0, 0}, {205, 0, 0}, {0, 205, 0}, {205, 205, 0},
	{0, 0, 238}, {205, 0, 205}, {0, 205, 205}, {229, 229, 229},
}

var brightColors = [8]Color{
	{127, 127, 127}, {255, 0, 0}, {0, 255, 0}, {255, 255, 0},
	{92, 92, 255}, {255, 0, 255}, {0, 255, 255}, {255, 255, 255},
}

// namedColors maps CSS/SVG color names to RGB.
var namedColors = map[string]Color{
	"black": {0, 0, 0}, "white": {255, 255, 255}, "red": {255, 0, 0},
	"green": {0, 128, 0}, "lime": {0, 255, 0}, "blue": {0, 0, 255},
	"yellow": {255, 255, 0}, "cyan": {0, 255, 255}, "aqua": {0, 255, 255},
	"magenta": {255, 0, 255}, "fuchsia": {255, 0, 255}, "silver": {192, 192, 192},
	"gray": {128, 128, 128}, "grey": {128, 128, 128}, "maroon": {128, 0, 0},
	"olive": {128, 128, 0}, "navy": {0, 0, 128}, "teal": {0, 128, 128},
	"purple": {128, 0, 128}, "orange": {255, 165, 0}, "pink": {255, 192, 203},
	"brown": {165, 42, 42}, "gold": {255, 215, 0}, "tomato": {255, 99, 71},
	"coral": {255, 127, 80}, "salmon": {250, 128, 114}, "indigo": {75, 0, 130},
	"violet": {238, 130, 238}, "khaki": {240, 230, 140}, "orchid": {218, 112, 214},
	"plum": {221, 160, 221}, "crimson": {220, 20, 60}, "chocolate": {210, 105, 30},
	"turquoise": {64, 224, 208}, "skyblue": {135, 206, 235}, "steelblue": {70, 130, 180},
	"dodgerblue": {30, 144, 255}, "royalblue": {65, 105, 225}, "slateblue": {106, 90, 205},
	"lavender": {230, 230, 250}, "mintcream": {245, 255, 250}, "mistyrose": {255, 228, 225},
	"seashell": {255, 245, 238}, "snow": {255, 250, 250}, "ivory": {255, 255, 240},
	"beige": {245, 245, 220}, "linen": {250, 240, 230}, "lavenderblush": {255, 240, 245},
	"lightgray": {211, 211, 211}, "lightgrey": {211, 211, 211}, "darkgray": {169, 169, 169},
	"darkgrey": {169, 169, 169}, "lightblue": {173, 216, 230}, "lightgreen": {144, 238, 144},
	"lightcoral": {240, 128, 128}, "lightcyan": {224, 255, 255}, "lightpink": {255, 182, 193},
	"lightsalmon": {255, 160, 122}, "lightseagreen": {32, 178, 170}, "lightskyblue": {135, 206, 250},
	"lightslategray": {119, 136, 153}, "lightsteelblue": {176, 196, 222}, "lightyellow": {255, 255, 224},
	"darkblue": {0, 0, 139}, "darkcyan": {0, 139, 139}, "darkgreen": {0, 100, 0},
	"darkkhaki": {189, 183, 107}, "darkmagenta": {139, 0, 139}, "darkolivegreen": {85, 107, 47},
	"darkorange": {255, 140, 0}, "darkorchid": {153, 50, 204}, "darkred": {139, 0, 0},
	"darksalmon": {233, 150, 122}, "darkseagreen": {143, 188, 143}, "darkslateblue": {72, 61, 139},
	"darkslategray": {47, 79, 79}, "darkturquoise": {0, 206, 209}, "darkviolet": {148, 0, 211},
	"deeppink": {255, 20, 147}, "deepskyblue": {0, 191, 255}, "dimgray": {105, 105, 105},
	"firebrick": {178, 34, 34}, "forestgreen": {34, 139, 34}, "gainsboro": {220, 220, 220},
	"ghostwhite": {248, 248, 255}, "greenyellow": {173, 255, 47}, "honeydew": {240, 255, 240},
	"hotpink": {255, 105, 180}, "indianred": {205, 92, 92}, "lawngreen": {124, 252, 0},
	"lemonchiffon": {255, 250, 205}, "mediumaquamarine": {102, 205, 170}, "mediumblue": {0, 0, 205},
	"mediumorchid": {186, 85, 211}, "mediumpurple": {147, 112, 219}, "mediumseagreen": {60, 179, 113},
	"mediumslateblue": {123, 104, 238}, "mediumspringgreen": {0, 250, 154}, "mediumturquoise": {72, 209, 204},
	"mediumvioletred": {199, 21, 133}, "midnightblue": {25, 25, 112}, "navajowhite": {255, 222, 173},
	"oldlace": {253, 245, 230}, "olivedrab": {107, 142, 35}, "orangered": {255, 69, 0},
	"palegoldenrod": {238, 232, 170}, "palegreen": {152, 251, 152}, "paleturquoise": {175, 238, 238},
	"palevioletred": {219, 112, 147}, "papayawhip": {255, 239, 213}, "peachpuff": {255, 218, 185},
	"peru": {205, 133, 63}, "powderblue": {176, 224, 230}, "rosybrown": {188, 143, 143},
	"saddlebrown": {139, 69, 19}, "sandybrown": {244, 164, 96}, "seagreen": {46, 139, 87},
	"sienna": {160, 82, 45}, "slategray": {112, 128, 144}, "springgreen": {0, 255, 127},
	"tan": {210, 180, 140}, "thistle": {216, 191, 216}, "wheat": {245, 222, 179},
	"whitesmoke": {245, 245, 245}, "yellowgreen": {154, 205, 50}, "rebeccapurple": {102, 51, 153},

	"hi-black":   brightColors[0],
	"hi-red":     brightColors[1],
	"hi-green":   brightColors[2],
	"hi-yellow":  brightColors[3],
	"hi-blue":    brightColors[4],
	"hi-magenta": brightColors[5],
	"hi-cyan":    brightColors[6],
	"hi-white":   brightColors[7],

	"bright-black":   brightColors[0],
	"bright-red":     brightColors[1],
	"bright-green":   brightColors[2],
	"bright-yellow":  brightColors[3],
	"bright-blue":    brightColors[4],
	"bright-magenta": brightColors[5],
	"bright-cyan":    brightColors[6],
	"bright-white":   brightColors[7],
}

// Named returns the RGB for a color name (case-insensitive).
func Named(name string) (Color, bool) {
	c, ok := namedColors[strings.ToLower(strings.TrimSpace(name))]
	return c, ok
}

func clamp8(v int) uint8 {
	switch {
	case v < 0:
		return 0
	case v > 255:
		return 255
	default:
		return uint8(v)
	}
}

func mod(a, n float64) float64 {
	r := mathMod(a, n)
	if r < 0 {
		r += n
	}
	return r
}

func mathMod(a, n float64) float64 {
	return a - n*float64(int64(a/n))
}

func absF(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

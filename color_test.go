package pigment

import (
	"bytes"
	"strings"
	"testing"
)

func withColor(t *testing.T) func() {
	t.Helper()
	old := Enabled
	Enabled = true
	return func() { Enabled = old }
}

func TestColorBasics(t *testing.T) {
	if c, err := Hex("#ff8800"); err != nil || c.R != 255 || c.G != 136 || c.B != 0 {
		t.Fatalf("Hex failed: %v %v", c, err)
	}
	if c := RGB(10, 20, 30); c.R != 10 || c.G != 20 || c.B != 30 {
		t.Fatalf("RGB failed: %v", c)
	}
	if c, ok := Named("red"); !ok || c.R != 255 || c.G != 0 || c.B != 0 {
		t.Fatalf("Named failed: %v %v", c, ok)
	}
	if c, ok := Named("nope"); ok || c != (Color{}) {
		t.Fatal("Named should return zero + false for unknown")
	}
	if Palette(16) != (Color{colorCube[0], colorCube[0], colorCube[0]}) {
		t.Fatalf("Palette cube off: %v", Palette(16))
	}
	if Palette(232) != (Color{8, 8, 8}) {
		t.Fatalf("Palette gray off: %v", Palette(232))
	}
}

func TestHiColorNames(t *testing.T) {
	tests := []struct {
		name string
		idx  int
	}{
		{"hi-black", 0}, {"hi-red", 1}, {"hi-green", 2}, {"hi-yellow", 3},
		{"hi-blue", 4}, {"hi-magenta", 5}, {"hi-cyan", 6}, {"hi-white", 7},
		{"bright-black", 0}, {"bright-red", 1}, {"bright-green", 2}, {"bright-yellow", 3},
		{"bright-blue", 4}, {"bright-magenta", 5}, {"bright-cyan", 6}, {"bright-white", 7},
	}
	for _, tt := range tests {
		c, ok := Named(tt.name)
		if !ok {
			t.Fatalf("Named(%q) not found", tt.name)
		}
		want := brightColors[tt.idx]
		if c != want {
			t.Fatalf("Named(%q) = %v, want %v", tt.name, c, want)
		}
	}
}

func TestMixAndHSL(t *testing.T) {
	a := RGB(0, 0, 0)
	b := RGB(255, 255, 255)
	m := a.Mix(b, 0.5)
	if m.R != 127 && m.R != 128 {
		t.Fatalf("Mix mid off: %v", m)
	}
	r := HSL(0, 1, 0.5)
	if r.R != 255 || r.G != 0 || r.B != 0 {
		t.Fatalf("HSL red off: %v", r)
	}
}

func TestStylePaint(t *testing.T) {
	defer withColor(t)()
	s := New().Fg("red").Bold()
	got := s.Paint("hi")
	want := "\x1b[1;38;2;255;0;0mhi\x1b[0m"
	if got != want {
		t.Fatalf("Paint = %q, want %q", got, want)
	}

	bg := New().BgRGB(0, 0, 0).Underline()
	got = bg.Paint("x")
	if !strings.Contains(got, "4;48;2;0;0;0") {
		t.Fatalf("Bg paint off: %q", got)
	}

	pal := New().FgPalette(196)
	if !strings.Contains(pal.Paint("y"), "38;5;196") {
		t.Fatalf("Palette paint off: %q", pal.Paint("y"))
	}
}

func TestStyleDisabled(t *testing.T) {
	old := Enabled
	Enabled = false
	defer func() { Enabled = old }()
	if got := New().Fg("red").Paint("hi"); got != "hi" {
		t.Fatalf("disabled paint should be plain, got %q", got)
	}
	s := New().Fg("red").Disable()
	if got := s.Paint("hi"); got != "hi" {
		t.Fatalf("per-style disable failed: %q", got)
	}
	s2 := New().Fg("red").Enable()
	if !strings.Contains(s2.Paint("hi"), "\x1b[") {
		t.Fatalf("per-style enable failed: %q", s2.Paint("hi"))
	}
}

func TestGradient(t *testing.T) {
	defer withColor(t)()
	out := Gradient("abc", RGB(255, 0, 0), RGB(0, 0, 255))
	if !strings.Contains(out, "\x1b[") {
		t.Fatalf("gradient should contain escapes: %q", out)
	}
	if len([]rune(out)) == 0 {
		t.Fatal("empty gradient")
	}
	sp := Gradient("a b", RGB(0, 0, 0), RGB(255, 255, 255))
	if !strings.Contains(sp, " ") {
		t.Fatalf("gradient dropped whitespace: %q", sp)
	}
	ms := Gradient("xyz", RGB(0, 0, 0), RGB(0, 0, 0), Stops(GradientStop{At: 0, Color: RGB(255, 0, 0)}, GradientStop{At: 0.5, Color: RGB(0, 255, 0)}, GradientStop{At: 1, Color: RGB(0, 0, 255)}))
	if !strings.Contains(ms, "38;2;0;255;0") {
		t.Fatalf("multi-stop gradient missing mid color: %q", ms)
	}
}

func TestRainbow(t *testing.T) {
	defer withColor(t)()
	out := Rainbow("hello")
	if !strings.Contains(out, "\x1b[38;2;") {
		t.Fatalf("rainbow should be colored: %q", out)
	}
	if !strings.Contains(out, "38;2;255;0;0") {
		t.Fatalf("rainbow start not red: %q", out)
	}
}

func TestMarkup(t *testing.T) {
	defer withColor(t)()
	got := Render("<red>hi</red> there")
	if !strings.Contains(got, "\x1b[38;2;255;0;0mhi\x1b[0m") {
		t.Fatalf("markup red off: %q", got)
	}
	if !strings.Contains(got, "there") {
		t.Fatalf("markup plain tail lost: %q", got)
	}
	nested := Render("<bold><blue>x</blue>y")
	if !strings.Contains(nested, "1;38;2;0;0;255mx\x1b[0m") {
		t.Fatalf("markup nesting off: %q", nested)
	}
	if !strings.HasSuffix(nested, "\x1b[0m") {
		t.Fatalf("markup not closed: %q", nested)
	}
	attr := Render("<bold underline>z</>")
	if !strings.Contains(attr, "1;4") {
		t.Fatalf("markup attrs off: %q", attr)
	}
	hb := Render("<fg=#00ff00 bg=blue>w</>")
	if !strings.Contains(hb, "38;2;0;255;0") || !strings.Contains(hb, "48;2;0;0;255") {
		t.Fatalf("markup hex/bg off: %q", hb)
	}
}

func TestTheme(t *testing.T) {
	defer withColor(t)()
	if Lookup("error") == nil {
		t.Fatal("default theme missing error")
	}
	got := Paint("success", "ok")
	if !strings.Contains(got, "\x1b[") {
		t.Fatalf("theme paint not colored: %q", got)
	}
	if Paint("nope", "x") != "x" {
		t.Fatal("unknown theme name should return plain")
	}
	Register("custom", New().Fg("pink"))
	if Lookup("custom") == nil {
		t.Fatal("Register failed")
	}
}

func TestLink(t *testing.T) {
	defer withColor(t)()
	got := Link("https://x.com", "x")
	if !strings.Contains(got, "\x1b]8;;https://x.com\x1b\\x\x1b]8;;\x1b\\") {
		t.Fatalf("link off: %q", got)
	}
	old := Enabled
	Enabled = false
	if Link("https://x.com", "x") != "x" {
		t.Fatal("link should be plain when disabled")
	}
	Enabled = old
	if Link("", "x") != "x" {
		t.Fatal("empty url should be plain")
	}
}

func TestEquals(t *testing.T) {
	a := New().Fg("red").Bold()
	b := New().Fg("red").Bold()
	c := New().Fg("red")
	if !a.Equals(b) {
		t.Fatal("equal styles should match")
	}
	if a.Equals(c) {
		t.Fatal("different styles should not match")
	}
	if a.Equals(nil) {
		t.Fatal("nil compare")
	}
}

func TestConvenienceStrings(t *testing.T) {
	defer withColor(t)()

	if !strings.Contains(RedString("test"), "\x1b[38;2;255;0;0m") {
		t.Fatalf("RedString should be red: %q", RedString("test"))
	}
	if !strings.Contains(GreenString("test"), "\x1b[38;2;0;128;0m") {
		t.Fatalf("GreenString should be green: %q", GreenString("test"))
	}
	if !strings.Contains(BlueString("test"), "\x1b[38;2;0;0;255m") {
		t.Fatalf("BlueString should be blue: %q", BlueString("test"))
	}
	if !strings.Contains(HiRedString("test"), "\x1b[38;2;255;0;0m") {
		t.Fatalf("HiRedString should be bright red: %q", HiRedString("test"))
	}
	if !strings.Contains(BgRedString("test"), "\x1b[48;2;255;0;0m") {
		t.Fatalf("BgRedString should have red bg: %q", BgRedString("test"))
	}
	if !strings.Contains(HiBgRedString("test"), "\x1b[48;2;255;0;0m") {
		t.Fatalf("HiBgRedString should have bright red bg: %q", HiBgRedString("test"))
	}

	formatted := RedString("hello %s", "world")
	if !strings.Contains(formatted, "hello world") {
		t.Fatalf("RedString with format should contain formatted text: %q", formatted)
	}

	plain := BlackString("text")
	if !strings.Contains(plain, "\x1b[38;2;0;0;0m") {
		t.Fatalf("BlackString should be black: %q", plain)
	}
}

func TestConveniencePrint(t *testing.T) {
	defer withColor(t)()

	var buf bytes.Buffer
	old := Output
	Output = &buf
	defer func() { Output = old }()

	Red("hello")
	if !strings.Contains(buf.String(), "\x1b[38;2;255;0;0m") || !strings.Contains(buf.String(), "hello\n") {
		t.Fatalf("Red print failed: %q", buf.String())
	}

	buf.Reset()
	Red("hello %s", "world")
	if !strings.Contains(buf.String(), "hello world") {
		t.Fatalf("Red print format failed: %q", buf.String())
	}

	buf.Reset()
	BgRed("bg test")
	if !strings.Contains(buf.String(), "\x1b[48;2;255;0;0m") {
		t.Fatalf("BgRed print failed: %q", buf.String())
	}

	buf.Reset()
	HiRed("hi text")
	if !strings.Contains(buf.String(), "\x1b[38;2;255;0;0m") {
		t.Fatalf("HiRed print failed: %q", buf.String())
	}

	buf.Reset()
	HiBgRed("hi bg text")
	if !strings.Contains(buf.String(), "\x1b[48;2;255;0;0m") {
		t.Fatalf("HiBgRed print failed: %q", buf.String())
	}
}

func TestConveniencePrintDisabled(t *testing.T) {
	old := Enabled
	Enabled = false
	defer func() { Enabled = old }()

	var buf bytes.Buffer
	oldOut := Output
	Output = &buf
	defer func() { Output = oldOut }()

	Red("plain text")
	if buf.String() != "plain text\n" {
		t.Fatalf("Red should be plain when disabled: %q", buf.String())
	}

	if RedString("plain") != "plain" {
		t.Fatalf("RedString should be plain when disabled: %q", RedString("plain"))
	}
}

func TestConvenienceColorsMap(t *testing.T) {
	if Colors["red"] == nil {
		t.Fatal("Colors map missing red")
	}
	if Colors["hi-red"] == nil {
		t.Fatal("Colors map missing hi-red")
	}
	if Colors["hi-black"] == nil {
		t.Fatal("Colors map missing hi-black")
	}
}

func TestStylePrintFunc(t *testing.T) {
	defer withColor(t)()

	var buf bytes.Buffer
	s := New().Fg("red")
	pf := s.FprintFunc()
	pf(&buf, "hello")
	if !strings.Contains(buf.String(), "\x1b[38;2;255;0;0mhello\x1b[0m") {
		t.Fatalf("FprintFunc failed: %q", buf.String())
	}
}

func TestStyleSprintFunc(t *testing.T) {
	defer withColor(t)()

	s := New().Fg("red").Bold()
	sf := s.SprintFunc()
	result := sf("hello")
	if !strings.Contains(result, "\x1b[1;38;2;255;0;0m") {
		t.Fatalf("SprintFunc failed: %q", result)
	}
}

func TestStyleSprintfFunc(t *testing.T) {
	defer withColor(t)()

	s := New().Fg("blue")
	sf := s.SprintfFunc()
	result := sf("hello %s", "world")
	if !strings.Contains(result, "hello world") {
		t.Fatalf("SprintfFunc formatting failed: %q", result)
	}
	if !strings.Contains(result, "\x1b[38;2;0;0;255m") {
		t.Fatalf("SprintfFunc color failed: %q", result)
	}
}

func TestPackageLevelPrintFunc(t *testing.T) {
	defer withColor(t)()

	var buf bytes.Buffer
	old := Output
	Output = &buf
	defer func() { Output = old }()

	redPrint := PrintFunc("red")
	redPrint("test")
	if !strings.Contains(buf.String(), "\x1b[38;2;255;0;0m") {
		t.Fatalf("PrintFunc failed: %q", buf.String())
	}
}

func TestPackageLevelSprintfFunc(t *testing.T) {
	defer withColor(t)()
	format := SprintfFunc("red")
	result := format("hello %s", "world")
	if !strings.Contains(result, "\x1b[38;2;255;0;0m") {
		t.Fatalf("SprintfFunc failed: %q", result)
	}
	if !strings.Contains(result, "hello world") {
		t.Fatalf("SprintfFunc content failed: %q", result)
	}
}

func TestConvenienceBlackWhite(t *testing.T) {
	if WhiteString("test") == "" {
		t.Fatal("WhiteString should not be empty")
	}
	if BlackString("test") == "" {
		t.Fatal("BlackString should not be empty")
	}
}

func TestFgBgConvenience(t *testing.T) {
	defer withColor(t)()
	s := Fg("red").Bold()
	if !strings.Contains(s.Paint("x"), "1;38;2;255;0;0") {
		t.Fatalf("Fg convenience failed: %q", s.Paint("x"))
	}

	s2 := HexStyle("#00ff00")
	if !strings.Contains(s2.Paint("x"), "38;2;0;255;0") {
		t.Fatalf("HexStyle failed: %q", s2.Paint("x"))
	}

	s3 := RGBStyle(255, 0, 255)
	if !strings.Contains(s3.Paint("x"), "38;2;255;0;255") {
		t.Fatalf("RGBStyle failed: %q", s3.Paint("x"))
	}

	s4 := PaletteStyle(196)
	if !strings.Contains(s4.Paint("x"), "38;5;196") {
		t.Fatalf("PaletteStyle failed: %q", s4.Paint("x"))
	}
}

func TestStyleSprintlnFunc(t *testing.T) {
	defer withColor(t)()

	s := New().Fg("cyan")
	slf := s.SprintlnFunc()
	result := slf("line")
	if !strings.HasSuffix(result, "\n") {
		t.Fatalf("SprintlnFunc should end with newline: %q", result)
	}
	if !strings.Contains(result, "\x1b[38;2;0;255;255m") {
		t.Fatalf("SprintlnFunc color failed: %q", result)
	}
}

func TestStyleFprintlnFunc(t *testing.T) {
	defer withColor(t)()

	var buf bytes.Buffer
	s := New().Fg("magenta")
	fplf := s.FprintlnFunc()
	fplf(&buf, "line")
	if !strings.Contains(buf.String(), "\x1b[38;2;255;0;255m") {
		t.Fatalf("FprintlnFunc failed: %q", buf.String())
	}
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Fatalf("FprintlnFunc should end with newline: %q", buf.String())
	}
}

func TestStyleFprintfFunc(t *testing.T) {
	defer withColor(t)()

	var buf bytes.Buffer
	s := New().Fg("yellow")
	fpf := s.FprintfFunc()
	fpf(&buf, "hello %s", "world")
	if !strings.Contains(buf.String(), "\x1b[38;2;255;255;0m") {
		t.Fatalf("FprintfFunc failed: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "hello world") {
		t.Fatalf("FprintfFunc content failed: %q", buf.String())
	}
}

func TestStylePrintfFunc(t *testing.T) {
	defer withColor(t)()

	var buf bytes.Buffer
	old := Output
	Output = &buf
	defer func() { Output = old }()

	s := New().Fg("green")
	pff := s.PrintfFunc()
	pff("hello %s", "world")
	if !strings.Contains(buf.String(), "\x1b[38;2;0;128;0m") {
		t.Fatalf("PrintfFunc failed: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "hello world") {
		t.Fatalf("PrintfFunc content failed: %q", buf.String())
	}
}

func TestPackageLevelPrintfFunc(t *testing.T) {
	defer withColor(t)()

	var buf bytes.Buffer
	old := Output
	Output = &buf
	defer func() { Output = old }()

	pff := PrintfFunc("blue")
	pff("test %d", 42)
	if !strings.Contains(buf.String(), "\x1b[38;2;0;0;255m") {
		t.Fatalf("Package-level PrintfFunc failed: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "test 42") {
		t.Fatalf("Package-level PrintfFunc content failed: %q", buf.String())
	}
}

func TestPackageLevelPrintlnFunc(t *testing.T) {
	defer withColor(t)()

	var buf bytes.Buffer
	old := Output
	Output = &buf
	defer func() { Output = old }()

	plf := PrintlnFunc("red")
	plf("test")
	if !strings.Contains(buf.String(), "\x1b[38;2;255;0;0m") {
		t.Fatalf("Package-level PrintlnFunc failed: %q", buf.String())
	}
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Fatalf("Package-level PrintlnFunc should end with newline: %q", buf.String())
	}
}

func TestPackageLevelSprintlnFunc(t *testing.T) {
	defer withColor(t)()
	slf := SprintlnFunc("red")
	result := slf("test")
	if !strings.Contains(result, "\x1b[38;2;255;0;0m") {
		t.Fatalf("Package-level SprintlnFunc failed: %q", result)
	}
	if !strings.HasSuffix(result, "\n") {
		t.Fatalf("Package-level SprintlnFunc should end with newline: %q", result)
	}
}

func TestDisabledConvenience(t *testing.T) {
	old := Enabled
	Enabled = false
	defer func() { Enabled = old }()

	if RedString("x") != "x" {
		t.Fatal("expected plain when disabled")
	}
	if GreenString("y") != "y" {
		t.Fatal("expected plain when disabled")
	}
	if BlueString("z") != "z" {
		t.Fatal("expected plain when disabled")
	}
}

func TestAllConveniencePrintFunctions(t *testing.T) {
	defer withColor(t)()

	tests := []struct {
		name string
		fn   func(string, ...any)
		fg   string
		bg   bool
	}{
		{"Black", Black, "black", false},
		{"Red", Red, "red", false},
		{"Green", Green, "green", false},
		{"Yellow", Yellow, "yellow", false},
		{"Blue", Blue, "blue", false},
		{"Magenta", Magenta, "magenta", false},
		{"Cyan", Cyan, "cyan", false},
		{"White", White, "white", false},
		{"HiBlack", HiBlack, "hi-black", false},
		{"HiRed", HiRed, "hi-red", false},
		{"HiGreen", HiGreen, "hi-green", false},
		{"HiYellow", HiYellow, "hi-yellow", false},
		{"HiBlue", HiBlue, "hi-blue", false},
		{"HiMagenta", HiMagenta, "hi-magenta", false},
		{"HiCyan", HiCyan, "hi-cyan", false},
		{"HiWhite", HiWhite, "hi-white", false},
		{"BgBlack", BgBlack, "black", true},
		{"BgRed", BgRed, "red", true},
		{"BgGreen", BgGreen, "green", true},
		{"BgYellow", BgYellow, "yellow", true},
		{"BgBlue", BgBlue, "blue", true},
		{"BgMagenta", BgMagenta, "magenta", true},
		{"BgCyan", BgCyan, "cyan", true},
		{"BgWhite", BgWhite, "white", true},
		{"HiBgBlack", HiBgBlack, "hi-black", true},
		{"HiBgRed", HiBgRed, "hi-red", true},
		{"HiBgGreen", HiBgGreen, "hi-green", true},
		{"HiBgYellow", HiBgYellow, "hi-yellow", true},
		{"HiBgBlue", HiBgBlue, "hi-blue", true},
		{"HiBgMagenta", HiBgMagenta, "hi-magenta", true},
		{"HiBgCyan", HiBgCyan, "hi-cyan", true},
		{"HiBgWhite", HiBgWhite, "hi-white", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			old := Output
			Output = &buf
			defer func() { Output = old }()

			tt.fn("test")
			if !strings.Contains(buf.String(), "\x1b[") {
				t.Fatalf("%s should produce escapes, got: %q", tt.name, buf.String())
			}
			if !strings.Contains(buf.String(), "test\n") {
				t.Fatalf("%s should contain newline, got: %q", tt.name, buf.String())
			}
		})
	}
}

func TestAllConvenienceStringFunctions(t *testing.T) {
	defer withColor(t)()

	tests := []struct {
		name string
		fn   func(string, ...any) string
		fg   string
		bg   bool
	}{
		{"BlackString", BlackString, "black", false},
		{"RedString", RedString, "red", false},
		{"GreenString", GreenString, "green", false},
		{"YellowString", YellowString, "yellow", false},
		{"BlueString", BlueString, "blue", false},
		{"MagentaString", MagentaString, "magenta", false},
		{"CyanString", CyanString, "cyan", false},
		{"WhiteString", WhiteString, "white", false},
		{"HiBlackString", HiBlackString, "hi-black", false},
		{"HiRedString", HiRedString, "hi-red", false},
		{"HiGreenString", HiGreenString, "hi-green", false},
		{"HiYellowString", HiYellowString, "hi-yellow", false},
		{"HiBlueString", HiBlueString, "hi-blue", false},
		{"HiMagentaString", HiMagentaString, "hi-magenta", false},
		{"HiCyanString", HiCyanString, "hi-cyan", false},
		{"HiWhiteString", HiWhiteString, "hi-white", false},
		{"BgBlackString", BgBlackString, "black", true},
		{"BgRedString", BgRedString, "red", true},
		{"BgGreenString", BgGreenString, "green", true},
		{"BgYellowString", BgYellowString, "yellow", true},
		{"BgBlueString", BgBlueString, "blue", true},
		{"BgMagentaString", BgMagentaString, "magenta", true},
		{"BgCyanString", BgCyanString, "cyan", true},
		{"BgWhiteString", BgWhiteString, "white", true},
		{"HiBgBlackString", HiBgBlackString, "hi-black", true},
		{"HiBgRedString", HiBgRedString, "hi-red", true},
		{"HiBgGreenString", HiBgGreenString, "hi-green", true},
		{"HiBgYellowString", HiBgYellowString, "hi-yellow", true},
		{"HiBgBlueString", HiBgBlueString, "hi-blue", true},
		{"HiBgMagentaString", HiBgMagentaString, "hi-magenta", true},
		{"HiBgCyanString", HiBgCyanString, "hi-cyan", true},
		{"HiBgWhiteString", HiBgWhiteString, "hi-white", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn("test")
			if !strings.Contains(result, "\x1b[") {
				t.Fatalf("%s should produce escapes, got: %q", tt.name, result)
			}
		})
	}
}

func TestSprintConvenience(t *testing.T) {
	defer withColor(t)()

	result := Sprint("red", "hello")
	if !strings.Contains(result, "\x1b[38;2;255;0;0m") {
		t.Fatalf("Sprint failed: %q", result)
	}
	if !strings.Contains(result, "hello") {
		t.Fatalf("Sprint content failed: %q", result)
	}

	result2 := Sprintf("red", "hello %s", "world")
	if !strings.Contains(result2, "hello world") {
		t.Fatalf("Sprintf failed: %q", result2)
	}
}

func TestFprintConvenience(t *testing.T) {
	defer withColor(t)()

	var buf bytes.Buffer
	Fprint(&buf, "red", "hello")
	if !strings.Contains(buf.String(), "\x1b[38;2;255;0;0m") {
		t.Fatalf("Fprint failed: %q", buf.String())
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Fatalf("Fprint content failed: %q", buf.String())
	}
}

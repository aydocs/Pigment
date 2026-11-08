package pigment

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestDetectFromEnvPrecedence(t *testing.T) {
	cases := []struct {
		name string
		env  map[string]string
		want bool
	}{
		{name: "force off", env: map[string]string{"FORCE_COLOR": "0", "CLICOLOR_FORCE": "1"}, want: false},
		{name: "force false off", env: map[string]string{"FORCE_COLOR": "false"}, want: false},
		{name: "force on", env: map[string]string{"FORCE_COLOR": "1", "NO_COLOR": "1"}, want: true},
		{name: "clicolor force", env: map[string]string{"CLICOLOR_FORCE": "1", "NO_COLOR": "1"}, want: true},
		{name: "no color", env: map[string]string{"NO_COLOR": "1"}, want: false},
		{name: "dumb terminal", env: map[string]string{"TERM": "dumb"}, want: false},
		{name: "clicolor off", env: map[string]string{"CLICOLOR": "0"}, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, key := range []string{"FORCE_COLOR", "CLICOLOR_FORCE", "NO_COLOR", "TERM", "CLICOLOR"} {
				t.Setenv(key, "")
				if err := os.Unsetenv(key); err != nil {
					t.Fatal(err)
				}
			}
			for key, value := range tc.env {
				t.Setenv(key, value)
			}
			if got := detectFromEnv(); got != tc.want {
				t.Fatalf("detectFromEnv() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStyleOutputAndErrors(t *testing.T) {
	withColors(t, true)
	style := New().FgRGB(1, 2, 3).BgPalette(999).Bold()
	if got := style.Sprint("x"); !strings.Contains(got, "38;2;1;2;3") || !strings.Contains(got, "48;5;255") {
		t.Fatalf("unexpected style output: %q", got)
	}
	if got := style.Sprintln("x"); !strings.HasSuffix(got, "\n") {
		t.Fatalf("Sprintln() = %q, want trailing newline", got)
	}
	var buf bytes.Buffer
	if _, err := style.Fprintln(&buf, "line"); err != nil || !strings.HasSuffix(buf.String(), "\n") {
		t.Fatalf("Fprintln() = %q, %v", buf.String(), err)
	}
	writeErr := errors.New("write failed")
	if _, err := style.Fprint(errorWriter{err: writeErr}, "x"); !errors.Is(err, writeErr) {
		t.Fatalf("Fprint() error = %v, want %v", err, writeErr)
	}
}

type errorWriter struct{ err error }

func (w errorWriter) Write([]byte) (int, error) { return 0, w.err }

func TestSetOutputAndGlobalTheme(t *testing.T) {
	withColors(t, false)
	var buf bytes.Buffer
	SetOutput(&buf)
	t.Cleanup(func() { SetOutput(nil) })
	if _, err := New().Println("hello"); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "hello\n" {
		t.Fatalf("output = %q", buf.String())
	}
	Register("release-test", New().Fg("red"))
	t.Cleanup(func() { Register("release-test", nil) })
	if got := Paint("missing", "plain"); got != "plain" || Paint("release-test", "x") != "x" {
		t.Fatalf("unexpected theme rendering")
	}
}

func TestLinkEscapeAndMarkupEscape(t *testing.T) {
	withColors(t, true)
	got := Link("https://example.test/a\x1b]52;;bad", "text")
	if strings.Contains(got, "\x1b]52") || !strings.Contains(got, "%1B]52;;") {
		t.Fatalf("unsafe URL was not encoded: %q", got)
	}
	literal := "<red> & \\ value"
	if got := Render(Escape(literal)); got != literal {
		t.Fatalf("Escape/Render = %q, want %q", got, literal)
	}
}

func TestGradientAndRainbowEdgeCases(t *testing.T) {
	withColors(t, true)
	if got := Gradient("", RGB(0, 0, 0), RGB(255, 255, 255)); got != "" {
		t.Fatalf("empty Gradient = %q", got)
	}
	got := Rainbow("abc", RainbowStart(90), RainbowEnd(0), RainbowSaturation(2), RainbowLightness(-1))
	if strings.Contains(got, "abc") || strings.Count(got, "\x1b[38;2;") != 3 {
		t.Fatalf("Rainbow edge case output = %q", got)
	}
}

func TestThemeConcurrentSnapshots(t *testing.T) {
	withColors(t, false)
	theme := NewTheme()
	style := New().Fg("red")
	theme.Register("x", style)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				theme.Register("x", style)
				_ = theme.Paint("x", "value")
				_ = theme.Names()
			}
		}()
	}
	wg.Wait()
}

func TestPublicHelpersDoNotRequireTerminal(t *testing.T) {
	withColors(t, false)
	if got := RedString("%s", "ok"); got != "ok" {
		t.Fatalf("RedString() = %q", got)
	}
	if got := Sprintf("red", "%s", "ok"); got != "ok" {
		t.Fatalf("Sprintf() = %q", got)
	}
	if got := New().Fg("no-such-color").Sprint("plain"); got != "plain" {
		t.Fatalf("invalid color changed output: %q", got)
	}
}

func FuzzRenderNeverPanics(f *testing.F) {
	for _, seed := range []string{"plain", "<red>x</red>", "if x < y > z", "\\<literal>"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		_ = Render(input)
	})
}

func BenchmarkRender(b *testing.B) {
	previous := IsEnabled()
	SetEnabled(true)
	b.Cleanup(func() { SetEnabled(previous) })
	input := "<bold><#ff5555>error</#ff5555>: check logs</bold>"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Render(input)
	}
}

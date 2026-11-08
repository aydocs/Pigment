package pigment

import (
	"strings"
	"testing"
)

func withColors(t *testing.T, enabledState bool) {
	t.Helper()
	previous := IsEnabled()
	SetEnabled(enabledState)
	t.Cleanup(func() { SetEnabled(previous) })
}

func TestHexNamedAndPalette(t *testing.T) {
	if got := MustHex("#f08"); got != (Color{255, 0, 136}) {
		t.Fatalf("MustHex(#f08) = %#v", got)
	}
	if got, ok := Named(" RebeccaPurple "); !ok || got != (Color{102, 51, 153}) {
		t.Fatalf("Named(rebeccapurple) = %#v, %v", got, ok)
	}
	if got := Palette(999); got != Palette(255) {
		t.Fatalf("Palette should clamp high indexes: %#v", got)
	}
}

func TestRenderPreservesInvalidTags(t *testing.T) {
	withColors(t, true)
	input := "for i < len(xs) { <fg=not-a-color>value }"
	if got := Render(input); got != input {
		t.Fatalf("Render changed invalid markup: %q", got)
	}
}

func TestRenderNestedTags(t *testing.T) {
	withColors(t, true)
	got := Render("<red><bold>x</bold> y</red>")
	if !strings.Contains(got, "x") || !strings.Contains(got, " y") {
		t.Fatalf("Render lost nested text: %q", got)
	}
	if strings.Contains(got, "<red>") {
		t.Fatalf("Render left valid markup behind: %q", got)
	}
}

func TestGradientWhitespaceAndDisabledOutput(t *testing.T) {
	withColors(t, false)
	if got := Gradient("a b\n", RGB(255, 0, 0), RGB(0, 0, 255)); got != "a b\n" {
		t.Fatalf("disabled Gradient = %q", got)
	}
	withColors(t, true)
	got := Gradient("a b\n", RGB(255, 0, 0), RGB(0, 0, 255), Stops())
	if strings.ContainsAny(got, " \n") == false || !strings.Contains(got, "\x1b[") {
		t.Fatalf("Gradient did not preserve whitespace or paint text: %q", got)
	}
}

func TestThemeSnapshotsAreIndependent(t *testing.T) {
	withColors(t, false)
	theme := NewTheme()
	original := New().Fg("red")
	theme.Register("error", original)
	original.Fg("blue")
	snapshot := theme.Style("error")
	if snapshot == nil || !snapshot.Equals(New().Fg("red")) {
		t.Fatalf("theme did not retain an independent snapshot")
	}
	snapshot.Fg("green")
	if !theme.Style("error").Equals(New().Fg("red")) {
		t.Fatalf("mutating a snapshot changed the theme")
	}
}

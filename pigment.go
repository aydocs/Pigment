package pigment

import (
	"fmt"
	"io"
	"strings"
)

// Colors maps color names to pre-built foreground styles.
var Colors = map[string]*Style{
	"black":   New().Fg("black"),
	"red":     New().Fg("red"),
	"green":   New().Fg("green"),
	"yellow":  New().Fg("yellow"),
	"blue":    New().Fg("blue"),
	"magenta": New().Fg("magenta"),
	"cyan":    New().Fg("cyan"),
	"white":   New().Fg("white"),

	"hi-black":   New().Fg("hi-black"),
	"hi-red":     New().Fg("hi-red"),
	"hi-green":   New().Fg("hi-green"),
	"hi-yellow":  New().Fg("hi-yellow"),
	"hi-blue":    New().Fg("hi-blue"),
	"hi-magenta": New().Fg("hi-magenta"),
	"hi-cyan":    New().Fg("hi-cyan"),
	"hi-white":   New().Fg("hi-white"),
}

// Fg returns a new Style with the foreground set from a flexible value.
func Fg(v any) *Style { return New().Fg(v) }

// RGBStyle returns a new Style with the given 24-bit RGB foreground.
func RGBStyle(r, g, b int) *Style { return New().FgRGB(r, g, b) }

// HexStyle returns a new Style with the given hex foreground.
func HexStyle(h string) *Style { return New().FgHex(h) }

// PaletteStyle returns a new Style with the given xterm 256 foreground.
func PaletteStyle(n int) *Style { return New().FgPalette(n) }

// Sprint is a shortcut for Fg(v).Sprint(a...).
func Sprint(v any, a ...any) string { return Fg(v).Sprint(a...) }

// Sprintf is a shortcut for Fg(v).Sprintf(format, a...).
func Sprintf(v any, format string, a ...any) string { return Fg(v).Sprintf(format, a...) }

// Fprint writes colored, formatted arguments to w.
func Fprint(w io.Writer, v any, a ...any) (int, error) {
	return fmt.Fprint(w, Fg(v).Paint(fmt.Sprint(a...)))
}

// ---------------------------------------------------------------------------
// Convenience print functions
// ---------------------------------------------------------------------------

func colorPrint(fg string, format string, a ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	if len(a) == 0 {
		New().Fg(fg).Print(format)
	} else {
		New().Fg(fg).Printf(format, a...)
	}
}

func colorString(fg string, format string, a ...any) string {
	if len(a) == 0 {
		return New().Fg(fg).Sprint(format)
	}
	return New().Fg(fg).Sprintf(format, a...)
}

func bgColorPrint(bg string, format string, a ...any) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	if len(a) == 0 {
		New().Bg(bg).Print(format)
	} else {
		New().Bg(bg).Printf(format, a...)
	}
}

func bgColorString(bg string, format string, a ...any) string {
	if len(a) == 0 {
		return New().Bg(bg).Sprint(format)
	}
	return New().Bg(bg).Sprintf(format, a...)
}

// Black prints text with black foreground.
func Black(format string, a ...any) { colorPrint("black", format, a...) }

// Red prints text with red foreground.
func Red(format string, a ...any) { colorPrint("red", format, a...) }

// Green prints text with green foreground.
func Green(format string, a ...any) { colorPrint("green", format, a...) }

// Yellow prints text with yellow foreground.
func Yellow(format string, a ...any) { colorPrint("yellow", format, a...) }

// Blue prints text with blue foreground.
func Blue(format string, a ...any) { colorPrint("blue", format, a...) }

// Magenta prints text with magenta foreground.
func Magenta(format string, a ...any) { colorPrint("magenta", format, a...) }

// Cyan prints text with cyan foreground.
func Cyan(format string, a ...any) { colorPrint("cyan", format, a...) }

// White prints text with white foreground.
func White(format string, a ...any) { colorPrint("white", format, a...) }

// HiBlack prints text with hi-intensity black foreground.
func HiBlack(format string, a ...any) { colorPrint("hi-black", format, a...) }

// HiRed prints text with hi-intensity red foreground.
func HiRed(format string, a ...any) { colorPrint("hi-red", format, a...) }

// HiGreen prints text with hi-intensity green foreground.
func HiGreen(format string, a ...any) { colorPrint("hi-green", format, a...) }

// HiYellow prints text with hi-intensity yellow foreground.
func HiYellow(format string, a ...any) { colorPrint("hi-yellow", format, a...) }

// HiBlue prints text with hi-intensity blue foreground.
func HiBlue(format string, a ...any) { colorPrint("hi-blue", format, a...) }

// HiMagenta prints text with hi-intensity magenta foreground.
func HiMagenta(format string, a ...any) { colorPrint("hi-magenta", format, a...) }

// HiCyan prints text with hi-intensity cyan foreground.
func HiCyan(format string, a ...any) { colorPrint("hi-cyan", format, a...) }

// HiWhite prints text with hi-intensity white foreground.
func HiWhite(format string, a ...any) { colorPrint("hi-white", format, a...) }

// BgBlack prints text with black background.
func BgBlack(format string, a ...any) { bgColorPrint("black", format, a...) }

// BgRed prints text with red background.
func BgRed(format string, a ...any) { bgColorPrint("red", format, a...) }

// BgGreen prints text with green background.
func BgGreen(format string, a ...any) { bgColorPrint("green", format, a...) }

// BgYellow prints text with yellow background.
func BgYellow(format string, a ...any) { bgColorPrint("yellow", format, a...) }

// BgBlue prints text with blue background.
func BgBlue(format string, a ...any) { bgColorPrint("blue", format, a...) }

// BgMagenta prints text with magenta background.
func BgMagenta(format string, a ...any) { bgColorPrint("magenta", format, a...) }

// BgCyan prints text with cyan background.
func BgCyan(format string, a ...any) { bgColorPrint("cyan", format, a...) }

// BgWhite prints text with white background.
func BgWhite(format string, a ...any) { bgColorPrint("white", format, a...) }

// HiBgBlack prints text with hi-intensity black background.
func HiBgBlack(format string, a ...any) { bgColorPrint("hi-black", format, a...) }

// HiBgRed prints text with hi-intensity red background.
func HiBgRed(format string, a ...any) { bgColorPrint("hi-red", format, a...) }

// HiBgGreen prints text with hi-intensity green background.
func HiBgGreen(format string, a ...any) { bgColorPrint("hi-green", format, a...) }

// HiBgYellow prints text with hi-intensity yellow background.
func HiBgYellow(format string, a ...any) { bgColorPrint("hi-yellow", format, a...) }

// HiBgBlue prints text with hi-intensity blue background.
func HiBgBlue(format string, a ...any) { bgColorPrint("hi-blue", format, a...) }

// HiBgMagenta prints text with hi-intensity magenta background.
func HiBgMagenta(format string, a ...any) { bgColorPrint("hi-magenta", format, a...) }

// HiBgCyan prints text with hi-intensity cyan background.
func HiBgCyan(format string, a ...any) { bgColorPrint("hi-cyan", format, a...) }

// HiBgWhite prints text with hi-intensity white background.
func HiBgWhite(format string, a ...any) { bgColorPrint("hi-white", format, a...) }

// ---------------------------------------------------------------------------
// Convenience string functions
// ---------------------------------------------------------------------------

// BlackString returns text with black foreground.
func BlackString(format string, a ...any) string { return colorString("black", format, a...) }

// RedString returns text with red foreground.
func RedString(format string, a ...any) string { return colorString("red", format, a...) }

// GreenString returns text with green foreground.
func GreenString(format string, a ...any) string { return colorString("green", format, a...) }

// YellowString returns text with yellow foreground.
func YellowString(format string, a ...any) string { return colorString("yellow", format, a...) }

// BlueString returns text with blue foreground.
func BlueString(format string, a ...any) string { return colorString("blue", format, a...) }

// MagentaString returns text with magenta foreground.
func MagentaString(format string, a ...any) string { return colorString("magenta", format, a...) }

// CyanString returns text with cyan foreground.
func CyanString(format string, a ...any) string { return colorString("cyan", format, a...) }

// WhiteString returns text with white foreground.
func WhiteString(format string, a ...any) string { return colorString("white", format, a...) }

// HiBlackString returns text with hi-intensity black foreground.
func HiBlackString(format string, a ...any) string { return colorString("hi-black", format, a...) }

// HiRedString returns text with hi-intensity red foreground.
func HiRedString(format string, a ...any) string { return colorString("hi-red", format, a...) }

// HiGreenString returns text with hi-intensity green foreground.
func HiGreenString(format string, a ...any) string { return colorString("hi-green", format, a...) }

// HiYellowString returns text with hi-intensity yellow foreground.
func HiYellowString(format string, a ...any) string { return colorString("hi-yellow", format, a...) }

// HiBlueString returns text with hi-intensity blue foreground.
func HiBlueString(format string, a ...any) string { return colorString("hi-blue", format, a...) }

// HiMagentaString returns text with hi-intensity magenta foreground.
func HiMagentaString(format string, a ...any) string { return colorString("hi-magenta", format, a...) }

// HiCyanString returns text with hi-intensity cyan foreground.
func HiCyanString(format string, a ...any) string { return colorString("hi-cyan", format, a...) }

// HiWhiteString returns text with hi-intensity white foreground.
func HiWhiteString(format string, a ...any) string { return colorString("hi-white", format, a...) }

// BgBlackString returns text with black background.
func BgBlackString(format string, a ...any) string { return bgColorString("black", format, a...) }

// BgRedString returns text with red background.
func BgRedString(format string, a ...any) string { return bgColorString("red", format, a...) }

// BgGreenString returns text with green background.
func BgGreenString(format string, a ...any) string { return bgColorString("green", format, a...) }

// BgYellowString returns text with yellow background.
func BgYellowString(format string, a ...any) string { return bgColorString("yellow", format, a...) }

// BgBlueString returns text with blue background.
func BgBlueString(format string, a ...any) string { return bgColorString("blue", format, a...) }

// BgMagentaString returns text with magenta background.
func BgMagentaString(format string, a ...any) string { return bgColorString("magenta", format, a...) }

// BgCyanString returns text with cyan background.
func BgCyanString(format string, a ...any) string { return bgColorString("cyan", format, a...) }

// BgWhiteString returns text with white background.
func BgWhiteString(format string, a ...any) string { return bgColorString("white", format, a...) }

// HiBgBlackString returns text with hi-intensity black background.
func HiBgBlackString(format string, a ...any) string { return bgColorString("hi-black", format, a...) }

// HiBgRedString returns text with hi-intensity red background.
func HiBgRedString(format string, a ...any) string { return bgColorString("hi-red", format, a...) }

// HiBgGreenString returns text with hi-intensity green background.
func HiBgGreenString(format string, a ...any) string { return bgColorString("hi-green", format, a...) }

// HiBgYellowString returns text with hi-intensity yellow background.
func HiBgYellowString(format string, a ...any) string { return bgColorString("hi-yellow", format, a...) }

// HiBgBlueString returns text with hi-intensity blue background.
func HiBgBlueString(format string, a ...any) string { return bgColorString("hi-blue", format, a...) }

// HiBgMagentaString returns text with hi-intensity magenta background.
func HiBgMagentaString(format string, a ...any) string { return bgColorString("hi-magenta", format, a...) }

// HiBgCyanString returns text with hi-intensity cyan background.
func HiBgCyanString(format string, a ...any) string { return bgColorString("hi-cyan", format, a...) }

// HiBgWhiteString returns text with hi-intensity white background.
func HiBgWhiteString(format string, a ...any) string { return bgColorString("hi-white", format, a...) }

// ---------------------------------------------------------------------------
// Convenience print func generators
// ---------------------------------------------------------------------------

// PrintFunc returns a function that prints with the given foreground color.
func PrintFunc(fg string) func(a ...any) { return New().Fg(fg).PrintFunc() }

// PrintfFunc returns a function that formats and prints with the given foreground.
func PrintfFunc(fg string) func(format string, a ...any) { return New().Fg(fg).PrintfFunc() }

// PrintlnFunc returns a function that prints lines with the given foreground.
func PrintlnFunc(fg string) func(a ...any) { return New().Fg(fg).PrintlnFunc() }

// SprintFunc returns a function that returns colorized strings.
func SprintFunc(fg string) func(a ...any) string { return New().Fg(fg).SprintFunc() }

// SprintfFunc returns a function that returns formatted colorized strings.
func SprintfFunc(fg string) func(format string, a ...any) string { return New().Fg(fg).SprintfFunc() }

// SprintlnFunc returns a function that returns colorized lines.
func SprintlnFunc(fg string) func(a ...any) string { return New().Fg(fg).SprintlnFunc() }

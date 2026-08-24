// Package pigment provides modern terminal styling for Go: truecolor
// foregrounds and backgrounds, text attributes, gradients, rainbow effects,
// inline markup, OSC 8 hyperlinks, semantic themes, HSL color math, mixing,
// and the complete CSS named-color set — with zero external dependencies.
//
// # Quick start
//
//	pigment.Red("error: something broke")            // prints a red line
//	pigment.New().Fg("navy").Bold().Println("hi")    // fluent styles
//	fmt.Println(pigment.Sprint("#ff8800", "orange")) // any color value
//
// A style value accepts colors from many sources: "#ff8800" or "ff8800"
// hex strings (3- and 6-digit), CSS color names ("rebeccapurple"),
// [pigment.RGB], [pigment.HSL], [pigment.Palette] xterm indices (0-255),
// Color values, *Color pointers, and [3]int/[3]uint8 triples.
//
// # Printing helpers
//
// For every one of the 16 base/bright ANSI colors there is a foreground
// printer (Red, HiRed, ...), a background printer (BgRed, HiBgRed, ...)
// and two string-returning counterparts (RedString, BgRedString, ...).
// Printers append a newline; String forms never do.
//
// # Detection and environment
//
// Whether color is enabled is decided once at startup (and re-runnable via
// [DetectEnabled]) from the environment, in this order:
//
//  1. FORCE_COLOR: "0" or "false" forces OFF, any other value forces ON
//  2. CLICOLOR_FORCE: set and != "0" forces ON
//  3. NO_COLOR: set non-empty forces OFF (https://no-color.org)
//  4. TERM == "dumb" forces OFF
//  5. CLICOLOR == "0" forces OFF
//  6. otherwise: stdout is a terminal (checked with a real tty probe)
//
// [SetEnabled] overrides globally; Style.Enable/Disable override per style.
//
// # Markup
//
// [Render] parses a small inline language:
//
//	<red>text</red>           named foreground
//	<#ff8800>text</>          hex foreground (empty close tag allowed)
//	<bg=blue>text</bg=blue>   named background
//	<palette=196>text</>      xterm 256 index
//	<bold underline red>x</>  multiple tokens in one tag
//
// Tags nest. Keys and color names are case-insensitive. A tag containing
// any unrecognized token is NOT consumed — it is emitted verbatim as plain
// text, so prose like "for i < len(x)" survives rendering untouched.
// Use [Escape] to guarantee literal output of arbitrary strings, and "\<"
// to embed a literal '<' inside markup ("\\\\" renders a single backslash).
//
// # Concurrency
//
// All global state (enabled flag, default writer, default theme) is safe
// for concurrent use. Styles returned by [Theme.Style]/[Lookup] are
// independent copies; mutating them does not affect the theme.
package pigment

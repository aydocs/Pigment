# pigment

> Modern terminal styling for Go. Zero dependencies.

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE.md)
[![Dependencies](https://img.shields.io/badge/dependencies-0-brightgreen)](go.mod)

---

## Install

```sh
go get github.com/aydocs/Pigment
```

```go
import "github.com/aydocs/Pigment" // package pigment
```

---

## About

pigment is a Go library for coloring terminal output. It provides truecolor
(24-bit) support, gradients, rainbow effects, inline markup, OSC 8 hyperlinks,
semantic themes, and a full color system with all 148 CSS Color Level 4 named
colors plus 16 bright variants, HSL, color mixing, and the xterm 256-color
palette.

It has zero external dependencies and uses only the Go standard library.

---

## Quick Start

```go
package main

import (
    "github.com/aydocs/Pigment"
)

func main() {
    pigment.Red("error: something broke")      // prints a red line (auto newline)
    pigment.Green("success: all good")
    pigment.HiBlue("info: update available")

    pigment.New().Fg("navy").Bold().Println("styled text")
    pigment.Render("<bold><#ff5555>Error:</#ff5555> check logs</bold>")
}
```

---

## Features

### Colors

Define colors in multiple ways:

```go
c := pigment.RGB(255, 128, 64)         // RGB components (clamped to 0-255)
c := pigment.MustHex("#ff8800")        // hex string ("#f00" also works)
c := pigment.HSL(0, 1, 0.5)            // hue 0-360, sat/lightness clamped to 0-1
c := pigment.Palette(196)              // xterm 256-color index (clamped to 0-255)
c, _ := pigment.Named("rebeccapurple") // 148 CSS Color Level 4 keywords
c, _ := pigment.Named("hi-red")        // 16 bright colors
```

All 148 CSS Color Level 4 extended keywords are supported (including every
gray/grey spelling and `rebeccapurple`), plus 16 bright variants with the
`hi-` and `bright-` prefixes — 164 names in total. Blend two colors with
`c.Mix(d, t)`; `t` is clamped to `[0, 1]` and results are rounded, not
truncated.

Invalid input never panics: out-of-range numbers are clamped, and setters
that receive an unparseable value simply leave the style unchanged.

### Styles

Fluent builder for foreground, background, and text attributes:

```go
s := pigment.New().Fg("red").Bg("navy").Bold().Italic().Underline()
s.Println("styled text")
s.Fprint(w, "to writer")
str := s.Sprint("as string")

w2 := s.Clone().Fg("#50fa7b") // Clone is independent: mutating w2 never touches s
```

Attributes: Bold, Faint, Italic, Underline, DoubleUnderline, BlinkSlow,
BlinkRapid, Reverse, Conceal, CrossedOut, Overline, Framed, Encircled.
Shortcut setters exist for the common ones; combine any attribute with
`Add(a)` / remove with `Del(a)`.

Per-style override beats the global switch: `pigment.SetEnabled(false)` mutes
everything except styles created with `.Enable()`; `.Disable()` does the
reverse.

### Print Helpers

32 print functions and 32 string functions:

```go
pigment.Red("text")           // foreground, appends "\n"
pigment.HiRed("text")         // bright foreground, appends "\n"
pigment.BgRed("text")         // background, appends "\n"
pigment.HiBgRed("text")       // bright background, appends "\n"
pigment.RedString("text")     // returns string, no newline added
pigment.Sprintf("red", "hello %s", "world")
```

The print helpers always append a trailing newline (like `fmt.Println`);
the String/Sprintf forms never do.

### Func Generators

Create reusable functions from any style:

```go
red := pigment.New().Fg("red").PrintfFunc()
red("Warning")
red("Error: %s", err)

cyan := pigment.New().Fg("cyan").SprintFunc()
fmt.Println(cyan("highlight"))
```

Style methods: PrintFunc, PrintfFunc, PrintlnFunc, SprintFunc, SprintfFunc,
SprintlnFunc, FprintFunc, FprintfFunc, FprintlnFunc. Package-level versions:
PrintFunc, PrintfFunc, PrintlnFunc, SprintfFunc, SprintlnFunc.

### Gradients

```go
// two-color gradient across the text
text := pigment.Gradient("hot to cold", pigment.RGB(255,60,60), pigment.RGB(60,120,255))

// custom multi-stop gradient
text := pigment.Gradient("sunset", first, last,
    pigment.Stops(
        pigment.GradientStop{At: 0.0, Color: pigment.MustHex("#ff5555")},
        pigment.GradientStop{At: 0.6, Color: pigment.MustHex("#f1fa8c")},
        pigment.GradientStop{At: 1.0, Color: pigment.MustHex("#bd93f9")},
    ),
)
```

Stops are sorted by position; missing endpoints are anchored automatically.
Whitespace (`space`, tab, newline, carriage return) is preserved uncolored so
layouts stay intact.

### Rainbow

```go
pigment.Rainbow("hello!")
pigment.Rainbow("pastel",
    pigment.RainbowSaturation(0.5),
    pigment.RainbowLightness(0.7),
    pigment.RainbowStart(90),   // start hue
    pigment.RainbowEnd(270),    // end hue (any value, including 0)
)
```

### Markup

```go
pigment.Render("<red>text</red>")                    // named foreground
pigment.Render("<bold underline red>important</>")   // multi-token tags
pigment.Render("<fg=#00ff00 bg=navy italic>styled</>")
pigment.Render("<palette=196>256-color</>")
```

Keys, attribute names, and color names are case-insensitive (`<FG=RED>` works).
Tags nest, auto-close at the end of the string, and close-tag matching is
stack-based: `</red>` closes the nearest open `<red>`.

**Prose safety:** a `<...>` sequence containing any unrecognized token is NOT
consumed — it is emitted verbatim as plain text. Strings like
`"for i < len(xs) { ... }"` render unchanged instead of losing characters.
Embed literal markup characters with escapes: `\<` renders `<`, `\>` renders
`>`, `\\` renders `\`. Use `pigment.Escape(s)` to make any arbitrary string
safe inside markup.

### Hyperlinks

```go
pigment.Link("https://example.com", "click here")
```

OSC 8 protocol; works in kitty, iTerm2, WezTerm, GNOME Terminal and Windows
Terminal. Control bytes in the URL are percent-encoded so they cannot break
out of the escape sequence. Never pass fully untrusted URLs regardless.

### Themes

```go
pigment.Paint("error", "something broke")
pigment.Paint("success", "all good")

theme := pigment.NewTheme()
theme.Register("critical", pigment.New().Fg("red").Bold())
```

DefaultTheme ships 9 Dracula-inspired semantic styles under 12 aliases
(`error`/`err`, `warn`/`warning`, `success`/`ok`, `info`, `debug`, `muted`,
`link`, `title`, `highlight`). Themes are safe for concurrent use, and
`Register`/`Style`/`Lookup` hand out independent snapshots: mutating what you
get back never affects the theme or other holders.

### Color Detection

Detection runs once at startup and can be re-run with `DetectEnabled()`.
Precedence, highest first:

| # | Signal | Effect |
|---|--------|--------|
| 1 | `FORCE_COLOR` set | `"0"`/`"false"` forces off, anything else forces on |
| 2 | `CLICOLOR_FORCE` set ≠ `"0"` | forces on |
| 3 | `NO_COLOR` non-empty | forces off |
| 4 | `TERM=dumb` | forces off |
| 5 | `CLICOLOR=0` | forces off |
| 6 | none of the above | enabled only when stdout is a real TTY |

TTY probing uses a genuine terminal check (ioctl on Unix, `GetConsoleMode`
on Windows) — pipes and `/dev/null` correctly report not-a-terminal.

```go
pigment.SetEnabled(false)          // global off
pigment.IsEnabled()                // current state
pigment.DetectEnabled()            // re-run env detection, returns result
pigment.New().Enable().Println("forced on")
pigment.New().Disable().Println("forced off")
```

---

## Compared to fatih/color

| Feature | fatih/color | pigment |
|---------|-------------|---------|
| External deps | 2 | 0 |
| Named colors | 8 | 164 |
| HSL | no | yes |
| 256-color palette | no | yes |
| Color mixing | no | yes |
| Gradients | no | yes |
| Rainbow | no | yes |
| Markup | no | yes |
| Hyperlinks | no | yes |
| Themes | no | yes |
| Background helpers | yes | yes |
| Per-style enable override | no | yes |

---

## Output Control

```go
pigment.SetOutput(os.Stderr)   // default writer for Print* methods

s := pigment.New().Fg("red")
s.Fprint(os.Stderr, "text")
s.Fprintln(w, "line")
s.Fprintf(w, "format %s", "text")
```

## Concurrency

Global state (enabled flag, output writer, default theme) is atomic /
mutex-protected and safe for concurrent use. A single `Style` value mutated by
multiple goroutines is *not* synchronized — clone it per goroutine if needed.

---

## License

MIT

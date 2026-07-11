# pigment

> Modern terminal styling for Go. Zero dependencies.

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE.md)
[![Dependencies](https://img.shields.io/badge/dependencies-0-brightgreen)](go.mod)

---

## Install

```sh
go get pigment
```

---

## About

pigment is a Go library for coloring terminal output. It provides truecolor (24-bit) support, gradients, rainbow effects, inline markup, OSC 8 hyperlinks, semantic themes, and a full color system with 163 named colors, HSL, and 256-color palette support.

It has zero external dependencies and uses only the Go standard library.

---

## Quick Start

```go
package main

import "pigment"

func main() {
    pigment.Red("error: something broke")
    pigment.Green("success: all good")
    pigment.HiBlue("info: update available")

    pigment.New().Fg("red").Bold().Println("styled text")
    pigment.Render("<bold><#ff5555>Error:</#ff5555> check logs</bold>")
    pigment.Gradient("smooth", pigment.RGB(255,0,0), pigment.RGB(0,0,255))
}
```

---

## Features

### Colors

Define colors in multiple ways:

```go
c := pigment.RGB(255, 128, 64)         // RGB components
c := pigment.MustHex("#ff8800")          // hex string
c := pigment.HSL(0, 1, 0.5)             // hue, saturation, lightness
c := pigment.Palette(196)               // xterm 256-color index
c, _ := pigment.Named("rebeccapurple")  // 147 CSS named colors
c, _ := pigment.Named("hi-red")         // 16 bright colors
```

147 CSS named colors + 16 bright variants with `hi-` and `bright-` prefixes. 256-color palette: 0-7 system colors, 8-15 bright, 16-231 color cube, 232-255 grayscale. Mix two colors with `c.Mix(d, 0.5)`.

### Styles

Fluent builder for foreground, background, and text attributes:

```go
s := pigment.New().Fg("red").Bg("navy").Bold().Italic().Underline()
s.Println("styled text")
s.Fprint(w, "to writer")
str := s.Sprint("as string")
```

Attributes: Bold, Faint, Italic, Underline, DoubleUnderline, BlinkSlow, BlinkRapid, Reverse, Conceal, CrossedOut, Overline, Framed, Encircled. Clone styles without affecting the original.

### Print Helpers

32 print functions and 32 string functions:

```go
pigment.Red("text")                      // foreground
pigment.HiRed("text")                    // bright foreground
pigment.BgRed("text")                    // background
pigment.HiBgRed("text")                  // bright background
pigment.RedString("text")                // returns string
pigment.Sprintf("red", "hello %s", "world")
```

### Func Generators

Create reusable functions from any style:

```go
red := pigment.New().Fg("red").PrintfFunc()
red("Warning")
red("Error: %s", err)

cyan := pigment.New().Fg("cyan").SprintFunc()
fmt.Println(cyan("highlight"))
```

9 generators: PrintFunc, PrintfFunc, PrintlnFunc, SprintFunc, SprintfFunc, SprintlnFunc, FprintFunc, FprintfFunc, FprintlnFunc. Package-level versions also available.

### Gradients

```go
pigment.Gradient("smooth", red, blue)
pigment.Gradient("multi", red, green, blue)
pigment.Gradient("custom", red, blue, pigment.Stops(...))
```

Multi-stop and custom-positioned stops. Whitespace preserved uncolored.

### Rainbow

```go
pigment.Rainbow("hello!")
pigment.Rainbow("pastel",
    pigment.RainbowSaturation(0.5),
    pigment.RainbowLightness(0.7),
)
```

Configurable start/end hue, saturation, lightness.

### Markup

```go
pigment.Render("<red>text</red>")
pigment.Render("<bold underline red>important</>")
pigment.Render("<fg=#00ff00 bg=navy italic>styled</>")
```

Supports named colors, hex, fg/bg prefixes, palette, hex attributes, and style attributes. Tags nest and auto-close.

### Hyperlinks

```go
pigment.Link("https://example.com", "click here")
```

OSC 8 protocol. Works in kitty, iTerm2, WezTerm, GNOME Terminal, Windows Terminal.

### Themes

```go
pigment.Paint("error", "something broke")
pigment.Paint("success", "all good")

theme := pigment.NewTheme()
theme.Register("critical", pigment.New().Fg("red").Bold().Blink())
```

Default theme with 9 Dracula-inspired semantic styles. Thread-safe with sync.RWMutex.

### Color Detection

Respects NO_COLOR, TERM=dumb, FORCE_COLOR, CLICOLOR_FORCE, CLICOLOR. Auto-detects TTY. Global and per-style override.

```go
pigment.SetEnabled(false)
pigment.New().Enable().Println("force on")
pigment.New().Disable().Println("force off")
```

---

## Compared to fatih/color

pigment is a superset. Same API, zero dependencies, more features.

| Feature | fatih/color | pigment |
|---------|-------------|---------|
| External deps | 2 | 0 |
| Named colors | 8 | 163 |
| HSL | no | yes |
| 256-color palette | no | yes |
| Color mixing | no | yes |
| Gradients | no | yes |
| Rainbow | no | yes |
| Markup | no | yes |
| Hyperlinks | no | yes |
| Themes | no | yes |
| Background helpers | no | yes |

---

## Output Control

```go
pigment.SetOutput(os.Stderr)

s := pigment.New().Fg("red")
s.Fprint(os.Stderr, "text")
s.Fprintln(w, "line")
s.Fprintf(w, "format %s", "text")
```

---

## License

MIT

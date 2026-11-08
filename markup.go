package pigment

import (
	"strconv"
	"strings"
)

// Render parses a small markup language and returns the colorized result.
//
//	<red>text</red>          named foreground
//	<#ff8800>text</>         hex foreground (empty close tag allowed)
//	<bg=blue>text</bg=blue>  named background
//	<fg=#00ff00>text</>      hex foreground
//	<palette=196>text</>     xterm 256 index
//	<bold> <italic> <u>      attributes (b, i, u, s, blink, reverse...)
//	<red bold>x</>           multiple tokens in one tag
//
// Tags nest and auto-close at the end of the string. A close tag unwinds to
// its nearest matching open tag; "</>" pops the innermost open tag. Any
// "<...>" sequence whose tokens do not form a valid tag is emitted verbatim,
// so plain text like "if x < y then z > w" survives untouched. In the input,
// "\<" and "\\" render a literal "<" and "\"; any other backslash passes
// through unchanged.
func Render(s string) string {
	p := &markupParser{src: s}
	return p.parse()
}

// Escape returns s with backslashes and '<' escaped so that Render emits them
// as literal text instead of interpreting them as markup.
func Escape(s string) string {
	return strings.NewReplacer(`\`, `\\`, `<`, `\<`).Replace(s)
}

type markupFrame struct {
	style *Style
	name  string // lowercased first token of the open tag; matched by close tags
}

type markupParser struct {
	src   string
	pos   int
	out   strings.Builder
	stack []markupFrame
}

func (p *markupParser) parse() string {
	p.stack = []markupFrame{{style: New()}}
	for p.pos < len(p.src) {
		i := strings.IndexAny(p.src[p.pos:], `<\`)
		if i < 0 {
			p.text(p.src[p.pos:])
			break
		}
		start := p.pos + i
		if i > 0 {
			p.text(p.src[p.pos:start])
		}
		if p.src[start] == '\\' {
			p.escape(start)
		} else {
			p.tagAt(start)
		}
	}
	p.stack = p.stack[:1]
	return p.out.String()
}

// escape consumes the backslash at i: "\<", "\>" and "\\" emit the escaped
// character literally; any other backslash is emitted unchanged.
func (p *markupParser) escape(i int) {
	if i+1 < len(p.src) {
		switch next := p.src[i+1]; next {
		case '<', '>', '\\':
			p.text(p.src[i+1 : i+2])
			p.pos = i + 2
			return
		}
	}
	p.text(`\`)
	p.pos = i + 1
}

// tagAt handles the '<' at i; a body that is not a valid tag is written out
// verbatim, delimiters included.
func (p *markupParser) tagAt(i int) {
	endRel := strings.IndexByte(p.src[i:], '>')
	if endRel < 0 {
		p.text(p.src[i:])
		p.pos = len(p.src)
		return
	}
	end := i + endRel
	body := strings.TrimSpace(p.src[i+1 : end])
	p.pos = end + 1
	if strings.HasPrefix(body, "/") {
		name := strings.ToLower(strings.TrimSpace(body[1:]))
		if p.close(name) || name != "" {
			return
		}
		p.out.WriteString(p.src[i : end+1])
		return
	}
	if !isValidTag(body) {
		p.out.WriteString(p.src[i : end+1])
		return
	}
	p.open(body)
}

func (p *markupParser) text(s string) {
	if s == "" {
		return
	}
	p.out.WriteString(p.stack[len(p.stack)-1].style.Paint(s))
}

func (p *markupParser) open(body string) {
	cur := p.stack[len(p.stack)-1]
	next := cur.style.Clone()
	name := ""
	for _, tok := range strings.Fields(body) {
		applyToken(next, tok)
		if name == "" {
			name = strings.ToLower(tok)
		}
	}
	p.stack = append(p.stack, markupFrame{style: next, name: name})
}

// close unwinds to the nearest frame whose name matches; an empty name
// ("</>") pops the innermost frame. Reports whether a frame was popped. An
// unmatched named close is swallowed silently; an anonymous "</>" with
// nothing open falls back to verbatim text in the caller.
func (p *markupParser) close(name string) bool {
	if name == "" {
		if len(p.stack) > 1 {
			p.stack = p.stack[:len(p.stack)-1]
			return true
		}
		return false
	}
	for i := len(p.stack) - 1; i > 0; i-- {
		if p.stack[i].name == name {
			p.stack = p.stack[:i]
			return true
		}
	}
	return false
}

// attrAliases maps lowercase attribute aliases to their Style setter.
var attrAliases = map[string]func(*Style) *Style{
	"b":          (*Style).Bold,
	"bold":       (*Style).Bold,
	"faint":      (*Style).Faint,
	"dim":        (*Style).Faint,
	"i":          (*Style).Italic,
	"italic":     (*Style).Italic,
	"u":          (*Style).Underline,
	"underline":  (*Style).Underline,
	"dunderline": func(s *Style) *Style { return s.Add(DoubleUnderline) },
	"double":     func(s *Style) *Style { return s.Add(DoubleUnderline) },
	"blink":      (*Style).Blink,
	"reverse":    (*Style).Reverse,
	"invert":     (*Style).Reverse,
	"conceal":    func(s *Style) *Style { return s.Add(Conceal) },
	"hide":       func(s *Style) *Style { return s.Add(Conceal) },
	"s":          (*Style).Strike,
	"strike":     (*Style).Strike,
	"crossed":    (*Style).Strike,
	"overline":   (*Style).Overline,
	"frame":      func(s *Style) *Style { return s.Add(Framed) },
	"framed":     func(s *Style) *Style { return s.Add(Framed) },
	"circle":     func(s *Style) *Style { return s.Add(Encircled) },
	"encircled":  func(s *Style) *Style { return s.Add(Encircled) },
}

// isValidTag reports whether every whitespace-separated token of body is a
// valid markup token: an attribute alias, a prefixed fg/bg/hex/palette form,
// a bare hex color, or a known named color.
func isValidTag(body string) bool {
	if strings.TrimSpace(body) == "" {
		return false
	}
	for _, tok := range strings.Fields(body) {
		if !isValidToken(tok) {
			return false
		}
	}
	return true
}

func isValidToken(tok string) bool {
	lower := strings.ToLower(tok)
	if _, ok := attrAliases[lower]; ok {
		return true
	}
	for _, prefix := range []string{"fg=", "bg=", "hex=", "fghex=", "bghex="} {
		if strings.HasPrefix(lower, prefix) {
			value := strings.TrimSpace(tok[len(prefix):])
			if value == "" {
				return false
			}
			if prefix == "hex=" || prefix == "fghex=" || prefix == "bghex=" {
				_, err := Hex(value)
				return err == nil
			}
			if strings.HasPrefix(strings.ToLower(value), "#") {
				_, err := Hex(value)
				return err == nil
			}
			_, ok := Named(value)
			return ok
		}
	}
	for _, prefix := range []string{"palette=", "bgpalette="} {
		if strings.HasPrefix(lower, prefix) {
			n, err := strconv.Atoi(strings.TrimSpace(tok[len(prefix):]))
			return err == nil && n >= 0 && n <= 255
		}
	}
	if strings.HasPrefix(lower, "#") {
		_, err := Hex(lower)
		return err == nil
	}
	_, ok := Named(lower)
	return ok
}

// applyToken applies one already-validated markup token to s. Keys match
// case-insensitively while values keep their original case.
func applyToken(s *Style, tok string) {
	lower := strings.ToLower(tok)
	if fn, ok := attrAliases[lower]; ok {
		fn(s)
		return
	}
	switch {
	case strings.HasPrefix(lower, "fg="):
		s.Fg(tok[len("fg="):])
	case strings.HasPrefix(lower, "bg="):
		s.Bg(tok[len("bg="):])
	case strings.HasPrefix(lower, "hex="), strings.HasPrefix(lower, "fghex="):
		s.FgHex(tok[strings.IndexByte(tok, '=')+1:])
	case strings.HasPrefix(lower, "bghex="):
		s.BgHex(tok[len("bghex="):])
	case strings.HasPrefix(lower, "palette="):
		if n, err := strconv.Atoi(strings.TrimSpace(tok[len("palette="):])); err == nil && n >= 0 && n <= 255 {
			s.FgPalette(n)
		}
	case strings.HasPrefix(lower, "bgpalette="):
		if n, err := strconv.Atoi(strings.TrimSpace(tok[len("bgpalette="):])); err == nil && n >= 0 && n <= 255 {
			s.BgPalette(n)
		}
	default:
		s.Fg(tok)
	}
}

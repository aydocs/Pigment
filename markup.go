package pigment

import "strings"

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
// Tags nest and are auto-closed at the end of the string.
func Render(s string) string {
	p := &markupParser{src: s}
	return p.parse()
}

type markupParser struct {
	src   string
	pos   int
	out   strings.Builder
	stack []*Style
}

func (p *markupParser) parse() string {
	p.stack = []*Style{New()}
	for p.pos < len(p.src) {
		i := strings.IndexByte(p.src[p.pos:], '<')
		if i < 0 {
			p.text(p.src[p.pos:])
			p.pos = len(p.src)
			break
		}
		start := p.pos + i
		if i > 0 {
			p.text(p.src[p.pos:start])
		}
		end := strings.IndexByte(p.src[start:], '>')
		if end < 0 {
			p.text(p.src[start:])
			p.pos = len(p.src)
			break
		}
		tag := p.src[start+1 : start+end]
		p.pos = start + end + 1
		p.tag(tag)
	}
	for len(p.stack) > 1 {
		p.stack = p.stack[:len(p.stack)-1]
	}
	return p.out.String()
}

func (p *markupParser) text(s string) {
	if s == "" {
		return
	}
	cur := p.stack[len(p.stack)-1]
	p.out.WriteString(cur.Paint(s))
}

func (p *markupParser) tag(tag string) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return
	}
	if strings.HasPrefix(tag, "/") {
		if len(p.stack) > 1 {
			p.stack = p.stack[:len(p.stack)-1]
		}
		return
	}
	cur := p.stack[len(p.stack)-1]
	next := cur.Clone()
	for _, tok := range strings.Fields(tag) {
		applyToken(next, tok)
	}
	p.stack = append(p.stack, next)
}

func applyToken(s *Style, tok string) {
	lower := strings.ToLower(tok)
	switch {
	case lower == "b" || lower == "bold":
		s.Bold()
	case lower == "faint" || lower == "dim":
		s.Faint()
	case lower == "i" || lower == "italic":
		s.Italic()
	case lower == "u" || lower == "underline":
		s.Underline()
	case lower == "dunderline" || lower == "double":
		s.Add(DoubleUnderline)
	case lower == "blink":
		s.Blink()
	case lower == "reverse" || lower == "invert":
		s.Reverse()
	case lower == "conceal" || lower == "hide":
		s.Add(Conceal)
	case lower == "s" || lower == "strike" || lower == "crossed":
		s.Strike()
	case lower == "overline":
		s.Overline()
	case lower == "frame" || lower == "framed":
		s.Add(Framed)
	case lower == "circle" || lower == "encircled":
		s.Add(Encircled)
	case strings.HasPrefix(lower, "fg="):
		s.Fg(strings.TrimPrefix(tok, "fg="))
	case strings.HasPrefix(lower, "bg="):
		s.Bg(strings.TrimPrefix(tok, "bg="))
	case strings.HasPrefix(lower, "palette="):
		if n, ok := atoiSafe(strings.TrimPrefix(tok, "palette=")); ok {
			s.FgPalette(n)
		}
	case strings.HasPrefix(lower, "bgpalette="):
		if n, ok := atoiSafe(strings.TrimPrefix(tok, "bgpalette=")); ok {
			s.BgPalette(n)
		}
	case strings.HasPrefix(lower, "hex=") || strings.HasPrefix(lower, "fghex="):
		s.FgHex(strings.TrimPrefix(strings.TrimPrefix(lower, "fghex="), "hex="))
	case strings.HasPrefix(lower, "bghex="):
		s.BgHex(strings.TrimPrefix(lower, "bghex="))
	default:
		s.Fg(tok)
	}
}

func atoiSafe(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	neg := false
	if s[0] == '-' {
		neg = true
		s = s[1:]
	}
	v := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, false
		}
		v = v*10 + int(c-'0')
	}
	if neg {
		v = -v
	}
	return v, true
}

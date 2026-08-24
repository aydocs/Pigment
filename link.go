package pigment

import "strings"

// Link wraps text in an OSC 8 hyperlink so it renders as a clickable link in
// supporting terminals. Returns plain text when color is off or url is empty.
//
// Control bytes in url (0x00-0x1F, 0x7F) are percent-encoded before
// embedding, so a URL can neither terminate the escape sequence nor smuggle
// terminal commands (e.g. OSC 52 clipboard overwrites). Never pass untrusted
// URLs: encoding covers control bytes only, every other byte is embedded as-is.
func Link(url, text string) string {
	if !colorEnabled() || url == "" {
		return text
	}
	var b strings.Builder
	b.WriteString(osc8Pre)
	b.WriteString(escapeOSCURL(url))
	b.WriteString(osc8Sep)
	b.WriteString(text)
	b.WriteString(osc8Pre)
	b.WriteString(osc8Sep)
	return b.String()
}

// escapeOSCURL percent-encodes every control byte (0x00-0x1F, 0x7F) as %XX
// with uppercase hex, defusing terminal-sequence injection via OSC 8.
func escapeOSCURL(u string) string {
	if !strings.ContainsFunc(u, func(r rune) bool { return r < 0x20 || r == 0x7f }) {
		return u
	}
	var b strings.Builder
	b.Grow(len(u) + 6)
	for i := 0; i < len(u); i++ {
		if c := u[i]; c < 0x20 || c == 0x7f {
			const hexDigits = "0123456789ABCDEF"
			b.WriteByte('%')
			b.WriteByte(hexDigits[c>>4])
			b.WriteByte(hexDigits[c&0x0f])
			continue
		}
		b.WriteByte(u[i])
	}
	return b.String()
}

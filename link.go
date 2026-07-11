package pigment

import "strings"

// Link wraps text in an OSC 8 hyperlink so it renders as a clickable link in
// supporting terminals. Returns plain text when color is off or url is empty.
func Link(url, text string) string {
	if !Enabled || url == "" {
		return text
	}
	var b strings.Builder
	b.WriteString(osc8Pre)
	b.WriteString(url)
	b.WriteString(osc8Sep)
	b.WriteString(text)
	b.WriteString(osc8Pre)
	b.WriteString(osc8Sep)
	return b.String()
}

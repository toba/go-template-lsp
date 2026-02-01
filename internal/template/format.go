package template

import (
	"regexp"
	"strings"
)

// voidElements lists HTML void elements that don't have closing tags.
var voidElements = map[string]bool{
	"area":     true,
	"base":     true,
	"br":       true,
	"col":      true,
	"embed":    true,
	"hr":       true,
	"img":      true,
	"input":    true,
	"link":     true,
	"meta":     true,
	"param":    true,
	"source":   true,
	"track":    true,
	"wbr":      true,
	"command":  true,
	"keygen":   true,
	"menuitem": true,
}

var (
	// Matches HTML open/close tags, capturing position for ordering.
	htmlTagRe = regexp.MustCompile(`<(/?)(\w+)`)

	// openTagRe matches the start of an HTML opening tag (not closing).
	openTagRe = regexp.MustCompile(`^<([a-zA-Z][a-zA-Z0-9]*)(\s|$|>|/)`)

	// attrTokenRe matches individual attribute chunks: name="value", name='value',
	// name, or template actions {{...}}.
	attrTokenRe = regexp.MustCompile(
		`(\{\{.*?\}\}(?:[^\s<>]*\{\{.*?\}\})*[^\s<>]*|\S+?="[^"]*"|\S+?='[^']*'|[^\s>]+)`,
	)
)

// FormatOptions holds configuration for the formatter.
type FormatOptions struct {
	TabSize      int
	InsertSpaces bool
	PrintWidth   int    // 0 = disabled, default 120
	AttrWrapMode string // "overflow" or "all"
}

// tagEvent represents an indent change at a position in a line.
type tagEvent struct {
	pos   int
	delta int // +1 for open, -1 for close
}

// Format re-indents the source based on HTML tag and Go template control block nesting.
// Only leading whitespace is changed; content within lines is never modified.
// When PrintWidth > 0, opening HTML tags that exceed the width are wrapped.
func Format(source []byte, opts FormatOptions) []byte {
	if len(source) == 0 {
		return source
	}

	indent := "\t"
	if opts.InsertSpaces {
		indent = strings.Repeat(" ", opts.TabSize)
	}

	lines := strings.Split(string(source), "\n")

	// Pre-pass: join multi-line tags into single logical lines
	lines = joinMultiLineTags(lines)

	level := 0
	var result []string

	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")

		if trimmed == "" {
			result = append(result, "")
			continue
		}

		before, after := computeLineDeltas(trimmed)

		// Apply before delta (dedent for closing tags/end appearing on this line)
		level += before
		if level < 0 {
			level = 0
		}

		indented := strings.Repeat(indent, level) + trimmed

		// Attribute wrapping
		if opts.PrintWidth > 0 && lineWidth(indented, opts.TabSize) > opts.PrintWidth {
			if wrapped, ok := wrapAttributes(trimmed, indent, level, opts); ok {
				result = append(result, wrapped...)
			} else {
				result = append(result, indented)
			}
		} else {
			result = append(result, indented)
		}

		// Apply after delta (indent increase for next lines)
		level += after
		if level < 0 {
			level = 0
		}
	}

	return []byte(strings.Join(result, "\n"))
}

// joinMultiLineTags joins opening HTML tags that span multiple lines into single lines.
// A tag spans lines when a line starts with <tagname but does not contain >.
func joinMultiLineTags(lines []string) []string {
	var result []string
	i := 0
	for i < len(lines) {
		trimmed := strings.TrimLeft(lines[i], " \t")
		if isOpenTagStart(trimmed) && !strings.Contains(trimmed, ">") {
			// Collect continuation lines until we find >
			joined := trimmed
			j := i + 1
			var joinedSb127 strings.Builder
			for j < len(lines) {
				cont := strings.TrimSpace(lines[j])
				if cont == "" {
					j++
					continue
				}
				joinedSb127.WriteString(" " + cont)
				if strings.Contains(cont, ">") {
					j++
					break
				}
				j++
			}
			joined += joinedSb127.String()
			result = append(result, joined)
			i = j
		} else {
			result = append(result, lines[i])
			i++
		}
	}
	return result
}

// isOpenTagStart checks if a trimmed line starts with an opening HTML tag (not closing).
func isOpenTagStart(trimmed string) bool {
	return openTagRe.MatchString(trimmed)
}

// lineWidth computes the display width of a line, counting tabs as tabSize characters.
func lineWidth(line string, tabSize int) int {
	width := 0
	for _, ch := range line {
		if ch == '\t' {
			width += tabSize
		} else {
			width++
		}
	}
	return width
}

// wrapAttributes wraps attributes of an HTML opening tag across multiple lines.
// Returns the wrapped lines and true if wrapping was performed, or nil and false
// if the line is not a wrappable opening tag.
func wrapAttributes(
	trimmed, indent string,
	level int,
	opts FormatOptions,
) ([]string, bool) {
	// Must be an opening tag
	m := openTagRe.FindStringSubmatch(trimmed)
	if m == nil {
		return nil, false
	}

	tagName := m[1]
	prefix := "<" + tagName

	// Find where attributes start (after tag name)
	attrStart := len(prefix)
	rest := trimmed[attrStart:]

	// Find the closing > or /> and any content after it
	closeIdx, closer, afterClose := findTagClose(rest)
	if closeIdx < 0 {
		return nil, false
	}

	attrPortion := rest[:closeIdx]

	// Tokenize attributes
	attrs := attrTokenRe.FindAllString(attrPortion, -1)
	if len(attrs) == 0 {
		return nil, false
	}

	baseIndent := strings.Repeat(indent, level)
	contIndent := strings.Repeat(indent, level+1)

	mode := opts.AttrWrapMode
	if mode == "" {
		mode = "overflow"
	}

	var lines []string

	if mode == "all" {
		// First line: just the tag name
		lines = append(lines, baseIndent+prefix)
		// Each attribute on its own line
		for i, attr := range attrs {
			line := contIndent + attr
			if i == len(attrs)-1 {
				line += closer + afterClose
			}
			lines = append(lines, line)
		}
	} else {
		// overflow mode: keep attrs on first line until one would overflow
		firstLine := baseIndent + prefix
		overflowed := false
		var remaining []string

		for i, attr := range attrs {
			candidate := firstLine + " " + attr
			if i == len(attrs)-1 {
				candidate += closer + afterClose
			}
			if !overflowed && lineWidth(candidate, opts.TabSize) <= opts.PrintWidth {
				firstLine = firstLine + " " + attr
			} else {
				overflowed = true
				remaining = append(remaining, attr)
			}
		}

		if !overflowed {
			// Everything fit — attach closer
			lines = append(lines, firstLine+closer+afterClose)
		} else {
			// First line gets closer only if no remaining attrs
			lines = append(lines, firstLine)
			for i, attr := range remaining {
				line := contIndent + attr
				if i == len(remaining)-1 {
					line += closer + afterClose
				}
				lines = append(lines, line)
			}
		}
	}

	return lines, true
}

// findTagClose finds the position and type of closing bracket in the attribute portion.
// Returns (index, closer, afterClose) where closer is ">" or " />" and afterClose
// is any content after the closing bracket.
func findTagClose(s string) (int, string, string) {
	depth := 0
	inDoubleQuote := false
	inSingleQuote := false

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			continue
		}
		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			continue
		}
		if inDoubleQuote || inSingleQuote {
			continue
		}

		if ch == '{' && i+1 < len(s) && s[i+1] == '{' {
			depth++
			i++ // skip second {
			continue
		}
		if ch == '}' && i+1 < len(s) && s[i+1] == '}' {
			depth--
			i++ // skip second }
			continue
		}

		if depth > 0 {
			continue
		}

		if ch == '>' {
			closer := ">"
			idx := i
			if i > 0 && s[i-1] == '/' {
				closer = " />"
				// Walk back to find start of the / excluding whitespace before it
				idx = i - 1
				for idx > 0 && s[idx-1] == ' ' {
					idx--
				}
			}
			afterClose := s[i+1:]
			return idx, closer, afterClose
		}
	}

	return -1, "", ""
}

// computeLineDeltas computes (beforeDelta, afterDelta) for a trimmed line by
// collecting all HTML tag events in source order.
//
// beforeDelta: the minimum running delta (always <= 0), applied before writing the line.
// afterDelta: the net delta minus beforeDelta, applied after writing the line.
func computeLineDeltas(line string) (before, after int) {
	var events []tagEvent

	// Collect HTML tag events
	matches := htmlTagRe.FindAllStringSubmatchIndex(line, -1)
	for _, match := range matches {
		isClosing := line[match[2]:match[3]] == "/"
		tagName := strings.ToLower(line[match[4]:match[5]])

		if voidElements[tagName] {
			continue
		}

		if !isClosing && isSelfClosing(match[5], line) {
			continue
		}

		delta := 1
		if isClosing {
			delta = -1
		}
		events = append(events, tagEvent{pos: match[0], delta: delta})
	}

	// Sort events by position (insertion sort since typically few events)
	for i := 1; i < len(events); i++ {
		for j := i; j > 0 && events[j].pos < events[j-1].pos; j-- {
			events[j], events[j-1] = events[j-1], events[j]
		}
	}

	// Compute running delta; track minimum for beforeDelta
	running := 0
	minRunning := 0
	for _, e := range events {
		running += e.delta
		if running < minRunning {
			minRunning = running
		}
	}

	// beforeDelta is the minimum (dedent this line), afterDelta is the remainder
	before = minRunning          // <= 0
	after = running - minRunning // >= 0

	return before, after
}

// isSelfClosing checks if a tag starting at tagEnd position is self-closing (ends with />).
func isSelfClosing(tagEnd int, line string) bool {
	for i := tagEnd; i < len(line); i++ {
		if line[i] == '>' {
			return i > 0 && line[i-1] == '/'
		}
		if line[i] == '<' {
			return false
		}
	}
	return false
}

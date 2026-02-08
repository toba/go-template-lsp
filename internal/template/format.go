package template

import (
	"cmp"
	"regexp"
	"slices"
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
		`(\{\{.*?\}\}(?:[^\s<>]*\{\{.*?\}\})*[^\s<>]*|\S+?="(?:\{\{.*?\}\}|[^"])*"|\S+?='(?:\{\{.*?\}\}|[^'])*'|[^\s>]+)`,
	)

	// tmplBlockRe matches template block-opening keywords.
	tmplBlockRe = regexp.MustCompile(`\{\{-?\s*(if|range|with|define|block)\b`)
	// tmplElseRe matches template else/else if/else with keywords.
	tmplElseRe = regexp.MustCompile(`\{\{-?\s*else\b`)
	// tmplEndRe matches template end keywords.
	tmplEndRe = regexp.MustCompile(`\{\{-?\s*end\b`)
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

// countStructuralKeywords counts template block-open, else, and end keywords on
// a line, cancelling any fully-inline blocks (matched open+end pairs) along with
// all elses that fall between them. Only unmatched (structural) keywords are counted.
func countStructuralKeywords(line string) (opens, elses, ends int) {
	type kw struct {
		pos  int
		kind byte // 'o' open, 'e' else, 'n' end
	}

	openMatches := tmplBlockRe.FindAllStringIndex(line, -1)
	elseMatches := tmplElseRe.FindAllStringIndex(line, -1)
	endMatches := tmplEndRe.FindAllStringIndex(line, -1)
	all := make([]kw, 0, len(openMatches)+len(elseMatches)+len(endMatches))
	for _, m := range openMatches {
		all = append(all, kw{m[0], 'o'})
	}
	for _, m := range elseMatches {
		all = append(all, kw{m[0], 'e'})
	}
	for _, m := range endMatches {
		all = append(all, kw{m[0], 'n'})
	}

	// Sort by position.
	slices.SortFunc(all, func(a, b kw) int {
		return cmp.Compare(a.pos, b.pos)
	})

	// Mark inline pairs: walk left-to-right with a stack of open positions.
	inline := make([]bool, len(all))
	var stack []int // indices into all
	for i, k := range all {
		switch k.kind {
		case 'o':
			stack = append(stack, i)
		case 'n':
			if len(stack) > 0 {
				openIdx := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				// Mark the open and end as inline.
				inline[openIdx] = true
				inline[i] = true
				// Mark all elses between them as inline too.
				for j := openIdx + 1; j < i; j++ {
					if all[j].kind == 'e' {
						inline[j] = true
					}
				}
			}
		}
	}

	// Count only structural (non-inline) keywords.
	for i, k := range all {
		if inline[i] {
			continue
		}
		switch k.kind {
		case 'o':
			opens++
		case 'e':
			elses++
		case 'n':
			ends++
		}
	}
	return
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

	// blockStack saves the HTML indent level at the start of each template
	// control block (if/range/with/define/block) so that {{else}} and {{end}}
	// can restore it, preventing drift from mutually exclusive branches.
	var blockStack []int

	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")

		if trimmed == "" {
			result = append(result, "")
			continue
		}

		// Count structural template block keywords (excluding inline pairs).
		opens, elses, ends := countStructuralKeywords(trimmed)

		// At {{else}}: restore level to block entry for indenting, then reset
		// for the new branch. At {{end}}: restore for indenting, but preserve
		// the last branch's net HTML delta for subsequent lines.
		var postLevel int
		hasRestore := false
		if elses > 0 && len(blockStack) > 0 {
			saved := blockStack[len(blockStack)-1]
			level = saved     // indent {{else}} at block entry level
			postLevel = saved // new branch starts fresh
			hasRestore = true
		}
		if ends > 0 && len(blockStack) > 0 {
			saved := blockStack[len(blockStack)-1]
			branchDelta := level - saved
			postLevel = saved + branchDelta // keep last branch's HTML delta
			level = saved                   // indent {{end}} at block entry level
			blockStack = blockStack[:len(blockStack)-1]
			hasRestore = true
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

		// After indenting: if we restored for a template keyword, set the
		// post-indent level; otherwise apply normal HTML after delta.
		if hasRestore {
			level = postLevel
		} else {
			level += after
			if level < 0 {
				level = 0
			}
		}

		// Push level at block open — after computing HTML deltas for this line.
		for range opens {
			blockStack = append(blockStack, level)
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
		// Check if first attr on tag line aligns with continuation indent
		firstOnTagLine := len(attrs) > 1 &&
			lineWidth(
				baseIndent+prefix+" ",
				opts.TabSize,
			) == lineWidth(
				contIndent,
				opts.TabSize,
			)
		startIdx := 0
		if firstOnTagLine {
			lines = append(lines, baseIndent+prefix+" "+attrs[0])
			startIdx = 1
		} else {
			lines = append(lines, baseIndent+prefix)
		}
		for i := startIdx; i < len(attrs); i++ {
			line := contIndent + attrs[i]
			if i == len(attrs)-1 {
				line += closer + afterClose
			}
			lines = append(lines, line)
		}
	} else {
		// overflow mode: pack as many attrs per line as fit within printWidth
		currentLine := baseIndent + prefix
		lineHasAttr := false

		for i, attr := range attrs {
			sep := " "
			if !lineHasAttr && currentLine == contIndent {
				sep = "" // no extra space at start of continuation line
			}
			candidate := currentLine + sep + attr
			if i == len(attrs)-1 {
				candidate += closer + afterClose
			}
			if lineWidth(candidate, opts.TabSize) <= opts.PrintWidth || !lineHasAttr {
				// Fits, or must take at least one attr per line
				currentLine += sep + attr
				lineHasAttr = true
			} else {
				// Start a new continuation line
				lines = append(lines, currentLine)
				currentLine = contIndent + attr
				lineHasAttr = true
			}
		}

		// Emit the last line with the closer
		lines = append(lines, currentLine+closer+afterClose)
	}

	// With only one attribute, wrapping just moves it to the next line without
	// reducing the overall width — skip wrapping in that case.
	if len(attrs) <= 1 {
		return nil, false
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

	// Sort events by position.
	slices.SortFunc(events, func(a, b tagEvent) int {
		return cmp.Compare(a.pos, b.pos)
	})

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

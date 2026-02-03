package template

import (
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     FormatOptions
		expected string
	}{
		{
			name:  "empty input",
			input: "",
			opts:  FormatOptions{TabSize: 2, InsertSpaces: true},
		},
		{
			name:     "single line no change",
			input:    "<div>hello</div>",
			opts:     FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: "<div>hello</div>",
		},
		{
			name:     "already formatted",
			input:    "<div>\n  <p>hello</p>\n</div>",
			opts:     FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: "<div>\n  <p>hello</p>\n</div>",
		},
		{
			name: "nested HTML divs",
			input: `<div>
<div>
<p>hello</p>
</div>
</div>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<div>
  <div>
    <p>hello</p>
  </div>
</div>`,
		},
		{
			name: "void elements no indent change",
			input: `<div>
<br>
<hr>
<img src="x">
<input type="text">
<p>text</p>
</div>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<div>
  <br>
  <hr>
  <img src="x">
  <input type="text">
  <p>text</p>
</div>`,
		},
		{
			name: "self-closing tags no indent change",
			input: `<div>
<br />
<Component />
</div>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<div>
  <br />
  <Component />
</div>`,
		},
		{
			name: "template if/else/end",
			input: `{{ if .Condition }}
<p>true</p>
{{ else }}
<p>false</p>
{{ end }}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{ if .Condition }}
<p>true</p>
{{ else }}
<p>false</p>
{{ end }}`,
		},
		{
			name: "template range",
			input: `{{ range .Items }}
<li>{{ . }}</li>
{{ end }}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{ range .Items }}
<li>{{ . }}</li>
{{ end }}`,
		},
		{
			name: "template with",
			input: `{{ with .Data }}
<span>{{ . }}</span>
{{ end }}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{ with .Data }}
<span>{{ . }}</span>
{{ end }}`,
		},
		{
			name: "template define",
			input: `{{ define "header" }}
<h1>Title</h1>
{{ end }}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{ define "header" }}
<h1>Title</h1>
{{ end }}`,
		},
		{
			name: "template block",
			input: `{{ block "content" . }}
<p>default</p>
{{ end }}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{ block "content" . }}
<p>default</p>
{{ end }}`,
		},
		{
			name: "combined HTML and template nesting",
			input: `<div>
{{ range .Items }}
<ul>
<li>{{ .Name }}</li>
</ul>
{{ end }}
</div>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<div>
  {{ range .Items }}
  <ul>
    <li>{{ .Name }}</li>
  </ul>
  {{ end }}
</div>`,
		},
		{
			name: "tabs instead of spaces",
			input: `<div>
<p>hello</p>
</div>`,
			opts:     FormatOptions{TabSize: 1, InsertSpaces: false},
			expected: "<div>\n\t<p>hello</p>\n</div>",
		},
		{
			name: "tab size 4",
			input: `<div>
<p>hello</p>
</div>`,
			opts:     FormatOptions{TabSize: 4, InsertSpaces: true},
			expected: "<div>\n    <p>hello</p>\n</div>",
		},
		{
			name:     "content within line preserved",
			input:    "  <p>  some   spaced   content  </p>",
			opts:     FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: "<p>  some   spaced   content  </p>",
		},
		{
			name: "inline if/end nets to zero",
			input: `<div>
{{ if .A }}{{ end }}
</div>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<div>
  {{ if .A }}{{ end }}
</div>`,
		},
		{
			name: "open and close tag on same line nets to zero",
			input: `<div>
<p>text</p>
</div>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<div>
  <p>text</p>
</div>`,
		},
		{
			name: "trimmed template actions",
			input: `{{- if .A }}
<p>yes</p>
{{- end }}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{- if .A }}
<p>yes</p>
{{- end }}`,
		},
		{
			name: "else if produces dedent then indent",
			input: `{{ if .A }}
<p>a</p>
{{ else if .B }}
<p>b</p>
{{ end }}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{ if .A }}
<p>a</p>
{{ else if .B }}
<p>b</p>
{{ end }}`,
		},
		{
			name: "else with produces dedent then indent",
			input: `{{ if .A }}
<p>a</p>
{{ else with .B }}
<p>b</p>
{{ end }}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{ if .A }}
<p>a</p>
{{ else with .B }}
<p>b</p>
{{ end }}`,
		},
		{
			name: "deeply nested mixed",
			input: `<html>
<body>
{{ range .Pages }}
<div>
{{ if .Show }}
<p>{{ .Text }}</p>
{{ end }}
</div>
{{ end }}
</body>
</html>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<html>
  <body>
    {{ range .Pages }}
    <div>
      {{ if .Show }}
      <p>{{ .Text }}</p>
      {{ end }}
    </div>
    {{ end }}
  </body>
</html>`,
		},
		{
			name: "define with HTML list and range",
			input: `{{ define "list" }}
<ul>
{{ range .Items }}
<li>{{ . }}</li>
{{ end }}
</ul>
{{ end }}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{ define "list" }}
<ul>
  {{ range .Items }}
  <li>{{ . }}</li>
  {{ end }}
</ul>
{{ end }}`,
		},
		{
			name:     "blank lines preserved",
			input:    "<div>\n\n<p>text</p>\n\n</div>",
			opts:     FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: "<div>\n\n  <p>text</p>\n\n</div>",
		},

		{
			name:  "tabs in input replaced with spaces",
			input: "{{define \"sidebar\" -}}\n<aside id=\"sidebar\">\n\t<header>\n\t\t<a href=\"/\" aria-label=\"Pacer home\">\n\t<img src=\"/web/images/logo.svg\" alt=\"\"/>\n\t\t\t<span></span>{{/* \"pacer\" SVG assigned in CSS */}}\n\t\t</a>\n\t</header>\n</aside>\n{{end}}",
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   100,
				AttrWrapMode: "overflow",
			},
			expected: "{{define \"sidebar\" -}}\n<aside id=\"sidebar\">\n   <header>\n      <a href=\"/\" aria-label=\"Pacer home\">\n         <img src=\"/web/images/logo.svg\" alt=\"\"/>\n         <span></span>{{/* \"pacer\" SVG assigned in CSS */}}\n      </a>\n   </header>\n</aside>\n{{end}}",
		},

		// Attribute wrapping tests
		{
			name:  "tag under printWidth no wrapping",
			input: `<button type="button" class="btn">`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   80,
				AttrWrapMode: "overflow",
			},
			expected: `<button type="button" class="btn">`,
		},
		{
			name:  "overflow mode wraps only overflowing attrs",
			input: `<button type="button" class="map-zoom-btn" title="Expand">`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   40,
				AttrWrapMode: "overflow",
			},
			expected: `<button type="button"
   class="map-zoom-btn" title="Expand">`,
		},
		{
			name:  "all mode wraps every attribute",
			input: `<button type="button" class="map-zoom-btn" title="Expand">`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   40,
				AttrWrapMode: "all",
			},
			expected: `<button
   type="button"
   class="map-zoom-btn"
   title="Expand">`,
		},
		{
			name:  "all mode keeps first attr on tag line when aligned",
			input: `<a href="/path" class="link" id="nav">`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   30,
				AttrWrapMode: "all",
			},
			expected: `<a href="/path"
   class="link"
   id="nav">`,
		},
		{
			name:  "all mode tag line alignment no match keeps tag alone",
			input: `<div class="container" id="main" role="presentation">`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   40,
				AttrWrapMode: "all",
			},
			expected: `<div
   class="container"
   id="main"
   role="presentation">`,
		},
		{
			name: "all mode nested tag keeps first attr on tag line when aligned",
			input: `<div>
<a href="/path" class="link" id="nav">
</a>
</div>`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   30,
				AttrWrapMode: "all",
			},
			expected: `<div>
   <a href="/path"
      class="link"
      id="nav">
   </a>
</div>`,
		},
		{
			name: "all mode nested tag gets extra indent",
			input: `<div>
<button type="button" class="map-zoom-btn" title="Expand">
</button>
</div>`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   50,
				AttrWrapMode: "all",
			},
			expected: `<div>
   <button
      type="button"
      class="map-zoom-btn"
      title="Expand">
   </button>
</div>`,
		},
		{
			name:  "self-closing tag wrapping preserves />",
			input: `<img src="logo.png" alt="Company Logo" width="200" height="100" />`,
			opts: FormatOptions{
				TabSize:      2,
				InsertSpaces: true,
				PrintWidth:   40,
				AttrWrapMode: "overflow",
			},
			expected: `<img src="logo.png" alt="Company Logo"
  width="200" height="100" />`,
		},
		{
			name:  "void element with attributes wrapping",
			input: `<input type="text" name="username" placeholder="Enter name" required>`,
			opts: FormatOptions{
				TabSize:      2,
				InsertSpaces: true,
				PrintWidth:   40,
				AttrWrapMode: "overflow",
			},
			expected: `<input type="text" name="username"
  placeholder="Enter name" required>`,
		},
		{
			name:  "template action as attribute",
			input: `<a href="/page" {{if .Active}}class="active"{{end}} title="Link">text</a>`,
			opts: FormatOptions{
				TabSize:      2,
				InsertSpaces: true,
				PrintWidth:   40,
				AttrWrapMode: "overflow",
			},
			expected: `<a href="/page"
  {{if .Active}}class="active"{{end}}
  title="Link">text</a>`,
		},
		{
			name: "already-wrapped multi-line tag gets re-wrapped",
			input: `<button
   type="button"
   class="btn"
   title="X">`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   80,
				AttrWrapMode: "overflow",
			},
			expected: `<button type="button" class="btn" title="X">`,
		},
		{
			name: "already-wrapped tag re-wrapped to different width",
			input: `<button
type="button"
class="map-zoom-btn"
title="Expand">`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   40,
				AttrWrapMode: "all",
			},
			expected: `<button
   type="button"
   class="map-zoom-btn"
   title="Expand">`,
		},
		{
			name:  "printWidth 0 disables wrapping",
			input: `<button type="button" class="map-zoom-btn" title="Expand">`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   0,
				AttrWrapMode: "overflow",
			},
			expected: `<button type="button" class="map-zoom-btn" title="Expand">`,
		},
		{
			name:  "wrapping with tab characters",
			input: `<button type="button" class="map-zoom-btn" title="Expand">`,
			// tab counts as tabSize chars for width measurement
			opts: FormatOptions{
				TabSize:      4,
				InsertSpaces: false,
				PrintWidth:   40,
				AttrWrapMode: "all",
			},
			expected: "<button\n" +
				"\ttype=\"button\"\n" +
				"\tclass=\"map-zoom-btn\"\n" +
				"\ttitle=\"Expand\">",
		},
		{
			name:  "overflow mode first attr stays on tag line",
			input: `<div class="very-long-class-name" id="my-id" data-value="something">`,
			opts: FormatOptions{
				TabSize:      2,
				InsertSpaces: true,
				PrintWidth:   40,
				AttrWrapMode: "overflow",
			},
			expected: `<div class="very-long-class-name"
  id="my-id" data-value="something">`,
		},
		{
			name:  "closing > on tag with content after attrs",
			input: `<a href="/long/path/to/page" class="nav-link active" title="Navigation">Home</a>`,
			opts: FormatOptions{
				TabSize:      2,
				InsertSpaces: true,
				PrintWidth:   50,
				AttrWrapMode: "overflow",
			},
			expected: `<a href="/long/path/to/page"
  class="nav-link active"
  title="Navigation">Home</a>`,
		},
		{
			name:  "closing tag is never wrapped",
			input: `</div>`,
			opts: FormatOptions{
				TabSize:      2,
				InsertSpaces: true,
				PrintWidth:   3,
				AttrWrapMode: "all",
			},
			expected: `</div>`,
		},
		{
			name:  "non-tag line is never wrapped",
			input: `Some very long text content that exceeds the print width significantly`,
			opts: FormatOptions{
				TabSize:      2,
				InsertSpaces: true,
				PrintWidth:   20,
				AttrWrapMode: "all",
			},
			expected: `Some very long text content that exceeds the print width significantly`,
		},
		{
			name:  "tag with single attribute no wrap needed",
			input: `<div class="short">`,
			opts: FormatOptions{
				TabSize:      2,
				InsertSpaces: true,
				PrintWidth:   80,
				AttrWrapMode: "all",
			},
			expected: `<div class="short">`,
		},
		{
			name:  "single long attr overflow mode skips wrapping",
			input: `<path d="M12 6V2H8M2 12h2m16 0h2m-2 4a2 2 0 0 1-2 2H8.828a2 2 0 0 0-1.414.586l-2.202 2.202A.71.71 0 0 1 4 20.286V8a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2z" />`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   100,
				AttrWrapMode: "overflow",
			},
			expected: `<path d="M12 6V2H8M2 12h2m16 0h2m-2 4a2 2 0 0 1-2 2H8.828a2 2 0 0 0-1.414.586l-2.202 2.202A.71.71 0 0 1 4 20.286V8a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2z" />`,
		},
		{
			name:  "single long attr all mode skips wrapping",
			input: `<path d="M12 6V2H8M2 12h2m16 0h2m-2 4a2 2 0 0 1-2 2H8.828a2 2 0 0 0-1.414.586l-2.202 2.202A.71.71 0 0 1 4 20.286V8a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2z" />`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   100,
				AttrWrapMode: "all",
			},
			expected: `<path d="M12 6V2H8M2 12h2m16 0h2m-2 4a2 2 0 0 1-2 2H8.828a2 2 0 0 0-1.414.586l-2.202 2.202A.71.71 0 0 1 4 20.286V8a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2z" />`,
		},
		{
			name:  "mix of short and long attrs still wraps",
			input: `<path class="icon" d="M12 6V2H8M2 12h2m16 0h2m-2 4a2 2 0 0 1-2 2H8.828a2 2 0 0 0-1.414.586l-2.202 2.202A.71.71 0 0 1 4 20.286V8a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2z" />`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   100,
				AttrWrapMode: "all",
			},
			expected: `<path
   class="icon"
   d="M12 6V2H8M2 12h2m16 0h2m-2 4a2 2 0 0 1-2 2H8.828a2 2 0 0 0-1.414.586l-2.202 2.202A.71.71 0 0 1 4 20.286V8a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2z" />`,
		},
		{
			name:  "overflow mode packs attrs on continuation lines",
			input: `<section class="portfolio" hx-boost:inherited="true" hx-target:inherited="main" hx-swap:inherited="innerHTML" hx-push-url:inherited="true">`,
			opts: FormatOptions{
				TabSize:      3,
				InsertSpaces: true,
				PrintWidth:   80,
				AttrWrapMode: "overflow",
			},
			expected: `<section class="portfolio" hx-boost:inherited="true" hx-target:inherited="main"
   hx-swap:inherited="innerHTML" hx-push-url:inherited="true">`,
		},
		{
			name: "if/else with HTML tags resets level at else",
			input: `<table>
<tr>
</tr><tr>
</tr>
{{if gt .Balance 0}}
<tr class="positive">
{{else}}
<tr class="negative">
{{end}}
</tr>
</table>`,
			opts: FormatOptions{TabSize: 3, InsertSpaces: true},
			expected: `<table>
   <tr>
   </tr><tr>
   </tr>
   {{if gt .Balance 0}}
   <tr class="positive">
   {{else}}
   <tr class="negative">
   {{end}}
   </tr>
</table>`,
		},
		{
			name: "nested template blocks with HTML",
			input: `<div>
{{range .Items}}
<ul>
{{if .Show}}
<li>{{.Name}}</li>
{{else}}
<li>hidden</li>
{{end}}
</ul>
{{end}}
</div>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<div>
  {{range .Items}}
  <ul>
    {{if .Show}}
    <li>{{.Name}}</li>
    {{else}}
    <li>hidden</li>
    {{end}}
  </ul>
  {{end}}
</div>`,
		},
		{
			name: "inline if/end no stack change",
			input: `<div>
{{if .A}}{{end}}
<p>text</p>
</div>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<div>
  {{if .A}}{{end}}
  <p>text</p>
</div>`,
		},
		{
			name: "range with no else restores level at end",
			input: `<table>
{{range .Rows}}
<tr>
<td>{{.}}</td>
</tr>
{{end}}
</table>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<table>
  {{range .Rows}}
  <tr>
    <td>{{.}}</td>
  </tr>
  {{end}}
</table>`,
		},
		{
			name: "else if resets level like else",
			input: `<div>
{{if .A}}
<span>a</span>
{{else if .B}}
<span>b</span>
{{else}}
<span>c</span>
{{end}}
</div>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<div>
  {{if .A}}
  <span>a</span>
  {{else if .B}}
  <span>b</span>
  {{else}}
  <span>c</span>
  {{end}}
</div>`,
		},
		{
			name: "inline if/end with HTML inside outer range",
			input: `{{range .Items}}
{{if .StreetAddress -}}<p>{{.StreetAddress}}</p>{{end}}
<time>{{.Date}}</time>
{{end}}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{range .Items}}
{{if .StreetAddress -}}<p>{{.StreetAddress}}</p>{{end}}
<time>{{.Date}}</time>
{{end}}`,
		},
		{
			name: "inline if/end in HTML attribute inside nested structure",
			input: `<section>
<div class="modal-body{{if and .Lat .Lng}} mapped{{end}}">
<p>content</p>
</div>
</section>`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<section>
  <div class="modal-body{{if and .Lat .Lng}} mapped{{end}}">
    <p>content</p>
  </div>
</section>`,
		},
		{
			name: "real template: define > dialog > section > inline if > time",
			input: `{{define "event-dialog"}}
<dialog>
<section>
{{if .StreetAddress -}}<p>{{.StreetAddress}}</p>{{end}}
<time>{{.Date}}</time>
</section>
</dialog>
{{end}}`,
			opts: FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `{{define "event-dialog"}}
<dialog>
  <section>
    {{if .StreetAddress -}}<p>{{.StreetAddress}}</p>{{end}}
    <time>{{.Date}}</time>
  </section>
</dialog>
{{end}}`,
		},
		{
			name:  "default opts no wrapping",
			input: `<button type="button" class="map-zoom-btn" title="Expand">`,
			// PrintWidth defaults to 0, so no wrapping
			opts:     FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: `<button type="button" class="map-zoom-btn" title="Expand">`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(Format([]byte(tt.input), tt.opts))
			if got != tt.expected {
				t.Errorf(
					"Format() mismatch\n--- got ---\n%s\n--- expected ---\n%s",
					got,
					tt.expected,
				)
			}
		})
	}
}

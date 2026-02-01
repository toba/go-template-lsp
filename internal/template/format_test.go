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
			name:     "blank lines preserved",
			input:    "<div>\n\n<p>text</p>\n\n</div>",
			opts:     FormatOptions{TabSize: 2, InsertSpaces: true},
			expected: "<div>\n\n  <p>text</p>\n\n</div>",
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
   class="map-zoom-btn"
   title="Expand">`,
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
  width="200"
  height="100" />`,
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
  placeholder="Enter name"
  required>`,
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
  id="my-id"
  data-value="something">`,
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

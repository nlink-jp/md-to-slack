package converter

import (
	"encoding/json"
	"testing"

	"github.com/nlink-jp/md-to-slack/internal/blocks"
)

// helper to marshal to JSON for comparison.
func toJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent: %v", err)
	}
	return string(b)
}

func TestHeadings(t *testing.T) {
	src := []byte("# Title\n\n## Subtitle\n\n### H3\n")
	blks := Convert(src)
	if len(blks) != 3 {
		t.Fatalf("expected 3 blocks, got %d", len(blks))
	}
	h1, ok := blks[0].(blocks.Header)
	if !ok {
		t.Fatalf("block[0] is %T, want Header", blks[0])
	}
	if h1.Text.Text != "Title" {
		t.Errorf("H1 text = %q, want %q", h1.Text.Text, "Title")
	}
	h2, ok := blks[1].(blocks.Header)
	if !ok {
		t.Fatalf("block[1] is %T, want Header", blks[1])
	}
	if h2.Text.Text != "Subtitle" {
		t.Errorf("H2 text = %q, want %q", h2.Text.Text, "Subtitle")
	}
	s3, ok := blks[2].(blocks.Section)
	if !ok {
		t.Fatalf("block[2] is %T, want Section", blks[2])
	}
	if s3.Text.Text != "*### H3*" {
		t.Errorf("H3 text = %q, want %q", s3.Text.Text, "*### H3*")
	}
}

func TestParagraph(t *testing.T) {
	src := []byte("Hello **world** and _italic_ text.\n")
	blks := Convert(src)
	if len(blks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blks))
	}
	s, ok := blks[0].(blocks.Section)
	if !ok {
		t.Fatalf("block[0] is %T, want Section", blks[0])
	}
	want := "Hello *world* and _italic_ text."
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestCodeBlock(t *testing.T) {
	src := []byte("```go\nfunc main() {}\n```\n")
	blks := Convert(src)
	if len(blks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blks))
	}
	s, ok := blks[0].(blocks.Section)
	if !ok {
		t.Fatalf("block[0] is %T, want Section", blks[0])
	}
	want := "```go\nfunc main() {}\n```"
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestFencedCodeBlockNoLang(t *testing.T) {
	src := []byte("```\nplain code\n```\n")
	blks := Convert(src)
	if len(blks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blks))
	}
	s := blks[0].(blocks.Section)
	want := "```\nplain code\n```"
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestThematicBreak(t *testing.T) {
	src := []byte("---\n")
	blks := Convert(src)
	if len(blks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blks))
	}
	_, ok := blks[0].(blocks.Divider)
	if !ok {
		t.Fatalf("block[0] is %T, want Divider", blks[0])
	}
}

func TestLink(t *testing.T) {
	src := []byte("[Slack](https://slack.com)\n")
	blks := Convert(src)
	s := blks[0].(blocks.Section)
	want := "<https://slack.com|Slack>"
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestLinkLabelEqualsURL(t *testing.T) {
	src := []byte("[https://slack.com](https://slack.com)\n")
	blks := Convert(src)
	s := blks[0].(blocks.Section)
	want := "<https://slack.com>"
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestAutoLink(t *testing.T) {
	src := []byte("<https://example.com>\n")
	blks := Convert(src)
	s := blks[0].(blocks.Section)
	want := "<https://example.com>"
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestStandaloneImage(t *testing.T) {
	src := []byte("![alt text](https://example.com/img.png)\n")
	blks := Convert(src)
	if len(blks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blks))
	}
	img, ok := blks[0].(blocks.Image)
	if !ok {
		t.Fatalf("block[0] is %T, want Image", blks[0])
	}
	if img.ImageURL != "https://example.com/img.png" {
		t.Errorf("ImageURL = %q", img.ImageURL)
	}
	if img.AltText != "alt text" {
		t.Errorf("AltText = %q", img.AltText)
	}
}

func TestInlineImage(t *testing.T) {
	// Image alongside text → section block with link fallback.
	src := []byte("See ![logo](https://example.com/logo.png) here.\n")
	blks := Convert(src)
	s, ok := blks[0].(blocks.Section)
	if !ok {
		t.Fatalf("block[0] is %T, want Section", blks[0])
	}
	if s.Text.Text == "" {
		t.Error("expected non-empty section text")
	}
}

func TestBlockquote(t *testing.T) {
	src := []byte("> quoted text\n")
	blks := Convert(src)
	s := blks[0].(blocks.Section)
	want := "> quoted text"
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestUnorderedList(t *testing.T) {
	src := []byte("- apple\n- banana\n- cherry\n")
	blks := Convert(src)
	s := blks[0].(blocks.Section)
	want := "• apple\n• banana\n• cherry"
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestOrderedList(t *testing.T) {
	src := []byte("1. first\n2. second\n3. third\n")
	blks := Convert(src)
	s := blks[0].(blocks.Section)
	want := "1. first\n2. second\n3. third"
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestNestedList(t *testing.T) {
	src := []byte("- a\n  - b\n  - c\n- d\n")
	blks := Convert(src)
	s := blks[0].(blocks.Section)
	if s.Text.Text == "" {
		t.Error("expected non-empty section for nested list")
	}
}

func TestStrikethrough(t *testing.T) {
	src := []byte("~~deleted~~\n")
	blks := Convert(src)
	s := blks[0].(blocks.Section)
	want := "~deleted~"
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestCodeSpan(t *testing.T) {
	src := []byte("Use `fmt.Println` to print.\n")
	blks := Convert(src)
	s := blks[0].(blocks.Section)
	want := "Use `fmt.Println` to print."
	if s.Text.Text != want {
		t.Errorf("text = %q, want %q", s.Text.Text, want)
	}
}

func TestTable(t *testing.T) {
	src := []byte("| Name | Age |\n|------|-----|\n| Alice | 30 |\n| Bob | 25 |\n")
	blks := Convert(src)
	if len(blks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blks))
	}
	s, ok := blks[0].(blocks.Section)
	if !ok {
		t.Fatalf("block[0] is %T, want Section", blks[0])
	}
	if !startsWith(s.Text.Text, "```") {
		t.Errorf("table should be in a code block, got: %q", s.Text.Text)
	}
}

func TestHTMLBlockDiscarded(t *testing.T) {
	src := []byte("<div>ignored</div>\n\nsome text\n")
	blks := Convert(src)
	// HTML block discarded; only the paragraph should remain.
	if len(blks) != 1 {
		t.Fatalf("expected 1 block, got %d: %v", len(blks), toJSON(t, blks))
	}
}

func TestConvertToJSON(t *testing.T) {
	src := []byte("# Hello\n\nWorld.\n")
	b, err := ConvertToJSON(src)
	if err != nil {
		t.Fatalf("ConvertToJSON: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(b, &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	blks, ok := payload["blocks"].([]any)
	if !ok || len(blks) != 2 {
		t.Fatalf("expected 2 blocks in JSON, got: %s", b)
	}
}

func TestEmptyInput(t *testing.T) {
	blks := Convert([]byte(""))
	if len(blks) != 0 {
		t.Errorf("expected 0 blocks for empty input, got %d", len(blks))
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

package converter

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"

	"github.com/nlink-jp/md-to-slack/internal/blocks"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
)

// Convert parses src as Markdown and returns the equivalent Slack Block Kit blocks.
func Convert(src []byte) []any {
	reader := text.NewReader(src)
	doc := md.Parser().Parse(reader)
	return walkDocument(doc, src)
}

// walkDocument iterates over the top-level children of the document node and
// converts each one to one or more Slack blocks.
func walkDocument(doc ast.Node, src []byte) []any {
	var result []any
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		result = append(result, blockOf(node, src)...)
	}
	return result
}

// blockOf converts a single block-level AST node to zero or more Slack blocks.
func blockOf(node ast.Node, src []byte) []any {
	switch n := node.(type) {

	case *ast.Heading:
		text := mrkdwnOf(n, src)
		if n.Level <= 2 {
			return []any{blocks.NewHeader(text)}
		}
		// H3–H6: bold prefix in a section block.
		prefix := strings.Repeat("#", n.Level) + " "
		return []any{blocks.NewSection("*" + prefix + text + "*")}

	case *ast.Paragraph:
		// Single standalone image → image block.
		if img, ok := singleImage(n); ok {
			url := string(img.Destination)
			alt := mrkdwnOf(img, src)
			if alt == "" {
				alt = url
			}
			title := string(img.Title)
			return []any{blocks.NewImage(url, alt, title)}
		}
		text := mrkdwnOf(n, src)
		if text == "" {
			return nil
		}
		return []any{blocks.NewSection(text)}

	case *ast.FencedCodeBlock:
		var sb strings.Builder
		sb.WriteString("```")
		if n.Info != nil {
			lang := string(n.Info.Segment.Value(src))
			// Strip trailing whitespace/newline from lang tag.
			lang = strings.TrimSpace(lang)
			if lang != "" {
				sb.WriteString(lang)
			}
		}
		sb.WriteByte('\n')
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			sb.Write(line.Value(src))
		}
		sb.WriteString("```")
		return []any{blocks.NewSection(sb.String())}

	case *ast.CodeBlock:
		var sb strings.Builder
		sb.WriteString("```\n")
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			sb.Write(line.Value(src))
		}
		sb.WriteString("```")
		return []any{blocks.NewSection(sb.String())}

	case *ast.Blockquote:
		// Render blockquote children as mrkdwn then prefix each line with "> ".
		inner := renderBlockquoteChildren(n, src)
		text := blockquoteLines(inner)
		if text == "" {
			return nil
		}
		return []any{blocks.NewSection(text)}

	case *ast.List:
		text := listToMrkdwn(n, src, 0)
		if text == "" {
			return nil
		}
		return []any{blocks.NewSection(text)}

	case *ast.ThematicBreak:
		return []any{blocks.NewDivider()}

	case *ast.HTMLBlock:
		// Discard raw HTML blocks — not renderable in Slack.
		return nil

	default:
		// GFM table.
		if node.Kind() == east.KindTable {
			text := tableToMrkdwn(node, src)
			if text == "" {
				return nil
			}
			return []any{blocks.NewSection(text)}
		}
		// Any other unknown block: try to extract inline text.
		text := mrkdwnOf(node, src)
		if text == "" {
			return nil
		}
		return []any{blocks.NewSection(text)}
	}
}

// singleImage reports whether a paragraph contains exactly one child that is
// an image, with no other visible content.
func singleImage(para *ast.Paragraph) (*ast.Image, bool) {
	var img *ast.Image
	for child := para.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *ast.Image:
			if img != nil {
				return nil, false // more than one image
			}
			img = c
		default:
			// Any other node (text, emphasis, etc.) means it is not standalone.
			_ = c
			return nil, false
		}
	}
	return img, img != nil
}

// renderBlockquoteChildren converts the children of a blockquote to a plain
// mrkdwn string. Each child block is rendered in turn, separated by newlines.
func renderBlockquoteChildren(bq ast.Node, src []byte) string {
	var parts []string
	for child := bq.FirstChild(); child != nil; child = child.NextSibling() {
		switch n := child.(type) {
		case *ast.Paragraph:
			parts = append(parts, mrkdwnOf(n, src))
		case *ast.List:
			parts = append(parts, listToMrkdwn(n, src, 0))
		case *ast.FencedCodeBlock, *ast.CodeBlock:
			blks := blockOf(n, src)
			for _, b := range blks {
				if s, ok := b.(blocks.Section); ok {
					parts = append(parts, s.Text.Text)
				}
			}
		case *ast.Blockquote:
			// Nested blockquote: recurse.
			inner := renderBlockquoteChildren(n, src)
			parts = append(parts, blockquoteLines(inner))
		default:
			t := mrkdwnOf(child, src)
			if t != "" {
				parts = append(parts, t)
			}
		}
	}
	return strings.Join(parts, "\n")
}

// ConvertToJSON parses src as Markdown and returns the Block Kit JSON bytes.
func ConvertToJSON(src []byte) ([]byte, error) {
	blks := Convert(src)
	payload := blocks.BlockKit{Blocks: blks}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

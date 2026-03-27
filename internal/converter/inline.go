package converter

import (
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

// mrkdwnOf converts an inline AST subtree to a Slack mrkdwn string.
// It is safe to call with a nil node (returns "").
func mrkdwnOf(node ast.Node, src []byte) string {
	if node == nil {
		return ""
	}
	var sb strings.Builder
	writeInline(&sb, node, src)
	return sb.String()
}

// writeInline writes the mrkdwn of every child of node into sb.
func writeInline(sb *strings.Builder, node ast.Node, src []byte) {
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		writeNode(sb, child, src)
	}
}

// writeNode writes one AST node (and its subtree) as mrkdwn.
func writeNode(sb *strings.Builder, node ast.Node, src []byte) {
	switch n := node.(type) {

	case *ast.Text:
		sb.Write(n.Segment.Value(src))
		if n.SoftLineBreak() || n.HardLineBreak() {
			sb.WriteByte('\n')
		}

	case *ast.String:
		sb.Write(n.Value)

	case *ast.CodeSpan:
		sb.WriteByte('`')
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if t, ok := c.(*ast.Text); ok {
				sb.Write(t.Segment.Value(src))
			}
		}
		sb.WriteByte('`')

	case *ast.Emphasis:
		marker := "_" // level 1 → italic
		if n.Level == 2 {
			marker = "*" // level 2 → bold
		}
		sb.WriteString(marker)
		writeInline(sb, n, src)
		sb.WriteString(marker)

	case *east.Strikethrough:
		sb.WriteByte('~')
		writeInline(sb, n, src)
		sb.WriteByte('~')

	case *ast.Link:
		url := string(n.Destination)
		label := mrkdwnOf(n, src)
		if label == "" || label == url {
			sb.WriteString("<" + url + ">")
		} else {
			sb.WriteString("<" + url + "|" + label + ">")
		}

	case *ast.AutoLink:
		url := string(n.URL(src))
		sb.WriteString("<" + url + ">")

	case *ast.Image:
		// Inline image (inside a paragraph alongside other content) → link fallback.
		url := string(n.Destination)
		alt := mrkdwnOf(n, src)
		if alt == "" {
			alt = url
		}
		sb.WriteString("<" + url + "|" + alt + ">")

	case *ast.RawHTML:
		// Discard raw HTML — not renderable in Slack.

	default:
		// Unknown node: recurse to preserve any nested text content.
		writeInline(sb, node, src)
	}
}

// blockquoteLines prefixes every line of text with "> ".
func blockquoteLines(text string) string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			out = append(out, ">")
		} else {
			out = append(out, "> "+line)
		}
	}
	// Trim trailing empty blockquote markers.
	for len(out) > 0 && out[len(out)-1] == ">" {
		out = out[:len(out)-1]
	}
	return strings.Join(out, "\n")
}

// listToMrkdwn converts a list node to a mrkdwn string at the given nesting depth.
func listToMrkdwn(list *ast.List, src []byte, depth int) string {
	indent := strings.Repeat("  ", depth)
	var sb strings.Builder
	orderedIdx := list.Start

	for item := list.FirstChild(); item != nil; item = item.NextSibling() {
		li, ok := item.(*ast.ListItem)
		if !ok {
			continue
		}

		var bullet string
		if list.IsOrdered() {
			bullet = fmt.Sprintf("%d.", orderedIdx)
			orderedIdx++
		} else {
			bullet = "•"
		}

		// Collect inline content from the first text block of the item.
		var itemText string
		for child := li.FirstChild(); child != nil; child = child.NextSibling() {
			switch child.(type) {
			case *ast.TextBlock, *ast.Paragraph:
				itemText = mrkdwnOf(child, src)
			}
		}

		sb.WriteString(indent + bullet + " " + itemText + "\n")

		// Nested lists.
		for child := li.FirstChild(); child != nil; child = child.NextSibling() {
			if nested, ok := child.(*ast.List); ok {
				sb.WriteString(listToMrkdwn(nested, src, depth+1))
				sb.WriteByte('\n')
			}
		}
	}

	return strings.TrimRight(sb.String(), "\n")
}

// tableToMrkdwn converts a GFM table node to a preformatted mrkdwn string.
// Slack does not natively render markdown tables in chat messages, so we
// output a plain-text approximation in a code block with aligned columns.
//
// goldmark GFM table structure:
//
//	Table
//	  TableHeader  (cells are direct children — no intermediate TableRow)
//	    TableCell
//	    TableCell ...
//	  TableRow     (body rows are direct children of Table)
//	    TableCell
//	    ...
func tableToMrkdwn(table ast.Node, src []byte) string {
	type tableRow struct {
		cells    []string
		isHeader bool
	}

	var rows []tableRow

	for child := table.FirstChild(); child != nil; child = child.NextSibling() {
		switch child.Kind() {
		case east.KindTableHeader:
			// Header cells are direct children of TableHeader.
			var cells []string
			for cell := child.FirstChild(); cell != nil; cell = cell.NextSibling() {
				cells = append(cells, mrkdwnOf(cell, src))
			}
			if len(cells) > 0 {
				rows = append(rows, tableRow{cells: cells, isHeader: true})
			}
		case east.KindTableRow:
			// Body rows: cells are direct children of TableRow.
			var cells []string
			for cell := child.FirstChild(); cell != nil; cell = cell.NextSibling() {
				cells = append(cells, mrkdwnOf(cell, src))
			}
			if len(cells) > 0 {
				rows = append(rows, tableRow{cells: cells, isHeader: false})
			}
		}
	}

	if len(rows) == 0 {
		return ""
	}

	// Calculate max column count and per-column widths.
	cols := 0
	for _, r := range rows {
		if len(r.cells) > cols {
			cols = len(r.cells)
		}
	}
	widths := make([]int, cols)
	for _, r := range rows {
		for i, cell := range r.cells {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("```\n")
	for i, r := range rows {
		parts := make([]string, cols)
		for ci := 0; ci < cols; ci++ {
			cell := ""
			if ci < len(r.cells) {
				cell = r.cells[ci]
			}
			parts[ci] = cell + strings.Repeat(" ", widths[ci]-len(cell))
		}
		sb.WriteString(strings.Join(parts, "  ") + "\n")
		// Separator after header row.
		if r.isHeader && i+1 < len(rows) {
			seps := make([]string, cols)
			for ci, w := range widths {
				seps[ci] = strings.Repeat("-", w)
			}
			sb.WriteString(strings.Join(seps, "  ") + "\n")
		}
	}
	sb.WriteString("```")
	return sb.String()
}

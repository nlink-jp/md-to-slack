# md-to-slack

Convert Markdown to [Slack Block Kit](https://api.slack.com/block-kit) JSON.

Reads Markdown from stdin and writes a `{"blocks": [...]}` JSON payload to stdout,
suitable for use with the Slack API or tools like [slackcat](https://github.com/bcicen/slackcat).

## Features

- GFM (GitHub Flavored Markdown) support
- H1/H2 → Slack header blocks (large font)
- H3–H6 → bold section blocks
- Paragraphs, blockquotes, ordered and unordered lists (nested)
- Fenced and indented code blocks
- GFM tables → aligned plain-text code blocks
- Standalone images → Slack image blocks
- Inline images → link fallback
- Strikethrough, bold, italic, inline code, links, autolinks
- Thematic breaks → Slack divider blocks
- Raw HTML discarded (not renderable in Slack)

## Installation

```bash
go install github.com/nlink-jp/md-to-slack/cmd/md-to-slack@latest
```

Or download a pre-built binary from the [Releases](https://github.com/nlink-jp/md-to-slack/releases) page.

## Usage

```bash
md-to-slack < README.md
echo "# Hello **world**" | md-to-slack
```

### Send to Slack

```bash
md-to-slack < message.md | curl -s \
  -X POST https://slack.com/api/chat.postMessage \
  -H "Authorization: Bearer $SLACK_TOKEN" \
  -H "Content-Type: application/json" \
  -d @- -d '{"channel":"#general"}'
```

### Flags

| Flag | Description |
|------|-------------|
| `--version`, `-V` | Print version and exit |
| `--help`, `-h` | Print usage and exit |

## Building

```bash
make build       # build for current platform
make build-all   # cross-compile for all platforms
make test        # run tests
make check       # vet + lint + test + build
```

## Output format

`md-to-slack` writes a Slack Block Kit payload:

```json
{
  "blocks": [
    { "type": "header", "text": { "type": "plain_text", "text": "Title", "emoji": true } },
    { "type": "section", "text": { "type": "mrkdwn", "text": "Body text." } }
  ]
}
```

## Documentation

- [日本語 README](README.ja.md)

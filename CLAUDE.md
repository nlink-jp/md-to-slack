# md-to-slack — CLAUDE.md

Project-specific instructions for Claude Code.

## Architecture

```
cmd/md-to-slack/main.go      CLI entry point (stdin → stdout, --version, --help)
internal/blocks/types.go     Slack Block Kit types (Header, Section, Divider, Image)
internal/converter/
  inline.go                  Inline/mrkdwn conversion: mrkdwnOf, listToMrkdwn, tableToMrkdwn
  converter.go               Block-level walker: Convert(), ConvertToJSON()
  converter_test.go          Unit tests
```

## Key decisions

- **goldmark GFM** for parsing: reliable, well-tested, extensible.
- **H1/H2 → Header block**, H3–H6 → bold Section block (Slack header only supports one size).
- **Standalone image** (sole child of paragraph) → Image block; mixed content → Section with link fallback.
- **GFM tables** → preformatted code block (Slack chat does not render Markdown tables).
- **Raw HTML** discarded silently — not renderable in Slack.
- `goldmark` table structure: `TableHeader` has `TableCell` as direct children (no `TableRow`
  wrapper); body rows are `TableRow` direct children of `Table`.

## Testing

```
make test        # go test ./...
make check       # vet + lint + test + build
```

## Shared conventions

See `../CONVENTIONS.md` (cli-series umbrella repo).

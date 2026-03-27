# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.1.1] - 2026-03-28

### Changed

- Updated README example to use `scli` instead of `slackcat` for Slack posting.

### Fixed

- Fixed `.gitignore` pattern `md-to-slack` that was inadvertently excluding the
  `cmd/md-to-slack/` source directory; changed to `/md-to-slack` to match only the
  root-level compiled binary.
- Added missing `cmd/md-to-slack/main.go` (CLI entry point) to the repository.

### Internal

- Added macOS-specific entries to `.gitignore`.

## [0.1.0] - 2026-03-27

### Added

- Initial release: Markdown to Slack Block Kit JSON converter.
- GFM support: headings, paragraphs, fenced code blocks, blockquotes, ordered and
  unordered lists (including nested), thematic breaks, tables, strikethrough,
  inline code, bold, italic, links, autolinks, and images.
- H1/H2 headings render as Slack header blocks (large font).
- H3–H6 headings render as bold section blocks.
- Standalone images render as Slack image blocks; inline images fall back to links.
- GFM tables render as preformatted plain-text code blocks with aligned columns.
- Raw HTML blocks and inline HTML are silently discarded.
- CLI reads Markdown from stdin and writes Block Kit JSON to stdout.
- `--version` / `-V` and `--help` / `-h` flags.

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.2.0] - 2026-07-12

### Added

- **`LICENSE` file (MIT).** The repository previously had no license file;
  it is now MIT-licensed and the license is bundled in every release archive.

### Removed

- **darwin/amd64 (Intel) pre-built binary.** macOS releases now ship
  **arm64 only**, per the org-wide policy (darwin is Apple-Silicon only; no
  universal binaries). Intel Mac users can build from source.

### Changed

- **Linux release archives are now `.tar.gz`** (darwin/windows remain `.zip`),
  per `nlink-jp/.github` CONVENTIONS.md §Release Archive Standard. Archives
  now bundle `README.md` + `LICENSE` alongside the canonical binary.
- **darwin code-signature identifier** is now the canonical `md-to-slack`
  (was `md-to-slack-darwin-arm64`), set via `codesign -i` so it stays stable
  after the archived binary is renamed to its canonical name.

No change to the binary's behaviour — a packaging / build-config release.

## [0.1.2] - 2026-05-22

### Added

- **`package` Makefile target.** Builds all 5 platforms, signs darwin
  binaries with Developer ID, zips each with README.md, and
  notarizes the darwin zips. Replaces the manual zip step the
  previous v0.1.1 release relied on.

### Changed

- **Darwin releases are now Developer ID signed and Apple-notarized.**
  `md-to-slack-v0.1.2-darwin-{amd64,arm64}.zip` carry full Apple
  Developer ID Application signatures and notarization tickets from
  Apple. End users on macOS no longer need to bypass Gatekeeper
  with right-click → Open or `xattr -d com.apple.quarantine` on
  first launch; local users who place `md-to-slack` under
  Dropbox-synced (or any other FileProvider-managed) paths are no
  longer killed by macOS's ad-hoc + provenance distrust policy.
  Pipeline: `scripts/codesign-darwin.sh` +
  `scripts/notarize-darwin.sh`, driven by `make package`. Adopts
  the org-wide convention in `nlink-jp/.github` CONVENTIONS.md
  §Code Signing.
- **Release zip filenames now embed the version**
  (`md-to-slack-vX.Y.Z-<os>-<arch>.zip`), aligning with the
  sibling chatops-series tools (scat, swrite, stail). Previous
  v0.1.1 assets used version-less names.

No behaviour change to the binary itself — feature-wise this is
identical to v0.1.1.

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

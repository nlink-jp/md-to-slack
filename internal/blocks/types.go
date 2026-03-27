// Package blocks defines the Slack Block Kit types produced by md-to-slack.
//
// Only the subset of Block Kit used for Markdown conversion is represented here:
// header, section, divider, and image blocks.
package blocks

// BlockKit is the top-level Slack Block Kit payload.
type BlockKit struct {
	Blocks []any `json:"blocks"`
}

// Header is a Slack header block (plain text, large font, no formatting).
// Used for H1 and H2 Markdown headings.
type Header struct {
	Type string    `json:"type"` // always "header"
	Text PlainText `json:"text"`
}

// Section is a Slack section block (mrkdwn or plain_text body).
// Used for paragraphs, H3–H6 headings, blockquotes, code blocks, and lists.
type Section struct {
	Type string `json:"type"` // always "section"
	Text Mrkdwn `json:"text"`
}

// Divider is a Slack divider block (horizontal rule).
// Used for Markdown thematic breaks (---).
type Divider struct {
	Type string `json:"type"` // always "divider"
}

// Image is a Slack image block.
// Used when a Markdown paragraph consists of a single image.
type Image struct {
	Type     string     `json:"type"` // always "image"
	ImageURL string     `json:"image_url"`
	AltText  string     `json:"alt_text"`
	Title    *PlainText `json:"title,omitempty"`
}

// PlainText is a text object with type "plain_text".
type PlainText struct {
	Type  string `json:"type"`  // always "plain_text"
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"` // true enables emoji shortcodes in Slack
}

// Mrkdwn is a text object with type "mrkdwn".
type Mrkdwn struct {
	Type string `json:"type"` // always "mrkdwn"
	Text string `json:"text"`
}

// NewHeader returns a header block with emoji enabled.
func NewHeader(text string) Header {
	return Header{
		Type: "header",
		Text: PlainText{Type: "plain_text", Text: text, Emoji: true},
	}
}

// NewSection returns a section block with mrkdwn body text.
func NewSection(text string) Section {
	return Section{
		Type: "section",
		Text: Mrkdwn{Type: "mrkdwn", Text: text},
	}
}

// NewDivider returns a divider block.
func NewDivider() Divider {
	return Divider{Type: "divider"}
}

// NewImage returns an image block. altText must not be empty (Slack API requirement).
func NewImage(url, altText, title string) Image {
	img := Image{
		Type:     "image",
		ImageURL: url,
		AltText:  altText,
	}
	if title != "" {
		img.Title = &PlainText{Type: "plain_text", Text: title, Emoji: true}
	}
	return img
}

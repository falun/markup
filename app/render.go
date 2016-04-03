package app

import (
	"bytes"

	"github.com/russross/blackfriday"
)

type FileResolver interface {
	Resolve(string) bool
}

type MarkupRenderer struct {
	blackfriday.Renderer
	extensions int
	res        FileResolver
}

var _ blackfriday.Renderer = MarkupRenderer{}

func NewRenderer(rslv FileResolver) MarkupRenderer {
	return MarkupRenderer{
		Renderer:   blackfriday.HtmlRenderer(blackfridayHTMLFlags, "", ""),
		extensions: blackfridayExtensions,
		res:        rslv,
	}
}

func (m MarkupRenderer) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	m.Renderer.AutoLink(out, link, kind)
}

func (m MarkupRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	m.Renderer.Link(out, link, title, content)
}

func (r MarkupRenderer) Render(in []byte) []byte {
	return blackfriday.MarkdownOptions(in, r, blackfridayOpts)
}

const (
	blackfridayHTMLFlags int = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	blackfridayExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS
)

var blackfridayOpts = blackfriday.Options{Extensions: blackfridayExtensions}

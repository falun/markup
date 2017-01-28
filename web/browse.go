package web

type MarkdownView struct {
	Host       string
	Title      string
	Style      string
	IndexPath  string
	FileString string
	IsMarkdown bool
	IsHTML     bool
}

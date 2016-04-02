package assets

import "fmt"

const pageTemplate = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
<title>%s</title>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<style type="text/css">
%s
</style>
<style>body{box-sizing:border-box;min-width:200px;max-width:980px;margin:0 auto;padding:45px;}</style>
</head>
<body class="markdown-body">
%s
</body>
</html>`

// GetPageHTML renders a given bit of markdown text as HTML.
//
//   markdown: the markdown to be rendered
//   titlep: If set this is the title of the page; otherwise the empty string is used
//   cssp: If set this is used as the CSS for this page; otherwise assets.Css is used
func GetPageHTML(markdown []byte, titlep, cssp *string) []byte {
	title := ""
	css := Css

	if cssp != nil {
		css = *cssp
	}

	if titlep != nil {
		title = *titlep
	}

	return []byte(fmt.Sprintf(pageTemplate, title, css, markdown))
}

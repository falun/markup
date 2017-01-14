package assets

import (
	"bytes"
	"fmt"
	"text/template"
)

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
<body>
<div class="banner"><a href="http://%s/%s">Index</a></div>
<div class="markdown-body">
%s
</div>
</body>
</html>`

// GetPageHTML renders a given bit of markdown text as HTML.
//
//   hosted: The host:port combination where markup is running
//   indexPath: 'index/<path>' to the directory hosting the file being rendered
//   markdown: the markdown to be rendered
//   titlep: If set this is the title of the page; otherwise the empty string is used
//   cssp: If set this is used as the CSS for this page; otherwise assets.Css is used
func GetPageHTML(hosted, indexPath string, markdown []byte, titlep, cssp *string) []byte {
	title := ""
	css := Css

	if cssp != nil {
		css = *cssp
	}

	if titlep != nil {
		title = *titlep
	}

	return []byte(fmt.Sprintf(pageTemplate, title, css, hosted, indexPath, markdown))
}

const indexTemplate = `
<html>
	<head>
		<title>Index Listing {{.Root}}</title>
	</head>

	<body>
		{{ with $dir := . }}
			<h2>Index: {{$dir.Root}}</h2>
			<ul>
				{{ if ne $dir.Root "/" }}<li><a href="/index/{{$dir.Root}}..">../</a></li>{{end}}
				{{ range $child := $dir.Children }}
					<li>
						{{ if $child.Dir }}
							<a href="/index/{{$dir.Root}}{{$child.Name}}">{{ $child.Name }}/</a>
						{{ else }}
							<a href="/browse/{{$dir.Root}}{{$child.Name}}">{{ $child.Name }}</a>
						{{ end }}
					</li>
				{{ end }}
			{{ end }}
		</hl>
	</body>
</html>
`

var indexTmpl = template.Must(template.New("index_page").Parse(indexTemplate))

// src should be an app.Dir; this is to avoid circular import dependency because
// I'm too lazy to factor things out correctly right now.
func RenderIndexPage(src interface{}) (string, error) {
	buffer := bytes.NewBuffer(nil)

	err := indexTmpl.Execute(buffer, src)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

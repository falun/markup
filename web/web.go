package web

import (
	"bytes"
	"os"
	"log"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"

	"github.com/falun/markup/markdown"
)

const (
	IndexTmpl    = "index"
	MarkdownTmpl = "browse"
)

type Renderer interface {
	Render(tmplName string, arg interface{}) ([]byte, error)
}

type renderer struct {
	assetDir string
}

func renderStr(str string) string {
	return string(MdRenderer.Render([]byte(str)))
}

func mkTmpl(n, tmpl string) (*template.Template, error) {
	t := template.New(string(n))
	t = t.Funcs(template.FuncMap{
		"markdown": renderStr,
	})

	return t.Parse(tmpl)
}

var (
	MdRenderer = markdown.NewRenderer()

	Templates = map[string]*template.Template{
		IndexTmpl:    template.Must(mkTmpl(IndexTmpl, IndexTmplString)),
		MarkdownTmpl: template.Must(mkTmpl(MarkdownTmpl, MarkdownTmplString)),
	}
)

func (r renderer) loadTmpl(name string) (*template.Template, error) {
	tmplPath := path.Join(r.assetDir, fmt.Sprintf("%s.tmpl", name))
	b, err := ioutil.ReadFile(tmplPath)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("could not load %s: %v", tmplPath, err)
	}

	var t *template.Template
	if b != nil {
		t, err = mkTmpl(name, string(b))
		if err != nil {
			t = nil
			log.Printf("could not parse %s, using built in fallback: %v", tmplPath, err)
		}
	}

	var ok bool
	if t == nil {
		t, ok = Templates[name]
		if !ok {
			return nil, fmt.Errorf("no built-in template found for %s", name)
		}
	}

	return t, nil
}

func (r renderer) Render(name string, arg interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	t, err := r.loadTmpl(name)
	if err != nil {
		return nil, fmt.Errorf("unable to load template %s: %v", name, err)
	}

	err = t.Execute(buf, arg)
	if err != nil {
		return nil, fmt.Errorf("could not render %v: %v", name, err)
	}

	return buf.Bytes(), nil
}

func NewRenderer(assetDir string) Renderer {
	return renderer{assetDir}
}

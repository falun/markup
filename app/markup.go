package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/falun/markup/assets"
)

type markupServer struct {
	root  string
	index string

	renderer MarkupRenderer
}

func (m markupServer) ResolvePath(p string) string {
	if p == "" {
		p = m.index
	}

	return path.Join(m.root, p)
}

type markupRequest struct{ *http.Request }

func (r markupRequest) IsRoot() bool {
	return "" == r.FilePath()
}

func (r markupRequest) IsHTML() bool {
	path := r.URL.Path
	return strings.HasSuffix(strings.ToLower(path), ".html") ||
		strings.HasSuffix(strings.ToLower(path), ".htm")
}

func (r markupRequest) IsMarkdown() bool {
	path := r.URL.Path
	return strings.HasSuffix(strings.ToLower(path), ".md")
}

func (r markupRequest) FilePath() string {
	p := r.URL.Path

	if p[0] == '/' {
		p = p[1:]
	}

	return p
}

func (srv markupServer) ServeHTTP(rw http.ResponseWriter, netReq *http.Request) {
	r := markupRequest{netReq}

	if r.IsRoot() {
		http.Redirect(rw, netReq, srv.index, http.StatusSeeOther)
		return
	}

	filepath := srv.ResolvePath(r.FilePath())
	fptr, err := os.Open(filepath)
	defer fptr.Close()

	if err != nil {
		rw.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(rw, "%s", err)
		return
	}

	b, err := ioutil.ReadAll(fptr)
	if err != nil {
		rw.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(rw, "%s", err)
		return
	}

	h := rw.Header()

	output := b
	ct := "text/plain"

	switch {
	case r.IsMarkdown():
		title := path.Base(filepath)
		ct = "text/html"
		md := srv.renderer.Render(b)
		output = assets.GetPageHTML(md, &title, nil)
	case r.IsHTML():
		ct = "text/html"
	}

	h["Content-Type"] = []string{ct}
	rw.WriteHeader(http.StatusOK)
	rw.Write(output)
}

func Main(cfg Config) {
	http.Handle("/", markupServer{
		root:     cfg.RootDir,
		index:    cfg.Index,
		renderer: NewRenderer(),
	})

	ingress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf(`
Listening on: %s
Serving from: %s
Default file: %s
`, ingress, cfg.RootDir, cfg.Index)

	err := http.ListenAndServe(ingress, nil)
	log.Fatal(err)
}

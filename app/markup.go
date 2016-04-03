package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/falun/markup/assets"
)

const browseToken = "browse"

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

func (r markupRequest) Browse() bool {
	return strings.HasPrefix(r.URL.Path, "/"+browseToken+"/")
}

func (r markupRequest) IsRoot() bool {
	uri := r.URL.Path
	return uri == "/"+browseToken ||
		uri == "/"+browseToken+"/" ||
		uri == "" ||
		uri == "/"
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

	cutlen := 2 + len(browseToken)
	p = p[cutlen:]

	return p
}

func (srv markupServer) ServeHTTP(rw http.ResponseWriter, netReq *http.Request) {
	r := markupRequest{netReq}

	if r.IsRoot() {
		http.Redirect(rw, netReq, "/"+browseToken+"/"+srv.index, http.StatusSeeOther)
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

type findFile struct {
	root string
}

func (ff findFile) Resolve(p string) bool {
	ue, err := url.QueryUnescape(p)
	if err != nil {
		log.Fatal(err)
	}

	fp := path.Join(ff.root, ue)
	fptr, err := os.Open(fp)
	defer fptr.Close()

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal(err)
	}

	return true
}

func Main(cfg Config) {
	http.Handle("/", markupServer{
		root:     cfg.RootDir,
		index:    cfg.Index,
		renderer: NewRenderer(findFile{cfg.RootDir}),
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

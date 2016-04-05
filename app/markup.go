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

const browseToken = "browse"

type browserHandler struct {
	root  string
	index string
	token string

	renderer MarkupRenderer
}

func (m browserHandler) ResolvePath(p string) string {
	return path.Join(m.root, p)
}

type browseRequest struct {
	*http.Request
	token string
}

func (r browseRequest) EndsWith(ext ...string) bool {
	path := strings.ToLower(r.URL.Path)

	for _, e := range ext {
		if strings.HasSuffix(path, strings.ToLower(e)) {
			return true
		}
	}

	return false
}

func (r browseRequest) IsHTML() bool {
	return r.EndsWith(".html", ".htm")
}

func (r browseRequest) IsMarkdown() bool {
	return r.EndsWith(".md")
}

func (r browseRequest) FilePath() string {
	p := r.URL.Path

	cutlen := 2 + len(r.token)
	p = p[cutlen:]

	return p
}

func (srv browserHandler) ServeHTTP(rw http.ResponseWriter, netReq *http.Request) {
	r := browseRequest{netReq, srv.token}

	// check for '/browse/' and serve index if requested
	reqPath := r.FilePath()
	if reqPath == "" {
		serveIndex("/"+srv.token+"/", srv.token, srv.index)(rw, netReq)
		return
	}

	// find the file being requested
	filepath := srv.ResolvePath(reqPath)
	fptr, err := os.Open(filepath)
	defer fptr.Close()

	if err != nil {
		if os.IsNotExist(err) {
			// not found if the file doesn't exist
			serveNotFound(rw, netReq)
		} else {
			// 503 on other errors
			serveError(rw, netReq, err)
		}
		return
	}

	// read contents of the file; return 500 on error
	b, err := ioutil.ReadAll(fptr)
	if err != nil {
		serveError(rw, netReq, err)
		return
	}

	h := rw.Header()

	output := b
	ct := "text/plain"

	// compute content type, render markdown, if applicable
	switch {
	case r.IsMarkdown():
		title := path.Base(filepath)
		ct = "text/html"
		md := srv.renderer.Render(b)
		output = assets.GetPageHTML(md, &title, nil)
	case r.IsHTML():
		ct = "text/html"
	}

	// set the mime type and write the results
	h["Content-Type"] = []string{ct}
	rw.WriteHeader(http.StatusOK)
	rw.Write(output)
}

func Main(cfg Config) {
	http.Handle("/", serveIndex("/", cfg.Token, cfg.Index))
	http.Handle(fmt.Sprintf("/%s/", cfg.Token), browserHandler{
		root:     cfg.RootDir,
		index:    cfg.Index,
		token:    cfg.Token,
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

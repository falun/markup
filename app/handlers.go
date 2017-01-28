package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/falun/markup/web"
)

func serveRoot(prefix, browseToken, index string) http.HandlerFunc {
	rootPage := fmt.Sprintf("/%s/%s", browseToken, index)
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == prefix {
			http.Redirect(rw, r, rootPage, http.StatusSeeOther)
		} else {
			serveNotFound(rw, r)
		}
	}
}

func serveIndex(
	renderer web.Renderer,
	indexToken string,
	browseToken string,
	root string,
	excluding []string,
) http.HandlerFunc {
	return indexHandler{
		renderer,
		indexToken,
		browseToken,
		root,
		excluding,
	}.ServeHTTP
}

func serveAsset(token, root string) http.HandlerFunc {
	i := 2 + len(token)
	asset := func(r *http.Request) string {
		return r.URL.Path[i:]
	}

	mime := web.NewMimeTypeProvider(root)

	return func(rw http.ResponseWriter, r *http.Request) {
		file := asset(r)
		if file == "" {
			serveError(rw, r, fmt.Errorf("No asset specified"))
			return
		}

		b, e := ioutil.ReadFile(path.Join(root, file))
		if e != nil {
			serveError(rw, r, e)
			return
		}
		h := rw.Header()
		h["Content-Type"] = []string{mime.ForFile(file)}
		rw.Write(b)
	}
}

var serveNotFound = http.NotFound

func serveError(rw http.ResponseWriter, r *http.Request, err error) {
	http.Error(rw, fmt.Sprintf("%v", err), http.StatusServiceUnavailable)
}

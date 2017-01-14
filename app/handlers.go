package app

import (
	"fmt"
	"net/http"
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
	indexToken string,
	browseToken string,
	root string,
	excluding []string,
) http.HandlerFunc {
	return indexHandler{
		indexToken, browseToken, root, excluding}.ServeHTTP
}

var serveNotFound = http.NotFound

func serveError(rw http.ResponseWriter, r *http.Request, err error) {
	http.Error(rw, fmt.Sprintf("%v", err), http.StatusServiceUnavailable)
}
